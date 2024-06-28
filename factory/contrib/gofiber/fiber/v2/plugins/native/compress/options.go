package compress

import (
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

// Options compress plugin for fiber options.
type Options struct {
	Enabled bool
	Level   int
}

// GetLevel returns current compress level.
func (o *Options) GetLevel() compress.Level {
	switch config.Int(level) {
	case -1:
		return compress.LevelDisabled
	case 1:
		return compress.LevelBestSpeed
	case 2:
		return compress.LevelBestCompression
	default:
		return compress.LevelDefault
	}
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return boost.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return boost.NewOptionsWithPath[Options](root, path)
}
