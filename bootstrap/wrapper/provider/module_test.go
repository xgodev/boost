package provider

import (
	"context"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/config"
	iglog "github.com/xgodev/boost/factory/local/wrapper/log"
	"go.uber.org/fx"
)

type WrapperModuleSuite struct {
	suite.Suite
}

func TestWrapperModuleSuite(t *testing.T) {
	suite.Run(t, new(WrapperModuleSuite))
}

func (s *WrapperModuleSuite) SetupSuite() {
	config.Load()
	iglog.New()
}

func (s *WrapperModuleSuite) TestWrapperEventModule() {

	eventModule := Module()
	var outs []*v2.Event

	app := fx.New(
		eventModule,
		fx.Invoke(func(s *EventWrapperProvider) error { return s.Publish(context.Background(), outs) }),
	)

	s.Assert().True(app.Err() == nil, "Error calling event module: %s", app.Err())
}
