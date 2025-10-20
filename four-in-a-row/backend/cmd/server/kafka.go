package main

import (
	"os"
	"strings"
)

func getKafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}

func getKafkaTopic() string {
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		return "game-analytics"
	}
	return topic
}
