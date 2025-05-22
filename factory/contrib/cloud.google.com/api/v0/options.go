package api

import (
	"strings"
	"time"
)

// Options holds shared API-level GCP client configuration.
type Options struct {
	ProjectID   string `config:"projectId"`
	Credentials struct {
		File string `config:"file"`
		JSON string `config:"json"`
	} `config:"credentials"`
	Endpoint     string        `config:"endpoint"`
	UseEmulator  bool          `config:"useEmulator"`
	EmulatorHost string        `config:"emulatorHost"`
	UserAgent    string        `config:"userAgent"`
	Scopes       []string      `config:"scopes"`
	Timeout      time.Duration `config:"timeout"`
	Proxy        string        `config:"proxy"`
	Retry        struct {
		MaxAttempts    int           `config:"maxAttempts"`
		InitialBackoff time.Duration `config:"initialBackoff"`
		MaxBackoff     time.Duration `config:"maxBackoff"`
		Multiplier     float64       `config:"multiplier"`
	} `config:"retry"`
}

// ParseScopes splits comma-separated scopes.
func (o *Options) ParseScopes() []string {
	var out []string
	for _, item := range o.Scopes {
		for _, part := range strings.Split(item, ",") {
			if s := strings.TrimSpace(part); s != "" {
				out = append(out, s)
			}
		}
	}
	return out
}
