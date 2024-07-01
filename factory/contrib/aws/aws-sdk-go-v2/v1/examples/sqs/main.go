package main

import (
	"context"
	asqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
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

	ilog.New()

	// get logrus instance from context
	logger := log.FromContext(ctx)

	// create default aws config
	awsConfig, err := aws.NewConfig(ctx)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	// create sns client
	sqsClient := asqs.NewFromConfig(awsConfig)
	client := sqs.NewClient(sqsClient)

	input := &asqs.SendMessageInput{
		MessageBody:             nil,
		QueueUrl:                nil,
		DelaySeconds:            0,
		MessageAttributes:       nil,
		MessageDeduplicationId:  nil,
		MessageGroupId:          nil,
		MessageSystemAttributes: nil,
	}

	// publish
	err = client.Publish(ctx, input)
	if err != nil {
		logger.Fatalf(err.Error())
	}

}
