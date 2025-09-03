package pubsub

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root                      = "boost.wrapper.publisher.driver.pubsub"
	logRoot                   = ".log"
	orderingKey               = ".orderingKey"
	level                     = logRoot + ".level"
	settings                  = ".settings"
	delayThreshold            = settings + ".delayThreshold"
	countThreshold            = settings + ".countThreshold"
	byteThreshold             = settings + ".byteThreshold"
	numGoroutines             = settings + ".numGoroutines"
	timeout                   = settings + ".timeout"
	bufferedByteLimit         = settings + ".bufferedByteLimit"
	flowControlSettings       = settings + ".flowControlSettings"
	maxOutstandingMessages    = flowControlSettings + ".maxOutstandingMessages"
	maxOutstandingBytes       = flowControlSettings + ".maxOutstandingBytes"
	limitExceededBehavior     = flowControlSettings + ".limitExceededBehavior"
	enableCompression         = settings + ".enableCompression"
	compressionBytesThreshold = settings + ".compressionBytesThreshold"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+level, "DEBUG", "defines log level")
	config.Add(path+orderingKey, false, "defines ordering key")

	config.Add(path+delayThreshold, 10*time.Millisecond, "the maximum duration to wait before sending a batch of messages")
	config.Add(path+countThreshold, 100, "the maximum number of messages to include in a batch")
	config.Add(path+byteThreshold, 1e6, "the maximum total size of messages to include in a batch")
	config.Add(path+numGoroutines, 1, "the number of goroutines that process batches of messages")
	config.Add(path+timeout, 60*time.Second, "the maximum duration to block Publish calls")
	config.Add(path+bufferedByteLimit, 10*1e7, "the maximum number of bytes that can be pending in memory across all topics")
	config.Add(path+maxOutstandingMessages, 1000, "the maximum number of messages that can be pending in memory for publishing")
	config.Add(path+maxOutstandingBytes, -1, "the maximum total size of messages that can be pending in memory for publishing")
	config.Add(path+limitExceededBehavior, 1, "behavior when flow control limits are exceeded: Block or Ignore")
	config.Add(path+enableCompression, false, "whether to compress messages before sending")
	config.Add(path+compressionBytesThreshold, 240, "the minimum size a message must be to be compressed before sending")
}
