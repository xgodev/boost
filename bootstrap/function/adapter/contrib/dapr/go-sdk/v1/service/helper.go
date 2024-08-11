package service

import (
	"context"
	"encoding/json"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
)

// Helper assists in creating event handlers.
type Helper[T any] struct {
	handler       function.Handler[T]
	subscriptions []common.Subscription
	service       common.Service
}

// NewHelperWithOptions returns a new Helper with options.
func NewHelperWithOptions[T any](service common.Service, handler function.Handler[T], options *Options) *Helper[T] {

	return &Helper[T]{
		handler:       handler,
		subscriptions: options.Subscriptions,
		service:       service,
	}
}

// NewHelper returns a new Helper with default options.
func NewHelper[T any](service common.Service, handler function.Handler[T]) *Helper[T] {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelperWithOptions(service, handler, opt)
}

func (h *Helper[T]) Start() {

	for _, sub := range h.subscriptions {
		if err := h.service.AddTopicEventHandler(&sub, h.eventHandler); err != nil {
			log.Fatalf("error adding topic subscription: %v", err)
		}

		log.Debugf("Added topic subscription: %v", sub)
	}

	if err := h.service.AddServiceInvocationHandler("/events", h.serviceHandler); err != nil {
		log.Fatalf("error adding service invocation handler: %v", err)
	}

	if err := h.service.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error listenning: %v", err)
	}

}

func (h *Helper[T]) serviceHandler(ctx context.Context, inv *common.InvocationEvent) (out *common.Content, err error) {
	logger := log.FromContext(ctx)
	if inv == nil {
		err = errors.New("nil inv parameter")
		return
	}
	logger.Tracef(
		"echo - ContentType:%s, Verb:%s, QueryString:%s, %s",
		inv.ContentType, inv.Verb, inv.QueryString, inv.Data,
	)

	in := event.New()
	in.SetSubject("changeme")     // TODO: set subject
	in.SetSource("changeme")      // TODO: set source
	in.SetSpecVersion("changeme") // TODO: set spec version
	//for key, value := range topicEvent.Metadata {
	//	in.SetExtension(key, value)
	//}
	//in.SetType(topicEvent.Type)
	err = in.SetData("", inv.Data)
	if err != nil {
		return nil, errors.Wrap(err, errors.New("could set data"))
	}

	ev, err := h.handler(ctx, in)
	if err != nil {
		return nil, err
	}

	var data []byte

	switch x := any(ev).(type) {
	case []*event.Event, *event.Event:
		data, err = json.Marshal(x)
	default:
		return nil, errors.New("unsupported handler type")
	}

	out = &common.Content{
		Data:        data,
		ContentType: "application/json",
		//DataTypeURL: inv.DataTypeURL,
	}
	return
}

func (h *Helper[T]) eventHandler(ctx context.Context, topicEvent *common.TopicEvent) (retry bool, err error) {

	logger := log.FromContext(ctx)

	in := event.New()
	in.SetSubject(topicEvent.Subject)
	in.SetSource(topicEvent.Source)
	in.SetSpecVersion(topicEvent.SpecVersion)
	for key, value := range topicEvent.Metadata {
		in.SetExtension(key, value)
	}
	in.SetType(topicEvent.Type)
	err = in.SetData(topicEvent.DataContentType, topicEvent.Data)
	if err != nil {
		return false, errors.Wrap(err, errors.New("could set data"))
	}

	logger.Tracef("dapr - event - PubsubName: %s, Topic: %s, ID: %s, Data: %s", topicEvent.PubsubName, topicEvent.Topic, topicEvent.ID, topicEvent.Data)

	_, err = h.handler(ctx, in)
	if err != nil {
		return false, err
	}

	return false, nil
}
