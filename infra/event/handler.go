package event

import (
	"bom-pedido-api/application/event"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type KafkaEventHandler struct {
	ConfigMap *kafka.ConfigMap
}

func NewKafkaEventHandler(server string) event.EventHandler {
	return &KafkaEventHandler{
		ConfigMap: &kafka.ConfigMap{
			"bootstrap.servers": server,
			"client.id":         "bom-pedido-api",
			"group.id":          "bom-pedido-api",
		},
	}
}

func (h *KafkaEventHandler) Consume(topic string, handler func(event event.Event) error) {
	consumer, err := kafka.NewConsumer(h.ConfigMap)
	defer consumer.Close()
	if err != nil {
		panic(err)
	}
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}
	for {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			var message event.Event
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				return
			}
			err = handler(message)
		}
	}
}
