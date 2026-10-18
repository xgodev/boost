package ftp

import (
	"context"
	"sync"

	"github.com/jlaffaye/ftp"
	f "github.com/xgodev/boost/factory/contrib/jlaffaye/ftp/v0"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for ftp server connection. The connection quits on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle) (*ftp.ServerConn, error) {
					conn, err := f.NewServerConn(ctx)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(conn.Quit))
					return conn, nil
				},
			),
		)
	})

	return options
}
