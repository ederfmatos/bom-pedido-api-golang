package event

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/infra/json"
	"bom-pedido-api/infra/telemetry"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/codes"
	"log/slog"
)

const exchange = "bompedido"

type RabbitMqAdapter struct {
	connection      *amqp.Connection
	producerChannel *amqp.Channel
	consumerChannel *amqp.Channel
}

func NewRabbitMqAdapter(server string) event.Handler {
	connection, err := amqp.Dial(server)
	if err != nil {
		panic(err)
	}
	slog.Info("Connected to rabbitmq successfully")
	producerChannel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	consumerChannel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	return &RabbitMqAdapter{
		connection:      connection,
		producerChannel: producerChannel,
		consumerChannel: consumerChannel,
	}
}

func (adapter *RabbitMqAdapter) Close() {
	adapter.producerChannel.Close()
	adapter.consumerChannel.Close()
	adapter.connection.Close()
}

func (adapter *RabbitMqAdapter) Emit(ctx context.Context, event *event.Event) error {
	ctx, span := telemetry.StartSpan(ctx, "RabbitMq.Emit")
	defer span.End()
	eventBytes, err := json.Marshal(ctx, event)
	if err != nil {
		slog.Error("Error on emit event", "event", event, "error", err)
		return err
	}
	err = adapter.producerChannel.PublishWithContext(
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

func (adapter *RabbitMqAdapter) Consume(options *event.ConsumerOptions, handler event.HandlerFunc) {
	_, err := adapter.consumerChannel.QueueDeclare(options.Queue, true, false, false, false, nil)
	if err != nil {
		slog.Error("Error on declare queue", "queue", options.Queue, "error", err)
	}
	messages, err := adapter.consumerChannel.Consume(
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
				adapter.handleMessage(message, handler)
			}
		}(messages)
	}
}

func (adapter *RabbitMqAdapter) handleMessage(message amqp.Delivery, handler event.HandlerFunc) {
	ctx, span := telemetry.StartSpan(context.Background(), "RabbitMq.Process")
	defer span.End()
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
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			messageEvent.Nack(ctx)
		}
	}()
	if err != nil {
		return
	}
	err = handler(ctx, messageEvent)
}
