package server

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	ms "github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/storj.io/drpc/v1/server"
	"github.com/xgodev/boost/fx/modules/extra/multiserver"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"storj.io/drpc/drpcmux"
)

// ran: Module and contextfx are once-guarded per process, so a rerun
// (-count>1) gets an empty graph.
var ran bool

type servers struct {
	fx.In
	Servers []ms.Server `group:"extra.multiserver.servers"`
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln.Close())
	return port
}

func TestModuleServesOnStartAndStopsOnStop(t *testing.T) {
	if ran {
		t.Skip("Module is once-guarded; the first run covers it")
	}
	ran = true
	port := freePort(t)
	t.Setenv("BOOST_FACTORY_DRPC_SERVER_PORT", strconv.Itoa(port))
	boost.Start()

	var (
		srv   *server.Server
		mux   *drpcmux.Mux
		group []ms.Server
	)
	app := fxtest.New(t,
		Module(),
		multiserver.Module(),
		fx.Invoke(func(s servers) { group = s.Servers }),
		fx.Populate(&srv, &mux),
	).RequireStart()

	assert.Equal(t, port, srv.ServerOptions().Port, "the port comes from config")
	assert.Same(t, srv.Mux(), mux)
	assert.Equal(t, []ms.Server{srv}, group)

	addr := "127.0.0.1:" + strconv.Itoa(port)
	require.Eventually(t, func() bool {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			return false
		}
		_ = c.Close()
		return true
	}, 5*time.Second, 10*time.Millisecond)

	app.RequireStop()

	require.Eventually(t, func() bool {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			return true
		}
		_ = c.Close()
		return false
	}, 5*time.Second, 10*time.Millisecond, "the server is down after the app stops")
	assert.Equal(t, fx.Options(), Module(), "a second Module is empty")
}
