package sns

import (
	"github.com/xgodev/boost/bootstrap/repository"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sns"
)

// NewEvent returns a initialized client
func NewEvent(c sns.Client) repository.Event {
	return NewClient(c)
}
