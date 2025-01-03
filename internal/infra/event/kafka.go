package event

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/json"
	"bom-pedido-api/internal/infra/retry"
	"bom-pedido-api/internal/infra/telemetry"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel/codes"
	"log/slog"
	"time"
)

type KafkaEventHandler struct {
	producer    *kafka.Producer
	consumers   []*kafka.Consumer
	environment *config.Environment
}

func NewKafkaEventHandler(environment *config.Environment) event.Handler {
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
	slog.Info("Connected to kafka successfully")
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
	body, err := json.Marshal(ctx, event)
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
	err = consumer.SubscribeTopics(options.Topics, nil)
	if err != nil {
		slog.Error("Error on subscribe topic", "topics", options.Topics, "error", err)
		panic(err)
	}
	handler.consumers = append(handler.consumers, consumer)
	go handler.processMessages(consumer, options, handlerFunc)
}

func (handler *KafkaEventHandler) processMessages(consumer *kafka.Consumer, options *event.ConsumerOptions, handlerFunc event.HandlerFunc) {
	for {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			slog.Error("Error on consume message", "error", err, "topic", options.Topics)
			continue
		}
		go handler.processMessage(message, consumer, handlerFunc)
	}
}

func (handler *KafkaEventHandler) processMessage(message *kafka.Message, consumer *kafka.Consumer, handlerFunc event.HandlerFunc) {
	ctx, span := telemetry.StartSpan(context.Background(), "KafkaEventEmitter.Process", "messageKey", string(message.Key), "topic", *message.TopicPartition.Topic)
	defer span.End()
	messageEvent, err := handler.createMessageEvent(ctx, message, consumer)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			messageEvent.Nack(ctx)
		}
	}()
	if err != nil {
		return
	}
	err = retry.Func(ctx, 6, time.Second, time.Second*30, func(ctx context.Context) error {
		return handlerFunc(ctx, messageEvent)
	})
	if err == nil {
		return
	}
	err = handler.sendMessageToDeadLetterTopic(message, err)
	if err == nil {
		_ = messageEvent.Ack(ctx)
	} else {
		messageEvent.Nack(ctx)
	}
}

func (handler *KafkaEventHandler) createMessageEvent(ctx context.Context, message *kafka.Message, consumer *kafka.Consumer) (*event.MessageEvent, error) {
	var messageEvent event.Event
	err := json.Unmarshal(ctx, message.Value, &messageEvent)
	return &event.MessageEvent{
		Event: &messageEvent,
		AckFn: func(ctx context.Context) error {
			_, err = consumer.CommitMessage(message)
			if err != nil {
				slog.Error("Error on commit message", "error", err, "topic", message.TopicPartition.Topic)
			}
			return nil
		},
		NackFn: func(ctx context.Context) {
			err = consumer.Seek(message.TopicPartition, 0)
			if err != nil {
				slog.Error("Error on seek offset", "error", err, "topic", message.TopicPartition.Topic)
			}
		},
	}, err
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
		_ = consumer.Close()
	}
	handler.consumers = make([]*kafka.Consumer, 0)
}
