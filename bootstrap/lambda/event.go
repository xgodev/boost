package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"github.com/xgodev/boost/wrapper/log"
)

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
	UserIdentity           interface{}                           `json:"userIdentity"`
}

type Event struct {
	ID         string    `json:"id"`
	Source     string    `json:"source"`
	Region     string    `json:"region"`
	DetailType string    `json:"detail-type"`
	Time       time.Time `json:"time"`
	Account    string    `json:"account"`
	Resources  []string  `json:"resources"`
	Records    []Record  `json:"Records"`
}

func convertEvent(ctx context.Context, event Event, from func(record Record) (*event.Event, error)) []*cloudevents.InOut {
	logger := log.FromContext(ctx)

	lc, _ := lambdacontext.FromContext(ctx)

	mu := &sync.Mutex{}
	var inouts []*cloudevents.InOut

	var wg sync.WaitGroup

	for _, record := range event.Records {
		wg.Add(1)
		go func(record Record) {
			defer wg.Done()
			j, _ := json.Marshal(record) //this should be avoided
			logger.Debug(string(j))
			in, err := from(record)
			if in.ID() == "" {
				in.SetID(record.EventID)
			}
			if in.Type() == "" {
				in.SetType(record.EventName)
			}
			if in.Source() == "" {
				in.SetSource(record.EventSource)
			}
			in.SetExtension("awsRequestID", lc.AwsRequestID)
			in.SetExtension("invokedFunctionArn", lc.InvokedFunctionArn)
			mu.Lock()
			inouts = append(inouts, &cloudevents.InOut{
				In:  in,
				Err: err,
			})
			mu.Unlock()
		}(record)
	}
	wg.Wait()
	return inouts
}

func transcode(in, out interface{}) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(in); err != nil {
		return err
	}

	if err := json.NewDecoder(buf).Decode(out); err != nil {
		return err
	}

	return nil
}
