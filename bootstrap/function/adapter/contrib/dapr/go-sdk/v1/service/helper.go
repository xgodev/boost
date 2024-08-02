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
type Helper struct {
	handler       function.Handler
	subscriptions []common.Subscription
	service       common.Service
}

// NewHelperWithOptions returns a new Helper with options.
func NewHelperWithOptions(service common.Service, handler function.Handler, options *Options) *Helper {

	return &Helper{
		handler:       handler,
		subscriptions: options.Subscriptions,
		service:       service,
	}
}

// NewHelper returns a new Helper with default options.
func NewHelper(service common.Service, handler function.Handler) *Helper {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelperWithOptions(service, handler, opt)
}

func (h *Helper) Start() {

	for _, sub := range h.subscriptions {
		if err := h.service.AddTopicEventHandler(&sub, h.eventHandler); err != nil {
			log.Fatalf("error adding topic subscription: %v", err)
		}

		log.Debugf("Added topic subscription: %v", sub)
	}

	if err := h.service.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error listenning: %v", err)
	}

}

func (h *Helper) eventHandler(ctx context.Context, topicEvent *common.TopicEvent) (retry bool, err error) {

	logger := log.FromContext(ctx)

	data, err := json.Marshal(topicEvent)
	if err != nil {
		return false, errors.Errorf("error parsing CloudEvent: %w", err)
	}

	in := event.New()
	err = json.Unmarshal(data, &in)
	if err != nil {
		return false, errors.Errorf("could set data: %w", err)
	}

	logger.Tracef("dapr - event - PubsubName: %s, Topic: %s, ID: %s, Data: %s", topicEvent.PubsubName, topicEvent.Topic, topicEvent.ID, topicEvent.Data)

	responseEvent, err := h.handler(ctx, in)
	if err != nil {
		return false, err
	}

	if responseEvent != nil {
		logger.Tracef("response event - ID: %s, Data: %s", responseEvent.ID(), responseEvent.Data())
	}

	return false, nil
}
