package azblob

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/stretchr/testify/assert"
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
	t.Setenv("BOOST_FACTORY_AZURE_AZBLOB_ACCOUNT__NAME", "orders")
	boost.Start()

	var c *azblob.Client
	fxtest.New(t, Module(), fx.Populate(&c)).RequireStart().RequireStop()

	assert.Equal(t, "https://orders.blob.core.windows.net/", c.URL())
	assert.Equal(t, fx.Options(), Module(), "a second Module is empty")
}
