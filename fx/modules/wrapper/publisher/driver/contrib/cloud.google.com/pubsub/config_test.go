package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
)

func defaults(key string) []any {
	var values []any
	for _, e := range config.Entries() {
		if e.Key == key {
			values = append(values, e.Value)
		}
	}
	return values
}

func TestKeysAreDeclaredOnce(t *testing.T) {
	assert.Equal(t, []any{pubsubcfg.V2}, defaults(VersionKey))
	assert.Equal(t, []any{true}, defaults(enabled))
}

func TestIsEnabled(t *testing.T) {
	for _, tt := range []struct {
		enabled string
		want    bool
	}{
		{enabled: "true", want: true},
		{enabled: "false", want: false},
	} {
		t.Run(tt.enabled, func(t *testing.T) {
			t.Setenv("BOOST_WRAPPER_PUBLISHER_DRIVER_PUBSUB_ENABLED", tt.enabled)
			config.Set(koanf.New())
			config.Load()

			assert.Equal(t, tt.want, IsEnabled())
		})
	}
}
