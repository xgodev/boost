package cloudevents

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/utils/strings"
	"github.com/xgodev/boost/wrapper/log"
)

func isSQSEvent(ctx context.Context, in *v2.Event) bool {
	return strings.SliceContains([]string{"com.amazon.sqs.message", "aws.sqs.message"}, in.Type())
}

func fromSQS(ctx context.Context, in *v2.Event) {

	logger := log.FromContext(ctx)

	var err error
	var sqsMessage events.SQSMessage

	err = json.Unmarshal(in.Data(), &sqsMessage)
	if err != nil {
		logger.Error(err)
	}

	var snsEntity events.SNSEntity

	err = json.Unmarshal([]byte(sqsMessage.Body), &snsEntity)
	if err != nil {
		logger.Error(err)
	}

	var data []byte

	if snsEntity.Message != "" {
		data = []byte(snsEntity.Message)
	} else {
		data = []byte(sqsMessage.Body)
	}

	ctype := v2.TextPlain
	var js map[string]interface{}

	if json.Unmarshal(data, &js) == nil {
		ctype = v2.ApplicationJSON
	}

	err = in.SetData(ctype, data)
	if err != nil {
		logger.Error(err)
	}

}
