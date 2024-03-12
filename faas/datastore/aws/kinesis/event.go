package kinesis

import (
	"github.com/xgodev/boost/faas/repository"
	"github.com/xgodev/boost/factory/aws/aws-sdk-go.v2/client/kinesis"
)

// NewEvent returns a initialized client
func NewEvent(c kinesis.Client, options *Options) repository.Event {
	return NewClient(c, options)
}
