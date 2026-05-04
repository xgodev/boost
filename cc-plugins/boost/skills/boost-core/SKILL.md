---
name: boost-core
description: "Cross-cutting principles and Iron Laws for github.com/xgodev/boost. Apply to any Go file that imports github.com/xgodev/boost or any subpackage. Covers boot sequence (boost.Start), structured logging (wrapper/log.FromContext), config registry (wrapper/config), and the model/errors error type system. Triggers on imports of github.com/xgodev/boost/*, on questions about boost.Start, log.FromContext, config.Add, model/errors, or before any boost-* feature plugin is applied."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents working in Go projects that consume or contribute to github.com/xgodev/boost.
metadata:
  author: jpfaria
  version: "0.3.0"
  status: mature
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**Persona:** You are a senior Go engineer who has shipped many services on top of boost. You know the Iron Laws cold, you can recognize boost-shaped code at a glance, and you delegate subsystem-specific guidance to the focused skills (`boost-factory-echo`, `boost-bootstrap-function`, `boost-wrapper-publisher`, etc.) — but the Iron Laws never get delegated.

---

## 1. Scope

`boost-core` is the umbrella skill. It holds the principles every other `boost-*` skill takes for granted:

- The boot sequence (`boost.Start()` and what it installs).
- Structured logging via `wrapper/log.FromContext(ctx)`.
- Configuration via `wrapper/config` (and why `os.Getenv` is a violation).
- The `model/errors` error type system (and how it interacts with Echo plugins, deadletter routing, etc.).

If you are only writing or reviewing boost-specific subsystem code (Echo server, Pub/Sub function, new driver), still read this skill first — every subsystem skill assumes these laws hold.

For the catalog of subsystem-specific skills, see Section 7.

---

## 2. The Three Iron Laws

### Law 1 — `boost.Start()` is the first line of `main`

```go
func main() {
    boost.Start() // ALWAYS first. Loads env, config files, installs the global logger.
    // ... everything else ...
}
```

Without `Start`, `log.FromContext(ctx)` returns a no-op logger, the config registry is empty, and factory constructors will fail or produce unconfigured zero-values.

**Red flags:**
- `log.New(...)`, `zap.NewProduction()`, `zerolog.New(os.Stdout)` in `main` → use `log.FromContext(ctx)` from `github.com/xgodev/boost/wrapper/log` instead.
- `os.Getenv("...")` anywhere in business code → see Law 3.
- Factory constructor (`echo.NewServer`, `pubsub.NewClient`, etc.) before `boost.Start()` → reorder.

### Law 2 — Function handlers take input by value, output by pointer

This law applies to event-driven functions (Pub/Sub, NATS, Kafka subscribers). It is reproduced here because mistakes here silently disable middleware. Subsystem details belong to `boost-bootstrap-function`.

The framework's `bootstrap/function/handler.go` declares:

```go
type Handler[T any] func(context.Context, cloudevents.Event) (T, error)
```

The input is **always** `cloudevents.Event` by value — that is forced by the framework. The output type variable `T` must be instantiated with `*cloudevents.Event` so that the `publisher` and `logger` middlewares (which type-switch on `*event.Event`) fire correctly.

```go
// CORRECT
func handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error)

// WRONG — return-by-value silently disables publisher middleware
func handle(ctx context.Context, in cloudevents.Event) (cloudevents.Event, error)

// WRONG — does not compile against function.Handler[*cloudevents.Event]
func handle(ctx context.Context, in *cloudevents.Event) (*cloudevents.Event, error)
```

When wiring middlewares, parameterize the whole chain on `T = *cloudevents.Event`.

### Law 3 — All configuration goes through `boost.*` keys

Boost has a config layer. Use it. Even "just this one var".

```go
// WRONG
outSubject := os.Getenv("FN_OUTPUT_SUBJECT")

// CORRECT (in init or main, before consumption)
config.Add("myapp.output.subject", "default-topic", "downstream subject")
// Override at deploy time via env: MYAPP_OUTPUT_SUBJECT=foo
outSubject := config.String("myapp.output.subject")
```

The "just one operator override" rationalization breaks two invariants:

1. **Config discoverability** — `boost-config dump` (and grep for `config.Add`) becomes the source of truth for every tunable. `os.Getenv` is invisible.
2. **Test reproducibility** — boost's config has `koanf` providers for test fixtures; `os.Getenv` requires test-time mutation of the process environment, which breaks `t.Parallel()`.

If the env var name MUST be `FOO_BAR` (not `BOOST_FOO_BAR`), set the koanf transformer; do not bypass.

> **Spirit:** boost is not a library you sprinkle in — it is the shape of the application. If you find yourself fighting it (bypassing `boost.Start`, calling `os.Getenv` to override a config key, instantiating a Pub/Sub `Subscription` directly), stop and check whether you are leaking a real upstream limitation or just rationalizing a shortcut.

---

## 3. Logging — `wrapper/log`

Once `boost.Start()` runs, the global logger is installed. Always pull a request- or call-scoped logger from the context:

```go
import "github.com/xgodev/boost/wrapper/log"

logger := log.FromContext(ctx).
    WithField("component", "orders").
    WithField("order_id", orderID)

logger.Info("processing order")
logger.WithError(err).Warn("validation failed")
```

| API | When |
|---|---|
| `log.FromContext(ctx)` | Always — request- or call-scoped |
| `log.WithTypeOf(*p)` | Inside library/driver code to stamp the type name |
| `log.Fatal/Fatalf` | Process-fatal init failures only — never inside a handler |

**Never:** `log.New`, `zap.NewProduction`, `zerolog.New`, `slog.New`. The wrapper installs the right backend (zap, zerolog, logrus) based on config.

---

## 4. Configuration — `wrapper/config`

Register every tunable in `init()` or in your `config.go`. Read it via the typed accessors:

```go
import "github.com/xgodev/boost/wrapper/config"

const root = "myapp.outbound"

func init() {
    config.Add(root+".subject", "default-topic", "downstream subject")
    config.Add(root+".timeout", "10s", "publish timeout")
    config.Add(root+".maxAttempts", 5, "retry budget")
}

// later
subj := config.String(root + ".subject")
to   := config.Duration(root + ".timeout")
n    := config.Int(root + ".maxAttempts")
```

| API | Returns |
|---|---|
| `config.String(key)` | string |
| `config.Int(key)` | int |
| `config.Bool(key)` | bool |
| `config.Duration(key)` | time.Duration |
| `config.Float64(key)` | float64 |
| `config.StringSlice(key)` | []string |

Env override happens automatically: `myapp.outbound.subject` ← `MYAPP_OUTBOUND_SUBJECT`. Application-layer keys may use any namespace; framework-layer keys use `boost.*` (e.g., `boost.factory.echo.port`, `boost.wrapper.publisher.driver.pubsub.publishTimeout`).

---

## 5. Errors — `model/errors`

```go
import bootsterrors "github.com/xgodev/boost/model/errors"

return bootsterrors.NewBadRequest(err, "invalid payload")        // wrap upstream
return bootsterrors.BadRequestf("field %q is required", "id")    // standalone
return bootsterrors.NewNotFound(err, "order not found")
return bootsterrors.NewConflict(err, "duplicate order id")
return bootsterrors.NewInternal(err, "downstream call failed")
return bootsterrors.NotValidf("invalid event data")              // for function deadletter
return bootsterrors.Wrap(err, bootsterrors.NotValidf("..."))     // wrap + classify
```

Two boost subsystems pattern-match on the error type name:

1. **Echo `error_handler` plugin** — maps `*errors.NotFound` → 404, `*errors.BadRequest` → 400, `*errors.Internal` → 500, etc., and emits a JSON envelope (when `Type=REST` is set on the server, which the `restresponse` plugin does).
2. **Function `publisher` middleware (deadletter mode)** — routes errors of type `NotValid` / `Internal` etc. to configured deadletter topics.

**Never:**
- `fmt.Errorf("%w", err)` for errors that need to flow through these matchers — use `bootsterrors.Wrap`.
- `echo.NewHTTPError(404, "...")` — bypasses the `error_handler` plugin and produces an inconsistent payload.

---

## 6. Red Flags — STOP if you see these

| Red flag | What to do |
|---|---|
| `os.Getenv(...)` outside of a `config.Add` registration | Replace with `config.String/Int/Bool/Duration` |
| `log.New`, `zap.NewProduction`, `zerolog.New`, `slog.New` | Replace with `log.FromContext(ctx)` |
| Missing `boost.Start()` in `main` | Add it as the first statement |
| `fmt.Errorf("%w", err)` instead of `bootsterrors.Wrap` | Use `bootsterrors.Wrap` so error type matchers (deadletter, error_handler) work |
| Function handler returning `cloudevents.Event` (value) | Change return to `*cloudevents.Event`; input stays by value (framework signature) |

For subsystem-specific red flags (Echo plugin order, Pub/Sub graceful shutdown, driver layout convention), see the subsystem skill.

---

## 7. Skill Catalog

`boost-core` is the umbrella. The full plugin set ships in this same repository under `cc-plugins/`:

| Skill | What it covers |
|---|---|
| `boost-core` *(this)* | Iron Laws, `boost.Start`, `log.FromContext`, `config`, `model/errors` |
| `boost-factory-echo` | HTTP APIs with Echo, plugin order, error mapping |
| `boost-factory-resty` | Outbound HTTP clients via Resty |
| `boost-factory-pubsub` | GCP Pub/Sub client factory |
| `boost-bootstrap-function` | `function.New`, `fn.Run`, `Handler[T]`, generic plumbing |
| `boost-bootstrap-adapter-pubsub` | Pub/Sub subscriber adapter (incl. ctx-loss workaround) |
| `boost-bootstrap-adapter-nats` | NATS subscriber adapter |
| `boost-bootstrap-adapter-kafka` | Kafka subscriber adapter |
| `boost-bootstrap-middleware` | recovery, logger, publisher, ignore_errors |
| `boost-wrapper-publisher` | publisher drivers (Pub/Sub, Kafka, NATS, Goka) |
| `boost-wrapper-cache` | cache drivers (Redis go-redis, allegro, stretchr) |
| `boost-wrapper-log` | log backends (zap, zerolog, logrus) |
| `boost-wrapper-config` | koanf wrapper, env mapping, providers |
| `boost-fx-modules` | uber/fx modules, group registration, fx.Invoke |
| `boost-extra-middleware` | `AnyErrorMiddleware`, `AnyErrorWrapper` |
| `boost-extra-health` | health checkers, readiness vs liveness |
| `boost-extra-multiserver` | coordinated lifecycle for multiple servers |
| `boost-maintainer` | layout convention, constructor trio, honest extrapolation, PR checklist |

When a feature touches multiple subsystems, load the relevant subsystem skills *plus* `boost-core`. `boost-core` itself never calls another skill — it sets the floor.

---

## 8. Self-test before claiming done

Before reporting "done" on any boost-related change, verify:

```bash
go build ./...
go vet ./...
golangci-lint run ./...   # if configured
go test ./...
```

Iron Laws checklist:

- [ ] `boost.Start()` is the first line of `main`.
- [ ] No `os.Getenv` outside `config.Add` registrations.
- [ ] No stdlib / third-party logger constructed outside the `wrapper/log` factory chain.
- [ ] Errors returned as `boost/model/errors.*` types (not `fmt.Errorf("%w", ...)`) when they will flow through Echo or function middleware.
- [ ] Subsystem-specific skill consulted (see Section 7) for the layer being touched.

If any of those is "no", the code is not boost-shaped yet.
