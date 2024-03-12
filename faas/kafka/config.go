package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/config"
	"time"
)

const (
	root             = "faas.kafka"
	topics           = root + ".topics"
	groupId          = root + ".groupId"
	brokers          = root + ".brokers"
	concurrency      = root + ".concurrency"
	queueCapacity    = root + ".queueCapacity"
	minBytes         = root + ".minBytes"
	maxBytes         = root + ".maxBytes"
	startOffset      = root + ".startOffset"
	readBatchTimeout = root + ".readBatchTimeout"
	maxWait          = root + ".maxWait"
)

func init() {
	config.Add(topics, []string{"changeme"}, "kafka listener topics")
	config.Add(brokers, []string{"localhost:9090"}, "kafka listener brokers")
	config.Add(groupId, "changeme", "kafka listener groupId")
	config.Add(concurrency, 10, "kafka listener concurrency")
	config.Add(queueCapacity, 100, "defines queue capacity")
	config.Add(minBytes, 1, "defines batch min bytes")
	config.Add(maxBytes, 10485760, "defines batch max bytes")
	config.Add(readBatchTimeout, 2*time.Second, "defines read batch timeout")
	config.Add(maxWait, 2*time.Second, "defines max wait")
	config.Add(startOffset, kafka.LastOffset, "defines start offset LastOffset=-1, FirstOffset=-2")
}
