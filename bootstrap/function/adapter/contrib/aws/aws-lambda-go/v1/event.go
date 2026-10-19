package lambda

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/model/errors"
)

// Event is the payload of one invocation: the records of an SQS, SNS,
// Kinesis, S3 or DynamoDB trigger, or an EventBridge (CloudWatch Events)
// event, which has no records.
type Event struct {
	ID         string          `json:"id"`
	Source     string          `json:"source"`
	Region     string          `json:"region"`
	DetailType string          `json:"detail-type"`
	Detail     json.RawMessage `json:"detail"`
	Time       time.Time       `json:"time"`
	Account    string          `json:"account"`
	Resources  []string        `json:"resources"`
	Records    []Record        `json:"Records"`
}

// Record is one record of a trigger; only the fields of its EventSource are set.
type Record struct {
	EventVersion           string                                `json:"eventVersion"`
	EventSubscriptionArn   string                                `json:"eventSubscriptionArn"`
	EventSource            string                                `json:"eventSource"`
	EventName              string                                `json:"eventName"`
	EventID                string                                `json:"eventID"`
	SNS                    events.SNSEntity                      `json:"sns"`
	S3                     events.S3Entity                       `json:"s3"`
	Kinesis                events.KinesisRecord                  `json:"kinesis"`
	DynamoDB               events.DynamoDBStreamRecord           `json:"dynamodb"`
	MessageId              string                                `json:"messageId"`
	ReceiptHandle          string                                `json:"receiptHandle"`
	Body                   string                                `json:"body"`
	Md5OfBody              string                                `json:"md5OfBody"`
	Md5OfMessageAttributes string                                `json:"md5OfMessageAttributes"`
	Attributes             map[string]string                     `json:"attributes"`
	MessageAttributes      map[string]events.SQSMessageAttribute `json:"messageAttributes"`
	EventSourceARN         string                                `json:"eventSourceARN"`
	AWSRegion              string                                `json:"awsRegion"`
	UserIdentity           *events.DynamoDBUserIdentity          `json:"userIdentity"`
}

// toCloudEvent converts a record by its EventSource, filling id, type and
// source from the record when its payload did not carry them.
func toCloudEvent(r Record) (event.Event, error) {
	var in event.Event
	var err error
	switch r.EventSource {
	case "aws:sqs":
		in, err = fromSQS(r)
	case "aws:sns":
		in, err = fromSNS(r)
	case "aws:kinesis":
		in, err = fromKinesis(r)
	case "aws:s3":
		in, err = fromS3(r)
	case "aws:dynamodb":
		in, err = fromDynamoDB(r)
	default:
		return event.Event{}, errors.NotImplementedf("lambda trigger %q", r.EventSource)
	}
	if err != nil {
		return in, err
	}
	if in.ID() == "" {
		in.SetID(r.EventID)
	}
	if in.Type() == "" {
		in.SetType(r.EventName)
	}
	if in.Type() == "" {
		in.SetType(r.EventSource)
	}
	if in.Source() == "" {
		in.SetSource(r.EventSource)
	}
	return in, nil
}

// fromEventBridge converts an EventBridge (CloudWatch Events) event, such as
// a schedule, whose detail becomes the data.
func fromEventBridge(e Event) (event.Event, error) {
	in := event.New()
	in.SetID(e.ID)
	in.SetType(e.DetailType)
	in.SetSource(e.Source)
	if !e.Time.IsZero() {
		in.SetTime(e.Time)
	}
	if len(e.Detail) > 0 {
		if err := in.SetData(event.ApplicationJSON, []byte(e.Detail)); err != nil {
			return in, errors.NewNotValid(err, "could not set the eventbridge detail as data")
		}
	}
	return in, nil
}

func fromSQS(r Record) (event.Event, error) {
	body := []byte(r.Body)
	if in, ok := cloudEvent(body); ok {
		return in, nil
	}

	// SNS delivers to SQS inside an envelope; its message is the payload.
	var envelope events.SNSEntity
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Message != "" {
		body = []byte(envelope.Message)
		if in, ok := cloudEvent(body); ok {
			return in, nil
		}
	}

	in, err := withData(body)
	in.SetID(r.MessageId)
	return in, err
}

func fromSNS(r Record) (event.Event, error) {
	message := []byte(r.SNS.Message)
	if in, ok := cloudEvent(message); ok {
		return in, nil
	}
	in, err := withData(message)
	in.SetID(r.SNS.MessageID)
	in.SetType(r.SNS.Type)
	return in, err
}

func fromKinesis(r Record) (event.Event, error) {
	if in, ok := cloudEvent(r.Kinesis.Data); ok {
		return in, nil
	}
	return withData(r.Kinesis.Data)
}

func fromS3(r Record) (event.Event, error) {
	in := event.New()
	in.SetID(r.S3.Object.Key)
	if err := in.SetData(event.ApplicationJSON, r.S3); err != nil {
		return in, errors.NewNotValid(err, "could not set the s3 entity as data")
	}
	return in, nil
}

func fromDynamoDB(r Record) (event.Event, error) {
	in := event.New()
	if r.UserIdentity != nil {
		in.SetExtension("useridentityprincipalid", r.UserIdentity.PrincipalID)
		in.SetExtension("useridentitytype", r.UserIdentity.Type)
	}
	if err := in.SetData(event.ApplicationJSON, r.DynamoDB); err != nil {
		return in, errors.NewNotValid(err, "could not set the dynamodb record as data")
	}
	return in, nil
}

// cloudEvent decodes payload as a CloudEvent in structured mode.
func cloudEvent(payload []byte) (event.Event, bool) {
	in := event.New()
	if err := json.Unmarshal(payload, &in); err != nil {
		return event.Event{}, false
	}
	return in, true
}

// withData returns a new event carrying payload as its data: JSON when it is
// JSON, plain text otherwise.
func withData(payload []byte) (event.Event, error) {
	in := event.New()
	contentType := event.TextPlain
	if json.Valid(payload) {
		contentType = event.ApplicationJSON
	}
	if err := in.SetData(contentType, payload); err != nil {
		return in, errors.NewNotValid(err, "could not set the record payload as data")
	}
	return in, nil
}
