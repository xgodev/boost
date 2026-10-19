package kinesis

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/xgodev/boost/model/errors"
)

// fakeKinesis fails every record of PutRecords; past 100 calls it gives up,
// so a client that retries forever still ends.
type fakeKinesis struct {
	calls atomic.Int32
}

func (f *fakeKinesis) PutRecord(context.Context, *kinesis.PutRecordInput, ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error) {
	return &kinesis.PutRecordOutput{SequenceNumber: aws.String("1"), ShardId: aws.String("shard-1")}, nil
}

func (f *fakeKinesis) PutRecords(_ context.Context, in *kinesis.PutRecordsInput, _ ...func(*kinesis.Options)) (*kinesis.PutRecordsOutput, error) {
	if f.calls.Add(1) > 100 {
		return nil, errors.New("too many calls")
	}
	records := make([]types.PutRecordsResultEntry, len(in.Records))
	for i := range records {
		records[i] = types.PutRecordsResultEntry{
			ErrorCode:    aws.String("ProvisionedThroughputExceededException"),
			ErrorMessage: aws.String("slow down"),
		}
	}
	return &kinesis.PutRecordsOutput{FailedRecordCount: aws.Int32(int32(len(records))), Records: records}, nil
}

func TestBulkPublish_GivesUpOnARecordThatAlwaysFails(t *testing.T) {
	f := &fakeKinesis{}
	c := &client{client: f}

	err := c.BulkPublish(context.Background(), []types.PutRecordsRequestEntry{{Data: []byte("x"), PartitionKey: aws.String("k")}}, "stream")
	if err == nil {
		t.Fatal("BulkPublish() error = nil, want an error for the record that never got in")
	}
	if got := f.calls.Load(); got > maxBulkAttempts {
		t.Errorf("PutRecords called %d times, want at most %d", got, maxBulkAttempts)
	}
}

// flakyKinesis rejects every record of the first PutRecords and takes them
// on the second.
type flakyKinesis struct {
	fakeKinesis
}

func (f *flakyKinesis) PutRecords(ctx context.Context, in *kinesis.PutRecordsInput, opts ...func(*kinesis.Options)) (*kinesis.PutRecordsOutput, error) {
	if f.calls.Load() == 0 {
		return f.fakeKinesis.PutRecords(ctx, in, opts...)
	}
	f.calls.Add(1)
	records := make([]types.PutRecordsResultEntry, len(in.Records))
	for i := range records {
		records[i] = types.PutRecordsResultEntry{SequenceNumber: aws.String("1"), ShardId: aws.String("shard-1")}
	}
	return &kinesis.PutRecordsOutput{FailedRecordCount: aws.Int32(0), Records: records}, nil
}

func TestBulkPublish_RetriesTheRejectedRecords(t *testing.T) {
	f := &flakyKinesis{}
	c := &client{client: f}

	if err := c.BulkPublish(context.Background(), []types.PutRecordsRequestEntry{{Data: []byte("x"), PartitionKey: aws.String("k")}}, "stream"); err != nil {
		t.Fatalf("BulkPublish() error = %v", err)
	}
	if got := f.calls.Load(); got != 2 {
		t.Errorf("PutRecords called %d times, want 2", got)
	}
}
