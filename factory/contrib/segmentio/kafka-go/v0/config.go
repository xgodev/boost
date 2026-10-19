package kafka

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = "boost.factory.kafkago"
	brokers = ".brokers"

	writer             = ".writer"
	writerBalancer     = writer + ".balancer"
	writerRequiredAcks = writer + ".requiredAcks"
	writerCompression  = writer + ".compression"
	writerMaxAttempts  = writer + ".maxAttempts"
	writerBatchSize    = writer + ".batchSize"
	writerBatchBytes   = writer + ".batchBytes"
	writerBatchTimeout = writer + ".batchTimeout"
	writerReadTimeout  = writer + ".readTimeout"
	writerWriteTimeout = writer + ".writeTimeout"
	writerAsync        = writer + ".async"
	writerAutoTopic    = writer + ".allowAutoTopicCreation"

	reader                 = ".reader"
	readerGroupID          = reader + ".groupId"
	readerTopic            = reader + ".topic"
	readerPartition        = reader + ".partition"
	readerQueueCapacity    = reader + ".queueCapacity"
	readerMinBytes         = reader + ".minBytes"
	readerMaxBytes         = reader + ".maxBytes"
	readerMaxWait          = reader + ".maxWait"
	readerReadBatchTimeout = reader + ".readBatchTimeout"
	readerCommitInterval   = reader + ".commitInterval"
	readerStartOffset      = reader + ".startOffset"

	conn          = ".conn"
	connType      = conn + ".type"
	connNetwork   = conn + ".network"
	connTopic     = conn + ".topic"
	connPartition = conn + ".partition"
	connTimeout   = conn + ".timeout"

	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd adds the kafka-go keys under path.
func ConfigAdd(path string) {
	config.Add(path+brokers, []string{"localhost:9092"}, "kafka broker addresses")

	// crc32 is the partitioner of librdkafka, so a key lands on the same
	// partition whether confluent-kafka-go or kafka-go produced it.
	config.Add(path+writerBalancer, "crc32", "partitioner: roundrobin, leastbytes, hash, crc32 or murmur2")
	config.Add(path+writerRequiredAcks, "all", "acks a produce request waits for: none, one or all")
	config.Add(path+writerCompression, "", "message compression: empty (none), gzip, snappy, lz4 or zstd")
	config.Add(path+writerMaxAttempts, 10, "attempts to deliver a batch before giving up")
	config.Add(path+writerBatchSize, 100, "messages buffered before a batch is sent")
	config.Add(path+writerBatchBytes, int64(1048576), "max bytes of a batch")
	// kafka-go waits 1s by default, which a synchronous publish of one event
	// would pay every time.
	config.Add(path+writerBatchTimeout, 10*time.Millisecond, "time an incomplete batch waits before being sent")
	config.Add(path+writerReadTimeout, 10*time.Second, "timeout of a read from the broker")
	config.Add(path+writerWriteTimeout, 10*time.Second, "timeout of a write to the broker")
	config.Add(path+writerAsync, false, "WriteMessages returns without waiting for the broker")
	config.Add(path+writerAutoTopic, false, "creates a missing topic on write")

	config.Add(path+readerGroupID, "", "consumer group; empty reads one partition")
	config.Add(path+readerTopic, "", "topic to read")
	config.Add(path+readerPartition, 0, "partition to read when there is no consumer group")
	config.Add(path+readerQueueCapacity, 100, "capacity of the internal message queue")
	config.Add(path+readerMinBytes, 1, "min bytes a fetch waits for")
	config.Add(path+readerMaxBytes, 10485760, "max bytes of a fetch")
	config.Add(path+readerMaxWait, 2*time.Second, "max time a fetch waits for min bytes")
	config.Add(path+readerReadBatchTimeout, 2*time.Second, "max time to read a batch")
	config.Add(path+readerCommitInterval, time.Duration(0), "interval of offset commits; 0 commits synchronously")
	config.Add(path+readerStartOffset, "last", "offset of a group with no commit: first or last")

	config.Add(path+connType, "leader", "connection: leader (of topic/partition), partition or server")
	config.Add(path+connNetwork, "tcp", "network of the connection")
	config.Add(path+connTopic, "", "topic of a leader or partition connection")
	config.Add(path+connPartition, 0, "partition of a leader or partition connection")
	config.Add(path+connTimeout, 10*time.Second, "dial timeout")
}
