// Package sqs is a publisher driver that sends each event to the SQS queue
// named by the event subject.
package sqs

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	v2 "github.com/cloudevents/sdk-go/v2"
	sqsclient "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/util"
)

// groupExtension is the event extension that becomes the MessageGroupId of a
// FIFO queue.
const groupExtension = "group"

type driver struct {
	client sqsclient.Client
}

// New returns a driver that publishes through client.
func New(client sqsclient.Client) publisher.Driver {
	return &driver{client: client}
}

// Publish sends the events concurrently, one SendMessage each, to the queue
// whose name is the event subject; a failure becomes the Error of that event's
// output. Retries are the AWS SDK retryer's.
func (d *driver) Publish(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {
	log.FromContext(ctx).WithTypeOf(driver{}).Tracef("publishing %d events to sqs", len(events))

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

	queueUrl, err := d.client.ResolveQueueUrl(ctx, event.Subject())
	if err != nil {
		return errors.Wrap(err, errors.Internalf("unable to resolve the url of sqs queue %s: %s", event.Subject(), err))
	}

	input := &awssqs.SendMessageInput{
		MessageBody: aws.String(string(payload)),
		QueueUrl:    queueUrl,
	}
	if group, ok := event.Extensions()[groupExtension]; ok {
		input.MessageGroupId = aws.String(fmt.Sprint(group))
	}

	log.FromContext(ctx).WithTypeOf(driver{}).
		WithField("subject", event.Subject()).
		WithField("id", event.ID()).
		Trace(string(payload))

	if err := d.client.Publish(ctx, input); err != nil {
		return errors.Wrap(err, errors.Internalf("unable to publish to sqs queue %s: %s", event.Subject(), err))
	}
	return nil
}
