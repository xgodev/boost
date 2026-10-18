package buildinfo

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stub(t *testing.T, version string, ok bool) {
	t.Helper()
	prev := readBuildInfo
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Main: debug.Module{Version: version}}, ok
	}
	t.Cleanup(func() { readBuildInfo = prev })
}

func TestVersion_MainModuleVersion(t *testing.T) {
	stub(t, "v1.4.2", true)
	assert.Equal(t, "v1.4.2", Version())
}

func TestVersion_DevelBuild(t *testing.T) {
	stub(t, "(devel)", true)
	assert.Equal(t, Unknown, Version())
}

func TestVersion_EmptyVersion(t *testing.T) {
	stub(t, "", true)
	assert.Equal(t, Unknown, Version())
}

func TestVersion_NoBuildInfo(t *testing.T) {
	stub(t, "v1.4.2", false)
	assert.Equal(t, Unknown, Version())
}
