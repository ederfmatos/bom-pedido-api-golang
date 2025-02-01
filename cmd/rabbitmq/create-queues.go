package main

import (
	"bom-pedido-api/pkg/log"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func main() {
	url := "amqp://guest:guest@localhost:5672/"

	if err := createQueues(url); err != nil {
		log.Error("Erro ao criar filas", err)
		return
	}

	log.Info("Recursos criados com sucesso")
}

type ResourceConfig struct {
	Exchanges map[string]Exchange `json:"exchanges"`
	Queues    map[string]Queue    `json:"queues"`
}

type Exchange struct {
	Type string `json:"type"`
}

type Queue struct {
	Bindings  []QueueBinding  `json:"bindings"`
	Arguments *QueueArguments `json:"arguments"`
}

type QueueBinding struct {
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routingKey"`
}

type QueueArguments struct {
	DeadLetterExchange   string        `json:"x-dead-letter-exchange"`
	DeadLetterRoutingKey string        `json:"x-dead-letter-routing-key"`
	MessageTtl           time.Duration `json:"x-message-ttl"`
}

var config ResourceConfig

func init() {
	defaultExchange := "bompedido"
	deadLetterExchange := "bompedido-dlx"
	config = ResourceConfig{
		Exchanges: map[string]Exchange{
			defaultExchange: {
				Type: amqp.ExchangeTopic,
			},
			deadLetterExchange: {
				Type: amqp.ExchangeDirect,
			},
		},
		Queues: map[string]Queue{
			"SEND_EMAIL": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "SEND_EMAIL"},
				},
			},
			"AWAIT_APPROVAL_ORDER": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "PIX_TRANSACTION_PAID"},
				},
			},
			"CANCEL_PIX_TRANSACTION": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "PIX_PAYMENT_CANCELLED"},
				},
			},
			"CREATE_PIX_TRANSACTION": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "PIX_PAYMENT_CREATED"},
				},
			},
			"PAY_PIX_TRANSACTION": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "PAYMENT_CALLBACK_RECEIVED"},
				},
			},
			"REFUND_PIX_TRANSACTION": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "PIX_PAYMENT_REFUNDED"},
				},
			},
			"ORDER_PAYMENT_FAILED": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "PIX_PAYMENT_CANCELLED"},
				},
			},
			"CREATE_PIX_PAYMENT": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "ORDER_CREATED"},
				},
			},
			"DELETE_SHOPPING_CART": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "ORDER_CREATED"},
				},
			},
			"REFUND_PIX_PAYMENT": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "ORDER_CANCELLED"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_REJECTED"},
				},
			},
			"WAIT_CHECK_PIX_PAYMENT_FAILED": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "PIX_PAYMENT_CREATED"},
				},
				Arguments: &QueueArguments{
					DeadLetterExchange:   deadLetterExchange,
					DeadLetterRoutingKey: "CHECK_PIX_PAYMENT_FAILED",
					MessageTtl:           time.Minute * 30,
				},
			},
			"CHECK_PIX_PAYMENT_FAILED": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "PAYMENT_CALLBACK_RECEIVED"},
					{Exchange: deadLetterExchange, RoutingKey: "CHECK_PIX_PAYMENT_FAILED"},
				},
			},
			"NOTIFY_CUSTOMER_ORDER_STATUS_CHANGED": {
				Bindings: []QueueBinding{
					{Exchange: defaultExchange, RoutingKey: "ORDER_AWAITING_APPROVAL"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_APPROVED"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_IN_PROGRESS"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_REJECTED"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_CANCELLED"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_DELIVERING"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_AWAITING_DELIVERY"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_AWAITING_WITHDRAW"},
					{Exchange: defaultExchange, RoutingKey: "ORDER_FINISHED"},
				},
			},
		},
	}
}

func createQueues(url string) error {
	connection, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("create connection: %v", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("create channel: %v", err)
	}

	for name, item := range config.Exchanges {
		err = channel.ExchangeDeclare(name, item.Type, true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("declare exchange %s: %v", name, err)
		}
		log.Info("Exchange criada com sucesso", "exchange", name)
	}

	for queueName, queue := range config.Queues {
		var arguments amqp.Table
		if queue.Arguments != nil {
			arguments = amqp.Table{
				"x-dead-letter-exchange":    queue.Arguments.DeadLetterExchange,
				"x-dead-letter-routing-key": queue.Arguments.DeadLetterRoutingKey,
				"x-message-ttl":             queue.Arguments.MessageTtl.Milliseconds(),
			}
		}
		_, err = channel.QueueDeclare(queueName, true, false, false, false, arguments)
		if err != nil {
			return fmt.Errorf("declare queue %s: %v", queueName, err)
		}
		log.Info("Queue criada com sucesso", "queue", queueName)

		for _, binding := range queue.Bindings {
			err = channel.QueueBind(queueName, binding.RoutingKey, binding.Exchange, false, nil)
			if err != nil {
				return fmt.Errorf("bind queue %s: %v", queueName, err)
			}
		}
	}
	return nil
}
