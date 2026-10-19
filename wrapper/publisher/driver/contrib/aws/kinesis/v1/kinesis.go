// Package kinesis is a publisher driver that puts each event in the Kinesis
// stream named by the event subject.
package kinesis

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awskinesis "github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	kinesisclient "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/util"
)

// partitionKeyExtension is the event extension that becomes the record's
// partition key.
const partitionKeyExtension = "partitionkey"

type driver struct {
	client kinesisclient.Client
}

// stream is what Publish puts in one stream: the records and the output index
// of the event of each one.
type stream struct {
	name    string
	records []types.PutRecordsRequestEntry
	outputs []int
}

// New returns a driver that publishes through client.
func New(client kinesisclient.Client) publisher.Driver {
	return &driver{client: client}
}

// Publish puts the events of each stream (the event subject) in one call, the
// streams concurrently: a PutRecord for a single event, a PutRecords (bulk)
// otherwise. A failure of a stream becomes the Error of every event of that
// stream, so a retry may duplicate the records the stream did accept.
func (d *driver) Publish(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {
	logger := log.FromContext(ctx).WithTypeOf(driver{})
	logger.Tracef("publishing %d events to kinesis", len(events))

	res := make([]publisher.PublishOutput, len(events))
	streams := map[string]*stream{}
	for i, event := range events {
		res[i].Event = event

		payload, err := util.Payload(event)
		if err != nil {
			res[i].Error = err
			continue
		}

		s, ok := streams[event.Subject()]
		if !ok {
			s = &stream{name: event.Subject()}
			streams[s.name] = s
		}
		pk := partitionKey(event)
		s.records = append(s.records, types.PutRecordsRequestEntry{Data: payload, PartitionKey: aws.String(pk)})
		s.outputs = append(s.outputs, i)

		logger.WithField("subject", event.Subject()).
			WithField("id", event.ID()).
			WithField("partitionKey", pk).
			Trace(string(payload))
	}

	var wg sync.WaitGroup
	for _, s := range streams {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := d.put(ctx, s); err != nil {
				for _, i := range s.outputs {
					res[i].Error = err
				}
			}
		}()
	}
	wg.Wait()

	return res, nil
}

func (d *driver) put(ctx context.Context, s *stream) error {
	var err error
	if len(s.records) == 1 {
		err = d.client.Publish(ctx, &awskinesis.PutRecordInput{
			Data:         s.records[0].Data,
			PartitionKey: s.records[0].PartitionKey,
			StreamName:   aws.String(s.name),
		})
	} else {
		err = d.client.BulkPublish(ctx, s.records, s.name)
	}
	if err != nil {
		return errors.Wrap(err, errors.Internalf("unable to publish to kinesis stream %s: %s", s.name, err))
	}
	return nil
}

// partitionKey returns the "partitionkey" extension of event or, without it,
// the event ID, so the records of different events spread over the shards.
func partitionKey(event *v2.Event) string {
	if pk, ok := event.Extensions()[partitionKeyExtension]; ok {
		if s := fmt.Sprint(pk); s != "" {
			return s
		}
	}
	if event.ID() != "" {
		return event.ID()
	}
	return uuid.NewString()
}
