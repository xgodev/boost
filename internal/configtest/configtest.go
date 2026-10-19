// Package configtest sets boost config values from tests.
package configtest

import (
	"os"
	"testing"

	"github.com/xgodev/boost"
)

// Set exports env=value and reloads the config with boost.Start. The config
// keeps a loaded value after its variable is gone, so on cleanup Set loads
// restore -- the key's default -- and unsets env again; a bad value set by
// one test never reaches the next.
func Set(t *testing.T, env, value, restore string) {
	t.Helper()
	if err := os.Setenv(env, value); err != nil {
		t.Fatalf("setenv %s: %v", env, err)
	}
	boost.Start()
	t.Cleanup(func() {
		if err := os.Setenv(env, restore); err != nil {
			t.Errorf("setenv %s: %v", env, err)
			return
		}
		boost.Start()
		if err := os.Unsetenv(env); err != nil {
			t.Errorf("unsetenv %s: %v", env, err)
		}
	})
}
