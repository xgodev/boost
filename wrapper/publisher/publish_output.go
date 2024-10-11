package publisher

import cloudevents "github.com/cloudevents/sdk-go/v2"

type PublishOutput struct {
	Event *cloudevents.Event
	Error error
}
