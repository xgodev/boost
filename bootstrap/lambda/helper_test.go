package lambda

import (
	"context"
	"reflect"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"github.com/xgodev/boost/config"
	iglog "github.com/xgodev/boost/factory/local/wrapper/log"
)

type LambdaHelperSuite struct {
	suite.Suite
}

func TestLambdaHelperSuite(t *testing.T) {
	suite.Run(t, new(LambdaHelperSuite))
}

func (s *LambdaHelperSuite) SetupSuite() {
	config.Load()
	iglog.New()
}

func (s *LambdaHelperSuite) TestLambdaNewHelper() {

	type args struct {
		handler *cloudevents.HandlerWrapper
		options *Options
	}

	defaultOptions, _ := DefaultOptions()
	handler := func(ctx context.Context, in v2.Event) (*v2.Event, error) { return nil, nil }
	hwOptions, _ := cloudevents.DefaultHandlerWrapperOptions()
	hw := cloudevents.NewHandlerWrapper(handler, hwOptions)

	tests := []struct {
		name string
		args args
		want *Helper
	}{
		{
			name: "success",
			args: args{
				handler: hw,
				options: defaultOptions,
			},
			want: NewHelper(hw, defaultOptions),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got := NewHelper(hw, tt.args.options)
			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewHelper() = %v, want %v", got, tt.want)
		})
	}
}

func (s *LambdaHelperSuite) TestLambdaNewDefaultHelper() {

	type args struct {
		handler *cloudevents.HandlerWrapper
		options *Options
	}

	handler := func(ctx context.Context, in v2.Event) (*v2.Event, error) { return nil, nil }
	hwOptions, _ := cloudevents.DefaultHandlerWrapperOptions()
	hw := cloudevents.NewHandlerWrapper(handler, hwOptions)

	tests := []struct {
		name string
		args args
		want *Helper
	}{
		{
			name: "success",
			args: args{
				handler: hw,
			},
			want: NewDefaultHelper(hw),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got := NewDefaultHelper(hw)
			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewDefaultHelper() = %v, want %v", got, tt.want)
		})
	}
}
