package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
)

const versionEnv = "BOOST_FACTORY_GCP_PUBSUB_VERSION"

// load sets the version through env and reloads the config. The koanf
// instance is global and keeps what earlier loads read, so every test sets
// the value it needs instead of relying on it being unset.
func load(t *testing.T, version string) {
	t.Helper()
	t.Setenv(versionEnv, version)
	config.Set(koanf.New())
	config.Load()
}

func TestVersionKeyIsDeclaredOnceWithV2(t *testing.T) {
	var defaults []any
	for _, e := range config.Entries() {
		if e.Key == VersionKey {
			defaults = append(defaults, e.Value)
		}
	}
	assert.Equal(t, []any{V2}, defaults)
}

func TestVersion(t *testing.T) {
	for _, v := range []string{V1, V2} {
		t.Run(v, func(t *testing.T) {
			load(t, v)

			got, err := Version(VersionKey)
			require.NoError(t, err)
			assert.Equal(t, v, got)
		})
	}
}

func TestVersionUnknownIsNotValid(t *testing.T) {
	load(t, "v3")

	got, err := Version(VersionKey)
	assert.Empty(t, got)
	assert.True(t, errors.IsNotValid(err), "got %v", err)
	assert.ErrorContains(t, err, VersionKey)
	assert.ErrorContains(t, err, `"v3"`)
}

func TestEnabled(t *testing.T) {
	tests := []struct {
		version string
		v1, v2  bool
	}{
		{version: V1, v1: true},
		{version: V2, v2: true},
		{version: "v3"},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			load(t, tt.version)

			assert.Equal(t, tt.v1, Enabled(VersionKey, V1))
			assert.Equal(t, tt.v2, Enabled(VersionKey, V2))
		})
	}
}
