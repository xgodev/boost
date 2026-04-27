---
name: golang-boost
description: "Idiomatic usage and maintenance of github.com/xgodev/boost — the Go application kernel used for HTTP APIs (Echo) and event-driven functions (Pub/Sub, NATS, Kafka via CloudEvents). Apply when writing or reviewing Go code that consumes boost (new APIs, new functions/subscribers, new handlers, new bootstrap wiring), and when contributing back to boost (new factory contribs, new wrapper drivers, new bootstrap adapters or middlewares). Triggers on imports of github.com/xgodev/boost/*, on files under bootstrap/, factory/, wrapper/, fx/, or extra/, and on questions about boost.Start, log.FromContext, function.New, fn.Run, echo.NewServer, publisher.New, or wrapper.config keys."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents working in Go projects that consume or contribute to github.com/xgodev/boost.
metadata:
  author: jpfaria
  version: "0.2.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**Persona:** You are a senior Go engineer who has shipped APIs and Cloud Functions on top of boost in production at Carrefour. You know which paths are canonical, which paths exist *because* of an upstream limitation, and where the line is between "consumer code" and "contribute back to the framework". You distrust shortcuts that "work locally" but break observability, lifecycle, or reproducibility.

---

## 1. Overview

`github.com/xgodev/boost` is an opinionated kernel that wires:

- **Bootstrap** — process startup (`boost.Start()`), serverless function adapters (CloudEvents over Pub/Sub, NATS, Kafka, HTTP), middleware stacks.
- **Factory** — constructor-style integrations (Echo HTTP server, Resty client, Pub/Sub client, MongoDB, BigQuery, etc.) configurable through `BOOST_FACTORY_*` env vars.
- **Wrapper** — provider-agnostic abstractions with pluggable drivers (logger, publisher, cache, config).
- **Fx** — uber/fx modules for dependency-injected wiring at scale.
- **Extra** — health checks, generic middleware (`AnyErrorWrapper`), multi-server lifecycle.
- **Model** — shared domain types (`model/errors`, `model/restresponse`).

The skill covers two audiences:

| Audience | Trigger | Output |
|---|---|---|
| **Consumer** | Building an app that imports boost | A `main.go` that *boots the boost way* |
| **Maintainer** | Adding to `bootstrap/`, `factory/`, `wrapper/`, `fx/`, or `extra/` | Code that mirrors existing contribs and won't surprise the next reviewer |

> **Spirit:** boost is not a library you sprinkle in — it is the shape of the application. If you find yourself fighting it (bypassing `boost.Start`, calling `os.Getenv` to override a config key, instantiating a Pub/Sub `Subscription` directly), stop and check whether you are leaking a real upstream limitation or just rationalizing a shortcut.

---

## 2. The Three Iron Laws

These are non-negotiable. Violations mean the code is not "boost-shaped" even if it compiles.

### Law 1 — `boost.Start()` is the first line of `main`

```go
func main() {
    boost.Start() // ALWAYS first. Loads env, files, installs the global logger.
    // ... everything else ...
}
```

Without `Start`, `log.FromContext(ctx)` returns a no-op logger, the config registry is empty, and factory constructors will fail or produce unconfigured zero-values.

**Red flags:**
- `log.New(...)`, `zap.NewProduction()`, `zerolog.New(os.Stdout)` in `main` → use `log.FromContext(ctx)` from `github.com/xgodev/boost/wrapper/log` instead.
- `os.Getenv("...")` anywhere in business code → use `config.String/Int/Bool/Duration` with a `boost.*` key.
- "I'll just read this one env var directly, it's an operator override" → see Law 3.

### Law 2 — Handlers take input by value, output by pointer

The framework's `bootstrap/function/handler.go` declares:

```go
type Handler[T any] func(context.Context, cloudevents.Event) (T, error)
```

The input is **always** `cloudevents.Event` by value — that is forced by the framework. The pointer concern is on the **output type `T`**: middlewares like `publisher` and `logger` type-switch on `*event.Event`, so the generic must be instantiated with `T = *cloudevents.Event`. Returning a value would silently disable downstream publishing.

```go
// CORRECT — input by value, output by pointer
func handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error)

// WRONG — return-by-value silently disables publisher middleware
func handle(ctx context.Context, in cloudevents.Event) (cloudevents.Event, error)

// WRONG — does not compile against function.Handler[*cloudevents.Event]
func handle(ctx context.Context, in *cloudevents.Event) (*cloudevents.Event, error)
```

When wiring middlewares, parameterize on the pointer: `function.New[*cloudevents.Event](...)`, `lm.NewAnyErrorMiddleware[*cloudevents.Event]()`, etc. The whole chain must agree on `T = *cloudevents.Event`.

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
2. **Test reproducibility** — boost's config has `koanf` providers for test fixtures; `os.Getenv` requires test-time mutation of the process environment.

If the env var name MUST be `FOO_BAR` (not `BOOST_FOO_BAR`), set the koanf transformer; do not bypass.

---

## 3. Consumer Path — HTTP API (Echo)

### Canonical shape

```go
package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    e "github.com/labstack/echo/v4"

    "github.com/xgodev/boost"
    echoserver "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
    "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
    logplugin "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
    restresponse "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/model/restresponse"
    recoverplugin "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/recover"
    "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/requestid"
    bootsterrors "github.com/xgodev/boost/model/errors"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    logger := log.FromContext(ctx).WithField("component", "api")

    srv, err := echoserver.NewServer(ctx,
        recoverplugin.Register, // panics → 500, no process death
        requestid.Register,     // X-Request-ID propagation
        logplugin.Register,     // request access log via boost logger
        restresponse.Register,  // sets Type=REST so error_handler emits JSON
        error_handler.Register, // boost model/errors → HTTP status mapping
    )
    if err != nil {
        logger.WithError(err).Fatal("failed to build echo server")
    }

    h := newOrderHandler(/* deps... */)
    srv.GET("/health", h.Health)
    srv.POST("/orders", h.Create)

    go srv.Serve(ctx)
    <-ctx.Done()

    shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancelShutdown()
    srv.Shutdown(shutdownCtx)
}

func (h *orderHandler) Create(c e.Context) error {
    ctx := c.Request().Context()
    var req CreateOrderRequest
    if err := c.Bind(&req); err != nil {
        return bootsterrors.NewBadRequest(err, "invalid order payload")
    }
    if req.ID == "" {
        return bootsterrors.BadRequestf("field 'id' is required")
    }
    return c.JSON(http.StatusCreated, req)
}
```

### Required plugin set (production)

| Plugin | Purpose | Skip when |
|---|---|---|
| `native/recover` | Panic → 500 | Never |
| `native/requestid` | X-Request-ID | Never |
| `local/wrapper/log` | Access log via boost logger | Never |
| `local/model/restresponse` | Sets `Type=REST` on context (required for `error_handler` JSON mode) | Never (for REST APIs) |
| `extra/error_handler` | `model/errors` → status code mapping | Never (without it, errors leak as plain strings) |
| `native/cors` | Browser clients | Internal-only services |
| `native/gzip` | Response compression | gRPC / streaming |

### Errors: always use `boost/model/errors`

```go
return bootsterrors.NewBadRequest(err, "invalid payload")        // wrap upstream
return bootsterrors.BadRequestf("field %q is required", "id")    // standalone
return bootsterrors.NewNotFound(err, "order not found")
return bootsterrors.NewConflict(err, "duplicate order id")
return bootsterrors.NewInternal(err, "downstream call failed")
```

The `error_handler` plugin pattern-matches the *type name* (`*errors.NotValid`, `*errors.NotFound`, `*errors.Internal`, etc.) and emits the right HTTP status with a JSON envelope. Do **not** use `echo.NewHTTPError` — it bypasses this and produces inconsistent payloads.

---

## 4. Consumer Path — Function (Pub/Sub Subscriber)

### Canonical shape (prototype / dev)

```go
package main

import (
    "context"
    cloudevents "github.com/cloudevents/sdk-go/v2"

    "github.com/xgodev/boost"
    "github.com/xgodev/boost/bootstrap/function"
    apubsub "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v1"
    lm "github.com/xgodev/boost/bootstrap/function/middleware/logger"
    pm "github.com/xgodev/boost/bootstrap/function/middleware/publisher"
    rm "github.com/xgodev/boost/bootstrap/function/middleware/recovery"
    fpubsub "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"
    "github.com/xgodev/boost/wrapper/log"
    "github.com/xgodev/boost/wrapper/publisher"
    drvpubsub "github.com/xgodev/boost/wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v1"
)

// Input is by value (forced by function.Handler[T]), output is *cloudevents.Event
// so the publisher/logger middlewares pick it up via their *event.Event type-switch.
func handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
    // ... derive event ...
    return &out, nil
}

func main() {
    boost.Start()
    ctx := context.Background()

    pb, err := fpubsub.NewClient(ctx)
    if err != nil { log.Fatalf("pubsub client: %v", err) }
    defer pb.Close()

    drv, err := drvpubsub.New(ctx, pb)
    if err != nil { log.Fatalf("publisher driver: %v", err) }
    pub := publisher.New(drv)

    rec := rm.NewRecovery[*cloudevents.Event]()
    lmi, _ := lm.NewAnyErrorMiddleware[*cloudevents.Event]()
    pmi, _ := pm.NewAnyErrorMiddleware[*cloudevents.Event](pub)

    fn, err := function.New[*cloudevents.Event](rec, lmi, pmi)
    if err != nil { log.Fatalf("function: %v", err) }

    if err := fn.Run(ctx, handle, apubsub.New[*cloudevents.Event](pb)); err != nil {
        log.Fatalf("run: %v", err)
    }
}
```

### Production caveat — known ctx-loss in adapter helpers

The canonical wiring **currently does not propagate signal cancellation into the receiver loop**. Inside `bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v1/helper.go:51` (and the analogous NATS / Kafka helpers), the call site is:

```go
if err := subscriber.Subscribe(context.Background()); err != nil { ... }
```

A `signal.NotifyContext` passed to `fn.Run(ctx, ...)` is used by the middleware wrapper but **not** by the subscription loop. SIGTERM will not gracefully drain in-flight messages.

**Production workaround** — bypass `fn.Run` and drive the subscriber directly with a signal-aware ctx:

```go
import "github.com/xgodev/boost/extra/middleware"

ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

opts, err := apubsub.DefaultOptions()
if err != nil { log.Fatalf("subscriber options: %v", err) }

wrp := middleware.NewAnyErrorWrapper[*cloudevents.Event](
    ctx, "bootstrap", rec, lmi, pmi,
)
wrappedHandler := function.Wrapper[*cloudevents.Event](wrp, handle)

var wg sync.WaitGroup
for _, sub := range opts.Subscriptions {
    wg.Add(1)
    go func(name string) {
        defer wg.Done()
        s := apubsub.NewSubscriber[*cloudevents.Event](pb, wrappedHandler, name, opts)
        if err := s.Subscribe(ctx); err != nil && !errors.Is(err, context.Canceled) {
            log.WithField("subscription", name).Errorf("subscriber exited: %v", err)
        }
    }(sub)
}
<-ctx.Done()
wg.Wait()
```

> **Required action when using the workaround:** add a `// TODO(boost-upstream):` comment with a link to a tracking issue. The workaround compensates for an upstream limitation, not a stylistic choice — when the helper is fixed, this code should collapse back to `fn.Run`.

### Middleware order

```
recovery → logger → publisher
```

- `recovery` outermost — a panic must not kill the worker goroutine before the logger sees it.
- `logger` middle — logs both raw and post-recovery errors.
- `publisher` innermost — fires on the handler's successful result; with deadletter config, also routes errors.

### Wrapping errors for deadletter routing

Boost's deadletter middleware matches on the unwrapped error *type name*. Wrap accordingly:

```go
// Malformed input → routable to "notvalid" deadletter
return nil, bootsterrors.Wrap(err, bootsterrors.NotValidf("invalid event data"))

// Transient internal failure → routable to "internal" deadletter (or retry)
return nil, bootsterrors.Wrap(err, bootsterrors.Internalf("downstream call failed"))
```

---

## 5. Maintainer Path — Adding a Driver / Adapter / Plugin

### Layout convention

| Adding | Path |
|---|---|
| Publisher driver | `wrapper/publisher/driver/contrib/<vendor>/<lib>/v<major>/` |
| Cache driver | `wrapper/cache/driver/contrib/<vendor>/<lib>/v<major>/` |
| Logger contrib | `wrapper/log/contrib/<vendor>/<lib>/v<major>/` |
| Function adapter | `bootstrap/function/adapter/contrib/<vendor>/<lib>/v<major>/` |
| Function middleware | `bootstrap/function/middleware/<name>/` |
| Factory contrib | `factory/contrib/<vendor>/<lib>/v<major>/` |
| Echo plugin | `factory/contrib/labstack/echo/v4/plugins/{native,extra,local}/<name>/` |
| Fx module | `fx/modules/<area>/<component>/` |

**Package name = the short library name** (`pubsub`, `nats`, `confluent`, `goka`, `sns`, `redis`), **not** the version directory leaf. Existing call sites alias if a name clash with the upstream SDK package would occur (e.g., `asns "github.com/aws/aws-sdk-go-v2/service/sns"`).

**Multi-service SDKs (AWS SDK v2, Azure SDK, GCP SDK):** for `wrapper/`, `bootstrap/`, and `extra/` — split per service: `<vendor>/<service>/v<major>/` (e.g., `aws/sns/v1/`, `aws/sqs/v1/`, `aws/dynamodb/v1/`). Do **not** nest under the SDK module dir (`aws/aws-sdk-go-v2/v1/`). Reason: a publisher driver, a cache driver, and a function adapter for SNS are independent contributions — coupling them under one umbrella dir creates merge contention. The umbrella SDK layout (`factory/contrib/aws/aws-sdk-go-v2/v1/client/<service>/`) is **exclusive to `factory/contrib/`** because the factory ships *clients*, which legitimately share an SDK version pin. Drivers do not.

### Constructor trio

Every driver/adapter exposes the same trio so call sites are interchangeable:

```go
// New — defaults from env-loaded config
func New(ctx context.Context, c *upstream.Client) (publisher.Driver, error) { ... }

// NewWithOptions — explicit options
func NewWithOptions(ctx context.Context, c *upstream.Client, opts *Options) publisher.Driver { ... }

// NewWithConfigPath — load options from a specific config root
func NewWithConfigPath(ctx context.Context, c *upstream.Client, path string) (publisher.Driver, error) { ... }
```

### Config registration

Register the config root from `init()` so it's discoverable without code execution:

```go
// config.go
package <pkg>

import "github.com/xgodev/boost/wrapper/config"

const root = "boost.wrapper.publisher.driver.<name>"

func init() {
    config.Add(root+".log.level", "INFO", "log level")
    config.Add(root+".publishTimeout", "10s", "per-event publish timeout")
    // ...
}
```

### Errors and logging

```go
import (
    "github.com/xgodev/boost/model/errors"
    "github.com/xgodev/boost/wrapper/log"
)

logger := log.FromContext(ctx).WithTypeOf(*p)
// ...
return nil, errors.Wrap(err, errors.Internalf("publish failed"))
```

Match the verbs in the closest existing driver — boost is consistent enough that diverging stands out in review.

### Honest extrapolation

When the new driver has no direct precedent (e.g., SNS's `TopicArn` does not map cleanly to Pub/Sub's topic name), mark the decision:

```go
// TODO(maintainer-review): SNS requires a full ARN. Falling back to options.TopicArn
// when ev.Subject() is not an ARN. Verify this is the wanted behavior, or require
// Subject to always be an ARN.
```

Then call out the marked decisions in the PR description. Do not "decide silently" — boost has 35+ contribs and consistency is the point.

### What NOT to add

- ❌ Code under `bootstrap/` for a wrapper-layer concern, or under `wrapper/` for an HTTP framework concern. The layering is enforced by directory.
- ❌ A new top-level interface alongside an existing one ("`Driver2` is a better Driver"). Extend the existing one or open an RFC issue first.
- ❌ A direct dependency on a third-party DI / config / log library. Use `wrapper/config`, `wrapper/log`, and `fx` modules.
- ❌ Tests that `os.Setenv` to inject config — use `config.Add` with explicit defaults in the test.

---

## 6. Quick Reference — boost imports cheatsheet

```go
// Boot
"github.com/xgodev/boost"                                 // boost.Start()
"github.com/xgodev/boost/wrapper/log"                     // log.FromContext(ctx)
"github.com/xgodev/boost/wrapper/config"                  // config.Add / config.String
"github.com/xgodev/boost/model/errors"                    // BadRequestf / NewNotFound / Wrap
"github.com/xgodev/boost/model/restresponse"              // REST envelopes

// HTTP
"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
".../plugins/native/{recover,requestid,cors,gzip}"
".../plugins/local/{wrapper/log, model/restresponse}"
".../plugins/extra/error_handler"

// Function
"github.com/xgodev/boost/bootstrap/function"
"github.com/xgodev/boost/bootstrap/function/middleware/{recovery,logger,publisher,ignore_errors}"
"github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v1"
"github.com/xgodev/boost/extra/middleware"                // NewAnyErrorWrapper

// Pub/Sub
"github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"           // client
"github.com/xgodev/boost/wrapper/publisher"                                    // publisher.New(driver)
"github.com/xgodev/boost/wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v1"
```

---

## 7. Red Flags — STOP if you see these

| Red flag | What to do |
|---|---|
| `os.Getenv(...)` outside of a `config.Add` registration | Replace with `config.String/Int/Bool/Duration` |
| `log.New`, `zap.NewProduction`, `zerolog.New` | Replace with `log.FromContext(ctx)` |
| `echo.New()` directly | Replace with `factory/.../echo/v4.NewServer(ctx, ...)` |
| `pubsub.Client.Subscription(...).Receive(...)` | Replace with `apubsub.NewSubscriber(...).Subscribe(ctx)` or `function.New + fn.Run` |
| `echo.NewHTTPError(404, "...")` | Replace with `bootsterrors.NewNotFound(err, "...")` + `error_handler` plugin |
| Handler returning `cloudevents.Event` (value) | Change return to `*cloudevents.Event`; input stays by value (framework signature) |
| Bypass of `function.New + fn.Run` without a `// TODO(boost-upstream):` comment | Add the comment + tracking issue, OR use the canonical path |
| Missing `boost.Start()` in `main` | Add it as the first statement |
| Wrapping with stdlib `fmt.Errorf("%w", err)` instead of `bootsterrors.Wrap` | Use `bootsterrors.Wrap` so error type matchers (deadletter, error_handler) work |

---

## 8. Common Rationalizations

| Rationalization | Reality |
|---|---|
| "It's just one env var, the operator needs to override it." | `config.Add` accepts a default and an env override. Use it. The "discoverability" cost compounds. |
| "I'm using stdlib log because boost's logger is overkill here." | Then `log.FromContext` is *also* not overkill — it returns a no-op when not configured, costs nothing. |
| "echo.New is simpler for a one-endpoint service." | You will copy the echo factory's plugin wiring back in 30 minutes. Start there. |
| "I'll wire the recovery middleware later." | Without it, a single panic kills the function worker mid-shift. |
| "function.New + fn.Run loses my ctx, so I'll roll my own loop." | True — but document it (`// TODO(boost-upstream):`) and use `apubsub.NewSubscriber` directly, not raw `client.Subscription().Receive()`. The middleware stack must still run. |
| "I'll fix the helper in boost later — for now I'll just not pass ctx." | Open the upstream issue NOW. Future-you and every other consumer will thank you. |
| "Mocking via `os.Setenv` in tests is fine." | It is not. Use `config.Add` with a test-specific default. Env mutation breaks test parallelism. |

---

## 9. When to Escalate

Stop and surface to the human if:

- The task asks you to **add a new top-level layer** (e.g., a peer to `bootstrap/`, `factory/`, `wrapper/`). That's an architectural change, not a contrib.
- You discover a **boost-internal bug** that makes the canonical path impossible (like the ctx-loss in adapter helpers). Document with `// TODO(boost-upstream):`, propose a fix, but don't unilaterally redesign the public API.
- The user wants a config key **outside** the `boost.*` namespace at the application layer. That's fine — but make sure `config.Add` is still the entry point, not `os.Getenv`.
- A consumer project asks for a **direct dependency on a non-boost framework** (gin, fiber, zap-direct, custom DI). Suggest the boost-equivalent first; only deviate with explicit user sign-off.

---

## 10. Self-test before claiming done

Before reporting "done" on any boost-related change, verify:

```bash
# Builds clean
go build ./...

# Vets clean
go vet ./...

# Lints clean (if golangci-lint is configured)
golangci-lint run ./...

# Tests pass (boost has table-driven tests under wrapper/, bootstrap/, factory/)
go test ./...

# For maintainer changes: vendor is consistent
make v 2>/dev/null || go mod vendor
```

For consumer changes, additionally:

- `boost.Start()` is the first line of `main`.
- No `os.Getenv` outside `config.Add` registrations.
- Handler signature: input `cloudevents.Event` by value, output `*cloudevents.Event` by pointer (matches `function.Handler[T]`).
- Generic instantiated as `function.New[*cloudevents.Event](...)` — and every middleware in the chain agrees on `T = *cloudevents.Event`.
- Graceful shutdown is wired via `signal.NotifyContext` (HTTP) or via the documented function workaround (subscribers).
- Errors are returned as `boost/model/errors.*` types (not `fmt.Errorf("%w", ...)`).

If any of those is "no", the code is not boost-shaped yet.

**Verification by example:** before claiming the skill or any boost code is correct, grep `bootstrap/function/handler.go` and confirm the actual `Handler[T]` type signature. Skills can drift from upstream when the framework evolves; the source-of-truth is the framework, not this document.
