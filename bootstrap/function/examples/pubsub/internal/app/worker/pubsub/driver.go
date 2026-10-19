// Package pubsub holds what the function of app worker needs besides its
// handler.
package pubsub

import (
	"context"

	gpubsub "cloud.google.com/go/pubsub/v2"
	"github.com/xgodev/boost/wrapper/publisher"
	drvpubsub "github.com/xgodev/boost/wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v2"
)

// Driver is the driver the publisher middleware publishes the replies with:
// each one goes to the pubsub topic named after its subject. client is the
// one the pubsub adapter receives with.
// @Provide(index=0)
func Driver(ctx context.Context, client *gpubsub.Client) (publisher.Driver, error) {
	return drvpubsub.New(ctx, client)
}
