// Package sns is a publisher driver that sends each event to the SNS topic
// whose ARN is the event subject.
package sns

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssns "github.com/aws/aws-sdk-go-v2/service/sns"
	v2 "github.com/cloudevents/sdk-go/v2"
	snsclient "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sns"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/util"
)

type driver struct {
	client snsclient.Client
}

// message is the SNS message with MessageStructure "json": every protocol
// receives Default.
type message struct {
	Default string `json:"default"`
}

// New returns a driver that publishes through client.
func New(client snsclient.Client) publisher.Driver {
	return &driver{client: client}
}

// Publish sends the events concurrently, one SNS Publish each; a failure
// becomes the Error of that event's output. Retries are the AWS SDK retryer's.
func (d *driver) Publish(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {
	log.FromContext(ctx).WithTypeOf(driver{}).Tracef("publishing %d events to sns", len(events))

	res := make([]publisher.PublishOutput, len(events))
	var wg sync.WaitGroup
	for i, event := range events {
		res[i].Event = event
		wg.Add(1)
		go func() {
			defer wg.Done()
			res[i].Error = d.publish(ctx, event)
		}()
	}
	wg.Wait()

	return res, nil
}

func (d *driver) publish(ctx context.Context, event *v2.Event) error {
	payload, err := util.Payload(event)
	if err != nil {
		return err
	}

	body, err := json.Marshal(message{Default: string(payload)})
	if err != nil {
		return errors.Wrap(err, errors.Internalf("unable to encode the sns message of event %s", event.ID()))
	}

	log.FromContext(ctx).WithTypeOf(driver{}).
		WithField("subject", event.Subject()).
		WithField("id", event.ID()).
		Trace(string(payload))

	err = d.client.Publish(ctx, &awssns.PublishInput{
		Message:          aws.String(string(body)),
		MessageStructure: aws.String("json"),
		TopicArn:         aws.String(event.Subject()),
	})
	if err != nil {
		return errors.Wrap(err, errors.Internalf("unable to publish to sns topic %s: %s", event.Subject(), err))
	}
	return nil
}
