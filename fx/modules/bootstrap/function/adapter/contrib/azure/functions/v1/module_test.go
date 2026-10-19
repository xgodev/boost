package azure

import (
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	fxfunction "github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestModuleAddsTheAzureAdapter(t *testing.T) {
	var adapters []function.CmdFunc[*event.Event]
	app := fx.New(
		fx.NopLogger,
		Module[*event.Event](),
		fx.Invoke(fx.Annotate(
			func(a []function.CmdFunc[*event.Event]) { adapters = a },
			fx.ParamTags(`group:"`+fxfunction.BSFunctionAdaptersGroupKey+`"`),
		)),
	)
	require.NoError(t, app.Err())

	require.Len(t, adapters, 1)
	assert.Equal(t, "azure", adapters[0](nil).Use)
}
