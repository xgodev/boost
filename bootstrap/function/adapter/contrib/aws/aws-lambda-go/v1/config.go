package lambda

import (
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root       = adapter.Root + ".lambda"
	skip       = root + ".skip"
	maxWorkers = root + ".maxWorkers"
)

func init() {
	config.Add(skip, false, "acknowledges every invocation without calling the handler")
	config.Add(maxWorkers, 10, "records of one invocation handled at the same time; 1 keeps their order")
}
