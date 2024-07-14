package event

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/infra/env"
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"time"
)

type KafkaEventHandler struct {
	producer    *kafka.Producer
	environment *env.Environment
}

func NewKafkaEventHandler(environment *env.Environment) event.Handler {
	configMapProducer := &kafka.ConfigMap{
		"bootstrap.servers":   environment.KafkaBootstrapServer,
		"delivery.timeout.ms": "0",
		"acks":                "all",
		"enable.idempotence":  "true",
		"client.id":           environment.KafkaClientId,
	}
	producer, err := kafka.NewProducer(configMapProducer)
	if err != nil {
		panic(err)
	}
	handler := &KafkaEventHandler{producer: producer, environment: environment}
	go handler.deliveryReport()
	return handler
}

func (handler *KafkaEventHandler) deliveryReport() {
	for e := range handler.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
			}
		}
	}
}

func (handler *KafkaEventHandler) Emit(_ context.Context, event *event.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &event.Name,
			Partition: kafka.PartitionAny,
		},
		Key:       []byte(event.CorrelationId),
		Value:     body,
		Timestamp: time.Now(),
	}
	err = handler.producer.Produce(message, nil)
	if err != nil {
		return err
	}
	return nil
}

func (handler *KafkaEventHandler) Consume(options *event.ConsumerOptions, handlerFunc event.HandlerFunc) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": handler.environment.KafkaBootstrapServer,
		"client.id":         handler.environment.KafkaClientId,
		"group.id":          fmt.Sprintf("%s_%s", handler.environment.KafkaClientId, options.Id),
		"auto.offset.reset": "earliest",
	}
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}
	err = consumer.SubscribeTopics([]string{options.TopicName}, nil)
	if err != nil {
		slog.Error("Error on subscribe topic", slog.String("topic", options.TopicName), "error", err)
		panic(err)
	}
	go func() {
		for {
			message, err := consumer.ReadMessage(-1)
			if err != nil {
				slog.Error("Error on consume message", "error", err, "topic", options.TopicName)
				continue
			}
			messageEvent := &KafkaMessageEvent{message: message, consumer: consumer}
			err = handlerFunc(messageEvent)
			if err != nil {
				messageEvent.Nack()
			}
		}
	}()
}

type KafkaMessageEvent struct {
	message  *kafka.Message
	consumer *kafka.Consumer
}

func (ev *KafkaMessageEvent) AckIfNoError(err error) error {
	if err == nil {
		return ev.Ack()
	}
	return err
}

func (ev *KafkaMessageEvent) Ack() error {
	_, err := ev.consumer.CommitMessage(ev.message)
	if err != nil {
		slog.Error("Error on commit message", "error", err, "topic", ev.message.TopicPartition.Topic)
		return err
	}
	return nil
}

func (ev *KafkaMessageEvent) Nack() {
	err := ev.consumer.Seek(ev.message.TopicPartition, 0)
	if err != nil {
		slog.Error("Error on seek offset", "error", err, "topic", ev.message.TopicPartition.Topic)
		return
	}
}

func (ev *KafkaMessageEvent) GetEvent() *event.Event {
	var event event.Event
	_ = json.Unmarshal(ev.message.Value, &event)
	return &event
}

func (handler *KafkaEventHandler) Close() {
	handler.producer.Close()
}
