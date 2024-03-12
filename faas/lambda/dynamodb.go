package lambda

import (
	"github.com/aws/aws-lambda-go/events"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
)

func fromDynamoDB(record Record) (*event.Event, error) {
	userIdentity := new(events.DynamoDBUserIdentity)
	if err := transcode(record.UserIdentity, userIdentity); err != nil {
		return nil, err
	}

	in := v2.NewEvent()
	in.SetExtension("userIdentityPrincipalID", userIdentity.PrincipalID)
	in.SetExtension("userIdentityType", userIdentity.Type)

	err := in.SetData("", record.DynamoDB)
	return &in, err
}
