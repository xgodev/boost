// Package redis keeps the counters on a redis client of boost's factory.
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/examples/client/internal/pkg/domain/repository/counter"
	"github.com/xgodev/boost/model/errors"
)

const prefix = "counter:"

type repository struct {
	client *redis.Client
}

// New provides the repository on the client the generator builds with the
// health plugin, which adds redis to the health checks.
// @Provide(index=0)
// @Inject(index=0, plugins=[health])
func New(client *redis.Client) counter.Repository {
	return &repository{client: client}
}

func (r *repository) Increment(ctx context.Context, name string) (int64, error) {
	count, err := r.client.Incr(ctx, prefix+name).Result()
	if err != nil {
		return 0, errors.Wrap(err, errors.Internalf("error incrementing counter %s", name))
	}
	return count, nil
}
