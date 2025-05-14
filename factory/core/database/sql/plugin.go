package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

// Plugin defines a go driver for oracle plugin signature.
type Plugin func(context.Context, *sql.DB, driver.Connector) (*sql.DB, error)
