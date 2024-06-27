package kinesis

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/bootstrap/repository"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis/mocks"
	iglog "github.com/xgodev/boost/factory/local/wrapper/log"
)

type EventSuite struct {
	suite.Suite
}

func (s *EventSuite) SetupSuite() {
	config.Load()
	iglog.New()
}

func (s *EventSuite) TestNewEvent() {

	client := new(mocks.Client)
	options, _ := DefaultOptions()

	type args struct {
		client  *mocks.Client
		options *Options
	}
	tests := []struct {
		name string
		args args
		want repository.Event
	}{
		{
			name: "Success",
			args: args{
				client:  client,
				options: options,
			},
			want: NewEvent(client, options),
		},
	}
	for _, tt := range tests {

		s.Run(tt.name, func() {
			got := NewEvent(tt.args.client, tt.args.options)
			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewEvent() = %v, want %v", got, tt.want)
		})
	}
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}
