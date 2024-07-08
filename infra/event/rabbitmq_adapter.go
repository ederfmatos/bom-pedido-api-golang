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

func (adapter *RabbitMqAdapter) Emit(context context.Context, event *event.Event) error {
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

func (adapter *RabbitMqAdapter) Consume(queue string, handler event.Handler) {
	channel := adapter.consumerChannel
	_, err := channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		slog.Error("Error on declare queue", "queue", queue, "error", err)
	}
	messages, err := channel.Consume(
		queue,
		"BOM_PEDIDO_API_"+queue,
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

	go func() {
		select {
		case message := <-messages:
			var messageEvent event.Event
			slog.Info("Message received", "exchange", message.Exchange, "routingKey", message.RoutingKey)
			err := json.Unmarshal(message.Body, &messageEvent)
			if err != nil {
				_ = message.Nack(false, true)
				return
			}

			err = handler(messageEvent)
			if err == nil {
				_ = message.Ack(false)
				return
			}
			_ = message.Nack(false, true)
		}
	}()
}
