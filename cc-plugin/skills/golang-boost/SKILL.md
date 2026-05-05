---
name: golang-boost
description: "Use when writing or reviewing Go code that imports github.com/xgodev/boost or any of its subpackages, including HTTP APIs (Echo factory), event-driven functions (CloudEvents over Pub/Sub, NATS, Kafka), publisher/cache/log/config wrappers, fx modules, or new contribs to the framework itself. Triggers on questions about boost.Start, log.FromContext, function.New, fn.Run, echo.NewServer, publisher.New, config.Add, model/errors, or layout decisions for new factory contribs / wrapper drivers / bootstrap adapters."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**Persona:** You are a senior Go engineer who has shipped APIs and Cloud Functions on top of `github.com/xgodev/boost` in production. You recognize boost-shaped code at a glance, know which paths are canonical, and can tell a real upstream limitation apart from a rationalized shortcut.

---

## 1. The Three Iron Laws

These are non-negotiable. Code that violates them is not boost-shaped, even when it compiles and tests pass.

### Law 1 — `boost.Start()` is the first line of `main`

```go
func main() {
    boost.Start() // ALWAYS first. Loads env+files into config, installs the global logger.
    // ... everything else ...
}
```

Without `Start`, `log.FromContext(ctx)` returns a no-op logger, the config registry is empty, and factory constructors fail or produce zero-valued junk.

### Law 2 — Function handlers take input by value, output by pointer

The framework's `bootstrap/function/handler.go` declares:

```go
type Handler[T any] func(context.Context, cloudevents.Event) (T, error)
```

Input is **always** `cloudevents.Event` by value (forced by the type). Instantiate `T = *cloudevents.Event` so the publisher and logger middlewares — which type-switch on `*event.Event` — fire correctly.

```go
// CORRECT
func handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error)

// WRONG — return-by-value silently disables publisher middleware
func handle(ctx context.Context, in cloudevents.Event) (cloudevents.Event, error)

// WRONG — does not compile against function.Handler[*cloudevents.Event]
func handle(ctx context.Context, in *cloudevents.Event) (*cloudevents.Event, error)
```

Wire the entire chain on `T = *cloudevents.Event`: `function.New[*cloudevents.Event](...)`, `lm.NewAnyErrorMiddleware[*cloudevents.Event]()`, `apubsub.New[*cloudevents.Event](pb)`.

### Law 3 — All configuration goes through `config.Add` + typed accessors

```go
// WRONG
outSubject := os.Getenv("FN_OUTPUT_SUBJECT")

// CORRECT (in init or main, before consumption)
config.Add("myapp.output.subject", "default-topic", "downstream subject")
// At deploy time the operator overrides via env: MYAPP_OUTPUT_SUBJECT=foo
outSubject := config.String("myapp.output.subject")
```

The "just one operator override" rationalization breaks two invariants:

1. **Discoverability** — the boot banner / `boost-config dump` enumerates every `config.Add` call. `os.Getenv` is invisible.
2. **Test reproducibility** — `wrapper/config` has koanf providers for fixtures. `os.Getenv` requires mutating the process environment, breaking `t.Parallel()`.

Application-layer keys can use any namespace; framework-layer keys are `boost.*` (e.g., `boost.factory.echo.port`).

> **Spirit:** boost is the shape of the application, not a library you sprinkle in. Fighting it (bypassing `boost.Start`, raw `os.Getenv`, instantiating `pubsub.Subscription` directly) means you're either leaking a real upstream limitation that deserves an issue/PR, or you're rationalizing a shortcut.

---

## 2. Logging — `wrapper/log`

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
| `log.Fatal/Fatalf` | Only for process-fatal init failures, never inside a handler |

**Never:** `log.New`, `zap.NewProduction`, `zerolog.New`, `slog.New`. The wrapper installs the configured backend (zap, zerolog, logrus) at `boost.Start`.

---

## 3. Errors — `model/errors`

```go
import bootsterrors "github.com/xgodev/boost/model/errors"

return bootsterrors.NewBadRequest(err, "invalid payload")
return bootsterrors.BadRequestf("field %q is required", "id")
return bootsterrors.NewNotFound(err, "order not found")
return bootsterrors.NewConflict(err, "duplicate order id")
return bootsterrors.NewInternal(err, "downstream call failed")
return bootsterrors.NotValidf("invalid event data")              // for function deadletter
return bootsterrors.Wrap(err, bootsterrors.NotValidf("..."))     // wrap + classify
```

Two boost subsystems pattern-match on the error type name:

1. **Echo `error_handler` plugin** — maps `*errors.NotFound` → 404, `*errors.BadRequest` → 400, `*errors.Internal` → 500, etc., emitting a JSON envelope (when `Type=REST` is set on the server, which the `restresponse` plugin does).
2. **Function `publisher` middleware (deadletter mode)** — routes errors of type `NotValid` / `Internal` / etc. to configured deadletter topics.

**Never:**
- `fmt.Errorf("%w", err)` for errors that flow through these matchers — use `bootsterrors.Wrap`.
- `echo.NewHTTPError(404, "...")` — bypasses `error_handler`, produces inconsistent payloads.

---

## 4. HTTP API — Echo Factory

### Canonical `main.go`

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
        error_handler.Register, // model/errors → HTTP status mapping
    )
    if err != nil {
        logger.WithError(err).Fatal("failed to build echo server")
    }

    h := newOrderHandler(/* deps */)
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

### Required plugin set

| Plugin | Purpose | Skip when |
|---|---|---|
| `native/recover` | Panic → 500 | Never |
| `native/requestid` | X-Request-ID | Never |
| `local/wrapper/log` | Access log via boost logger | Never |
| `local/model/restresponse` | Sets `Type=REST` for `error_handler` JSON mode | Never (REST APIs) |
| `extra/error_handler` | `model/errors` → status code mapping | Never (without it, errors leak as plain strings) |
| `native/cors` | Browser clients | Internal-only services |
| `native/gzip` | Response compression | gRPC / streaming |

`local/model/restresponse` MUST be registered before `extra/error_handler` so `Type=REST` is set when the handler picks `ErrorHandlerJSON` over `ErrorHandlerString`.

---

## 5. Function — Pub/Sub Subscriber

### Prototype shape (canonical)

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

// Input by value (framework signature), output by pointer (Law 2).
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

The canonical wiring **does not propagate signal cancellation into the receiver loop**. Inside `bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v1/helper.go:51` (and the analogous NATS / Kafka helpers), the call is:

```go
if err := subscriber.Subscribe(context.Background()); err != nil { ... }
```

A `signal.NotifyContext` passed to `fn.Run(ctx, ...)` reaches the middleware wrapper but **not** the subscription loop. SIGTERM does not gracefully drain in-flight messages.

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

**Required action when using the workaround:** add a `// TODO(boost-upstream):` comment naming `helper.go:51` and stating the collapse-back condition. The workaround compensates for an upstream limitation, not a stylistic choice — when the helper is fixed, this code should collapse back to `fn.Run`.

### Middleware order

```
recovery → logger → publisher
```

- `recovery` outermost — a panic must not kill the worker before the logger sees it.
- `logger` middle — sees both raw and post-recovery errors.
- `publisher` innermost — fires on the handler's successful result; with deadletter config, also routes errors.

### Wrapping errors for deadletter routing

The publisher middleware matches on the unwrapped error type name:

```go
return nil, bootsterrors.Wrap(err, bootsterrors.NotValidf("invalid event data"))
return nil, bootsterrors.Wrap(err, bootsterrors.Internalf("downstream call failed"))
```

---

## 6. Maintainer — Adding a Driver / Adapter / Plugin to boost

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

**Package name = the short library name** (`pubsub`, `nats`, `confluent`, `goka`, `sns`, `redis`), **not** the version directory leaf. Alias the upstream SDK at the import site if a name clash would occur (e.g., `asns "github.com/aws/aws-sdk-go-v2/service/sns"`).

**Multi-service SDKs (AWS SDK v2, Azure SDK, GCP SDK):** for `wrapper/`, `bootstrap/`, and `extra/` — split per service: `<vendor>/<service>/v<major>/` (e.g., `aws/sns/v1/`, `aws/sqs/v1/`, `aws/dynamodb/v1/`). Do **not** nest under the SDK module dir. The umbrella SDK layout (`factory/contrib/aws/aws-sdk-go-v2/v1/client/<service>/`) is **exclusive to `factory/contrib/`** — factories ship clients that legitimately share an SDK version pin; drivers don't.

### Constructor trio

Every driver/adapter exposes the same trio so call sites are interchangeable:

```go
func New(ctx context.Context, c *upstream.Client) (publisher.Driver, error)
func NewWithOptions(ctx context.Context, c *upstream.Client, opts *Options) publisher.Driver
func NewWithConfigPath(ctx context.Context, c *upstream.Client, path string) (publisher.Driver, error)
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
}
```

### Errors and logging

Match the verbs in the closest existing driver. Boost is consistent enough that diverging stands out in review.

```go
import (
    "github.com/xgodev/boost/model/errors"
    "github.com/xgodev/boost/wrapper/log"
)

logger := log.FromContext(ctx).WithTypeOf(*p)
return nil, errors.Wrap(err, errors.Internalf("publish failed"))
```

### Honest extrapolation

When the new driver has no direct precedent (e.g., SNS's `TopicArn` does not map cleanly to Pub/Sub's topic name), mark every guess:

```go
// TODO(maintainer-review): SNS requires a full ARN. Falling back to options.TopicArn
// when ev.Subject() is not an ARN. Verify this is the wanted behavior, or require
// Subject to always be an ARN.
```

Then call out the marked decisions in the PR description.

### What NOT to add

- ❌ Code under `bootstrap/` for a wrapper-layer concern, or under `wrapper/` for an HTTP framework concern. Layering is enforced by directory.
- ❌ A new top-level interface alongside an existing one (`Driver2`). Extend the existing one or open an RFC issue first.
- ❌ A direct dependency on a third-party DI / config / log library. Use `wrapper/config`, `wrapper/log`, and `fx` modules.
- ❌ Tests that `os.Setenv` to inject config — use `config.Add` with explicit defaults.

---

## 7. Quick Reference — Imports

```go
// Boot
"github.com/xgodev/boost"                                 // boost.Start()
"github.com/xgodev/boost/wrapper/log"                     // log.FromContext(ctx)
"github.com/xgodev/boost/wrapper/config"                  // config.Add / config.String
"github.com/xgodev/boost/model/errors"                    // BadRequestf / NewNotFound / Wrap

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

## 8. Red Flags — STOP if you see these

| Red flag | Fix |
|---|---|
| `os.Getenv(...)` outside a `config.Add` registration | Replace with `config.String/Int/Bool/Duration` |
| `log.New`, `zap.NewProduction`, `zerolog.New`, `slog.New` | Replace with `log.FromContext(ctx)` |
| `echo.New()` directly | Replace with `factory/.../echo/v4.NewServer(ctx, ...)` |
| `pubsub.Client.Subscription(...).Receive(...)` | Use `apubsub.NewSubscriber(...).Subscribe(ctx)` or `function.New + fn.Run` |
| `echo.NewHTTPError(404, "...")` | Use `bootsterrors.NewNotFound(err, "...")` + `error_handler` plugin |
| Function handler returning `cloudevents.Event` (value) | Change return to `*cloudevents.Event` |
| Bypass of `function.New + fn.Run` without `// TODO(boost-upstream):` | Add the comment + tracking issue, OR use the canonical path |
| Missing `boost.Start()` in `main` | Add it as the first statement |
| `fmt.Errorf("%w", err)` instead of `bootsterrors.Wrap` | Use `bootsterrors.Wrap` so error type matchers (deadletter, error_handler) work |

---

## 9. Self-test before claiming done

```bash
go build ./...
go vet ./...
golangci-lint run ./...   # if configured
go test ./...
```

Iron Laws checklist:

- [ ] `boost.Start()` is the first line of `main`.
- [ ] No `os.Getenv` outside `config.Add` registrations.
- [ ] No stdlib / third-party logger constructed outside `wrapper/log`.
- [ ] Function handler signature: input `cloudevents.Event` by value, output `*cloudevents.Event` by pointer.
- [ ] Generic instantiated as `function.New[*cloudevents.Event](...)`; whole middleware chain agrees on `T = *cloudevents.Event`.
- [ ] Errors returned as `boost/model/errors.*` types when they flow through Echo or function middleware.
- [ ] Graceful shutdown wired via `signal.NotifyContext` (HTTP) or via the documented function workaround (subscribers).

If any answer is "no", the code is not boost-shaped yet.

**Verification by example:** before claiming the skill or any boost code is correct, grep `bootstrap/function/handler.go` and confirm the actual `Handler[T]` type signature. Skills can drift from upstream when the framework evolves; the source-of-truth is the framework, not this document.
