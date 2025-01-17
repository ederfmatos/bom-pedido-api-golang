package event

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/infra/json"
	"bom-pedido-api/internal/infra/retry"
	"bom-pedido-api/internal/infra/telemetry"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"os"
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

	eventHandler := &RabbitMqEventHandler{
		connection:      connection,
		producerChannel: producerChannel,
		consumerChannel: consumerChannel,
	}
	if err = eventHandler.createQueues(); err != nil {
		return nil, fmt.Errorf("create queues: %v", err)
	}

	return eventHandler, nil
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
		slog.Error("Error on emit event", "event", event, "error", err)
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
		slog.Error("Error on publish event", "event", event, "exchange", exchange, "error", err)
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
		slog.Error("Error on consume messages", "error", err)
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
	slog.Info("Mensagem recebida", "consumer", message.ConsumerTag, "routingKey", message.RoutingKey)
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
			slog.Info("Mensagem consumida com sucesso", "consumer", message.ConsumerTag, "routingKey", message.RoutingKey)
			return
		}
		span.RecordError(err)
		messageEvent.Nack(ctx)
		slog.Error("Ocorreu um erro no consumo da mensagem", "error", err, "consumer", message.ConsumerTag, "routingKey", message.RoutingKey)
	}()
	if err != nil {
		return
	}
	err = retry.Func(ctx, 5, time.Second, time.Second*30, func(ctx context.Context) error {
		return handler(ctx, messageEvent)
	})
}

type ResourceConfig struct {
	Exchanges []Exchange `json:"exchanges"`
	Queues    []Queue    `json:"queues"`
}

type Exchange struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Queue struct {
	Name     string `json:"name"`
	Bindings []struct {
		Exchange   string `json:"exchange"`
		RoutingKey string `json:"routingKey"`
	} `json:"bindings"`
	Arguments *struct {
		DeadLetterExchange   string `json:"x-dead-letter-exchange"`
		DeadLetterRoutingKey string `json:"x-dead-letter-routing-key"`
		MessageTtl           int32  `json:"x-message-ttl"`
	} `json:"arguments"`
}

// TODO: Criar script para criar fila fora da aplicação
func (r *RabbitMqEventHandler) createQueues() error {
	file, err := os.ReadFile(".resources/rabbitmq.json")
	if err != nil {
		return err
	}
	var resourceConfig ResourceConfig
	err = json.Unmarshal(context.Background(), file, &resourceConfig)
	if err != nil {
		return err
	}
	for _, item := range resourceConfig.Exchanges {
		err = r.consumerChannel.ExchangeDeclare(item.Name, item.Type, true, false, false, false, nil)
		if err != nil {
			return err
		}
	}
	for _, queue := range resourceConfig.Queues {
		var arguments amqp.Table
		if queue.Arguments != nil {
			arguments = amqp.Table{
				"x-dead-letter-exchange":    queue.Arguments.DeadLetterExchange,
				"x-dead-letter-routing-key": queue.Arguments.DeadLetterRoutingKey,
				"x-message-ttl":             queue.Arguments.MessageTtl,
			}
		}
		_, err = r.consumerChannel.QueueDeclare(queue.Name, true, false, false, false, arguments)
		if err != nil {
			return err
		}
		for _, binding := range queue.Bindings {
			err = r.consumerChannel.QueueBind(queue.Name, binding.RoutingKey, binding.Exchange, false, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
