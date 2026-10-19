package middlewares

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/builder"
)

// TestGeneratedCodeIsUpToDate fails when gen/ or a main under cmd/ differs
// from what go generate writes for the annotations of this example.
func TestGeneratedCodeIsUpToDate(t *testing.T) {
	stale, err := builder.Stale(context.Background(), ".", "gen", "cmd")
	require.NoError(t, err)
	assert.Empty(t, stale, "run go generate ./... in bootstrap/function/examples/middlewares")
}
