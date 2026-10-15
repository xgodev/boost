package otelsql_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/core/database/sql"
	otelsql "github.com/xgodev/boost/factory/core/database/sql/plugins/contrib/xsam/otelsql/v0"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

const instrumentationScope = "github.com/XSAM/otelsql"

// In an fx app the *sql.DB is built before any other component starts the
// boost OTel providers, so the plugin must not depend on that ordering.
func TestOTel_DBBuiltBeforeProvidersStartReportsToThem(t *testing.T) {
	boost.Start()
	ctx := context.Background()

	db, err := sql.NewDB(ctx, fakeConnector{}, otelsql.NewOTel())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() }) // intentional: fake connections hold no resources

	opts, err := otelboost.NewOptions()
	require.NoError(t, err)

	spans := tracetest.NewInMemoryExporter()
	otelboost.StartTracerProviderWithOptions(ctx, opts, sdktrace.WithSyncer(spans))

	reader := sdkmetric.NewManualReader()
	otelboost.StartMetricProviderWithOptions(ctx, opts, sdkmetric.WithReader(reader))

	_, err = db.ExecContext(ctx, "SELECT 1")
	require.NoError(t, err)

	t.Run("spans", func(t *testing.T) {
		var got []string
		for _, s := range spans.GetSpans() {
			if s.InstrumentationScope.Name == instrumentationScope {
				got = append(got, s.Name)
			}
		}
		require.NotEmpty(t, got, "no otelsql span reached the tracer provider started after the DB")
	})

	t.Run("db stats metrics", func(t *testing.T) {
		var rm metricdata.ResourceMetrics
		require.NoError(t, reader.Collect(ctx, &rm))

		var got []string
		for _, sm := range rm.ScopeMetrics {
			if sm.Scope.Name != instrumentationScope {
				continue
			}
			for _, m := range sm.Metrics {
				got = append(got, m.Name)
			}
		}
		require.NotEmpty(t, got, "no otelsql metric reached the meter provider started after the DB")
	})
}

var errNotSupported = errors.New("not supported by fake driver")

type fakeConnector struct{}

func (fakeConnector) Connect(context.Context) (driver.Conn, error) { return fakeConn{}, nil }
func (fakeConnector) Driver() driver.Driver                        { return fakeDriver{} }

type fakeDriver struct{}

func (fakeDriver) Open(string) (driver.Conn, error) { return fakeConn{}, nil }

type fakeConn struct{}

func (fakeConn) Prepare(string) (driver.Stmt, error) { return nil, errNotSupported }
func (fakeConn) Close() error                        { return nil }
func (fakeConn) Begin() (driver.Tx, error)           { return nil, errNotSupported }

func (fakeConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return driver.RowsAffected(0), nil
}
