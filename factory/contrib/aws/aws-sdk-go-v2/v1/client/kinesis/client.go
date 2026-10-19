package kinesis

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Client knows how to bulk publish on kinesis.
type Client interface {
	BulkPublish(ctx context.Context, messages []types.PutRecordsRequestEntry, resource string) error
	Publish(ctx context.Context, input *kinesis.PutRecordInput) error
}

// maxBulkAttempts bounds the PutRecords calls for one lot: the records
// Kinesis rejects are sent again until they get in or the attempts run out.
const maxBulkAttempts = 5

// api is the part of the Kinesis client the boost client calls.
type api interface {
	PutRecord(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
	PutRecords(ctx context.Context, params *kinesis.PutRecordsInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordsOutput, error)
}

// Client holds client and resource name.
type client struct {
	client api
}

// NewClient returns an initialized client.
func NewClient(c *kinesis.Client) Client {
	return &client{c}
}

// Publish publishes message on kinesis.
func (c *client) Publish(ctx context.Context, input *kinesis.PutRecordInput) error {

	logger := log.FromContext(ctx).
		WithTypeOf(*c).
		WithField("resource", input.StreamName).
		WithField("partition_key", input.PartitionKey)

	logger.Trace("sending message to kinesis")

	response, err := c.client.PutRecord(ctx, input)
	if err != nil {
		return errors.Wrap(err, errors.New("error publishing message on kinesis"))
	}

	logger.
		WithField("sequence_number", *response.SequenceNumber).
		WithField("shard_id", *response.ShardId).
		Debug("message sent to kinesis")

	return nil
}

// BulkPublish publishes an array of messages on kinesis, in lots of 500.
// It fails when a lot cannot be sent or some of its records are still
// rejected after maxBulkAttempts.
func (c *client) BulkPublish(ctx context.Context, entries []types.PutRecordsRequestEntry, resource string) error {

	logger := log.FromContext(ctx).
		WithTypeOf(*c).
		WithField("resource", resource)

	for _, lot := range c.splitInputs(entries, 500) {
		if err := c.putLot(ctx, logger, lot, resource); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) putLot(ctx context.Context, logger log.Logger, lot []types.PutRecordsRequestEntry, resource string) error {
	total := len(lot)

	for attempt := 1; ; attempt++ {

		logger.Trace("sending bulk message to kinesis")

		response, err := c.client.PutRecords(ctx, c.buildPutRecordsInput(lot, resource))
		if err != nil {
			return errors.Wrap(err, errors.New("error publishing message on kinesis"))
		}

		var retry []types.PutRecordsRequestEntry

		for i, r := range response.Records {

			if r.ErrorMessage != nil {
				logger.
					WithField("cause", aws.ToString(r.ErrorMessage)).
					WithField("code", aws.ToString(r.ErrorCode)).
					Warn("error in kinesis bulk record")
				retry = append(retry, lot[i])
				continue
			}

			logger.
				WithField("sequence_number", aws.ToString(r.SequenceNumber)).
				WithField("shard_id", aws.ToString(r.ShardId)).
				Debug("message sent to kinesis")

		}

		if len(retry) == 0 {
			return nil
		}

		if attempt == maxBulkAttempts {
			return errors.Internalf("%d of %d kinesis records rejected after %d attempts", len(retry), total, maxBulkAttempts)
		}

		logger.Warnf("retrying %d of %d kinesis records", len(retry), total)
		lot = retry
	}
}

func (c *client) buildPutRecordsInput(messages []types.PutRecordsRequestEntry,
	resource string) *kinesis.PutRecordsInput {
	return &kinesis.PutRecordsInput{
		Records:    messages,
		StreamName: aws.String(resource),
	}
}

func (c *client) splitInputs(inputs []types.PutRecordsRequestEntry, chunkSize int) (chunks [][]types.PutRecordsRequestEntry) {
	for chunkSize < len(inputs) {
		inputs, chunks = inputs[chunkSize:], append(chunks, inputs[0:chunkSize:chunkSize])
	}

	return append(chunks, inputs)
}
