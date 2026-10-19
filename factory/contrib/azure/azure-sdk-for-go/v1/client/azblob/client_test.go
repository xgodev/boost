package azblob

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// token is a credential that is never asked for a token in these tests.
type token struct{}

func (token) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{}, errors.New("not used")
}

func TestNewClientWithOptionsUsesTheAccountURL(t *testing.T) {
	c, err := NewClientWithOptions(context.Background(), token{}, &Options{AccountName: "orders"})

	require.NoError(t, err)
	assert.Equal(t, "https://orders.blob.core.windows.net/", c.URL())
}

func TestNewClientWithOptionsUsesTheConfiguredURL(t *testing.T) {
	c, err := NewClientWithOptions(context.Background(), token{}, &Options{AccountName: "orders", URL: "https://blob.example.com/orders/"})

	require.NoError(t, err)
	assert.Equal(t, "https://blob.example.com/orders/", c.URL())
}

func TestNewClientWithOptionsUsesTheConnectionString(t *testing.T) {
	cs := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"

	c, err := NewClientWithOptions(context.Background(), nil, &Options{AccountName: "orders", ConnectionString: cs})

	require.NoError(t, err)
	assert.Equal(t, "http://127.0.0.1:10000/devstoreaccount1/", c.URL())
}

func TestPluginsCustomizeTheClientOptions(t *testing.T) {
	var got *azblob.ClientOptions
	plugin := func(_ context.Context, o *azblob.ClientOptions) error { got = o; return nil }

	_, err := NewClientWithOptions(context.Background(), token{}, &Options{AccountName: "orders"}, plugin)

	require.NoError(t, err)
	assert.NotNil(t, got)
}

func TestFailingPluginStopsTheBuild(t *testing.T) {
	boom := errors.New("boom")
	failing := func(context.Context, *azblob.ClientOptions) error { return boom }

	c, err := NewClientWithOptions(context.Background(), token{}, &Options{AccountName: "orders"}, failing)

	assert.Nil(t, c)
	assert.ErrorIs(t, err, boom)
}

func TestNewOptionsDefaults(t *testing.T) {
	opts, err := NewOptions()

	require.NoError(t, err)
	assert.Equal(t, &Options{AccountName: "unknown"}, opts)
}

func TestNewOptionsWithPathReadsTheEnv(t *testing.T) {
	t.Cleanup(config.Load)
	t.Setenv("APP_BLOB_ACCOUNT__NAME", "orders")
	t.Setenv("APP_BLOB_URL", "https://blob.example.com/")

	opts, err := NewOptionsWithPath("app.blob")

	require.NoError(t, err)
	assert.Equal(t, "orders", opts.AccountName)
	assert.Equal(t, "https://blob.example.com/", opts.URL)
}
