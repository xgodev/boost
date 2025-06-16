package pgx

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	sqll "github.com/xgodev/boost/factory/core/database/sql"
	"github.com/xgodev/boost/wrapper/log"
)

// NewDBWithConfigPath returns a new sql DB with options from config path.
func NewDBWithConfigPath(ctx context.Context, path string, plugins ...sqll.Plugin) (*sql.DB, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewDBWithOptions(ctx, opts, plugins...)
}

// NewDBWithOptions returns a new sql DB with options.
func NewDBWithOptions(ctx context.Context, o *Options, plugins ...sqll.Plugin) (db *sql.DB, err error) {

	logger := log.FromContext(ctx)

	// Parse DSN into a Connector
	cfg, err := pgx.ParseConfig(o.ConnectString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	if cfg.User == "" {
		cfg.User = o.Username
	}

	if cfg.Password == "" {
		cfg.Password = o.Password
	}

	connector := stdlib.GetConnector(*cfg)

	db, err = sqll.NewDB(ctx, connector, plugins...)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Connected to PostgreSQL (pgx) server")

	return db, err
}

// NewDB returns a new DB.
func NewDB(ctx context.Context, plugins ...sqll.Plugin) (*sql.DB, error) {

	logger := log.FromContext(ctx)

	o, err := NewOptions()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	return NewDBWithOptions(ctx, o, plugins...)
}
