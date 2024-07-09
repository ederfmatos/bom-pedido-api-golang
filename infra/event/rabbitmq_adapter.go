package event

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/domain/events"
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

func (adapter *RabbitMqAdapter) Close() {
	adapter.producerChannel.Close()
	adapter.consumerChannel.Close()
	adapter.connection.Close()
}

func NewRabbitMqAdapter(server string) event.Dispatcher {
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

func (adapter *RabbitMqAdapter) Emit(context context.Context, event *events.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		slog.Error("Error on emit event", "event", event, "error", err)
		return err
	}
	slog.Info("Emitting event", "event", event)
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
	slog.Info("Event emitted", "event", event)
	return nil
}

func (adapter *RabbitMqAdapter) Consume(options *event.ConsumerOptions, handler event.Handler) {
	_, err := adapter.consumerChannel.QueueDeclare(options.Id, true, false, false, false, nil)
	if err != nil {
		slog.Error("Error on declare queue", "queue", options.Id, "error", err)
	}
	messages, err := adapter.consumerChannel.Consume(
		options.Id,
		"BOM_PEDIDO_API_"+options.Id,
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

func (adapter *RabbitMqAdapter) handleMessages(messages <-chan amqp.Delivery, handler event.Handler) {
	for message := range messages {
		messageEvent := &RabbitMqMessageEvent{message}
		err := handler(messageEvent)
		if err != nil {
			messageEvent.Nack()
		}
	}
}

type RabbitMqMessageEvent struct {
	message amqp.Delivery
}

func (ev *RabbitMqMessageEvent) AckIfNoError(err error) error {
	if err == nil {
		return ev.Ack()
	}
	return err
}

func (ev *RabbitMqMessageEvent) Ack() error {
	return ev.message.Ack(false)
}

func (ev *RabbitMqMessageEvent) Nack() {
	_ = ev.message.Nack(false, true)
}

func (ev *RabbitMqMessageEvent) ParseData(event interface{}) {
	var messageEvent events.Event
	_ = json.Unmarshal(ev.message.Body, &messageEvent)
	data, _ := json.Marshal(messageEvent.Data)
	_ = json.Unmarshal(data, event)
}
