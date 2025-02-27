package event

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/infra/retry"
	"bom-pedido-api/pkg/log"
	"bom-pedido-api/pkg/telemetry"
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
	"time"
)

const (
	exchange                  = "bompedido"
	_rabbitMQEventHandlerName = "RABBITMQ"
)

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
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error("Error on emit event", err, "event", event)
		return err
	}

	err = telemetry.StartSpanReturningError(ctx, "RabbitMQ.Emit", func(ctx context.Context) error {
		headers := telemetry.GetPropagationHeaders(ctx)

		return r.producerChannel.PublishWithContext(
			ctx,
			exchange,
			string(event.Name),
			false,
			false,
			amqp.Publishing{
				ContentType: "text/json",
				Body:        eventBytes,
				Headers:     NewAmqpHeaders(headers),
			},
		)
	}, "event.name", string(event.Name))
	if err != nil {
		log.Error("Error on publish event", err, "event", event, "exchange", exchange)
		return err
	}
	return nil
}

func (r *RabbitMqEventHandler) OnEvent(eventName string, handlerFunc event.HandlerFunc) {
	if r.consumerChannel.IsClosed() {
		consumerChannel, err := r.connection.Channel()
		if err != nil {
			log.Error("Error on consume messages", err)
			return
		}
		r.consumerChannel = consumerChannel
	}
	messages, err := r.consumerChannel.Consume(
		eventName,
		"BOM_PEDIDO_API_"+eventName,
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

	for range 3 {
		go func(messages <-chan amqp.Delivery) {
			for message := range messages {
				r.handleMessage(message, handlerFunc)
			}
		}(messages)
	}
}

func (r *RabbitMqEventHandler) handleMessage(message amqp.Delivery, handler event.HandlerFunc) {
	ctx := telemetry.InjectPropagationHeaders(context.Background(), ParseAmqpTable(message.Headers))

	_ = telemetry.StartSpanReturningError(ctx, fmt.Sprintf("%s.HandleMessage", message.ConsumerTag), func(ctx context.Context) error {
		var applicationEvent event.Event
		err := json.Unmarshal(message.Body, &applicationEvent)
		messageEvent := &event.MessageEvent{
			Topic: strings.ReplaceAll(message.ConsumerTag, "BOM_PEDIDO_API_", ""),
			Event: &applicationEvent,
			AckFn: func(ctx context.Context) error {
				return message.Ack(false)
			},
			NackFn: func(ctx context.Context) {
				_ = message.Nack(false, true)
			},
		}
		defer func() {
			if err != nil {
				messageEvent.Nack(ctx)
				log.Error("Ocorreu um erro no consumo da mensagem", err, "consumer", message.ConsumerTag, "routingKey", message.RoutingKey)
			}
		}()
		if err != nil {
			return nil
		}

		err = retry.Func(ctx, 5, time.Second, time.Second*30, func(ctx context.Context) error {
			return handler(ctx, messageEvent)
		})

		return err
	})
}

func (r *RabbitMqEventHandler) Name() string {
	return _rabbitMQEventHandlerName
}

func NewAmqpHeaders(headers map[string]string) amqp.Table {
	table := amqp.Table{}
	for key, value := range headers {
		table[key] = value
	}
	return table
}

func ParseAmqpTable(headers amqp.Table) map[string]string {
	values := map[string]string{}
	for key, value := range headers {
		values[key] = value.(string)
	}
	return values
}
