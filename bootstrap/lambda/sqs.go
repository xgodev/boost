package lambda

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/model/errors"
)

func fromSQS(record Record) (*event.Event, error) {
	in := v2.NewEvent()
	body := []byte(record.Body)
	if err := json.Unmarshal(body, &in); err != nil {
		snsEvent := events.SNSEntity{}
		if err = json.Unmarshal(body, &snsEvent); err == nil && snsEvent.Message != "" {
			body = []byte(snsEvent.Message)
		}
		if err = in.SetData(v2.ApplicationJSON, body); err != nil {
			return nil, errors.NewNotValid(err, "could not set data in event")
		}
	}
	if in.ID() == "" {
		in.SetID(record.MessageId)
	}
	return &in, nil
}
