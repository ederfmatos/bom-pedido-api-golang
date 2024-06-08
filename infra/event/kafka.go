package event

import (
	"bom-pedido-api/application/event"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaEventEmitter struct {
	ConfigMap *kafka.ConfigMap
}

func NewKafkaEventEmitter(server string) *KafkaEventEmitter {
	return &KafkaEventEmitter{
		ConfigMap: &kafka.ConfigMap{
			"bootstrap.servers":   server,
			"delivery.timeout.ms": "0",
			"enable.idempotence":  "true",
		},
	}
}

func (handler *KafkaEventEmitter) Emit(event *event.Event) error {
	producer, err := kafka.NewProducer(handler.ConfigMap)
	if err != nil {
		return err
	}

	eventJson, err := json.Marshal(event)
	if err != nil {
		return err
	}

	topic := event.Name
	key := event.Id
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          eventJson,
		Key:            []byte(key),
	}
	return producer.Produce(message, nil)
}
