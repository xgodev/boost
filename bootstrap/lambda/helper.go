package lambda

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xgodev/boost/bootstrap/cloudevents"
)

// Helper assists in creating event handlers.
type Helper struct {
	handler *Handler
}

// NewHelper returns a new Helper with options.
func NewHelper(handler *cloudevents.HandlerWrapper, options *Options) *Helper {

	h := NewHandler(handler, options)

	return &Helper{
		handler: h,
	}
}

// NewHelper returns a new Helper with default options.
func NewDefaultHelper(handler *cloudevents.HandlerWrapper) *Helper {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelper(handler, opt)
}

// Start starts HTTP client for handle events.
func (h *Helper) Start() {
	lambda.Start(h.handler.Handle)
}
