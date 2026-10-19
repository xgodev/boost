package local

import (
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/config/examples/env_file/internal/pkg/domain/model/entity"
)

func init() {
	// env APP_APPLICATION_NAME
	config.Add("app.application.name", "app_test", "name of the application")
	// env APP_APPLICATION_MY__NAME
	config.Add("app.application.myName", "my_name_test", "my name of the application")
}

// Config is the config under app.
type Config struct {
	Application struct {
		Name   string
		MyName string
	}
}

func (c Config) entity() *entity.Application {
	return &entity.Application{Name: c.Application.Name, MyName: c.Application.MyName}
}
