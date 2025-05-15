package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/xgodev/boost/wrapper/log"
)

// Plugin define os hooks de “antes” e “depois” de sql.OpenDB.
type Plugin interface {
	WrapConnector(ctx context.Context, connector driver.Connector) (driver.Connector, error)
	InitDB(ctx context.Context, db *sql.DB) error
}

// NewDBWithOptions abre o *sql.DB* aplicando os plugins.
// (é o corpo que você já tem)
func NewDBWithOptions(
	ctx context.Context,
	connector driver.Connector,
	options *Options,
	plugins ...Plugin,
) (*sql.DB, error) {

	logger := log.FromContext(ctx)

	// 1) Wrap do connector
	var err error
	for _, pl := range plugins {
		if pl != nil {
			connector, err = pl.WrapConnector(ctx, connector)
			if err != nil {
				return nil, fmt.Errorf("connector plugin: %w", err)
			}
		}
	}

	logger.Debugf("after WrapConnector, driver is %T", connector.Driver())

	// 2) Abre o DB e configura pool
	db := sql.OpenDB(connector)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)
	db.SetConnMaxIdleTime(options.ConnMaxIdletime)
	db.SetMaxIdleConns(options.MaxIdletime)
	db.SetMaxOpenConns(options.MaxOpenConns)

	logger.Debugf("db.Driver() = %T", db.Driver())

	// 3) Init no DB
	for _, pl := range plugins {
		if pl != nil {
			if err := pl.InitDB(ctx, db); err != nil {
				return nil, fmt.Errorf("db plugin: %w", err)
			}
		}
	}

	// 4) Ping & retorna
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

// NewDBWithConfigPath carrega as Options de um arquivo e chama NewDBWithOptions.
func NewDBWithConfigPath(
	ctx context.Context,
	connector driver.Connector,
	configPath string,
	plugins ...Plugin,
) (*sql.DB, error) {
	opts, err := NewOptionsWithPath(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return NewDBWithOptions(ctx, connector, opts, plugins...)
}

// NewDB carrega as Options padrão e chama NewDBWithOptions.
func NewDB(
	ctx context.Context,
	connector driver.Connector,
	plugins ...Plugin,
) (*sql.DB, error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to load default options: %w", err)
	}
	return NewDBWithOptions(ctx, connector, opts, plugins...)
}
