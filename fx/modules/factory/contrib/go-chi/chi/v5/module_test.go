package chi

import (
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	server "github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/fx/modules/extra/multiserver"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// ran: Module and contextfx are once-guarded per process, so a rerun
// (-count>1) gets an empty graph.
var ran bool

type servers struct {
	fx.In
	Servers []server.Server `group:"extra.multiserver.servers"`
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
	t.Setenv("BOOST_FACTORY_CHI_PORT", strconv.Itoa(port))
	boost.Start()

	var (
		srv   *chi.Server
		group []server.Server
	)
	app := fxtest.New(t,
		Module(),
		multiserver.Module(),
		fx.Invoke(func(s *chi.Server) {
			s.Get("/ping", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("pong")) })
		}),
		fx.Invoke(func(s servers) { group = s.Servers }),
		fx.Populate(&srv),
	).RequireStart()

	assert.Equal(t, port, srv.ServerOptions().Port, "the port comes from config")
	assert.Equal(t, []server.Server{srv}, group)

	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	url := "http://127.0.0.1:" + strconv.Itoa(port) + "/ping"
	require.Eventually(t, func() bool {
		resp, err := client.Get(url)
		if err != nil {
			return false
		}
		defer func() { _ = resp.Body.Close() }()
		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 10*time.Millisecond)

	app.RequireStop()

	_, err := client.Get(url)
	assert.Error(t, err, "the server is down after the app stops")
	assert.Equal(t, fx.Options(), Module(), "a second Module is empty")
}
