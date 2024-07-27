package event

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/retry"
	"bom-pedido-api/infra/telemetry"
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"time"
)

type KafkaEventHandler struct {
	producer    *kafka.Producer
	consumers   []*kafka.Consumer
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
	handler := &KafkaEventHandler{
		producer:    producer,
		environment: environment,
		consumers:   make([]*kafka.Consumer, 0),
	}
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

func (handler *KafkaEventHandler) Emit(ctx context.Context, event *event.Event) error {
	ctx, span := telemetry.StartSpan(ctx, "KafkaEventEmitter.Emit",
		"eventId", event.Id, "eventName", event.Name,
	)
	defer span.End()
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
	return handler.producer.Produce(message, nil)
}

func (handler *KafkaEventHandler) Consume(options *event.ConsumerOptions, handlerFunc event.HandlerFunc) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":  handler.environment.KafkaBootstrapServer,
		"client.id":          handler.environment.KafkaClientId,
		"group.id":           fmt.Sprintf("%s_%s", handler.environment.KafkaClientId, options.Id),
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
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
	handler.consumers = append(handler.consumers, consumer)
	go handler.processMessages(consumer, options, handlerFunc)
}

func (handler *KafkaEventHandler) processMessages(consumer *kafka.Consumer, options *event.ConsumerOptions, handlerFunc event.HandlerFunc) {
	for {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			slog.Error("Error on consume message", "error", err, "topic", options.TopicName)
			continue
		}
		go handler.processMessage(message, consumer, handlerFunc)
	}
}

func (handler *KafkaEventHandler) processMessage(message *kafka.Message, consumer *kafka.Consumer, handlerFunc event.HandlerFunc) {
	_, span := telemetry.StartSpan(context.Background(), "KafkaEventEmitter.Process", "messageKey", string(message.Key))
	defer span.End()
	messageEvent := &event.MessageEvent{
		AckFn: func() error {
			_, err := consumer.CommitMessage(message)
			if err != nil {
				slog.Error("Error on commit message", "error", err, "topic", message.TopicPartition.Topic)
				return err
			}
			return nil
		},
		NackFn: func() {
			err := consumer.Seek(message.TopicPartition, 0)
			if err != nil {
				slog.Error("Error on seek offset", "error", err, "topic", message.TopicPartition.Topic)
			}
		},
		GetEventFn: func() *event.Event {
			var event event.Event
			_ = json.Unmarshal(message.Value, &event)
			return &event
		},
	}
	err := retry.Func(6, time.Second, time.Second*30, func() error {
		return handlerFunc(messageEvent)
	})
	if err == nil {
		return
	}
	err = handler.sendMessageToDeadLetterTopic(message, err)
	if err == nil {
		messageEvent.Ack()
	} else {
		messageEvent.Nack()
	}
}

func (handler *KafkaEventHandler) sendMessageToDeadLetterTopic(originalMessage *kafka.Message, err error) error {
	topic := "DEAD_LETTER_TOPIC"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   originalMessage.Key,
		Value: originalMessage.Value,
		Headers: []kafka.Header{
			{Key: "error", Value: []byte(err.Error())},
			{Key: "timestamp", Value: []byte(originalMessage.Timestamp.String())},
			{Key: "topic", Value: []byte(*originalMessage.TopicPartition.Topic)},
		},
		Timestamp: time.Now(),
	}
	return handler.producer.Produce(message, nil)
}

func (handler *KafkaEventHandler) Close() {
	handler.producer.Close()
	for _, consumer := range handler.consumers {
		consumer.Close()
	}
	handler.consumers = make([]*kafka.Consumer, 0)
}
