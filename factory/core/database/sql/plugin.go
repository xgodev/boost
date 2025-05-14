package sql

import (
	"context"
	"database/sql"
)

// Plugin defines a go driver for oracle plugin signature.
type Plugin func(context.Context, *sql.DB) error
