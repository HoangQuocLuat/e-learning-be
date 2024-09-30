package main

import (
	"context"
	kafka_config "e-learning/src/kafka"
	kafka_consumer "e-learning/src/kafka/consumer"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Worker is starting...")
	ctx, cancel := context.WithCancel(context.Background())

	sendMailConsumer := kafka_config.NewKafkaConsumer()
	orderHandler := kafka_consumer.NewSendMailConsumer()
	go kafka_consumer.ConsumeTopic(ctx, sendMailConsumer, orderHandler.Consumer)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)

	stop := false
	for !stop {
		s := <-terminateSignals
		log.Println("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
		cancel()
		stop = true

	}

	time.Sleep(5 * time.Second)
}
