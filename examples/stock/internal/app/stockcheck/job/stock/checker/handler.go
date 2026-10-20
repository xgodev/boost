// Package checker checks the stock of the skus Options lists in app
// stock-check, a job: the app runs the check once and exits, the shape of a
// Kubernetes CronJob.
package checker

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/finder"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Handler checks the stock of skus.
type Handler struct {
	svc  finder.Service
	skus []string
}

// New provides the handler; the generated job door gets it from here.
// @Provide(index=0)
func New(svc finder.Service, o *Options) *Handler {
	return &Handler{svc: svc, skus: o.SKUs}
}

// Check logs the units of each sku and warns of the sold-out ones. A sku the
// stock does not have stops the check with its error, and the app exits 1.
// @Job(app=stock-check)
func (h *Handler) Check(ctx context.Context) error {
	logger := log.FromContext(ctx)
	for _, sku := range h.skus {
		s, err := h.svc.Find(ctx, sku)
		if err != nil {
			return errors.Annotatef(err, "checking %s", sku)
		}
		if s.Units == 0 {
			logger.Warnf("%s is sold out", s.SKU)
			continue
		}
		logger.Infof("%s has %d units", s.SKU, s.Units)
	}
	return nil
}
