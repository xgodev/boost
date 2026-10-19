package sqs

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/mock"
)

// TestResolveQueueUrl_Concurrent resolves queues from many goroutines, as the
// publisher does for a batch; run with -race.
func TestResolveQueueUrl_Concurrent(t *testing.T) {
	m := &mockSqsClient{}
	m.On("GetQueueUrl", mock.Anything, mock.Anything).Return(func(_ context.Context, in *sqs.GetQueueUrlInput, _ ...func(*sqs.Options)) *sqs.GetQueueUrlOutput {
		return &sqs.GetQueueUrlOutput{QueueUrl: aws.String("https://sqs/" + *in.QueueName)}
	}, nil)
	c := &client{client: m, queueUrls: map[string]*string{}}

	var wg sync.WaitGroup
	for i := range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			name := fmt.Sprintf("queue-%d", i%4)
			url, err := c.ResolveQueueUrl(context.Background(), name)
			if err != nil || *url != "https://sqs/"+name {
				t.Errorf("ResolveQueueUrl(%s) = %v, %v", name, url, err)
			}
		}()
	}
	wg.Wait()
}
