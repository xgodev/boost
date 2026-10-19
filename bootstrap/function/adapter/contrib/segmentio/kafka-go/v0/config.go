package kafka

import (
	"time"

	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = adapter.Root + ".kafkago"
	topics      = root + ".topics"
	buffer      = root + ".buffer"
	backoff     = root + ".backoff"
	backoffBase = root + ".backoffBase"
	maxBackoff  = root + ".maxBackoff"
	retryLimit  = root + ".retryLimit"
)

func init() {
	config.Add(topics, []string{"changeme"}, "kafka topics read by the consumer group")
	config.Add(buffer, 100, "messages queued per partition while its handler is busy")
	config.Add(backoff, true, "waits between retries of a failed message")
	config.Add(backoffBase, 1*time.Second, "wait before the first retry, doubled on each retry")
	config.Add(maxBackoff, 20*time.Second, "max wait between retries")
	config.Add(retryLimit, -1, "retries of a failed message before it is committed and skipped; -1 retries forever")
}
