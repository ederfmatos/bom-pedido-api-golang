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
	}
}

type RabbitMqAdapter struct {
	connection *amqp.Connection
}

func (adapter *RabbitMqAdapter) Close() {
	adapter.connection.Close()
}

func NewRabbitMqAdapter(server string) *RabbitMqAdapter {
	connection, err := amqp.Dial(server)
	if err != nil {
		panic(err)
	}
	return &RabbitMqAdapter{connection: connection}
}

func (adapter *RabbitMqAdapter) Emit(context context.Context, event *event.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		slog.Error("Error on emit event", "event", event, "error", err)
		return err
	}
	slog.Info("Emitting event", "event", event)
	exchange := eventExchanges[event.Name]
	channel, err := adapter.connection.Channel()
	if err != nil {
		slog.Error("Error on open rabbitmq channel", err)
		return err
	}
	err = channel.PublishWithContext(
		context,
		exchange,
		event.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        eventBytes,
		},
	)
	if err != nil {
		slog.Error("Error on publish event", "event", event, "exchange", exchange, "error", err)
		return err
	}
	slog.Info("Event emitted", "event", event)
	return nil
}

func (adapter *RabbitMqAdapter) Consume(queue string, handler func(event event.Event) error) {
	channel, err := adapter.connection.Channel()
	defer channel.Close()
	if err != nil {
		slog.Error("Error on open rabbitmq channel", err)
		return
	}
	_, err = channel.QueueDeclarePassive(queue, true, false, false, false, nil)
	if err != nil {
		slog.Error("Error on declare queue", "queue", queue, "error", err)
	}
	messages, err := channel.Consume(
		queue,
		"",
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

	for message := range messages {
		var messageEvent event.Event
		slog.Info("Message received", "exchange", message.Exchange, "routingKey", message.RoutingKey)
		err := json.Unmarshal(message.Body, &messageEvent)
		if err != nil {
			return
		}
		go func() {
			err := handler(messageEvent)
			if err := message.Ack(err == nil); err != nil {
				slog.Error("Error on handle message", err)
			}
		}()
	}
}
