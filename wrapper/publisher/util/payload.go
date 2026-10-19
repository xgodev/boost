package util

import (
	"encoding/json"
	"fmt"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/model/errors"
)

// targetData is the value of the "target" extension that makes a driver send
// only the event data.
const targetData = "data"

// Payload returns what a driver sends for event: the whole event as JSON or,
// when its "target" extension is "data", only its data as JSON.
func Payload(event *v2.Event) ([]byte, error) {
	if target, ok := event.Extensions()["target"]; !ok || fmt.Sprint(target) != targetData {
		return JSONBytes(*event)
	}

	var data interface{}
	if err := event.DataAs(&data); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("unable to decode the data of event %s: %s", event.ID(), err))
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("unable to encode the data of event %s: %s", event.ID(), err))
	}
	return payload, nil
}
