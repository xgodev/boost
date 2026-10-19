package log

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestDisabledAddsNothing(t *testing.T) {
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, NewLogWithOptions(&Options{}).Register)
	require.NoError(t, err)
	assert.Empty(t, srv.Mux().Middlewares())
}

func TestRegisterPutsTheRequestLoggerInTheContext(t *testing.T) {
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, Register)
	require.NoError(t, err)
	var got log.Logger
	srv.Get("/x", func(w http.ResponseWriter, r *http.Request) {
		got = log.FromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	srv.Mux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.NotNil(t, got)
}

func TestRegisterKeepsTheHandlerResponse(t *testing.T) {
	for _, level := range []string{"TRACE", "DEBUG", "INFO"} {
		t.Run(level, func(t *testing.T) {
			srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{},
				NewLogWithOptions(&Options{Enabled: true, Level: level}).Register)
			require.NoError(t, err)
			srv.Post("/x", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte("made"))
			})

			rec := httptest.NewRecorder()
			srv.Mux().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/x", nil))

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "made", rec.Body.String())
		})
	}
}
