package main

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/xgodev/boost/config"
	igaws "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	iglog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
)

type Message struct {
	Default string `json:"default"`
}

func main() {

	config.Load()
	iglog.New()

	var err error

	ctx := context.Background()

	awsConfig, _ := igaws.NewConfig(ctx)
	// if you have already set aws credentials in your system environment variables,
	// ignore the two lines below
	awsConfig.Region = "YOUR_AWS_REGION"
	awsConfig.Credentials = credentials.
		NewStaticCredentialsProvider("YOUR_AWS_ACCESS_KEY_ID", "YOUR_AWS_SECRET_ACCESS_KEY", "")

	client := sns.NewFromConfig(awsConfig)

	topic := "arn:aws:sns:us-east-1:000000000000:changeme"

	var b []byte
	b, err = ioutil.ReadFile("examples/simple/client/example-sns.json")
	if err != nil {
		log.Fatal(err)
	}

	msg := Message{
		Default: string(b),
	}
	msgBytes, _ := json.Marshal(msg)
	msgStr := string(msgBytes)

	res, err := client.Publish(ctx, &sns.PublishInput{
		Message:          aws.String(msgStr),
		MessageStructure: aws.String("json"),
		TopicArn:         aws.String(topic),
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("published group message on topic [%s]", topic)

	resJSON, _ := json.Marshal(res)

	log.Info(string(resJSON))
}
