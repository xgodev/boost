package kinesis

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	client "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis"
	"go.uber.org/fx"
)

func TestModuleProvidesTheSDKAndTheBoostClient(t *testing.T) {
	empty := filepath.Join(t.TempDir(), "empty")
	require.NoError(t, os.WriteFile(empty, nil, 0o600))
	t.Setenv("AWS_CONFIG_FILE", empty)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", empty)
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("BOOST_FACTORY_AWS_DEFAULT__REGION", "sa-east-1")
	boost.Start()

	var sdk *kinesis.Client
	var c client.Client
	app := fx.New(fx.NopLogger, Module(), fx.Populate(&sdk, &c))
	require.NoError(t, app.Err())

	require.NotNil(t, sdk)
	assert.Equal(t, "sa-east-1", sdk.Options().Region)
	assert.NotNil(t, c)
}
