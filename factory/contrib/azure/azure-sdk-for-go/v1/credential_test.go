package azure

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestNewCredentialWithOptionsUsesClientSecretWhenComplete(t *testing.T) {
	cred, err := NewCredentialWithOptions(context.Background(), &Options{
		ClientId: "11111111-1111-1111-1111-111111111111", ClientSecret: "secret", TenantId: "22222222-2222-2222-2222-222222222222",
	})

	require.NoError(t, err)
	assert.IsType(t, &azidentity.ClientSecretCredential{}, cred)
}

func TestNewCredentialWithOptionsFallsBackToDefaultChain(t *testing.T) {
	for _, opts := range []*Options{
		{},
		{ClientId: "id", ClientSecret: "secret"},
		{ClientId: "id", TenantId: "tenant"},
	} {
		cred, err := NewCredentialWithOptions(context.Background(), opts)

		require.NoError(t, err)
		assert.IsType(t, &azidentity.DefaultAzureCredential{}, cred)
	}
}

func TestNewOptionsDefaults(t *testing.T) {
	opts, err := NewOptions()

	require.NoError(t, err)
	assert.Equal(t, &Options{}, opts)
}

func TestNewOptionsWithPathReadsTheEnv(t *testing.T) {
	t.Cleanup(config.Load)
	t.Setenv("APP_AZ_CLIENT__ID", "id")
	t.Setenv("APP_AZ_CLIENT__SECRET", "secret")
	t.Setenv("APP_AZ_TENANT__ID", "tenant")

	opts, err := NewOptionsWithPath("app.az")

	require.NoError(t, err)
	assert.Equal(t, &Options{ClientId: "id", ClientSecret: "secret", TenantId: "tenant"}, opts)
}
