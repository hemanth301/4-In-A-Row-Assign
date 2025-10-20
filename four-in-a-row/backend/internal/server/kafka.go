package server

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
		Async: true,
	}
	return &KafkaProducer{writer: w}
}

func (k *KafkaProducer) Emit(key, value string) error {
	return k.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})
}

func (k *KafkaProducer) Close() error {
	return k.writer.Close()
}
