package cache

import (
	"context"
)

type Driver interface {
	Set(ctx context.Context, key string, data []byte) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (data []byte, err error)
}
