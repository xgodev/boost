package main

import (
	"context"
	akinesis "github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

const Bucket = "aws.s3.bucket"

func init() {
	config.Add(Bucket, "example", "s3 example bucket")
}

func main() {

	boost.Start()

	// create background context
	ctx := context.Background()

	// get logrus instance from context
	logger := log.FromContext(ctx)

	// create default aws config
	awsConfig, err := aws.NewConfig(ctx)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	// create sns client
	sqsClient := akinesis.NewFromConfig(awsConfig)
	client := kinesis.NewClient(sqsClient)

	input := &akinesis.PutRecordInput{
		Data:                      nil,
		PartitionKey:              nil,
		StreamName:                nil,
		ExplicitHashKey:           nil,
		SequenceNumberForOrdering: nil,
	}

	// publish
	err = client.Publish(ctx, input)
	if err != nil {
		logger.Fatalf(err.Error())
	}

}
