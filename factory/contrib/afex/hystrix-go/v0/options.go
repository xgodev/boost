package hystrix

import (
	"github.com/xgodev/boost"
	"strings"
)

type Options struct {
	Enabled                bool
	Timeout                int
	RequestVolumeThreshold int
	ErrorPercentThreshold  int
	MaxConcurrentRequests  int
	SleepWindow            int
}

// NewOptionsFromCommand unmarshals options based a given key path.
func NewOptionsFromCommand(cmd string) (*Options, error) {
	path := strings.Join([]string{cmdRoot, cmd}, ".")
	return boost.NewOptionsWithPath[Options](path)
}
