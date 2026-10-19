package client

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"storj.io/drpc/drpcconn"
)

// ran: Module and contextfx are once-guarded per process, so a rerun
// (-count>1) gets an empty graph.
var ran bool

func TestModuleDialsOnStartAndClosesOnStop(t *testing.T) {
	if ran {
		t.Skip("Module is once-guarded; the first run covers it")
	}
	ran = true

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() { _ = c.Close() })
		}
	}()
	t.Setenv("BOOST_FACTORY_DRPC_CLIENT_HOST", "127.0.0.1")
	t.Setenv("BOOST_FACTORY_DRPC_CLIENT_PORT", strconv.Itoa(ln.Addr().(*net.TCPAddr).Port))
	boost.Start()

	var conn *drpcconn.Conn
	app := fxtest.New(t, Module(), fx.Populate(&conn)).RequireStart()
	require.NotNil(t, conn)

	app.RequireStop()

	select {
	case <-conn.Closed():
	case <-time.After(5 * time.Second):
		t.Fatal("conn not closed on stop")
	}
	assert.Equal(t, fx.Options(), Module(), "a second Module is empty")
}
