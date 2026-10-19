package probe

import (
	"context"

	"github.com/labstack/echo/v4"
	grpcserver "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Plugins builds the admin server: every response says which server served
// it. A []Plugin provided under the name of a server reaches only that
// server, rest and grpc included; rest and grpc also take the unnamed one.
// @Provide(name=admin, index=0)
func Plugins() []boostecho.Plugin {
	return []boostecho.Plugin{servedBy("admin")}
}

// RestPlugins builds the rest server as Plugins builds admin.
// @Provide(name=rest, index=0)
func RestPlugins() []boostecho.Plugin {
	return []boostecho.Plugin{servedBy("rest")}
}

// GrpcPlugins builds the grpc server: every response carries the
// x-served-by header.
// @Provide(name=grpc, index=0)
func GrpcPlugins() []grpcserver.Plugin {
	return []grpcserver.Plugin{func(context.Context) []grpc.ServerOption {
		return []grpc.ServerOption{grpc.ChainUnaryInterceptor(
			func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
				if err := grpc.SetHeader(ctx, metadata.Pairs("x-served-by", "grpc")); err != nil {
					return nil, err
				}
				return handler(ctx, req)
			})}
	}}
}

func servedBy(server string) boostecho.Plugin {
	return func(_ context.Context, srv *boostecho.Server) error {
		srv.Instance().Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set("X-Served-By", server)
				return next(c)
			}
		})
		return nil
	}
}
