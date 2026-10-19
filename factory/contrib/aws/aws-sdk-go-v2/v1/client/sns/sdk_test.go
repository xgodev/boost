package sns

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	awsfactory "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
)

func TestNewSDKClientWithOptions(t *testing.T) {
	empty := filepath.Join(t.TempDir(), "empty")
	if err := os.WriteFile(empty, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AWS_CONFIG_FILE", empty)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", empty)
	t.Setenv("AWS_PROFILE", "")

	c, err := NewSDKClientWithOptions(context.Background(), &awsfactory.Options{DefaultRegion: "sa-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Options().Region; got != "sa-east-1" {
		t.Errorf("Region = %q, want sa-east-1", got)
	}
}
