package kinesis

import (
	"context"
	"errors"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis/mocks"
	iglog "github.com/xgodev/boost/factory/local/wrapper/log"
)

type ClientSuite struct {
	suite.Suite
}

func (s *ClientSuite) SetupSuite() {
	config.Load()
	iglog.New()
}

func (s *ClientSuite) TestClient_Publish() {

	event := v2.NewEvent()
	event.SetID("changeme")
	event.SetSubject("changeme")
	event.SetSource("changeme")
	event.SetType("changeme")
	event.SetExtension("partitionkey", "changeme")
	event.SetData("", nil)

	options, _ := DefaultOptions()

	type fields struct {
		client  *mocks.Client
		options *Options
	}

	type args struct {
		ctx      context.Context
		events   []*v2.Event
		resource string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func(eventPublisher *mocks.Client)
	}{
		{
			name: "when push one message successfully",
			fields: fields{
				client:  new(mocks.Client),
				options: options,
			},
			args: args{
				ctx:      context.Background(),
				events:   []*v2.Event{&event},
				resource: "subject",
			},
			wantErr: false,
			mock: func(client *mocks.Client) {
				client.On("Publish", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "when push messages successfully",
			fields: fields{
				client:  new(mocks.Client),
				options: options,
			},
			args: args{
				ctx:      context.Background(),
				events:   []*v2.Event{&event, &event},
				resource: "subject",
			},
			wantErr: false,
			mock: func(client *mocks.Client) {
				client.On("BulkPublish", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "when push messages failed",
			fields: fields{
				client:  new(mocks.Client),
				options: options,
			},
			args: args{
				ctx:      context.Background(),
				events:   []*v2.Event{&event},
				resource: "subject",
			},
			wantErr: true,
			mock: func(client *mocks.Client) {
				client.On("Publish", mock.Anything, mock.Anything).Return(errors.New("Error")).Times(5)
			},
		},
	}
	for _, tt := range tests {

		s.Run(tt.name, func() {

			tt.mock(tt.fields.client)

			p := NewClient(tt.fields.client, tt.fields.options)

			err := p.Publish(tt.args.ctx, tt.args.events)

			s.Assert().True((err != nil) == tt.wantErr, "Publish() error = %v, wantErr %v", err, tt.wantErr)

			tt.fields.client.AssertExpectations(s.T())

		})
	}
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
