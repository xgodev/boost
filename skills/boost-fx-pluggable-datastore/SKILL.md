---
name: boost-fx-pluggable-datastore
description: "Use when a Go service needs the concrete implementation behind a domain port (repository / datastore / cache / provider) to be selectable at boot by configuration — e.g. mongodb in prod, memory in local/dev — without the choice leaking into application or domain layers, in a boost + uber-go/fx service. Covers the config-key-in-the-port pattern, per-implementation fx.Module providers, and the switch dispatcher that returns one fx.Option. Triggers on swapping mongo for in-memory via env, pluggable repository/provider, multiple implementations of one interface chosen by config, datastore selection, or fx.As/ResultTags gymnastics to bind an interface."
license: MIT
metadata:
  author: jpfaria
  version: "0.2.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-fx-modules` (Module() shape), `boost-wrapper-config` (config.Add / env override).

## The pattern

One domain port, N implementations, the live one chosen at boot by a single config key. Selection is `config string → switch → fx.Option`. The `application`/`domain` layers only ever see the port interface — they never learn which concrete answered, and never import mongo/redis/memory.

```
domain/.../rolloutfilterer/
  repository.go        # the port (interface) — pure, imports nothing infra
  config.go            # config.Add(<key>, "memory") + DataStore() string  ← choice point lives HERE
fx/modules/datastore/.../rolloutfilterer/
  memory/module.go     # fx.Module(name, fx.Provide(memory.New))
  mongodb/module.go    # fx.Module(name, mongodb.Module(), fx.Provide(mongodb.New))   ← carries its OWN conn
fx/modules/domain/repository/.../rolloutfilterer/
  module.go            # the dispatcher: switch DataStore() { ... }
```

## The three rules that make it work

**1. The config key + accessor live in the domain port package.** The port owns the identity of the choice point.

```go
package rolloutfilterer

import "github.com/xgodev/boost/wrapper/config"

const (
	root      = "internal.pkg.domain.repository.categorizer.rollout-filterer"
	datastore = root + ".datastore"
)

func init() { config.Add(datastore, "memory", "defines datastore type") }

func DataStore() string { return config.String(datastore) }
```

**2. Each concrete constructor returns the PORT INTERFACE, never the concrete struct.** This is what lets fx bind it directly — no `fx.As`, no `ResultTags`/`ParamTags`, no named-instance rebinding.

```go
func New() rolloutfilterer.Repository                 { return &datastore{} }          // memory: no deps
func New(db *mongo.Database) rolloutfilterer.Repository { return &datastore{db: db} }   // mongodb: shared conn
```

**3. Each implementation's fx.Module carries its own infra dependency** (the mongo sub-module pulls in the connection's `Module()`), so deps resolve *inside the chosen case* — never in the dispatcher.

```go
// memory/module.go
func Module() fx.Option {
	return fx.Module("rolloutfilterer_datastore_memory", fx.Provide(memory.New))
}

// mongodb/module.go
func Module() fx.Option {
	return fx.Module("rolloutfilterer_datastore_mongodb",
		mongodb.Module(),          // ← the *mongo.Database provider rides along
		fx.Provide(datastore.New),
	)
}
```

## The dispatcher

A plain switch returning the chosen sub-module's `fx.Option`. The default arm is the safe local impl.

```go
func Module() fx.Option {
	switch rolloutfilterer.DataStore() {
	case "mongodb":
		return mongodb.Module()
	default:
		return memory.Module()
	}
}
```

`main` composes `rolloutfilterer.Module()` blindly. When `datastore=memory`, the mongo module is **never added to the graph** — fx never resolves `*mongo.Database`, so local/dev opens no connection.

## Adding a provider (e.g. redis) — closed, mechanical

1. New leaf `redis/` package: `New(rdb *redis.Client) rolloutfilterer.Repository` (returns the interface) + `module.go` with `fx.Module(name, redisfact.Module(), fx.Provide(redis.New))`.
2. One `case "redis":` arm in the dispatcher.

Zero change to `domain/`, `application/`, the other impls, or `main`.

## Dev default = memory (boot with zero external infra)

Default the provider to `memory` so the service boots standalone in dev — **production overrides to the real backend via env** (`MYAPP_REDIS_PROVIDER=boost`, …). The payoff is real only if the `memory` arm wires *nothing* that touches the network:

- The dispatcher provides **only** the in-memory impl for `memory`. The real factory (`mongofact.NewConn`, `redisfact.NewClient`) is never invoked → no eager Mongo ping that blocks boot, no Redis client spamming `connection refused` and downing the process.
- **Cache port:** consumers depend on `cache.Driver` (the interface), not the concrete `*cacheredis.Client`. The selector wires the Redis driver for `boost`, a noop / in-memory driver for `memory` — same `Manager[T]` either way.
- **Health checkers** take their client/conn as an **optional fx dep** and noop when it's nil, so selecting `memory` doesn't force the real connection to exist just to satisfy `/health`:

```go
type RedisHealthParams struct {
	fx.In
	Client *goredis.Client `optional:"true"` // nil when redis.provider=memory
}
func NewRedisChecker(p RedisHealthParams) health.Checker {
	if p.Client == nil { return health.Noop{} }
	return &redisChecker{client: p.Client}
}
```

## Red flags

| Red flag | Fix |
|---|---|
| `os.Getenv` / reading the provider in `main` | `config.Add(<key>, default)` in the **domain port package** + `DataStore()` accessor |
| The `switch` living in `main.go` or in `application/` | It lives only in the fx dispatcher module; `main` calls `Module()` blindly |
| `New` returns `*datastore` (concrete) + `fx.As`/`fx.ResultTags`/`fx.ParamTags` to rebind it to the interface | Return the interface straight from `New` — no tag gymnastics, no named instances |
| The shared `*mongo.Database` provided globally even when `memory` is selected | Let each impl module carry its own conn `Module()`, so it's only wired when that case is chosen |
| `//go:build mongodb` build tags to pick the impl | One binary, runtime `config` switch — env-tunable per environment, no rebuild |
| Dispatcher panics / silently falls back on an unknown value | `default:` = the safe local impl, OR return `fx.Error(...)` to fail-fast at boot |
| The `memory` arm still constructs the real client (eager ping / dial) | The memory case provides ONLY the in-memory impl — never call the boost factory, so there is no network at boot |
| Production silently boots in-memory because the default is `memory` | Deploy MUST set `<key>=boost` (real) via env; state it in the `config.Add` description |
| Cache consumers depend on the concrete `*cacheredis.Client` | Depend on `cache.Driver` (interface) so the `memory` arm can swap a noop driver |
| Health checker hard-requires the real client/conn | Take it as an `optional:"true"` fx dep and noop when nil |
