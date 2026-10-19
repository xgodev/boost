// Package memory keeps the greeting in memory and reports its health.
package memory

import (
	"context"

	"github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/health/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/health/internal/pkg/domain/repository/greeting"
)

type datastore struct {
	greeting entity.Greeting
}

// New provides the greeting and adds its checker to extra/health, as the
// health plugins of boost's factories do for the clients they build.
// @Provide(index=0)
func New() greeting.Repository {
	d := &datastore{greeting: entity.Greeting{Message: "hello world"}}
	health.Add(health.NewHealthChecker("greetings", "the greeting in memory", d, true, true))
	return d
}

func (d *datastore) Get(context.Context) (*entity.Greeting, error) {
	g := d.greeting
	return &g, nil
}

// Check is always healthy: the greeting is in memory.
func (d *datastore) Check(context.Context) error {
	return nil
}
