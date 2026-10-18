// Package buildinfo reads the application version Go stamps into the binary.
package buildinfo

import "runtime/debug"

// Unknown is the version when the binary carries none (go run, go test, no
// VCS tag); it matches the otel and datadog factories' version default.
const Unknown = "0.0.0"

var readBuildInfo = debug.ReadBuildInfo

// Version returns the main module version of the running binary, e.g. v1.4.2
// for a build from a tagged commit.
func Version() string {
	info, ok := readBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return Unknown
	}
	return info.Main.Version
}
