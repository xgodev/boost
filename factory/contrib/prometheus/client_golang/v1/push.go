package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/log"
)

func Push(ctx context.Context) {
	if PushGatewayEnabled() {
		if PushGatewayAsync() {
			go pushMetrics(ctx)
		} else {
			pushMetrics(ctx)
		}
	}
}

func pushMetrics(ctx context.Context) {
	if err := push.New(PushGatewayURL(), boost.ApplicationName()).
		Gatherer(prometheus.DefaultGatherer).
		Push(); err != nil {
		logger := log.FromContext(ctx)
		logger.WithError(err).Warnf("error on push metrics")
	}
}
