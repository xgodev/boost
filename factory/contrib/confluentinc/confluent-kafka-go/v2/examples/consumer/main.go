package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2"
	"github.com/xgodev/boost/wrapper/log"
	"time"
)

func main() {

	boost.Start()

	ctx := context.Background()

	consumer, err := confluent.NewConsumer(ctx)
	if err != nil {
		panic(err)
	}

	logger := log.FromContext(ctx)

	if err := consumer.SubscribeTopics([]string{"topic"}, nil); err != nil {
		panic(err)
	}

	errorCount := 0
	for {

		msg, err := consumer.ReadMessage(10 * time.Second)
		if err != nil {
			if err.(kafka.Error).IsTimeout() {
				logger.Warnf("Consumer error: %v (%v)", err, msg)
				continue
			}
			continue
		}

		for {
			logger.Infof("Processing message on %s: %s", msg.TopicPartition, string(msg.Value))

			if errorCount < 5 {
				errorCount++
				logger.Errorf("Simulated error. Error count: %d", errorCount)
				time.Sleep(1 * time.Second) // Simulando um tempo de espera para retry
				// Não comitar o offset aqui, para que a mensagem continue sendo reprocessada
				continue
			}

			errorCount = 0

			if _, err := consumer.CommitMessage(msg); err != nil {
				logger.Errorf("Failed to commit message: %v", err)
				// Não comitar, o que garantirá que a mensagem seja reprocessada
				continue
			}

			logger.Infof("Message successfully processed and committed: %s", string(msg.Value))
			break
		}

	}
}
