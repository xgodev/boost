package main

import (
	"context"
	asns "github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sns"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

const Bucket = "aws.s3.bucket"

func init() {
	config.Add(Bucket, "example", "s3 example bucket")
}

func main() {

	config.Load()

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
	snsClient := asns.NewFromConfig(awsConfig)
	client := sns.NewClient(snsClient)

	input := &asns.PublishInput{
		Message:                nil,
		MessageAttributes:      nil,
		MessageDeduplicationId: nil,
		MessageGroupId:         nil,
		MessageStructure:       nil,
		PhoneNumber:            nil,
		Subject:                nil,
		TargetArn:              nil,
		TopicArn:               nil,
	}

	// publish
	err = client.Publish(ctx, input)
	if err != nil {
		logger.Fatalf(err.Error())
	}

}
