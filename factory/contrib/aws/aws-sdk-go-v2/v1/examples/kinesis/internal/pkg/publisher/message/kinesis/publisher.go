// Package kinesis puts the messages on a stream, through the kinesis client
// of boost's factory.
package kinesis

import (
	"context"

	akinesis "github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/kinesis/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/kinesis/internal/pkg/domain/publisher/message"
)

type publisher struct {
	client kinesis.Client
	o      *Options
}

// New loads the aws config and provides the publisher on the kinesis client:
// the generator does not build aws clients, so the constructor calls the
// factory.
// @Provide(index=0)
func New(ctx context.Context, o *Options) (message.Publisher, error) {
	cfg, err := aws.NewConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &publisher{client: kinesis.NewClient(akinesis.NewFromConfig(cfg)), o: o}, nil
}

func (p *publisher) Publish(ctx context.Context, msg *entity.Message) error {
	return p.client.Publish(ctx, &akinesis.PutRecordInput{
		StreamName:   &p.o.Stream,
		PartitionKey: &p.o.PartitionKey,
		Data:         []byte(msg.Text),
	})
}
