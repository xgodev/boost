package confluent

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root            = "boost.factory.confluent"
	brokers         = ".brokers"
	producer        = ".producer"
	logRoot         = ".log"
	level           = logRoot + ".level"
	logEnabled      = logRoot + ".enabled"
	acks            = producer + ".acks"
	timeout         = producer + ".timeout"
	request         = timeout + ".request"
	message         = timeout + ".message"
	batch           = producer + ".batch"
	numMessages     = batch + ".numMessages"
	size            = batch + ".size"
	consumer        = ".consumer"
	topics          = consumer + ".topics"
	groupId         = consumer + ".groupId"
	autoOffsetReset = consumer + ".autoOffsetReset"
	autoCommit      = consumer + ".autoCommit"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+brokers, "localhost:9092", "defines brokers addresses")
	config.Add(path+topics, []string{"changeme"}, "defines topics")
	config.Add(path+groupId, "changeme", "defines consumer groupid")
	config.Add(path+autoOffsetReset, "earliest", "defines consumer auto offset reset")
	config.Add(path+autoCommit, false, "defines consumer auto commit")
	config.Add(path+numMessages, 10000, "Maximum number of messages batched in one MessageSet. The total MessageSet size is also limited by batch.size and message.max.bytes")
	config.Add(path+size, 1000000, "Maximum size (in bytes) of all messages batched in one MessageSet, including protocol framing overhead. This limit is applied after the first message has been added to the batch, regardless of the first message's size, this is to ensure that messages that exceed batch.size are produced. The total MessageSet size is also limited by batch.num.messages and message.max.bytes")
	config.Add(path+acks, -1, "This field indicates the number of acknowledgements the leader broker must receive from ISR brokers before responding to the request: 0=Broker does not send any response/ack to client, -1 or all=Broker will block until message is committed by all in sync replicas (ISRs). If there are less than min.insync.replicas (broker configuration) in the ISR set the produce request will fail")
	config.Add(path+request, 30000, "The ack timeout of the producer request in milliseconds. This value is only enforced by the broker and relies on request.required.acks being != 0")
	config.Add(path+message, 300000, "Local message timeout. This value is only enforced locally and limits the time a produced message waits for successful delivery. A time of 0 is infinite. This is the maximum time librdkafka may use to deliver a message (including retries). Delivery error occurs when either the retry count or the message timeout are exceeded. The message timeout is automatically adjusted to transaction.timeout.ms if transactional.id is configured")
	config.Add(path+level, "DEBUG", "defines log level")
	config.Add(path+logEnabled, true, "defines log enabled")
}
