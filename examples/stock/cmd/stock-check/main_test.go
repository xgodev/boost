package main

import (
	"context"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/job"
	stockcheck "github.com/xgodev/boost/examples/stock/gen/builder/stock-check"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/finder"
	"go.uber.org/fx"
)

// TestMain boots boost: the stock and the check read their config.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// TestStockCheck builds app stock-check from the generated Modules() and
// runs the job handler it provides, as the job module does once on start.
func TestStockCheck(t *testing.T) {
	var (
		check job.Handler
		svc   finder.Service
	)
	app := fx.New(append(stockcheck.Modules(), fx.NopLogger, fx.Populate(&check, &svc))...)
	require.NoError(t, app.Err())

	assert.Contains(t, reflect.TypeOf(svc).Elem().PkgPath(), "/gen/builder/middleware/",
		"the check gets the finder wrapped in the middlewares of @Middleware")
	require.NoError(t, check(context.Background()), "the memory stock has tshirt and mug")
}

const (
	// mainEnv makes the process the TestMain_* tests start run the generated
	// main instead of the tests.
	mainEnv = "STOCK_EXAMPLE_MAIN"
	// skusEnv sets the skus the check looks at
	// (internal.app.stockcheck.job.stock.checker.skus).
	skusEnv = "INTERNAL_APP_STOCKCHECK_JOB_STOCK_CHECKER_SKUS"
)

// TestMain_ExitsZeroWhenTheCheckPasses runs the generated main of app
// stock-check on the seeded memory stock: the check returns nil and the
// process exits 0 by itself.
func TestMain_ExitsZeroWhenTheCheckPasses(t *testing.T) {
	out, err := runMain(t)

	require.NoError(t, err, "the check did not exit 0: %s", out)
	assert.Contains(t, out, "tshirt has 3 units")
	assert.Contains(t, out, "mug has 1 units")
}

// TestMain_ExitsOneWhenTheCheckFails runs the generated main of app
// stock-check on a sku the stock does not have: the check returns the
// NotFound, which the app logs before exiting 1.
func TestMain_ExitsOneWhenTheCheckFails(t *testing.T) {
	out, err := runMain(t, skusEnv+"=tshirt,cap")

	var exit *exec.ExitError
	require.ErrorAs(t, err, &exit, "the check must fail: %s", out)
	assert.Equal(t, 1, exit.ExitCode())
	assert.Contains(t, out, "job failed: checking cap")
}

// runMain runs the generated main of app stock-check with env and returns
// its output and exit, failing the test if it does not exit by itself.
func runMain(t *testing.T, env ...string) (string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestMainChild$")
	cmd.Env = append(append(os.Environ(), env...), mainEnv+"=1")
	out, err := cmd.CombinedOutput()
	require.NoError(t, ctx.Err(), "the check did not exit by itself: %s", out)
	return string(out), err
}

// TestMainChild runs the generated main, with no argument as a deployed
// binary gets, in the process the TestMain_* tests start.
func TestMainChild(t *testing.T) {
	if os.Getenv(mainEnv) != "1" {
		t.Skip("run by the TestMain_* tests")
	}
	os.Args = os.Args[:1]
	main()
}
