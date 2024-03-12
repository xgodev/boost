package sns

import (
	"github.com/xgodev/boost/faas/repository"
	"github.com/xgodev/boost/factory/aws/aws-sdk-go.v2/client/sns"
)

// NewEvent returns a initialized client
func NewEvent(c sns.Client) repository.Event {
	return NewClient(c)
}
