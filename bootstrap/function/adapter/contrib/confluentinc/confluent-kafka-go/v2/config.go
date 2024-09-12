package confluent

import (
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

const (
	root         = adapter.Root + ".kafka_confluent"
	topics       = root + ".topics"
	timeOut      = root + ".timeOut"
	manualCommit = root + ".manualCommit"
	useSemaphore = root + ".useSemaphore"
	maxWorkers   = root + ".maxWorkers"
	backoff      = root + ".backoff"
	backoffBase  = root + ".backoffBase"
	maxBackoff   = root + ".maxBackoff"
	retryLimit   = root + ".retryLimit"
)

func init() {
	config.Add(topics, []string{"changeme"}, "kafka topics")
	config.Add(timeOut, 2*time.Second, "kafka timeout")
	config.Add(manualCommit, true, "kafka manual commit")
	config.Add(useSemaphore, true, "kafka use semaphore")
	config.Add(maxWorkers, int64(10), "kafka max workers")
	config.Add(backoff, true, "kafka backoff")
	config.Add(backoffBase, 1*time.Second, "kafka backoff base")
	config.Add(maxBackoff, 20*time.Second, "kafka max backoff")
	config.Add(retryLimit, -1, "kafka retry limit")
}
