package confluent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"github.com/xgodev/boost/wrapper/config/model"
)

// Issue #48: every key under this adapter's root must be reachable from env
// vars. The old root ".kafka_confluent" produced a segment with a literal
// underscore, which the env loader can never emit, so BOOST_..._CONFLUENT_*
// silently did nothing.
func TestConfig_EnvVarReachesTopics(t *testing.T) {
	t.Setenv("BOOST_BOOTSTRAP_FUNCTION_ADAPTER_CONFLUENT_TOPICS", "orders")

	config.Set(koanf.New())
	config.Load()

	opts, err := DefaultOptions()
	assert.NoError(t, err)
	assert.Equal(t, []string{"orders"}, opts.Topics)
}

// Backward compat: values set under the legacy ".kafka_confluent" root via a
// config file must keep working after the rename to ".confluent".
func TestConfig_LegacyFileNamespaceStillWorks(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "legacy.json")
	err := os.WriteFile(f, []byte(
		`{"boost":{"bootstrap":{"function":{"adapter":{"kafka_confluent":{"topics":["legacy"]}}}}}}`), 0o644)
	assert.NoError(t, err)
	t.Setenv(model.ConfEnvironment, f)

	config.Set(koanf.New())
	config.Load()

	opts, err := DefaultOptions()
	assert.NoError(t, err)
	assert.Equal(t, []string{"legacy"}, opts.Topics)
}
