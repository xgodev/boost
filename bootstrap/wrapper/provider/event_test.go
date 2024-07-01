package provider

import (
	"context"
	"errors"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/bootstrap/repository/mocks"
	"github.com/xgodev/boost/config"
	iglog "github.com/xgodev/boost/factory/local/wrapper/log"
)

type WrapperEventSuite struct {
	suite.Suite
}

func (s *WrapperEventSuite) SetupSuite() {
	config.Load()
	iglog.New()
}

func (s *WrapperEventSuite) TestWrapperEvent_Publish() {

	event := v2.NewEvent()
	event.SetID("changeme")
	event.SetSubject("changeme")
	event.SetSource("changeme")
	event.SetType("changeme")
	event.SetExtension("partitionkey", "changeme")
	event.SetData("", nil)

	type fields struct {
		events *mocks.Event
		pkg    string
		impl   string
	}

	type args struct {
		ctx    context.Context
		events []*v2.Event
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func(eventPublisher *mocks.Event)
	}{
		{
			name: "when push one message successfully",
			fields: fields{
				events: new(mocks.Event),
				pkg:    "changeme",
				impl:   "changeme",
			},
			args: args{
				ctx:    context.Background(),
				events: []*v2.Event{&event},
			},
			wantErr: false,
			mock: func(events *mocks.Event) {
				events.On("Publish", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "when push messages failed",
			fields: fields{
				events: new(mocks.Event),
				pkg:    "changeme",
				impl:   "changeme",
			},
			args: args{
				ctx:    context.Background(),
				events: []*v2.Event{&event},
			},
			wantErr: true,
			mock: func(events *mocks.Event) {
				events.On("Publish", mock.Anything, mock.Anything).Return(errors.New("Error")).Times(1)
			},
		},
	}
	for _, tt := range tests {

		s.Run(tt.name, func() {

			tt.mock(tt.fields.events)

			p := NewEventWrapperProvider(tt.fields.events)

			err := p.Publish(tt.args.ctx, tt.args.events)

			s.Assert().True((err != nil) == tt.wantErr, "Publish() error = %v, wantErr %v", err, tt.wantErr)

			tt.fields.events.AssertExpectations(s.T())
		})
	}
}

func TestWrapperEventSuite(t *testing.T) {
	suite.Run(t, new(WrapperEventSuite))
}
