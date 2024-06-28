package main

import (
	"context"
	"github.com/xgodev/boost/wrapper/config"
	"os"

	"github.com/xgodev/boost/factory/contrib/storj.io/drpc/v0/server"
	"github.com/xgodev/boost/factory/contrib/storj.io/drpc/v0/server/examples/examplesvc/pb"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	alog "github.com/xgodev/boost/wrapper/log"
)

func init() {
	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "TRACE")
}

func main() {

	ctx := context.Background()

	config.Load()

	ilog.New()

	srv, _ := server.NewServer(ctx)
	m := srv.Mux()
	if err := pb.DRPCRegisterExample(m, &Service{}); err != nil {
		panic(err)
	}

	srv.Serve(ctx)
}

type Service struct {
	pb.DRPCExampleUnimplementedServer
}

func (h *Service) Test(ctx context.Context, request *pb.TestRequest) (*pb.TestResponse, error) {

	logger := alog.FromContext(ctx)

	logger.Infof(request.Message)

	return &pb.TestResponse{Message: "hello world"}, nil
}

func NewService() pb.DRPCExampleServer {
	return &Service{}
}
