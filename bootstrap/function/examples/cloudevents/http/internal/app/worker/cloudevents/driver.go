// Package cloudevents holds what the function of app worker needs besides its
// handler.
package cloudevents

import (
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/driver/extra/noop"
)

// Driver is the driver the publisher middleware publishes the replies with:
// noop only logs them; an app that delivers them provides the nats, kafka or
// pubsub driver here.
// @Provide(index=0)
func Driver() publisher.Driver {
	return noop.New()
}
