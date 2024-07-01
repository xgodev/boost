package kinesis

import (
	"github.com/xgodev/boost/bootstrap/repository"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis"
)

// NewEvent returns a initialized client
func NewEvent(c kinesis.Client, options *Options) repository.Event {
	return NewClient(c, options)
}
