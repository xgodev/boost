package server

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	ms "github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type servers struct {
	fx.In
	All []ms.Server `group:"extra.multiserver.servers"`
}

// ran: Module and contextfx are once-guarded per process, so a rerun
// (-count>1) gets an empty graph.
var ran bool

func TestModule(t *testing.T) {
	if ran {
		t.Skip("Module is once-guarded; the first run covers it")
	}
	ran = true
	boost.Start()

	var (
		srv  *server.Server
		inst *mcp.Server
		in   servers
	)
	fxtest.New(t, Module(), fx.Populate(&srv, &inst, &in)).RequireStart().RequireStop()

	require.NotNil(t, srv)
	assert.Same(t, srv.Instance(), inst)
	require.Len(t, in.All, 1)
	assert.Same(t, srv, in.All[0])
	assert.Equal(t, fx.Options(), Module(), "once: a second Module adds nothing")
}
