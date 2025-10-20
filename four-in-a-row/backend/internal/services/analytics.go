package services

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type AnalyticsProducer struct {
	Writer *kafka.Writer
}

func NewAnalyticsProducer(brokers []string, topic string) *AnalyticsProducer {
	return &AnalyticsProducer{
		Writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: topic,
		},
	}
}

func (a *AnalyticsProducer) Emit(eventType, payload string) error {
	return a.Writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(eventType),
		Value: []byte(payload),
	})
}
