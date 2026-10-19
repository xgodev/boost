// Package orders manages the Pub/Sub topics of orders through the "orders"
// client.
package orders

import (
	"context"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

// Topics creates Pub/Sub topics.
type Topics struct {
	client *pubsub.Client
}

// NewTopics provides the topics on the orders client, which inject builds
// with the log plugin.
// @Provide(index=0)
// @Inject(index=0, name=orders, plugins=[log])
func NewTopics(client *pubsub.Client) *Topics {
	return &Topics{client: client}
}

// Create creates the topic id in the client's project.
func (t *Topics) Create(ctx context.Context, id string) error {
	_, err := t.client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: "projects/" + t.client.Project() + "/topics/" + id,
	})
	return err
}
