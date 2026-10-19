// Package grpc replies to messages through the Test rpc of the Example
// service, on a grpc client connection of boost's factory
// (boost.factory.grpc.client).
package grpc

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/examples/examplesvc/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/examples/examplesvc/internal/pkg/domain/service/message/replier"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/pb"
	"github.com/xgodev/boost/model/errors"
)

type replierClient struct {
	example pb.ExampleClient
}

// New opens the client connection with the log plugin and provides the
// replier on it: the generator does not build grpc clients, so the
// constructor calls the factory.
// @Provide(index=0)
func New(ctx context.Context) replier.Service {
	return &replierClient{example: pb.NewExampleClient(client.NewClientConn(ctx, log.Register))}
}

func (c *replierClient) Reply(ctx context.Context, msg *entity.Message) (*entity.Message, error) {
	res, err := c.example.Test(ctx, &pb.TestRequest{Message: msg.Text})
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error calling the Example service"))
	}
	return &entity.Message{Text: res.GetMessage()}, nil
}
