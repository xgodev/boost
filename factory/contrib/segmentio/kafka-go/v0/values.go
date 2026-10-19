package kafka

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

func balancer(name string) (kafka.Balancer, error) {
	switch strings.ToLower(name) {
	case "roundrobin":
		return &kafka.RoundRobin{}, nil
	case "leastbytes":
		return &kafka.LeastBytes{}, nil
	case "hash":
		return &kafka.Hash{}, nil
	case "crc32":
		return &kafka.CRC32Balancer{}, nil
	case "murmur2":
		return &kafka.Murmur2Balancer{}, nil
	}
	return nil, errors.NotValidf("kafka-go balancer %q", name)
}

func requiredAcks(name string) (kafka.RequiredAcks, error) {
	switch strings.ToLower(name) {
	case "none":
		return kafka.RequireNone, nil
	case "one":
		return kafka.RequireOne, nil
	case "all":
		return kafka.RequireAll, nil
	}
	return 0, errors.NotValidf("kafka-go requiredAcks %q", name)
}

func compression(name string) (kafka.Compression, error) {
	switch strings.ToLower(name) {
	case "", "none":
		return 0, nil
	case "gzip":
		return kafka.Gzip, nil
	case "snappy":
		return kafka.Snappy, nil
	case "lz4":
		return kafka.Lz4, nil
	case "zstd":
		return kafka.Zstd, nil
	}
	return 0, errors.NotValidf("kafka-go compression %q", name)
}

func startOffset(name string) (int64, error) {
	switch strings.ToLower(name) {
	case "first":
		return kafka.FirstOffset, nil
	case "last":
		return kafka.LastOffset, nil
	}
	return 0, errors.NotValidf("kafka-go startOffset %q", name)
}

// loggers adapts the boost logger to kafka-go: its chatty info log goes to
// trace, its errors to error.
func loggers(ctx context.Context) (kafka.Logger, kafka.Logger) {
	logger := log.FromContext(ctx)
	return kafka.LoggerFunc(logger.Tracef), kafka.LoggerFunc(logger.Errorf)
}
