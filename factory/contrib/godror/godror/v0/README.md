# godror Factory

## Overview

Builds a `*sql.DB` for Oracle through [godror](https://github.com/godror/godror),
with the `database/sql` plugins of `factory/core/database/sql/plugins`
(`health`, `otel`, `datadog`). godror needs cgo to build and the Oracle
client libraries (Instant Client) at run time.

## Usage

```go
boost.Start()

db, err := godror.NewDB(ctx)                                 // root boost.factory.godror
db, err := godror.NewDBWithConfigPath(ctx, "app.db.orders")  // another root, same keys
```

In a builder app, `@Inject(index=N, drivers=[godror])` builds it (see
[`inject/README.md`](../../../../../inject/README.md#sql-databases)).

## Configuration

| Key | Default | Description |
|---|---|---|
| `boost.factory.godror.connectString` | `localhost:1521/database?connect_timeout=2` | connection string |
| `boost.factory.godror.username` | `""` | database user (used only with `password`) |
| `boost.factory.godror.password` | `""` | database password (used only with `username`) |
| `boost.factory.godror.maxLifetime` | `1h` | maximum time a session may be reused |
| `boost.factory.godror.sessionTimeout` | `5m` | session timeout |
| `boost.factory.godror.waitTimeout` | `30s` | wait timeout |
| `boost.factory.godror.maxSessions` | `1000` | maximum sessions |
| `boost.factory.godror.minSessions` | `1` | minimum sessions |
| `boost.factory.godror.sessionIncrement` | `1` | session increment |
| `boost.factory.godror.maxSessionsPerShard` | `0` | maximum sessions per shard |
| `boost.factory.godror.ping` | `true` | pings the database while building the `*sql.DB`; `false` skips it (see [Ping](#ping)) |

The pool keys of `database/sql` stay under `boost.factory.sql`.

### Ping

By default `NewDB*` pings the database and returns its error, so a database
that is down or a wrong password fails the build of the `*sql.DB`. With
`ping=false` it does not ping: the pool is lazy, so the first connection,
and its error, happen on the first query (or `db.PingContext`).

In a builder app the `*sql.DB` is built while the app boots, so a failed ping
stops the boot before a `@Job` or handler runs. Set `ping=false` when the app
must handle that error itself (for example, to send an alert):

```
BOOST_FACTORY_GODROR_CLIENTS_GOLD_PING=false
```

The key works on the root and under `clients.<name>`, like every other key.
