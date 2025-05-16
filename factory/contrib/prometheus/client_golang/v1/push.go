package prometheus

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/log"
)

var (
	once sync.Once
)

// startPusherInternal inicia o ticker + hook de shutdown
func startPusherInternal() {
	if !PushGatewayEnabled() {
		return
	}

	interval := PushInterval()
	ticker := time.NewTicker(interval)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer ticker.Stop()
		logger := log.WithField("interval", interval.String())
		logger.Debugf("starting Prometheus PushGateway pusher")

		for {
			select {
			case <-ticker.C:
				pushOnce()

			case sig := <-sigCh:
				logger.Tracef("received signal %s, doing final push", sig)
				pushOnce()
				logger.Debugf("stopping Prometheus PushGateway pusher")
				return
			}
		}
	}()
}

// StartPusher garante que o background pusher só comece 1×
func StartPusher() {
	once.Do(startPusherInternal)
}

// FlushMetrics força um push síncrono **uma única vez**, quando você chamar
func FlushMetrics() {
	pushOnce()
}

// pushOnce faz o push, se habilitado
func pushOnce() {
	if !PushGatewayEnabled() {
		return
	}
	if err := push.
		New(PushGatewayURL(), boost.ApplicationName()).
		Gatherer(prometheus.DefaultGatherer).
		Push(); err != nil {
		log.WithError(err).Warnf("error pushing metrics to PushGateway")
	}
}
