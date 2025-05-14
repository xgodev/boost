package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/xgodev/boost/wrapper/log"
)

// NewDBWithConfigPath returns a new sql DB with options from config path.
func NewDBWithConfigPath(ctx context.Context, connector driver.Connector, path string, plugins ...Plugin) (*sql.DB, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewDBWithOptions(ctx, connector, opts, plugins...)
}

// NewDBWithOptions returns a new sql DB.
func NewDBWithOptions(ctx context.Context, connector driver.Connector, options *Options, plugins ...Plugin) (db *sql.DB, err error) {

	logger := log.FromContext(ctx)

	db = sql.OpenDB(connector)

	db.SetConnMaxLifetime(options.ConnMaxLifetime)
	db.SetConnMaxIdleTime(options.ConnMaxIdletime)
	db.SetMaxIdleConns(options.MaxIdletime)
	db.SetMaxOpenConns(options.MaxOpenConns)

	for _, plugin := range plugins {
		db, err = plugin(ctx, db, connector)
		if err != nil {
			return db, err
		}
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Connected to sql connector driver")

	return db, err
}

// NewDB returns a new DB.
func NewDB(ctx context.Context, connector driver.Connector, plugins ...Plugin) (*sql.DB, error) {

	logger := log.FromContext(ctx)

	o, err := NewOptions()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	return NewDBWithOptions(ctx, connector, o, plugins...)
}
