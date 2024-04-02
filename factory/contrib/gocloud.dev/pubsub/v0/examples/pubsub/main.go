package main

import (
	"context"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
	p "gocloud.dev/pubsub"
)

func main() {

	config.Load()

	ctx := context.Background()

	ilog.New()

	logger := log.FromContext(ctx)

	topic, err := pubsub.NewTopic(ctx)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	meta := map[string]string{}

	data := []byte("Hello, World!")

	message := &p.Message{
		Body:     data,
		Metadata: meta,
	}

	if err := topic.Send(ctx, message); err != nil {
		logger.Fatalf(err.Error())
	}

	defer topic.Shutdown(ctx)

	logger.Infof("sucesss message send")

	// Don't works using memory
	// subscription, err := gocloud.NewSubscription(ctx)
	// if err != nil {
	// 	logger.Fatalf(err.Error())
	// }

	// Loop on received messages.
	// for {
	// 	m, err := subscription.Receive(ctx)
	// 	if err != nil {
	// 		logger.Info("Receiving message: %v", err)
	// 		break
	// 	}
	// 	logger.Info("Got message: ", string(m.Body))
	// 	m.Ack()
	// }

	// defer subscription.Shutdown(ctx)
}
