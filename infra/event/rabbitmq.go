package event

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/infra/json"
	"bom-pedido-api/infra/retry"
	"bom-pedido-api/infra/telemetry"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/codes"
	"log"
	"log/slog"
	"os"
	"time"
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
	adapter := &RabbitMqAdapter{
		connection:      connection,
		producerChannel: producerChannel,
		consumerChannel: consumerChannel,
	}
	adapter.createQueues()
	return adapter
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
		span.SetStatus(codes.Error, err.Error())
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

func (adapter *RabbitMqAdapter) createQueues() {
	file, err := os.ReadFile(".resources/rabbitmq.json")
	handleError(err)
	var resourceConfig ResourceConfig
	err = json.Unmarshal(context.Background(), file, &resourceConfig)
	handleError(err)
	for _, item := range resourceConfig.Exchanges {
		err = adapter.consumerChannel.ExchangeDeclare(item.Name, item.Type, true, false, false, false, nil)
		handleError(err)
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
		_, err = adapter.consumerChannel.QueueDeclare(queue.Name, true, false, false, false, arguments)
		handleError(err)
		for _, binding := range queue.Bindings {
			err = adapter.consumerChannel.QueueBind(queue.Name, binding.RoutingKey, binding.Exchange, false, nil)
			handleError(err)
		}
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
