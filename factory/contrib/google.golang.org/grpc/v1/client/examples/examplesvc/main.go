package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/plugins/local/wrapper/log"
	"os"
	"time"

	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/pb"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	alog "github.com/xgodev/boost/wrapper/log"
)

func init() {
	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "TRACE")
}

func main() {

	ctx := context.Background()

	boost.Start()

	ilog.New()

	request := &pb.TestRequest{
		Message: "mensagem da requisição",
	}

	options, _ := client.NewOptions()
	options.Host = "localhost"
	options.Port = 8080
	options.TLS.Enabled = true
	options.TLS.InsecureSkipVerify = true

	conn := client.NewClientConnWithOptions(ctx, options, log.Register)
	defer conn.Close()

	c := pb.NewExampleClient(conn)

	rctx, _ := context.WithTimeout(ctx, 1*time.Minute)

	test, err := c.Test(rctx, request)
	if err != nil {
		alog.Fatalf("%v.Call(_) = _, %v", c, err)
	}

	alog.Infof(test.Message)

	alog.Infof(conn.GetState().String())
}
