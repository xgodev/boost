package fiber

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

// Options fiber server options.
type Options struct {
	Port   int
	Type   string
	Config Config
}

// Config holds the fiber.Config settings read from config.
type Config struct {
	Prefork                   bool
	ServerHeader              string
	StrictRouting             bool
	CaseSensitive             bool
	Immutable                 bool
	UnescapePath              bool
	ETag                      bool `config:"etag"`
	BodyLimit                 int
	Concurrency               int
	ReadTimeout               time.Duration
	WriteTimeout              time.Duration
	IdleTimeout               time.Duration
	ReadBufferSize            int
	WriteBufferSize           int
	CompressedFileSuffix      string
	ProxyHeader               string
	GETOnly                   bool `config:"getOnly"`
	DisableKeepalive          bool
	DisableDefaultDate        bool
	DisableDefaultContentType bool
	DisableHeaderNormalizing  bool
	DisableStartupMessage     bool
	ReduceMemoryUsage         bool
	Network                   string
	EnableTrustedProxyCheck   bool
	TrustedProxies            []string
}

// fiberConfig returns the fiber.Config the settings describe.
func (c Config) fiberConfig() fiber.Config {
	return fiber.Config{
		Prefork:                   c.Prefork,
		ServerHeader:              c.ServerHeader,
		StrictRouting:             c.StrictRouting,
		CaseSensitive:             c.CaseSensitive,
		Immutable:                 c.Immutable,
		UnescapePath:              c.UnescapePath,
		ETag:                      c.ETag,
		BodyLimit:                 c.BodyLimit,
		Concurrency:               c.Concurrency,
		ReadTimeout:               c.ReadTimeout,
		WriteTimeout:              c.WriteTimeout,
		IdleTimeout:               c.IdleTimeout,
		ReadBufferSize:            c.ReadBufferSize,
		WriteBufferSize:           c.WriteBufferSize,
		CompressedFileSuffix:      c.CompressedFileSuffix,
		ProxyHeader:               c.ProxyHeader,
		GETOnly:                   c.GETOnly,
		DisableKeepalive:          c.DisableKeepalive,
		DisableDefaultDate:        c.DisableDefaultDate,
		DisableDefaultContentType: c.DisableDefaultContentType,
		DisableHeaderNormalizing:  c.DisableHeaderNormalizing,
		DisableStartupMessage:     c.DisableStartupMessage,
		ReduceMemoryUsage:         c.ReduceMemoryUsage,
		Network:                   c.Network,
		EnableTrustedProxyCheck:   c.EnableTrustedProxyCheck,
		TrustedProxies:            c.TrustedProxies,
	}
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	ConfigAdd(path)
	config.Load()
	return config.NewOptionsWithPath[Options](root, path)
}
