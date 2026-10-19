package local

import (
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/config/examples/env/internal/pkg/domain/model/entity"
)

// The key of the name of the application, env APP_APPLICATION_NAME.
const name = "app.application.name"

func init() {
	config.Add(name, "app_test", "name of the application")
}

// Config is the config under app.
type Config struct {
	Application struct {
		Name string
	}
}

func (c Config) entity() *entity.Application {
	return &entity.Application{Name: c.Application.Name}
}
