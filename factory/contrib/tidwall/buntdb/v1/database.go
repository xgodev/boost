package buntdb

import (
	"context"

	"github.com/tidwall/buntdb"
	"github.com/xgodev/boost/wrapper/log"
)

// NewDBWithConfigPath returns DB with options from config path.
func NewDBWithConfigPath(ctx context.Context, path string) (*buntdb.DB, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewDBWithOptions(ctx, options)
}

// NewDBWithOptions returns DB with options.
func NewDBWithOptions(ctx context.Context, o *Options) (db *buntdb.DB, err error) {

	logger := log.FromContext(ctx)

	var syncPolicy buntdb.SyncPolicy

	switch o.SyncPolicy {
	case 1:
		syncPolicy = 1
	case 2:
		syncPolicy = 2
	default:
		syncPolicy = 0
	}

	config := buntdb.Config{
		SyncPolicy:           syncPolicy,
		AutoShrinkPercentage: o.AutoShrink.Percentage,
		AutoShrinkMinSize:    o.AutoShrink.MinSize,
		AutoShrinkDisabled:   o.AutoShrink.Disabled,
	}

	db, err = buntdb.Open(o.Path)
	if err != nil {
		return nil, err
	}

	err = db.SetConfig(config)
	if err != nil {
		return nil, err
	}

	logger.Infof("Connected to buntdb: %s", o.Path)

	return db, err

}

// NewDB returns DB with default options.
func NewDB(ctx context.Context) (*buntdb.DB, error) {

	logger := log.FromContext(ctx)

	o, err := NewOptions()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	return NewDBWithOptions(ctx, o)
}
