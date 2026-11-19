# pgx Factory

## Overview

Builds a `*sql.DB` for PostgreSQL through [pgx/v5](https://github.com/jackc/pgx)
(`stdlib` connector), with the `database/sql` plugins of
`factory/core/database/sql/plugins` (`health`, `otel`, `datadog`).

## Usage

```go
boost.Start()

db, err := pgx.NewDB(ctx)                                  // root boost.factory.pgx
db, err := pgx.NewDBWithConfigPath(ctx, "app.db.orders")   // another root, same keys
```

In a builder app, `@Inject(index=N, drivers=[pgx])` builds it (see
[`inject/README.md`](../../../../../inject/README.md#sql-databases)).

## Configuration

| Key | Default | Description |
|---|---|---|
| `boost.factory.pgx.connectString` | `postgres://localhost:5432/postgres?sslmode=disable` | connection string (URL or `key=value`) |
| `boost.factory.pgx.username` | `""` | database user |
| `boost.factory.pgx.password` | `""` | database password |
| `boost.factory.pgx.ping` | `true` | pings the database while building the `*sql.DB`; `false` skips it (see [Ping](#ping)) |

### Username and password

A non-empty `username` or `password` **overrides** the one in the
connection string, and also what `pgx.ParseConfig` falls back to when the
connection string has none (`PGUSER`/`PGPASSWORD` and the OS user: `root`
in a container running as root). An empty one keeps what the connection string
resolved. So the user and password can come from a secret as env vars
while the connection string carries only host, port and database:

```
BOOST_FACTORY_PGX_CONNECTSTRING=postgres://db.example.com:5432/orders?sslmode=require
BOOST_FACTORY_PGX_USERNAME=<from the secret>
BOOST_FACTORY_PGX_PASSWORD=<from the secret>
```

`NewConnConfig(o)` returns the `*pgx.ConnConfig` this resolution builds,
without opening a connection.

### Ping

By default `NewDB*` pings the database and returns its error, so a database
that is down or a wrong password fails the build of the `*sql.DB`. With
`ping=false` it does not ping: the pool is lazy, so the first connection,
and its error, happen on the first query (or `db.PingContext`).

In a builder app the `*sql.DB` is built while the app boots, so a failed ping
stops the boot before a `@Job` or handler runs. Set `ping=false` when the app
must handle that error itself (for example, to send an alert):

```
BOOST_FACTORY_PGX_CLIENTS_GESPROM_PING=false
```

The key works on the root and under `clients.<name>`, like every other key.
