package pubsub

import (
	"fmt"
	"time"

	pb "cloud.google.com/go/pubsub/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/xgodev/boost/model/errors"
)

func generateCloudEvent(msg *pb.Message, subscription string) (event.Event, error) {
	in := event.New()

	ce := false
	contentType := "application/json"

	// Checks attributes and transforms into a CloudEvent if applicable
	for key, value := range msg.Attributes {
		switch key {
		case "content-type":
			in.SetDataContentType(value)
			contentType = value
		case "ce_specversion":
			in.SetSpecVersion(value)
			ce = true
		case "ce_id":
			in.SetID(value)
			ce = true
		case "ce_source":
			in.SetSource(value)
			ce = true
		case "ce_type":
			in.SetType(value)
			ce = true
		case "ce_time":
			ce = true
			if t, err := time.Parse(time.RFC3339, value); err == nil {
				in.SetTime(t)
			}
		case "ce_subject":
			ce = true
			in.SetSubject(value)
		default:
			in.SetExtension(key, value)
		}
	}

	// If the event does not have a time, populate it with the time the message was published
	if in.Time().IsZero() {
		in.SetTime(msg.PublishTime)
	}

	// If it's not a CloudEvent, create one manually
	if !ce {
		in.SetID(uuid.NewString())
		in.SetSource(fmt.Sprintf("pubsub://%s", subscription))
		in.SetType("pubsub.message")
	}

	// Set the message body as CloudEvent data
	if err := in.SetData(contentType, msg.Data); err != nil {
		return event.Event{}, errors.Wrap(err, errors.Internalf("could not set data from pubsub message: %s", err.Error()))
	}

	return in, nil
}
