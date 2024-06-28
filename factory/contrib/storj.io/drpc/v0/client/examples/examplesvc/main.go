package main

import (
	"context"
	"github.com/xgodev/boost/wrapper/config"
	"os"
	"time"

	"github.com/xgodev/boost/factory/contrib/storj.io/drpc/v0/client"
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

	request := &pb.TestRequest{
		Message: "mensagem da requisição",
	}

	conn, _ := client.NewClientConn(ctx)
	defer conn.Close()

	c := pb.NewDRPCExampleClient(conn)

	rctx, _ := context.WithTimeout(ctx, 1*time.Minute)

	test, err := c.Test(rctx, request)
	if err != nil {
		alog.Fatalf("%v.Call(_) = _, %v", c, err)
	}

	alog.Infof(test.Message)
}
