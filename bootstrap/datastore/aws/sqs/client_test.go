package sqs

import (
	"context"
	iglogrus "github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus/v1"
	"reflect"
	"testing"

	"github.com/xgodev/boost/model/errors"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs"
	sqsmock "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs/mocks"
)

type SqsClientSuite struct {
	suite.Suite
}

func TestSqsClientSuite(t *testing.T) {
	suite.Run(t, new(SqsClientSuite))
}

func (s *SqsClientSuite) SetupSuite() {
	iglogrus.NewLogger()
}

func (s *SqsClientSuite) TestNewClient() {

	var sqsCli sqs.Client

	tt := []struct {
		name   string
		client sqs.Client
		want   *Client
	}{
		{
			name:   "success",
			client: sqsCli,
			want:   &Client{sqsCli},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewClient(t.client)
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *SqsClientSuite) Test_rawMessage() {
	tt := []struct {
		name    string
		in      func() *v2.Event
		cli     *Client
		want    string
		wantErr func(error) bool
	}{
		{
			name: "when extensions target is not present",
			in: func() *v2.Event {
				e := v2.NewEvent()
				e.SetSubject("blah")
				e.SetID("123")
				e.SetSource("/home/blah")
				e.SetType("order")
				return &e
			},
			cli:  &Client{},
			want: `{"specversion":"1.0","id":"123","source":"/home/blah","type":"order","subject":"blah"}`,
			wantErr: func(e error) bool {
				return e == nil
			},
		},
		{
			name: "when extensions target is present and is data",
			in: func() *v2.Event {
				e := v2.NewEvent()
				e.SetSubject("blah")
				e.SetID("123")
				e.SetSource("/home/blah")
				e.SetType("order")
				e.SetExtension("target", "data")
				e.SetData("", map[string]string{
					"id":   "123",
					"name": "xablau",
				})
				return &e
			},
			cli:  &Client{},
			want: `{"id":"123","name":"xablau"}`,
			wantErr: func(e error) bool {
				return e == nil
			},
		},
		{
			name: "when extensions target is present and is not data",
			in: func() *v2.Event {
				e := v2.NewEvent()
				e.SetSubject("blah")
				e.SetID("123")
				e.SetSource("/home/blah")
				e.SetType("order")
				e.SetExtension("target", "all")
				e.SetData("", map[string]string{
					"id":   "123",
					"name": "xablau",
				})
				return &e
			},
			cli:  &Client{},
			want: `{"specversion":"1.0","id":"123","source":"/home/blah","type":"order","subject":"blah","data":{"id":"123","name":"xablau"},"target":"all"}`,
			wantErr: func(e error) bool {
				return e == nil
			},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got, err := t.cli.rawMessage(t.in())
			s.Assert().True(string(got) == t.want, "got  %v\nwant %v", string(got), t.want)
			s.Assert().True(t.wantErr(err), "unexpected error %v", err)
		})
	}
}

func (s *SqsClientSuite) Test_send() {
	tt := []struct {
		name    string
		in      func() []*v2.Event
		mock    func(*sqsmock.Client)
		wantErr func(error) bool
	}{
		{
			name: "when ResolveQueueUrl pops an error",
			in: func() []*v2.Event {
				e := v2.NewEvent()
				e.SetSubject("blah")
				e.SetID("123")
				e.SetSource("/home/blah")
				e.SetType("order")
				return []*v2.Event{&e}
			},
			mock: func(m *sqsmock.Client) {
				m.On("ResolveQueueUrl", mock.Anything, mock.Anything).Once().Return(nil, errors.Internalf("Ops!"))
			},
			wantErr: errors.IsInternal,
		},
		{
			name: "when Publish pops an error",
			in: func() []*v2.Event {
				e := v2.NewEvent()
				e.SetSubject("blah")
				e.SetID("123")
				e.SetSource("/home/blah")
				e.SetType("order")
				return []*v2.Event{&e}
			},
			mock: func(m *sqsmock.Client) {
				r := "blah"
				m.On("ResolveQueueUrl", mock.Anything, mock.Anything).Once().Return(&r, nil)
				m.On("Publish", mock.Anything, mock.Anything).Times(5).Return(errors.BadRequestf("Wth.."))
			},
			wantErr: errors.IsInternal,
		},
		{
			name: "when success",
			in: func() []*v2.Event {
				e := v2.NewEvent()
				e.SetSubject("blah")
				e.SetID("123")
				e.SetSource("/home/blah")
				e.SetType("order")
				e.SetExtension("group", "g2")
				return []*v2.Event{&e}
			},
			mock: func(m *sqsmock.Client) {
				r := "blah"
				m.On("ResolveQueueUrl", mock.Anything, mock.Anything).Once().Return(&r, nil)
				m.On("Publish", mock.Anything, mock.Anything).Once().Return(nil)
			},
			wantErr: func(e error) bool { return e == nil },
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			m := new(sqsmock.Client)
			t.mock(m)
			cli := NewClient(m)
			err := cli.send(context.TODO(), t.in())
			s.Assert().True(t.wantErr(err), "unexpected error %v", err)
			m.AssertExpectations(s.T())
		})
	}
}

func (s *SqsClientSuite) TestPublish() {
	tt := []struct {
		name    string
		in      func() []*v2.Event
		mock    func(*sqsmock.Client)
		wantErr func(error) bool
	}{
		{
			name: "when theres are no events to send",
			in: func() []*v2.Event {
				return []*v2.Event{}
			},
			mock: func(m *sqsmock.Client) {
			},
			wantErr: func(e error) bool { return e == nil },
		},
		{
			name: "when success",
			in: func() []*v2.Event {
				e := v2.NewEvent()
				e.SetSubject("blah")
				e.SetID("123")
				e.SetSource("/home/blah")
				e.SetType("order")
				e.SetExtension("group", "g2")
				return []*v2.Event{&e}
			},
			mock: func(m *sqsmock.Client) {
				r := "blah"
				m.On("ResolveQueueUrl", mock.Anything, mock.Anything).Once().Return(&r, nil)
				m.On("Publish", mock.Anything, mock.Anything).Once().Return(nil)
			},
			wantErr: func(e error) bool { return e == nil },
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			m := new(sqsmock.Client)
			t.mock(m)
			cli := NewClient(m)
			err := cli.Publish(context.TODO(), t.in())
			s.Assert().True(t.wantErr(err), "unexpected error %v", err)
			m.AssertExpectations(s.T())
		})
	}
}
