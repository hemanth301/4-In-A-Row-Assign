package services

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type AnalyticsConsumer struct {
	Reader *kafka.Reader
}

func NewAnalyticsConsumer(brokers []string, topic, groupID string) *AnalyticsConsumer {
	return &AnalyticsConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (a *AnalyticsConsumer) StartConsuming(ctx context.Context, handle func(eventType, payload string)) {
	for {
		m, err := a.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka consumer error: %v", err)
			continue
		}
		handle(string(m.Key), string(m.Value))
	}
}
