# inject: factory clients from `@Inject` -- design

Issue #85. Piece 4 of the builder boilerplate roadmap
(`docs/builder-boilerplate/README.md`), after `@Middleware` (#72), the
implementation switch (#75) and `@Config` (#80).

## Intent

A builder app that talks to an HTTP API, Redis or Mongo writes by hand the
fx module that builds the boost factory client with its plugins
(claude-plugin layout: `internal/pkg/fx/modules/lib/boost/factory/<factory>/module.go`).
The maintainer wants that generated: the constructor that receives the
client says which client it wants and with which plugins, and inject writes
the rest.

## Decided (by the maintainer)

```go
// internal/pkg/client/quote/resty/client.go
// @Provide(index=0)
// @Inject(index=0, name=quote, plugins=[otel, prometheus])
func New(c *resty.Client) quote.Repository { ... }
```

- Each client has its own config (one client per service, not one shared
  client per app).
- The client is identified by the `name` of `@Inject`: two clients of the
  same factory (two resty, for two APIs) are two names. Every constructor
  asking for the same name receives the same instance.
- Plugins are chosen in the annotation; config only turns each one off
  (every boost plugin registers `.enabled = true`, so wiring them all would
  turn on datadog and otel together).

## Behaviour

- **Trigger.** A parameter of a `@Provide`/`@Invoke` function whose type is
  a supported factory client (table below) and that carries `@Inject`. If
  the app itself provides that value (same type and `name`) through a
  `@Provide`, the app's provider is used and nothing is generated. A client
  parameter without `@Inject` is left to the app, as today.
- **Name.** `@Inject(index=N, name=x)` is the client `x`;
  `@Inject(index=N)` is the default (unnamed) client of that factory. A name
  matches `^[a-z][a-zA-Z0-9]*$` so an env var reaches its keys.
- **Plugins.** `plugins=[a, b]` on the same `@Inject`, in order; names come
  from the factory's registry (table below). Several consumers of one client
  may each list plugins only if the lists are equal; a consumer without
  `plugins` takes the client as declared elsewhere; no list anywhere means a
  client without plugins.
- **Config.** The default client keeps the factory's root
  (`boost.factory.resty.*`, env `BOOST_FACTORY_RESTY_*`). A named client
  lives under `<factory root>.clients.<name>` (`boost.factory.resty.clients.quote.host`,
  env `BOOST_FACTORY_RESTY_CLIENTS_QUOTE_HOST`), with every factory key
  registered there with the factory's default. Plugin keys sit under the
  client: `<client root>.plugins.<plugin>.*`
  (`boost.factory.resty.clients.quote.plugins.otel.enabled`, default `true`).
- **Generated file.** inject writes
  `gen/client/<client type import path>/<name>_client.go` (`default_client.go`
  for the unnamed one): `init()` registers the client and plugin keys through
  the factory's and each plugin's `ConfigAdd(path)`, and
  `<Name>ClientModule()` (once-guarded, with `contextfx.Module()`) provides
  the client -- with the fx `name` tag when named -- built by the factory's
  `...WithConfigPath(ctx, path, plugins...)`, each plugin built by its own
  `...WithConfigPath(<client root>.plugins.<plugin>)`. Every consumer's
  module includes it; so does `gen.Module()`.
- **Implementation switch.** A client used only by members of an
  implementation set (#75) is built only when such a member is chosen,
  because only the members' modules include it.

## Supported clients (first registry)

| Client type | Factory | Plugins |
|---|---|---|
| `*resty.Client` (`github.com/go-resty/resty/v2`) | `factory/contrib/go-resty/resty/v2` | `log`, `health`, `otel`, `datadog`, `prometheus`, `retry`, `requestid` |
| `*redis.Client` (`github.com/redis/go-redis/v9`) | `factory/contrib/redis/go-redis/v9` | `health`, `otel`, `datadog`, `prometheus` |
| `*redis.ClusterClient` (`github.com/redis/go-redis/v9`) | same | the plugins that have a cluster constructor |
| `*mongo.Conn` (boost `factory/contrib/go.mongodb.org/mongo-driver/v1` and `/v2`) | same | `health`, `otel`, `datadog`, `prometheus` |

The plugin name is the last segment of the plugin's config root
(`resty.PluginsRoot + ".otel"` -> `otel`). Adding a factory is one registry
entry. `*sql.DB` stays out: several drivers return it (pgx, godror), so the
type alone does not say which factory to use.

## Failures (generation stops, naming file and function)

- `plugins` on a parameter that is not a supported client, or on a client
  the app provides itself.
- An unknown plugin name (the error lists the valid ones) or a repeated one.
- Two consumers of one client with different plugin lists.
- `group` on a client parameter with `plugins`.
- A name that does not match `^[a-z][a-zA-Z0-9]*$`.

## Out of scope

- `*sql.DB` (pgx, godror) and the other factories (elasticsearch, gocql,
  pubsub, nats, ...): later registry entries.
- Plugin options beyond config (a hand-built `retry` policy stays
  hand-written).

## Calls made by the builder of this design (not the maintainer's)

- A client param needs `@Inject` to be generated: unannotated client params
  keep today's behaviour ("left to the app"), so an inject-only app that
  includes `restyfx.Module()` by hand does not get a duplicate provider.
- Named-client path `<factory root>.clients.<name>` and the
  `gen/client/...` file layout, mirroring `gen/implementation/...`.
- Every generated client registers its keys with `ConfigAdd(path)`, the
  default one included, so a plugin whose `init` registers under the wrong
  root (mongo's `prometheus` uses `resty.PluginsRoot`) still gets its keys
  under the client.
- Example: `builder/examples/simple` gains a `stock/remote` implementation
  that asks a warehouse API through `@Inject(index=0, name=warehouse,
  plugins=[log])`; its test points
  `BOOST_FACTORY_RESTY_CLIENTS_WAREHOUSE_HOST` at an `httptest` server.

## Testing

- inject unit tests: trigger, default vs named, shared instance, plugin
  list rules, app provider wins, every failure above, rendered output
  compared with golden files.
- The simple example boots with `stock.implementation=remote` against the
  `httptest` server and answers through the generated client; `gen/` stale
  check as today.
