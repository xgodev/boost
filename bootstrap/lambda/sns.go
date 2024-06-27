package lambda

import (
	"encoding/json"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/model/errors"
)

func fromSNS(record Record) (*event.Event, error) {

	var err error

	in := v2.NewEvent()
	message := []byte(record.SNS.Message)
	if err = json.Unmarshal(message, &in); err != nil {

		var data interface{}

		if err = json.Unmarshal(message, &data); err != nil {
			err = errors.NewNotValid(err, "could not decode SNS record")
		} else {
			if err = in.SetData(v2.ApplicationJSON, data); err != nil {
				err = errors.NewNotValid(err, "could not set data in event")
			}
		}

	}

	if in.ID() == "" {
		in.SetID(record.SNS.MessageID)
	}

	if in.Type() == "" {
		in.SetType(record.SNS.Type)
	}

	return &in, err
}
