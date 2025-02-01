package event

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/infra/json"
	"bom-pedido-api/internal/infra/retry"
	"bom-pedido-api/internal/infra/telemetry"
	"bom-pedido-api/pkg/log"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

const exchange = "bompedido"

type RabbitMqEventHandler struct {
	connection      *amqp.Connection
	producerChannel *amqp.Channel
	consumerChannel *amqp.Channel
}

func NewRabbitMqEventHandler(server string) (*RabbitMqEventHandler, error) {
	connection, err := amqp.Dial(server)
	if err != nil {
		return nil, fmt.Errorf("create connection: %v", err)
	}

	producerChannel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("create producer channel: %v", err)
	}

	consumerChannel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("create consumer channel: %v", err)
	}

	return &RabbitMqEventHandler{
		connection:      connection,
		producerChannel: producerChannel,
		consumerChannel: consumerChannel,
	}, nil
}

func (r *RabbitMqEventHandler) Close() {
	_ = r.producerChannel.Close()
	_ = r.consumerChannel.Close()
	_ = r.connection.Close()
}

func (r *RabbitMqEventHandler) Emit(ctx context.Context, event *event.Event) error {
	ctx, span := telemetry.StartSpan(ctx, "RabbitMq.Emit")
	defer span.End()
	eventBytes, err := json.Marshal(ctx, event)
	if err != nil {
		log.Error("Error on emit event", err, "event", event)
		return err
	}
	err = r.producerChannel.PublishWithContext(
		ctx,
		exchange,
		event.Name,
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: eventBytes},
	)
	if err != nil {
		log.Error("Error on publish event", err, "event", event, "exchange", exchange)
		return err
	}
	return nil
}

func (r *RabbitMqEventHandler) Consume(options *event.ConsumerOptions, handler event.HandlerFunc) {
	messages, err := r.consumerChannel.Consume(
		options.Queue,
		"BOM_PEDIDO_API_"+options.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("Error on consume messages", err)
		return
	}

	for range options.WorkerPoolSize {
		go func(messages <-chan amqp.Delivery) {
			for message := range messages {
				r.handleMessage(message, handler, options.Id)
			}
		}(messages)
	}
}

func (r *RabbitMqEventHandler) handleMessage(message amqp.Delivery, handler event.HandlerFunc, name string) {
	ctx, span := telemetry.StartSpan(context.Background(), "RabbitMq.Process::"+name)
	defer span.End()
	log.Info("Mensagem recebida", "consumer", message.ConsumerTag, "routingKey", message.RoutingKey)
	var applicationEvent event.Event
	err := json.Unmarshal(ctx, message.Body, &applicationEvent)
	messageEvent := &event.MessageEvent{
		Event: &applicationEvent,
		AckFn: func(ctx context.Context) error {
			return message.Ack(false)
		},
		NackFn: func(ctx context.Context) {
			_ = message.Nack(false, true)
		},
	}
	defer func() {
		if err == nil {
			log.Info("Mensagem consumida com sucesso", "consumer", message.ConsumerTag, "routingKey", message.RoutingKey)
			return
		}
		span.RecordError(err)
		messageEvent.Nack(ctx)
		log.Error("Ocorreu um erro no consumo da mensagem", err, "consumer", message.ConsumerTag, "routingKey", message.RoutingKey)
	}()
	if err != nil {
		return
	}
	err = retry.Func(ctx, 5, time.Second, time.Second*30, func(ctx context.Context) error {
		return handler(ctx, messageEvent)
	})
}
