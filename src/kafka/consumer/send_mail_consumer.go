package kafka_consumer

import (
	"log"

	"github.com/segmentio/kafka-go"
	"gopkg.in/gomail.v2"
)

type SendMailConsumer struct {
}

func NewSendMailConsumer() *SendMailConsumer {
	return &SendMailConsumer{}
}

func (c *SendMailConsumer) Consumer(mess *kafka.Message) error {
	email := string(mess.Value)
	log.Println("Received email:", email)

	abc := gomail.NewMessage()
	abc.SetHeader("From", "hoangquocluatspak@gmail.com")
	abc.SetHeader("To", email)
	abc.SetHeader("Subject", "Lịch đã được thêm mới")
	abc.SetBody("text/plain", "Lịch đã được thêm mới, hãy kiểm tra lại lịch")

	// Sử dụng mật khẩu ứng dụng nếu xác thực 2 bước được bật
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "hoangquocluatspak@gmail.com", "tyyk yafp tpdr qgio")

	// Gửi email
	if err := dialer.DialAndSend(abc); err != nil {
		log.Printf("consumer err: %v", err)
		return err
	}

	return nil
}
