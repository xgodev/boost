package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/pb"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/plugins/local/wrapper/log"
	"os"

	alog "github.com/xgodev/boost/wrapper/log"
)

func init() {
	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "TRACE")
}

func main() {

	ctx := context.Background()

	boost.Start()

	options, _ := server.NewOptions()
	options.Port = 8080

	srv := server.NewServerWithOptions(ctx, options, log.Register)

	pb.RegisterExampleServer(srv.ServiceRegistrar(), NewService())

	srv.Serve(ctx)
}

type Service struct {
	pb.UnimplementedExampleServer
}

func (h *Service) Test(ctx context.Context, request *pb.TestRequest) (*pb.TestResponse, error) {

	logger := alog.FromContext(ctx)

	logger.Infof(request.Message)

	return &pb.TestResponse{Message: "hello world"}, nil
}

func NewService() pb.ExampleServer {
	return &Service{}
}
