package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// ran: Module and contextfx are once-guarded per process, so a rerun
// (-count>1) gets an empty graph.
var ran bool

func TestModuleProvidesTheClient(t *testing.T) {
	if ran {
		t.Skip("Module is once-guarded; the first run covers it")
	}
	ran = true
	boost.Start()

	var c *fasthttp.Client
	fxtest.New(t, Module(), fx.Populate(&c)).RequireStart().RequireStop()

	assert.Equal(t, "fasthttp-client", c.Name)
	assert.Equal(t, fx.Options(), Module(), "a second Module is empty")
}
