package kafka_config

import (
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

var KafkaProducer *kafka.Writer

func InitKafkaProducer() {
	KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Topic:    "send-email",
		Balancer: &kafka.LeastBytes{},
	}
}

func CLoseKafka() {
	if err := KafkaProducer.Close(); err != nil {
		log.Fatalf("Failed to close kafka producer: %v", err)
	}
}

func NewKafkaConsumer() *kafka.Reader {
	brokers := strings.Split("localhost:29092", ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        "send-mail-group",
		Topic:          "send-email",
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})
}
