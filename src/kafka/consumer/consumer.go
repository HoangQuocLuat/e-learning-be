package kafka_consumer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type ConsumerHandler func(message *kafka.Message) error

func ConsumeTopic(ctx context.Context, reader *kafka.Reader, handler ConsumerHandler) {
	run := true

	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			message, err := reader.FetchMessage(ctx)
			if err == nil {
				err := handler(&message)
				if err != nil {
					log.Println("Failed to process message")
				} else {
					if err := reader.CommitMessages(ctx, message); err != nil {
						log.Println("Failed to commit message")
					}
				}
			} else if !isTimeoutError(err) {
				log.Printf("Consumer error: %v", err)
			}
		}
	}
	if err := reader.Close(); err != nil {
		panic(err)
	}
}

func isTimeoutError(err error) bool {
	if kafkaErr, ok := err.(kafka.Error); ok {
		return kafkaErr.Timeout()
	}
	return false
}
