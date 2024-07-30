package nats

import (
	"context"
	"fmt"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
	"reflect"
	"testing"

	n "github.com/nats-io/nats.go"

	"github.com/stretchr/testify/suite"
)

type NatsHelperSuite struct {
	suite.Suite
}

func TestNatsHelperSuite(t *testing.T) {
	suite.Run(t, new(NatsHelperSuite))
}

func (s *NatsHelperSuite) SetupSuite() {
	boost.Start()
}

func (s *NatsHelperSuite) TestNatsNewHelper() {

	ctx := context.Background()
	defaultOptions, _ := DefaultOptions()

	sUrl := fmt.Sprintf("nats://127.0.0.1:%d", TestPort)
	options, _ := nats.NewOptions()
	options.Url = sUrl
	conn, _ := nats.NewConnWithOptions(ctx, options)

	type args struct {
		ctx     context.Context
		conn    *n.Conn
		options *Options
		handler function.Handler
	}
	tests := []struct {
		name string
		args args
		want *Helper
	}{
		{
			name: "success",
			args: args{
				ctx:     ctx,
				conn:    conn,
				options: defaultOptions,
				handler: nil,
			},
			want: &Helper{nil, "changeme", []string{"changeme"}, conn},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := NewHelper(tt.args.conn, tt.args.options, tt.args.handler)
			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewHelper() = %v, want %v")
		})
	}
}

func (s *NatsHelperSuite) TestNatsNewDefaultHelper() {

	ctx := context.Background()
	defaultOptions, _ := DefaultOptions()

	sUrl := fmt.Sprintf("nats://127.0.0.1:%d", TestPort)
	options, _ := nats.NewOptions()
	options.Url = sUrl
	conn, _ := nats.NewConnWithOptions(ctx, options)

	type args struct {
		ctx     context.Context
		conn    *n.Conn
		options *Options
		handler function.Handler
	}
	tests := []struct {
		name string
		args args
		want *Helper
	}{
		{
			name: "success",
			args: args{
				ctx:     ctx,
				conn:    conn,
				options: defaultOptions,
				handler: nil,
			},
			want: &Helper{nil, "changeme", []string{"changeme"}, conn},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := NewDefaultHelper(tt.args.conn, tt.args.handler)
			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewHelper() = %v, want %v")
		})
	}
}
