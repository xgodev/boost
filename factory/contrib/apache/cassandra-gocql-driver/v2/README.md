# gocql v2 factory

Builds a `*gocql.Session` with
[`github.com/apache/cassandra-gocql-driver/v2`](https://github.com/apache/cassandra-gocql-driver),
the Apache home of `gocql/gocql`. The v1 factory (`factory/contrib/gocql/gocql/v1`)
stays; both read the same keys, so moving a service to v2 changes the import
and nothing in its config.

```go
import (
	"github.com/xgodev/boost"
	gocql "github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2"
	"github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2/plugins/local/extra/health"
)

boost.Start()

session, err := gocql.NewSession(ctx, health.Register)
if err != nil {
	return err
}
defer session.Close()
```

| Constructor | Options from |
|---|---|
| `NewSession(ctx, plugins...)` | `boost.factory.gocql` |
| `NewSessionWithConfigPath(ctx, path, plugins...)` | `path`, falling back to `boost.factory.gocql` for keys it leaves out |
| `NewSessionWithOptions(ctx, options, plugins...)` | the `*Options` given |

A plugin is `func(context.Context, *gocql.Session) error` (`factory.Plugin`).
The fx module is `fx/modules/factory/contrib/apache/cassandra-gocql-driver/v2`:
it provides the session, takes the `[]Plugin` in the graph, and closes the
session when the app stops.

## Config

| Key (`boost.factory.gocql.`) | Default | |
|---|---|---|
| `hosts` | `[127.0.0.1]` | addresses for the initial connections |
| `port` | `9042` | |
| `dc` | | local DC: routes to its nodes first and filters out the others |
| `username` / `password` | | password authentication when either is set |
| `CQLVersion` | `3.0.0` | |
| `protoVersion` | `0` | native protocol; `0` negotiates |
| `timeout` / `connecttimeout` | `600ms` | |
| `keyspace` | | |
| `numConns` | `2` | connections per host |
| `consistency` | `QUORUM` | `ANY`, `ONE`, `TWO`, `THREE`, `QUORUM`, `ALL`, `LOCAL_QUORUM`, `EACH_QUORUM`, `LOCAL_ONE` |
| `socketKeepAlive` | `0` | keepalive period, on when > 0 |
| `maxPreparedStmts` / `maxRoutingKeyInfo` | `1000` | |
| `pageSize` | `5000` | |
| `defaultTimestamp` | `true` | client-side timestamp on every request |
| `reconnectInterval` | `10ms` | |
| `maxWaitSchemaAgreement` | `60s` | |
| `disableInitialHostLookup` | `true` | skip `system.peers` |
| `writeCoalesceWaitTime` | `200µs` | |

A zero value keeps the driver's default.

## Health plugin

`plugins/local/extra/health` registers a checker that runs
`SELECT release_version FROM system.local`. Keys under
`boost.factory.gocql.plugins.health.`: `name` (`cassandra`), `description`
(`default connection`), `required` (`true`), `enabled` (`true`).

## Moving from v1

- **Errors, not panics.** A plugin that fails closes the session and
  `NewSession…` returns `gocql plugin: <err>`; v1 panicked. An unknown
  `consistency` is an error; v1 panicked inside the driver.
- **`health.NewHealth()` returns `(*Health, error)`** and `health.Register`
  returns a bad health config as its error; v1 called `log.Fatalf` and
  `Register` skipped it.
- **The health check is a real query.** v1 ran `void`, which Cassandra rejects.
- **Driver changes.** The password authenticator answers any server
  authenticator unless `AllowedAuthenticators` is set, protocols 1 and 2 are
  gone, and the driver's own defaults moved (`Timeout`/`ConnectTimeout` 11s,
  `ReconnectInterval` 60s); boost's defaults above still apply.

## Tests

`go test ./...` runs without Cassandra. The tests that talk to a server read
`BOOST_TEST_CASSANDRA=host:port` and skip without it:

```bash
docker run -d --rm --name cassandra -p 19042:9042 cassandra:5.0
BOOST_TEST_CASSANDRA=127.0.0.1:19042 go test ./factory/contrib/apache/... ./fx/modules/factory/contrib/apache/...
```
