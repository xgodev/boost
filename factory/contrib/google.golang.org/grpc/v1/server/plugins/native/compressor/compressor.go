package compressor

import (
	"context"

	"github.com/xgodev/boost/wrapper/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

// Register registers a new compressor plugin for grpc server.
func Register(ctx context.Context) []grpc.ServerOption {
	o, err := NewOptions()
	if err != nil {
		return nil
	}
	h := NewCompressorWithOptions(o)
	return h.Register(ctx)
}

// Compressor represents compressor plugin for grpc server.
type Compressor struct {
	options *Options
}

// NewCompressorWithOptions returns a new compressor plugin with options.
func NewCompressorWithOptions(options *Options) *Compressor {
	return &Compressor{options: options}
}

// NewCompressorWithConfigPath returns a new compressor plugin with options from config path.
func NewCompressorWithConfigPath(path string) (*Compressor, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewCompressorWithOptions(o), nil
}

// NewCompressor returns a new compressor plugin with default options.
func NewCompressor() *Compressor {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewCompressorWithOptions(o)
}

// Register registers this compressor plugin for grpc server.
func (i *Compressor) Register(ctx context.Context) []grpc.ServerOption {

	logger := log.FromContext(ctx)
	logger.Debug("compressor successfully enabled in grpc server")

	err := gzip.SetLevel(i.options.Level)
	if err != nil {
		logger.Fatalf("could not set level: %s", err.Error())
	}

	return nil
}
