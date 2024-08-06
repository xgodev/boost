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

	type args[T any] struct {
		ctx     context.Context
		conn    *n.Conn
		options *Options
		handler function.Handler[T]
	}
	tests := []struct {
		name string
		args args[string]
		want *Helper[string]
	}{
		{
			name: "success",
			args: args[string]{
				ctx:     ctx,
				conn:    conn,
				options: defaultOptions,
				handler: nil,
			},
			want: &Helper[string]{nil, "changeme", []string{"changeme"}, conn},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := NewHelperWithOptions[string](tt.args.conn, tt.args.handler, tt.args.options)
			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewHelperWithOptions() = %v, want %v")
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

	type args[T any] struct {
		ctx     context.Context
		conn    *n.Conn
		options *Options
		handler function.Handler[T]
	}
	tests := []struct {
		name string
		args args[string]
		want *Helper[string]
	}{
		{
			name: "success",
			args: args[string]{
				ctx:     ctx,
				conn:    conn,
				options: defaultOptions,
				handler: nil,
			},
			want: &Helper[string]{nil, "changeme", []string{"changeme"}, conn},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := NewHelper(tt.args.conn, tt.args.handler)
			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewHelperWithOptions() = %v, want %v")
		})
	}
}
