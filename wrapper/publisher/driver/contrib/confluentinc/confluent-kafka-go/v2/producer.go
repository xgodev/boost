package confluent

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

// producer is the part of *kafka.Producer the driver uses. Produce only
// enqueues the message; its outcome arrives later on deliveryChan.
type producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}
