package event

import (
	"bom-pedido-api/application/event"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

var eventExchanges map[string]string

func init() {
	eventExchanges = map[string]string{
		"PRODUCT_CREATED": "PRODUCTS",
		"ORDER_CREATED":   "ORDERS",
		"ORDER_APPROVED":  "ORDERS",
	}
}

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

func (adapter *RabbitMqAdapter) Emit(context context.Context, event *event.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		slog.Error("Error on emit event", "event", event, "error", err)
		return err
	}
	exchange := eventExchanges[event.Name]
	err = adapter.producerChannel.PublishWithContext(
		context,
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
		slog.Error("Error on consume messages", err)
		return
	}

	for range options.WorkerPoolSize {
		go adapter.handleMessages(messages, handler)
	}
}

func (adapter *RabbitMqAdapter) handleMessages(messages <-chan amqp.Delivery, handler event.HandlerFunc) {
	for message := range messages {
		messageEvent := &event.MessageEvent{
			AckFn: func() error {
				return message.Ack(false)
			},
			NackFn: func() {
				_ = message.Nack(false, true)
			},
			GetEventFn: func() *event.Event {
				var event event.Event
				_ = json.Unmarshal(message.Body, &event)
				return &event
			},
		}
		err := handler(messageEvent)
		messageEvent.NackIfError(err)
	}
}
