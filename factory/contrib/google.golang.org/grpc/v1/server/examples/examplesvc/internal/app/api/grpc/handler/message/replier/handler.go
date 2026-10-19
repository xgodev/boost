// Package replier serves the Test rpc of the Example service.
package replier

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/internal/pkg/domain/service/message/replier"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/pb"
)

// Handler serves the Example service.
type Handler struct {
	pb.UnimplementedExampleServer
	svc replier.Service
}

// New provides the handler, which the generated grpc door registers with
// pb.RegisterExampleServer.
// @Provide(index=0)
// @GRPC(register=pb.RegisterExampleServer)
func New(svc replier.Service) *Handler {
	return &Handler{svc: svc}
}

// Test replies to the message of the request.
func (h *Handler) Test(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	reply, err := h.svc.Reply(ctx, &entity.Message{Text: req.GetMessage()})
	if err != nil {
		return nil, server.Error(err)
	}
	return &pb.TestResponse{Message: reply.Text}, nil
}
