package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2"
	"github.com/xgodev/boost/wrapper/log"
	"strconv"
)

func main() {

	boost.Start()

	ctx := context.Background()

	consumer, err := confluent.NewProducer(ctx)
	if err != nil {
		panic(err)
	}

	logger := log.FromContext(ctx)

	topic := "topic"

	for i := 0; i < 10; i++ {

		value := "Hello Go! " + strconv.Itoa(i)

		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
		}
		err := consumer.Produce(msg, nil)
		if err != nil {
			logger.Errorf("Failed to produce message: %v", err)
		}

		logger.Infof("Produced message: %s", value)

	}

	consumer.Flush(15 * 1000)
}
