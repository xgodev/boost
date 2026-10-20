// Package stockaudit checks the stock of a list of items in app stock-audit,
// a job: the app runs the audit once and exits, the shape of a CronJob.
package stockaudit

import (
	"context"

	"github.com/xgodev/boost/builder/examples/simple/stock"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Options configures the audit. inject generates its config keys
// (builder.examples.simple.stockaudit.*) and provides *Options to New.
// @Config
type Options struct {
	Items []string `default:"book,pen,sold-out" desc:"items the audit checks"`
}

// Auditor checks the stock of the items Options lists.
type Auditor struct {
	stock stock.Service
	items []string
}

// New provides the auditor with the stock service of the app, wrapped in
// the middlewares @Middleware lists on it.
// @Provide(index=0)
func New(svc stock.Service, o *Options) *Auditor {
	return &Auditor{stock: svc, items: o.Items}
}

// Run asks the stock for each item and logs the ones that are out. A stock
// that fails stops the audit with its error, and the app exits 1.
// @Job(app=stock-audit)
func (a *Auditor) Run(ctx context.Context) error {
	logger := log.FromContext(ctx)
	out := 0
	for _, item := range a.items {
		available, err := a.stock.Available(ctx, item)
		if err != nil {
			return errors.Annotatef(err, "auditing %s", item)
		}
		if !available {
			out++
			logger.Warnf("%s is out of stock", item)
		}
	}
	logger.Infof("audited %d items, %d out of stock", len(a.items), out)
	return nil
}
