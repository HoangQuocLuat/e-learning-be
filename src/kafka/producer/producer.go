package kafka_producer

// import (
// 	"context"
// 	kafka_event "e-learning/src/kafka/event"
// 	"encoding/json"
// 	"log"

// 	"github.com/segmentio/kafka-go"
// )

// type Producer[T kafka_event.Event] struct {
// 	Producer *kafka.Writer
// 	Topic    string
// }

// func (p *Producer[T]) GetTopic() *string {
// 	return &p.Topic
// }

// func (p *Producer[T]) Send(event T) error {
// 	value, err := json.Marshal(event)
// 	if err != nil {
// 		log.Println("failed to marshal event")
// 		return err
// 	}

// 	message := kafka.Message{
// 		Topic: "send",
// 		Key:   []byte("event"),
// 		Value: value,
// 	}

// 	err = p.Producer.WriteMessages(context.Background(), message)

// 	if err != nil {
// 		log.Println("failed to produce message:", err)
// 		return err
// 	}
// 	return nil
// }
