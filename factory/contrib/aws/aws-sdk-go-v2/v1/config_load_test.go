package aws

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfigWithOptions_ReturnsTheLoadError(t *testing.T) {
	empty := filepath.Join(t.TempDir(), "empty")
	if err := os.WriteFile(empty, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AWS_CONFIG_FILE", empty)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", empty)
	t.Setenv("AWS_PROFILE", "boost-profile-that-does-not-exist")

	if _, err := NewConfigWithOptions(context.Background(), &Options{}); err == nil {
		t.Fatal("NewConfigWithOptions() error = nil, want the error loading a profile that does not exist")
	}
}
