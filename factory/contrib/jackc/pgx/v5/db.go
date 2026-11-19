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

	cfg, err := NewConnConfig(o)
	if err != nil {
		return nil, err
	}

	connector := stdlib.GetConnector(*cfg)

	sqlOpts, err := sqll.NewOptions()
	if err != nil {
		return nil, err
	}
	sqlOpts.SkipPing = !o.Ping

	db, err = sqll.NewDBWithOptions(ctx, connector, sqlOpts, plugins...)
	if err != nil {
		return nil, err
	}

	if !o.Ping {
		logger.Info("PostgreSQL (pgx) pool created; ping skipped (ping=false)")
		return db, nil
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

// NewConnConfig parses the connect string and applies the username and password from options.
// A non-empty Username or Password overrides the one in the connect string, and also the
// PGUSER/PGPASSWORD env vars and OS user that pgx.ParseConfig falls back to; empty ones keep
// what the connect string resolved.
func NewConnConfig(o *Options) (*pgx.ConnConfig, error) {
	cfg, err := pgx.ParseConfig(o.ConnectString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	if o.Username != "" {
		cfg.User = o.Username
	}

	if o.Password != "" {
		cfg.Password = o.Password
	}

	return cfg, nil
}
