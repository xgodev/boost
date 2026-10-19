package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/examples/examplesvc/gen/builder/api"
	grpcserver "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/pb"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

type example struct {
	pb.UnimplementedExampleServer
}

func (example) Test(_ context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{Message: "hello " + req.GetMessage()}, nil
}

// TestMain points the client at an Example service of the test, served over
// TLS like server/examples/examplesvcautotls.
func TestMain(m *testing.M) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	os.Setenv("CONF", "../../configs/config.yaml")
	os.Setenv("BOOST_FACTORY_GRPC_CLIENT_PORT", strconv.Itoa(lis.Addr().(*net.TCPAddr).Port))
	boost.Start()

	srv := grpcserver.NewServerWithOptions(context.Background(), &grpcserver.Options{
		TLS: grpcserver.TLSOptions{Enabled: true, Type: "AUTO", Auto: grpcserver.TLSAutoOptions{Host: "localhost"}},
	})
	pb.RegisterExampleServer(srv.ServiceRegistrar(), example{})
	go func() { _ = srv.Server().Serve(lis) }()

	code := m.Run()
	srv.Server().Stop()
	os.Exit(code)
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(`{"message":"world"}`))
	req.Header.Set("Content-Type", "application/json")
	srv.Instance().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"hello world"}`, rec.Body.String(), "the handler called the Test rpc over TLS")
}
