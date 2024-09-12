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

	partitions := make(map[int32]chan *kafka.Message)

	for {
		msg, err := consumer.ReadMessage(10 * time.Second)
		if err != nil {
			if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
				logger.Warnf("Consumer timeout: %v", kafkaErr)
				continue
			}
			logger.Errorf("Failed to read message: %v", err)
			continue
		}

		partition := msg.TopicPartition.Partition

		if _, exists := partitions[partition]; !exists {
			partitions[partition] = make(chan *kafka.Message, 100)

			go processPartition(ctx, consumer, partitions[partition], partition)
		}

		partitions[partition] <- msg
	}
}

func processPartition(ctx context.Context, consumer *kafka.Consumer, msgs chan *kafka.Message, partition int32) {
	logger := log.FromContext(ctx)
	errorCount := 0

	for msg := range msgs {
		for {
			logger.Infof("Processing message from partition %d: %s", partition, string(msg.Value))

			// Simula um erro no processamento
			if errorCount < 5 {
				errorCount++
				logger.Errorf("Simulated error in partition %d. Error count: %d", partition, errorCount)
				time.Sleep(1 * time.Second) // Simulando um tempo de espera para retry
				continue
			}

			// Após sucesso, resetar o contador de erros
			errorCount = 0

			if _, err := consumer.CommitMessage(msg); err != nil {
				logger.Errorf("Failed to commit message from partition %d: %v", partition, err)
				continue // Retry no commit
			}

			logger.Infof("Message from partition %d successfully processed and committed: %s", partition, string(msg.Value))
			break
		}
	}
}
