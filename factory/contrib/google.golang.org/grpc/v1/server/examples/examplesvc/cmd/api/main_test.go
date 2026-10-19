package main

import (
	"context"

	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	grpcserver "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/gen/builder/api"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/pb"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// TestAPI calls the Test rpc on the generated grpc server.
func TestAPI(t *testing.T) {
	var srv *grpcserver.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	lis := bufconn.Listen(1 << 20)
	go func() { _ = srv.Server().Serve(lis) }()
	t.Cleanup(srv.Server().Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	res, err := pb.NewExampleClient(conn).Test(context.Background(), &pb.TestRequest{Message: "hi"})
	require.NoError(t, err)
	assert.Equal(t, "hello world", res.GetMessage())
}
