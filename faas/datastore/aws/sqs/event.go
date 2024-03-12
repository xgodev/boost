package sqs

import (
	"github.com/xgodev/boost/faas/repository"
	giawsclientsqs "github.com/xgodev/boost/factory/aws/aws-sdk-go.v2/client/sqs"
)

// NewEvent returns a initialized client
func NewEvent(c giawsclientsqs.Client) repository.Event {
	return NewClient(c)
}
