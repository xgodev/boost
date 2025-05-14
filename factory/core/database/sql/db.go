package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/xgodev/boost/wrapper/log"
)

// NewDB returns a new sql DB.
func NewDB(ctx context.Context, connector driver.Connector, plugins ...Plugin) (db *sql.DB, err error) {

	logger := log.FromContext(ctx)

	db = sql.OpenDB(connector)

	for _, plugin := range plugins {
		db, err = plugin(ctx, db, connector)
		if err != nil {
			return db, err
		}
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Connected to Oracle (godror) server")

	return db, err
}
