package kafka_producer

import (
	"github.com/segmentio/kafka-go"
)

func NewKafkaProducer() *kafka.Writer {
	return &kafka.Writer{
		// Addr:     kafka.TCP(config.Get().KafkaAddr),
		Addr:     kafka.TCP("127.0.0.1:29092"),
		Topic:    "send",
		Balancer: &kafka.LeastBytes{},
	}
}

// type SendMailProducer struct {
// 	Producer[*kafka_event.SendMailEvent]
// }

// func NewSendEmailProducer(producer *kafka.Writer) *SendMailProducer {
// 	return &SendMailProducer{
// 		Producer: Producer[*kafka_event.SendMailEvent]{
// 			Producer: producer,
// 			Topic:    "send",
// 		},
// 	}
// }
