package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/matryer/try"
	"github.com/xgodev/boost/bootstrap/util"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/sync/errgroup"
)

// Client represents a sqs client.
type Client struct {
	client sqs.Client
}

// NewClient creates a new sqs client.
func NewClient(c sqs.Client) *Client {
	return &Client{client: c}
}

// Publish publishes an event slice.
func (p *Client) Publish(ctx context.Context, events []*v2.Event) error {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Info("publishing to awssqs")

	if len(events) > 0 {

		return p.send(ctx, events)

	}

	logger.Warnf("no messages were reported for posting")

	return nil
}

func (p *Client) send(ctx context.Context, events []*v2.Event) (err error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	g, gctx := errgroup.WithContext(ctx)
	defer gctx.Done()

	for _, e := range events {

		event := e

		g.Go(func() (err error) {

			var rawMessage []byte

			rawMessage, err = p.rawMessage(event)
			if err != nil {
				return errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
			}

			queueUrl, err := p.client.ResolveQueueUrl(ctx, event.Subject())
			if err != nil {
				return err
			}

			input := &awssqs.SendMessageInput{
				MessageBody: aws.String(string(rawMessage)),
				QueueUrl:    queueUrl,
			}

			if group, ok := event.Extensions()["group"]; ok {
				input.MessageGroupId = aws.String(fmt.Sprintf("%v", group))
			}

			logger.WithField("subject", event.Subject()).
				WithField("id", event.ID()).
				Info(string(rawMessage))

			err = try.Do(func(attempt int) (bool, error) {
				err := p.client.Publish(gctx, input)
				if err != nil {
					return attempt < 5, errors.NewInternal(err, "could not be published in awssqs")
				}
				return false, nil
			})

			return err

		})

	}

	return g.Wait()
}

func (p *Client) rawMessage(out *v2.Event) (rawMessage []byte, err error) {
	exts := out.Extensions()

	source, ok := exts["target"]

	if ok {

		s := source.(string)

		if s == "data" {
			var data interface{}

			err = out.DataAs(&data)
			if err != nil {
				return nil, errors.Wrap(err, errors.Internalf("error on data as. %s", err.Error()))
			}

			rawMessage, err = json.Marshal(data)

		} else {
			rawMessage, err = util.JSONBytes(*out)
		}
	} else {
		rawMessage, err = util.JSONBytes(*out)
	}

	return rawMessage, err
}
