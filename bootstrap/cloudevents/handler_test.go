package cloudevents

import (
	"context"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus/v1"
	"reflect"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/config"
)

type CloudEventsHandlerSuite struct {
	suite.Suite
}

func TestCloudEventsHandlerSuite(t *testing.T) {
	suite.Run(t, new(CloudEventsHandlerSuite))
}

func (s *CloudEventsHandlerSuite) SetupSuite() {
	config.Load()
	logrus.NewLogger()
}

func (s *CloudEventsHandlerSuite) TestCloudEventsNewHandler() {

	type args struct {
		handler Handler
	}

	handler := func(ctx context.Context, in v2.Event) (*v2.Event, error) { return nil, nil }
	hwOptions, _ := DefaultHandlerWrapperOptions()
	hw := NewHandlerWrapper(handler, hwOptions)

	tests := []struct {
		name string
		want *Handler
	}{
		{
			name: "success",
			want: NewHandler(hw),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got := NewHandler(hw)
			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewHandler() = %v, want %v", got, tt.want)
		})
	}
}

func (s *CloudEventsHandlerSuite) TestCloudEventsHandler_Handle() {

	handler := func(ctx context.Context, in v2.Event) (*v2.Event, error) { return nil, nil }
	hwOptions, _ := DefaultHandlerWrapperOptions()
	hw := NewHandlerWrapper(handler, hwOptions)

	type fields struct {
		handler *HandlerWrapper
	}

	type args struct {
		ctx context.Context
		in  func() v2.Event
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "on kinesis success event",
			fields: fields{
				handler: hw,
			},
			args: args{
				ctx: context.Background(),
				in: func() v2.Event {
					e := v2.NewEvent()
					e.SetSubject("changeme")
					e.SetSource("changeme")
					e.SetType("changeme")
					e.SetData("", "changeme")
					return e
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			h := NewHandler(tt.fields.handler)

			_, err := h.Handle(tt.args.ctx, tt.args.in())
			s.Assert().True((err != nil) == tt.wantErr, "Handle() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
