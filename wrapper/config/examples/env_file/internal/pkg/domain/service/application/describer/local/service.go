// Package local describes the application from its config.
package local

import (
	"context"

	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/config/examples/env_file/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/wrapper/config/examples/env_file/internal/pkg/domain/service/application/describer"
)

type service struct{}

// New provides the describer.
// @Provide(index=0)
func New() describer.Service {
	return &service{}
}

// Describe reads the keys under app on every call, so a config.Load after an
// env var changed shows in the next answer.
func (s *service) Describe(context.Context) (*entity.Application, error) {
	var c Config
	if err := config.UnmarshalWithPath("app", &c); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error reading the config of the application"))
	}
	return c.entity(), nil
}
