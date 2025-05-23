package godror

import (
	"context"
	"database/sql"
	"github.com/godror/godror"
	sqll "github.com/xgodev/boost/factory/core/database/sql"
	"github.com/xgodev/boost/factory/core/time"
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

	var params godror.ConnectionParams
	params.ConnectString = o.ConnectString
	if o.Username != "" && o.Password != "" {
		params.Username, params.Password = o.Username, godror.NewPassword(o.Password)
	}
	params.SessionTimeout = o.SessionTimeout
	params.MaxLifeTime = o.MaxLifetime
	params.MaxSessions = o.MaxSessions
	params.MinSessions = o.MinSessions
	params.MaxSessionsPerShard = o.MaxSessionsPerShard
	params.Timezone = time.Location()
	params.WaitTimeout = o.WaitTimeout
	params.SessionIncrement = o.SessionIncrement
	// params.SetSessionParamOnInit("NLS_NUMERIC_CHARACTERS", ",.")
	// params.SetSessionParamOnInit("NLS_LANGUAGE", "FRENCH")

	connector := godror.NewConnector(params)

	db, err = sqll.NewDB(ctx, connector, plugins...)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Connected to Oracle (godror) server")

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
