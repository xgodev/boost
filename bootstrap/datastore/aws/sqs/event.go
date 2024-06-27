package sqs

import (
	"github.com/xgodev/boost/bootstrap/repository"
	giawsclientsqs "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs"
)

// NewEvent returns a initialized client
func NewEvent(c giawsclientsqs.Client) repository.Event {
	return NewClient(c)
}
