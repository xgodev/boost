package lambda

import (
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
)

func fromS3(record Record) (*event.Event, error) {

	in := v2.NewEvent()
	in.SetID(record.S3.Object.Key)
	err := in.SetData("", record.S3)

	return &in, err
}
