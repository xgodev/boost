package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type Plugin interface {
	// WrapConnector é chamado antes de abrir o *sql.DB*,
	// pra você trocar o driver/connector se quiser.
	WrapConnector(ctx context.Context, connector driver.Connector) (driver.Connector, error)
	// InitDB é chamado logo depois de abrir o *sql.DB*,
	// pra você configurar métricas ou outros hooks.
	InitDB(ctx context.Context, db *sql.DB) error
}

func NewDBWithOptions(
	ctx context.Context,
	connector driver.Connector,
	options *Options,
	plugins []Plugin,
) (*sql.DB, error) {
	// 1) wrap do connector
	var err error
	for _, pl := range plugins {
		if pl != nil {
			connector, err = pl.WrapConnector(ctx, connector)
			if err != nil {
				return nil, fmt.Errorf("connector plugin: %w", err)
			}
		}
	}

	// 2) abre o DB
	db := sql.OpenDB(connector)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)
	db.SetConnMaxIdleTime(options.ConnMaxIdletime)
	db.SetMaxIdleConns(options.MaxIdletime)
	db.SetMaxOpenConns(options.MaxOpenConns)

	// 3) init no DB
	for _, pl := range plugins {
		if err := pl.InitDB(ctx, db); err != nil {
			return nil, fmt.Errorf("db plugin: %w", err)
		}
	}

	// 4) ping
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
