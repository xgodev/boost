package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
)

// Helper assists in creating event handlers.
type Helper struct {
	handler function.Handler
	topics  []string
	service common.Service
}

// NewHelper returns a new Helper with options.
func NewHelper(service common.Service, options *Options, handler function.Handler) *Helper {

	return &Helper{
		handler: handler,
		topics:  options.Topics,
		service: service,
	}
}

// NewDefaultHelper returns a new Helper with default options.
func NewDefaultHelper(service common.Service, handler function.Handler) *Helper {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelper(service, opt, handler)
}

func (h *Helper) Start() {

	// add some topic subscriptions
	sub := &common.Subscription{
		PubsubName: "messages",
		Topic:      "topic1",
		Route:      "/events",
	}
	if err := h.service.AddTopicEventHandler(sub, h.eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	if err := h.service.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error listenning: %v", err)
	}

}

func (h *Helper) eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	data, err := json.Marshal(e.Data)
	if err != nil {
		return false, fmt.Errorf("error parsing CloudEvent: %w", err)
	}

	var event cloudevents.Event

	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("error parsing CloudEvent: %v", err)
		return false, fmt.Errorf("error parsing CloudEvent: %w", err)
	}

	log.Printf("event - PubsubName: %s, Topic: %s, ID: %s, Data: %s", e.PubsubName, e.Topic, e.ID, e.Data)

	responseEvent, err := h.handler(ctx, event)
	if err != nil {
		return false, err
	}

	if responseEvent != nil {
		// Handle response event if needed
		log.Printf("response event - ID: %s, Data: %s", responseEvent.ID(), responseEvent.Data())
	}

	return false, nil
}
