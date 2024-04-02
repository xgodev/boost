package cache

import "context"

type Cacheable[T any] func(ctx context.Context) (T, error)
