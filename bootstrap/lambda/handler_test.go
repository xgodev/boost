package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xgodev/boost/bootstrap/cloudevents/plugins/local/wrapper/log"
	"io/ioutil"
	"reflect"
	"testing"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"github.com/xgodev/boost/config"
	igcloudevents "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
	iglog "github.com/xgodev/boost/factory/local/wrapper/log"
)

type HandlerSuite struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) SetupSuite() {
	config.Load()
	iglog.New()
}

func (s *HandlerSuite) TestHandler_Handle() {

	lc := new(lambdacontext.LambdaContext)
	ctx := lambdacontext.NewContext(context.Background(), lc)

	var middlewares []cloudevents.Middleware

	middlewares = append(middlewares, log.NewLogger())

	options, _ := DefaultOptions()

	type fields struct {
		handler     func(*assert.Assertions) igcloudevents.Handler
		middlewares []cloudevents.Middleware
		options     *Options
	}

	tests := []struct {
		name          string
		eventFilename string
		fields        fields
		wantErr       func(error) bool
		want          v2.Event
	}{
		{
			name: "on kinesis success event",
			fields: fields{
				handler: func(a *assert.Assertions) igcloudevents.Handler {
					return func(ctx context.Context, in v2.Event) (*v2.Event, error) {
						a.Equal("6aee8d49-76fd-4900-a49b-eecac3552c5a", in.ID())
						a.Equal("receiver/internal/app/receiver/domain/model/item", in.Type())
						a.Equal("receiver-receiver", in.Source())
						a.Equal(lc.AwsRequestID, in.Extensions()["awsrequestid"])
						a.Equal(lc.InvokedFunctionArn, in.Extensions()["invokedfunctionarn"])
						var evt map[string]interface{}
						in.DataAs(&evt)
						a.Equal("123", evt["id"])
						a.Equal("skyhub", evt["source"])
						a.Equal("Iphone 11 XS", evt["name"])
						return &in, nil
					}
				},
				middlewares: middlewares,
				options:     options,
			},
			eventFilename: "kinesis_success.json",
			wantErr:       func(err error) bool { return err == nil },
		},
		{
			name: "on s3 success event",
			fields: fields{
				handler: func(a *assert.Assertions) igcloudevents.Handler {
					return func(ctx context.Context, in v2.Event) (*v2.Event, error) {
						a.Equal("b21b84d653bb07b05b1e6b33684dc11b", in.ID())
						a.Equal("ObjectCreated:Put", in.Type())
						a.Equal("aws:s3", in.Source())
						a.Equal(lc.AwsRequestID, in.Extensions()["awsrequestid"])
						a.Equal(lc.InvokedFunctionArn, in.Extensions()["invokedfunctionarn"])
						var evt awsevents.S3Entity
						in.DataAs(&evt)
						a.Equal("1.0", evt.SchemaVersion)
						a.Equal("828aa6fc-f7b5-4305-8584-487c791949c1", evt.ConfigurationID)
						a.Equal("lambda-artifacts-deafc19498e3f2df", evt.Bucket.Name)
						a.Equal("A3I5XTEXAMAI3E", evt.Bucket.OwnerIdentity.PrincipalID)
						a.Equal("arn:aws:s3:::lambda-artifacts-deafc19498e3f2df", evt.Bucket.Arn)
						a.Equal("b21b84d653bb07b05b1e6b33684dc11b", evt.Object.Key)
						a.Equal(int64(1305107), evt.Object.Size)
						a.Equal("b21b84d653bb07b05b1e6b33684dc11b", evt.Object.ETag)
						a.Equal("0C0F6F405D6ED209E1", evt.Object.Sequencer)
						return &in, nil
					}
				},
				middlewares: middlewares,
				options:     options,
			},
			eventFilename: "s3_success.json",
			wantErr:       func(err error) bool { return err == nil },
		},
		{
			name: "on dynamodb success event",
			fields: fields{
				handler: func(a *assert.Assertions) igcloudevents.Handler {
					return func(ctx context.Context, in v2.Event) (*v2.Event, error) {
						a.Equal("7de3041dd709b024af6f29e4fa13d34c", in.ID())
						a.Equal("INSERT", in.Type())
						a.Equal("aws:dynamodb", in.Source())
						a.Equal(lc.AwsRequestID, in.Extensions()["awsrequestid"])
						a.Equal(lc.InvokedFunctionArn, in.Extensions()["invokedfunctionarn"])
						a.Equal("Service", in.Extensions()["useridentitytype"])
						a.Equal("dynamodb.amazonaws.com", in.Extensions()["useridentityprincipalid"])

						var evt awsevents.DynamoDBStreamRecord
						in.DataAs(&evt)

						a.Equal(int64(1479499740), evt.ApproximateCreationDateTime.Time.Unix())

						a.Equal("2016-11-18:12:09:36", evt.Keys["Timestamp"].String())
						a.Equal("John Doe", evt.Keys["Username"].String())

						a.Equal("2016-11-18:12:09:36", evt.NewImage["Timestamp"].String())
						a.Equal("This is a bark from the Woofer social network", evt.NewImage["Message"].String())
						a.Equal("John Doe", evt.NewImage["Username"].String())

						a.Equal("13021600000000001596893679", evt.SequenceNumber)
						a.Equal(int64(112), evt.SizeBytes)
						a.Equal("NEW_IMAGE", evt.StreamViewType)
						return &in, nil
					}
				},
				middlewares: middlewares,
				options:     options,
			},
			eventFilename: "dynamodb_success.json",
			wantErr:       func(err error) bool { return err == nil },
		},
		{
			name: "on sns success event",
			fields: fields{
				handler: func(a *assert.Assertions) igcloudevents.Handler {
					return func(ctx context.Context, in v2.Event) (*v2.Event, error) {
						a.Equal("95df01b4-ee98-5cb9-9903-4c221d41eb5e", in.ID())
						a.Equal("Notification", in.Type())
						a.Equal("aws:sns", in.Source())
						a.Equal(lc.AwsRequestID, in.Extensions()["awsrequestid"])
						a.Equal(lc.InvokedFunctionArn, in.Extensions()["invokedfunctionarn"])
						var evt map[string]interface{}
						in.DataAs(&evt)
						a.Equal("123", evt["test"])
						return &in, nil
					}
				},
				middlewares: middlewares,
				options:     options,
			},
			eventFilename: "sns_success.json",
			wantErr:       func(err error) bool { return err == nil },
		},
		{
			name: "on sqs plain message success event",
			fields: fields{
				handler: func(a *assert.Assertions) igcloudevents.Handler {
					return func(ctx context.Context, in v2.Event) (*v2.Event, error) {
						a.Equal("aws:sqs", in.Source())
						var evt map[string]interface{}
						in.DataAs(&evt)
						a.Equal("abc", evt["id"])
						return &in, nil
					}
				},
				middlewares: middlewares,
				options:     options,
			},
			eventFilename: "sqs_plain_success.json",
			wantErr:       func(err error) bool { return err == nil },
		},
		{
			name: "on sqs success event",
			fields: fields{
				handler: func(a *assert.Assertions) igcloudevents.Handler {
					return func(ctx context.Context, in v2.Event) (*v2.Event, error) {
						a.Equal("123456", in.ID())
						a.Equal("myType", in.Type())
						a.Equal("mySource", in.Source())
						a.Equal(lc.AwsRequestID, in.Extensions()["awsrequestid"])
						a.Equal(lc.InvokedFunctionArn, in.Extensions()["invokedfunctionarn"])
						var evt map[string]interface{}
						in.DataAs(&evt)
						a.Equal("abc", evt["id"])
						return &in, nil
					}
				},
				middlewares: middlewares,
				options:     options,
			},
			eventFilename: "sqs_success.json",
			wantErr:       func(err error) bool { return err == nil },
		},
		{
			name: "on sqs wrapping sns success event",
			fields: fields{
				handler: func(a *assert.Assertions) igcloudevents.Handler {
					return func(ctx context.Context, in v2.Event) (*v2.Event, error) {
						a.Equal("2e1424d4-f796-459a-8184-9c92662be6da", in.ID())
						var evt map[string]interface{}
						in.DataAs(&evt)
						a.Equal("123", evt["test"])
						a.Equal("id-123", evt["id"])
						return &in, nil
					}
				},
				middlewares: middlewares,
				options:     options,
			},
			eventFilename: "sqs_wraps_sns_success.json",
			wantErr:       func(err error) bool { return err == nil },
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			hwOptions, _ := cloudevents.DefaultHandlerWrapperOptions()
			hw := cloudevents.NewHandlerWrapper(tt.fields.handler(s.Assertions), hwOptions, tt.fields.middlewares...)
			h := NewHandler(hw, tt.fields.options)
			var evt Event
			b, _ := ioutil.ReadFile(fmt.Sprintf("testdata/%s", tt.eventFilename))
			//fmt.Println(string(b))
			json.Unmarshal(b, &evt)
			//str, _ := json.MarshalIndent(evt, "", "  ")
			//fmt.Println(string(str))
			err := h.Handle(ctx, evt)
			s.Assert().True(tt.wantErr(err), "Handle() unexpected error = %v", err)
		})
	}
}

func (s *HandlerSuite) TestNewHandler() {

	type args struct {
		handler igcloudevents.Handler
		options *Options
	}

	handler := func(ctx context.Context, in v2.Event) (*v2.Event, error) { return nil, nil }
	options, _ := DefaultOptions()
	hwOptions, _ := cloudevents.DefaultHandlerWrapperOptions()
	hw := cloudevents.NewHandlerWrapper(handler, hwOptions)

	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "success",
			args: args{
				handler: handler,
				options: options,
			},
			want: NewHandler(hw, options),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			got := NewHandler(hw, tt.args.options)

			s.Assert().True(reflect.DeepEqual(got, tt.want), "NewHandler() = %v, want %v", got, tt.want)

		})
	}
}
