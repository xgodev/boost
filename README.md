<p align="center">
  <img src="docs/site/logo.svg" width="68" alt="boost logo: Go gopher by Renée French, CC BY 4.0" title="Gopher by Renée French, CC BY 4.0">
</p>

<h1 align="center">boost</h1>

<p align="center"><b>Enterprise Go services. Built fast.</b></p>

<p align="center">
Boot, config, logs, 58 factories, middleware, health and code generation in one framework. You write the business logic.
</p>

<p align="center">
  <a href="https://github.com/xgodev/boost/releases/latest"><img src="https://img.shields.io/github/v/release/xgodev/boost?label=release" alt="release"></a>
  <a href="https://github.com/xgodev/boost/actions/workflows/go.yml"><img src="https://github.com/xgodev/boost/actions/workflows/go.yml/badge.svg" alt="Go"></a>
  <img src="https://img.shields.io/badge/go-1.18%2B-00ADD8" alt="go 1.18+">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT"></a>
</p>

<p align="center">
  <a href="https://xgodev.github.io/boost/">Site</a> ·
  <a href="#from-annotation-to-binary-in-two-commands">Get started</a> ·
  <a href="https://github.com/xgodev/boost/tree/main/examples">Examples</a>
</p>

```bash
go get github.com/xgodev/boost
```

| 58 | 46 | 107 | 221 |
|:---:|:---:|:---:|:---:|
| factories | libraries | plugins | releases since 2024 |

## A complete service

A complete service: one route, a service, a repository and two datastores. The
generated files and the run are the real output. Full source:
[`docs/site/example/shop`](docs/site/example/shop).

### You write (9 files)

```text
github.com/acme/shop
├── doc.go                                        go:generate
└── internal/
    ├── app/api/rest/
    │   └── handler/stock/reserver/handler.go     @REST
    └── pkg/
        ├── server/rest/doc.go                    @Server
        ├── domain/
        │   ├── service/stock/reserver/
        │   │   ├── local/service.go              @Middleware
        │   │   └── service.go                    interface
        │   └── repository/stock/repository.go    interface
        └── datastore/stock/
            ├── memory/datastore.go               @Default
            └── mongodb/
                ├── config.go                     @Config
                └── datastore.go                  @Inject
```

`internal/app/api/rest/handler/stock/reserver/handler.go`

```go
// Package reserver serves POST /stock/:sku/reserve.
package reserver

import (
	"net/http"

	"github.com/acme/shop/internal/pkg/domain/service/stock/reserver"
	"github.com/labstack/echo/v4"
)

// Handler reserves one unit of an item.
type Handler struct {
	svc reserver.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc reserver.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle reserves one unit of the sku in the path.
// @REST(method=POST, path=/stock/:sku/reserve)
func (h *Handler) Handle(c echo.Context) error {
	if err := h.svc.Reserve(c.Request().Context(), c.Param("sku")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
```

<details>
<summary><code>internal/pkg/server/rest/doc.go</code></summary>

```go
// Package rest holds what the rest server of app api needs besides its
// handlers: errorHandler turns boost errors into HTTP statuses, so a
// Conflict answers 409.
// @Server(type=rest, plugins=[recover, errorHandler])
package rest
```

</details>

<details>
<summary><code>internal/pkg/domain/service/stock/reserver/local/service.go</code></summary>

```go
// Package local reserves stock in process, on the stock repository.
package local

import (
	"context"

	"github.com/acme/shop/internal/pkg/domain/repository/stock"
	"github.com/acme/shop/internal/pkg/domain/service/stock/reserver"
)

type service struct {
	repo stock.Repository
}

// New provides the service, wrapped in the log and prometheus middlewares.
// @Provide(index=0)
// @Middleware(log, prometheus)
func New(repo stock.Repository) reserver.Service {
	return &service{repo: repo}
}

func (s *service) Reserve(ctx context.Context, sku string) error {
	return s.repo.Take(ctx, sku)
}
```

</details>

<details>
<summary><code>internal/pkg/domain/service/stock/reserver/service.go</code></summary>

```go
// Package reserver reserves one unit of an item.
package reserver

import "context"

// Service reserves stock.
type Service interface {
	Reserve(ctx context.Context, sku string) error
}
```

</details>

<details>
<summary><code>internal/pkg/domain/repository/stock/repository.go</code></summary>

```go
// Package stock is the persistence of the units in stock.
package stock

import "context"

// Repository keeps the units in stock of each item.
type Repository interface {
	// Take removes one unit of sku; a sold-out sku is a Conflict.
	Take(ctx context.Context, sku string) error
}
```

</details>

<details>
<summary><code>internal/pkg/datastore/stock/memory/datastore.go</code></summary>

```go
// Package memory keeps the stock in memory: the default, no infra needed.
package memory

import (
	"context"
	"sync"

	"github.com/acme/shop/internal/pkg/domain/repository/stock"
	"github.com/xgodev/boost/model/errors"
)

type datastore struct {
	mu    sync.Mutex
	units map[string]int
}

// New provides the stock wired when no other is configured.
// @Provide(index=0)
// @Default
func New() stock.Repository {
	return &datastore{units: map[string]int{"tshirt": 3}}
}

func (d *datastore) Take(_ context.Context, sku string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.units[sku] == 0 {
		return errors.Conflictf("%s is sold out", sku)
	}
	d.units[sku]--
	return nil
}
```

</details>

<details>
<summary><code>internal/pkg/datastore/stock/mongodb/config.go</code></summary>

```go
package mongodb

// Options configures the mongo stock; inject generates its config keys.
// @Config
type Options struct {
	Collection string `default:"stock" desc:"collection of the units in stock"`
}
```

</details>

<details>
<summary><code>internal/pkg/datastore/stock/mongodb/datastore.go</code></summary>

```go
// Package mongodb keeps the stock in MongoDB, on boost's mongo factory.
package mongodb

import (
	"context"

	"github.com/acme/shop/internal/pkg/domain/repository/stock"
	mongofact "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
	"github.com/xgodev/boost/model/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type datastore struct {
	coll *mongo.Collection
}

// New keeps the stock on the conn inject builds from boost.factory.mongo.*,
// which the app closes when it stops.
// @Provide(index=0)
// @Inject(index=0, plugins=[health, prometheus])
func New(conn *mongofact.Conn, o *Options) stock.Repository {
	return &datastore{coll: conn.Database.Collection(o.Collection)}
}

func (d *datastore) Take(ctx context.Context, sku string) error {
	res, err := d.coll.UpdateOne(ctx,
		bson.M{"_id": bson.M{"$eq": sku}, "units": bson.M{"$gt": 0}},
		bson.M{"$inc": bson.M{"units": -1}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.Conflictf("%s is sold out", sku)
	}
	return nil
}
```

</details>

<details>
<summary><code>doc.go</code></summary>

```go
// Package shop is the example app of the boost site.
package shop

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
```

</details>


### Generated (12 files)

```text
cmd/api/main.go                                   main
gen/
├── builder/
│   ├── api/
│   │   ├── modules.go                            app api
│   │   └── rest_module.go                        @REST
│   └── middleware/…/local/new_middleware.go      @Middleware
├── client/…/mongo-driver/v2/
│   └── default_conn.go                           @Inject
├── implementation/…/stock/
│   └── repository_implementation.go              @Default
├── internal/
│   ├── app/api/rest/
│   │   └── handler/stock/reserver/new_module.go
│   └── pkg/
│       ├── domain/service/…/local/new_module.go
│       └── datastore/stock/
│           ├── memory/new_module.go
│           └── mongodb/
│               ├── new_module.go
│               └── options_module.go             @Config
└── module.go                                     all modules
```

`cmd/api/main.go`

```go
// Code generated by builder; DO NOT EDIT.

// Command api runs app api: the doors builder generates in
// github.com/acme/shop/gen/builder/api.
package main

import (
	app "github.com/acme/shop/gen/builder/api"
	"github.com/xgodev/boost"
	fxfact "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1"
)

func main() {
	boost.Start()

	fxfact.NewApp(app.Modules()...).Run()
}
```

<details>
<summary><code>gen/builder/api/modules.go</code></summary>

```go
// Code generated by builder; DO NOT EDIT.

package api

import (
	"go.uber.org/fx"

	mw_local_54c3d838 "github.com/acme/shop/gen/builder/middleware/internal_/pkg/domain/service/stock/reserver/local"
	multiserverfx "github.com/xgodev/boost/fx/modules/extra/multiserver"
)

// Modules returns the doors of app api and the multiserver that runs its servers:
//
//	fxfact.NewApp(api.Modules()...).Run()
func Modules() []fx.Option {
	return []fx.Option{
		RestModule(),
		mw_local_54c3d838.DecorateNew(),
		multiserverfx.Module(),
	}
}
```

</details>

<details>
<summary><code>gen/builder/api/rest_module.go</code></summary>

```go
// Code generated by builder; DO NOT EDIT.

package api

import (
	"sync"

	"go.uber.org/fx"

	gen_reserver_5a8b66f7 "github.com/acme/shop/gen/internal/app/api/rest/handler/stock/reserver"
	src_reserver_5a8b66f7 "github.com/acme/shop/internal/app/api/rest/handler/stock/reserver"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	echoerrorhandler "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
	echorecover "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/recover"
	echofx "github.com/xgodev/boost/fx/modules/factory/contrib/labstack/echo/v4"
)

var restOnce sync.Once

// RestModule serves the REST routes of server rest on echo.
// It returns an empty option after the first call.
func RestModule() fx.Option {
	options := fx.Options()

	restOnce.Do(func() {
		options = fx.Module("api.rest",
			echofx.Module(),
			fx.Provide(func() []boostecho.Plugin {
				return []boostecho.Plugin{
					echorecover.Register,
					echoerrorhandler.Register,
				}
			}),
			gen_reserver_5a8b66f7.NewModule(),
			fx.Invoke(func(srv *boostecho.Server, h0 *src_reserver_5a8b66f7.Handler) {
				srv.Add("POST", "/stock/:sku/reserve", h0.Handle)
			}),
		)
	})

	return options
}
```

</details>

<details>
<summary><code>gen/builder/middleware/internal_/pkg/domain/service/stock/reserver/local/new_middleware.go</code></summary>

```go
// Code generated by builder; DO NOT EDIT.

package local

import (
	"go.uber.org/fx"

	context "context"
	src_reserver_399f2048 "github.com/acme/shop/internal/pkg/domain/service/stock/reserver"
	middleware "github.com/xgodev/boost/extra/middleware"
	fmiddleware "github.com/xgodev/boost/factory/local/extra/middleware"
	prometheusmw "github.com/xgodev/boost/factory/local/extra/middleware/plugins/contrib/prometheus/client_golang/v1"
	logmw "github.com/xgodev/boost/factory/local/extra/middleware/plugins/local/wrapper/log"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
)

// DecorateNew decorates the reserver.Service that local.New provides:
// every consumer in the app receives a newDecorator, which runs the middlewares
// log, prometheus around each method; config can disable each one.
func DecorateNew() fx.Option {
	return fx.Options(
		contextfx.Module(),
		fx.Decorate(decorateNew),
	)
}

// decorateNew builds the wrappers of every method of reserver.Service once.
func decorateNew(ctx context.Context, next src_reserver_399f2048.Service) (src_reserver_399f2048.Service, error) {
	d := &newDecorator{next: next}
	var err error

	d.reserveWrapper, err = fmiddleware.NewErrorWrapperWithOptions(ctx, &fmiddleware.Options{Name: "reserver.service.reserve"},
		logmw.NewError,
		prometheusmw.NewError,
	)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// newDecorator implements reserver.Service: each method runs next's in its wrapper.
type newDecorator struct {
	next           src_reserver_399f2048.Service
	reserveWrapper *middleware.ErrorWrapper
}

// Reserve runs next.Reserve in the middlewares of reserver.service.reserve.
func (d *newDecorator) Reserve(ctx context.Context, a1 string) error {
	return d.reserveWrapper.Exec(ctx, "reserver.service.reserve",
		func(ctx context.Context) error {
			return d.next.Reserve(ctx, a1)
		},
		func(_ context.Context, err error) error {
			return err
		},
	)
}
```

</details>

<details>
<summary><code>gen/client/github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2/default_conn.go</code></summary>

```go
// Code generated by inject; DO NOT EDIT.

package mongo

import (
	"context"
	"sync"

	"go.uber.org/fx"

	factory "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
	prometheusplugin "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2/plugins/contrib/prometheus/client_golang/v1"
	healthplugin "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2/plugins/local/health"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
)

// defaultConnRoot is the config root of the default client (*github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2.Conn);
// env BOOST_FACTORY_MONGO_*.
const defaultConnRoot = "boost.factory.mongo"

func init() {
	factory.ConfigAdd(defaultConnRoot)
	healthplugin.ConfigAdd(defaultConnRoot + ".plugins.health")
	prometheusplugin.ConfigAdd(defaultConnRoot + ".plugins.prometheus")
}

var defaultConnOnce sync.Once

// DefaultConnModule provides the default client (*github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2.Conn), built from the keys under
// boost.factory.mongo with the plugins health, prometheus, and closes it when the app stops.
// It returns an empty option after the first call, so every consumer can
// include it.
func DefaultConnModule() fx.Option {
	options := fx.Options()

	defaultConnOnce.Do(func() {
		options = fx.Module("github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2.DefaultConn",
			contextfx.Module(),
			fx.Provide(newDefaultConn),
		)
	})

	return options
}

func newDefaultConn(ctx context.Context, lc fx.Lifecycle) (*factory.Conn, error) {
	health, err := healthplugin.NewHealthWithConfigPath(defaultConnRoot + ".plugins.health")
	if err != nil {
		return nil, err
	}
	prometheus, err := prometheusplugin.NewPrometheusWithConfigPath(defaultConnRoot + ".plugins.prometheus")
	if err != nil {
		return nil, err
	}
	c, err := factory.NewConnWithConfigPath(ctx, defaultConnRoot, health.Register, prometheus.Register)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.StopHook(c.Close))
	return c, nil
}
```

</details>

<details>
<summary><code>gen/implementation/internal_/pkg/domain/repository/stock/repository_implementation.go</code></summary>

```go
// Code generated by inject; DO NOT EDIT.

package stock

import (
	"sync"

	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
	"go.uber.org/fx"

	gen_memory_6a1a5ba1 "github.com/acme/shop/gen/internal/pkg/datastore/stock/memory"
	gen_mongodb_1bacb721 "github.com/acme/shop/gen/internal/pkg/datastore/stock/mongodb"
)

// repositoryImplementationKey chooses the implementation of github.com/acme/shop/internal/pkg/domain/repository/stock.Repository.
const repositoryImplementationKey = "internal.pkg.domain.repository.stock.implementation"

func init() {
	config.Add(repositoryImplementationKey, "memory", "implementation of github.com/acme/shop/internal/pkg/domain/repository/stock.Repository: memory, mongodb")
}

// RepositoryImplementation returns the name of the implementation of
// github.com/acme/shop/internal/pkg/domain/repository/stock.Repository that RepositoryImplementationModule wires.
func RepositoryImplementation() string {
	return config.String(repositoryImplementationKey)
}

var repositoryImplementationOnce sync.Once

// RepositoryImplementationModule wires the implementation of github.com/acme/shop/internal/pkg/domain/repository/stock.Repository
// named by internal.pkg.domain.repository.stock.implementation (memory, mongodb; default memory) and only that one.
// An unknown name fails the boot. It returns an empty option after the first
// call, so every consumer can include it.
func RepositoryImplementationModule() fx.Option {
	options := fx.Options()

	repositoryImplementationOnce.Do(func() {
		switch name := RepositoryImplementation(); name {
		case "memory":
			options = gen_memory_6a1a5ba1.NewModule()
		case "mongodb":
			options = gen_mongodb_1bacb721.NewModule()
		default:
			options = fx.Error(errors.NotValidf("%s: unknown implementation %q; valid: %s", repositoryImplementationKey, name, "memory, mongodb"))
		}
	})

	return options
}
```

</details>

<details>
<summary><code>gen/internal/pkg/datastore/stock/mongodb/options_module.go</code></summary>

```go
// Code generated by inject; DO NOT EDIT.

package mongodb

import (
	"sync"

	"github.com/xgodev/boost/wrapper/config"
	"go.uber.org/fx"

	src_mongodb_1bacb721 "github.com/acme/shop/internal/pkg/datastore/stock/mongodb"
)

// optionsRoot is the config key holding github.com/acme/shop/internal/pkg/datastore/stock/mongodb.Options.
const optionsRoot = "internal.pkg.datastore.stock.mongodb"

func init() {
	// env INTERNAL_PKG_DATASTORE_STOCK_MONGODB_COLLECTION
	config.Add(optionsRoot+".collection", "stock", "collection of the units in stock")
}

var optionsOnce sync.Once

// OptionsModule provides the Options of github.com/acme/shop/internal/pkg/datastore/stock/mongodb,
// filled from the config keys under internal.pkg.datastore.stock.mongodb.
// It returns an empty option after the first call, so every consumer can
// include it.
func OptionsModule() fx.Option {
	options := fx.Options()

	optionsOnce.Do(func() {
		options = fx.Module("github.com/acme/shop/internal/pkg/datastore/stock/mongodb.Options",
			fx.Provide(func() (*src_mongodb_1bacb721.Options, error) {
				return config.NewOptionsWithPath[src_mongodb_1bacb721.Options](optionsRoot)
			}),
		)
	})

	return options
}
```

</details>


The other generated files are in
[`docs/site/example/shop/gen`](docs/site/example/shop/gen).

### Run

```console
$ go generate ./...
5:35PM INF scanning ~/shop (module github.com/acme/shop)
5:35PM INF generated 5 inject modules and 1 apps in ~/shop/gen, mains in ~/shop/cmd
$ go run ./cmd/api &
5:36PM INF starting Echo Server. https://echo.labstack.com/
# memory is the @Default datastore: 3 t-shirts in stock
$ for i in 1 2 3 4; do curl -s -w ' %{http_code}\n' -X POST localhost:8080/stock/tshirt/reserve; done
 204
 204
 204
{"httpStatusCode":409,"errorCode":"409","message":"tshirt is sold out"}
 409
# same binary on MongoDB, through boost's mongo factory:
# INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION=mongodb
```

## One boot. Everything ready after it.

`boost.Start()` runs first. When it returns, config is loaded and the logger is
set, so every factory reads its own keys and logs the same way.

| | | |
|---|---|---|
| 01 | `boost.Start()` | Reads flags, env vars and config files. |
| 02 | `config` | Every key declared once, with a default and a description. |
| 03 | `log` | zap, zerolog or logrus behind one wrapper. |
| 04 | `factory` | One constructor per component, plugins applied. |
| 05 | `fx` | Generated modules wire the dependency graph. |
| 06 | `serve` | REST, gRPC, MCP or an event function. |

### One key, three ways to set it

Keys live under `boost.factory.<component>.*` and reach the app from wherever
you deploy.

| | |
|---|---|
| key | `boost.factory.zap.console.level` |
| env | `BOOST_FACTORY_ZAP_CONSOLE_LEVEL=DEBUG` |
| flag | `--boost.factory.zap.console.level=DEBUG` |

```go
// cmd/api/main.go, written by go generate
func main() {
	boost.Start()
	fxfact.NewApp(app.Modules()...).Run()
}
```

[`bootstrap`](bootstrap/README.md)

## Every piece has one job.

Factories build, wrappers abstract, generators remove the wiring. Use one
package or all of them.

### Factories

`factory/contrib/labstack/echo/v4` · [docs](factory/README.md)

One constructor per component. You never instantiate the SDK by hand: the
factory reads its config keys, applies the plugins you list and returns the
ready client.

```go
srv, _ := echo.NewServer(ctx,
	cors.Register,
	requestid.Register,
	gzip.Register,
	log.Register,
	health.Register,
)
srv.GET("/stock/:sku", Get)
srv.Serve(ctx)
```

**58** factories · **46** upstream libraries · **107** plugins

### Plugins

`…/echo/v4/plugins`

Behaviour you opt into per factory, applied in the order you list it.

`cors` `requestid` `gzip` `log` `health` `recover` `bodylimit` `prometheus`
`otel` `datadog` `swaggo` `semaphore`

### Wrappers

`wrapper/{log,config,cache,publisher}` · [docs](wrapper/README.md)

Swap zap for zerolog by changing one line in main. App code only ever calls the
wrapper.

```go
log.FromContext(ctx).
	Infof("reserved %s", sku)
```

### Middleware

`extra/middleware` · [docs](extra/README.md)

Wraps any call in log, cache, hystrix, datadog, prometheus or fallback. With
`@Middleware` the decorator is generated and the method body stays clean.

```text
log
└── prometheus
    └── hystrix
        └── Service.Reserve()
```

`@Middleware(log, prometheus, hystrix)` · first is outermost

### Inject

`inject` · [docs](inject/README.md)

`@Provide`, `@Inject` and `@Invoke` become one fx module per function and a root
`Module()`.

```text
@Default     memory.New | mongodb.New   stock.Repository
    ↓
@Middleware  local.New(repo)            reserver.Service
    ↓
@Provide     reserver.New(svc)          *Handler
    ↓
@REST        Handler.Handle             POST /stock/:sku/reserve
```

- cycles and missing providers fail at `go generate`, not in production
- two providers of the same value are rejected

### Annotation

`annotation` · [docs](annotation/README.md)

Reads `// @Name(key=value)` from doc comments into typed structs. Every
generator here is built on it.

```go
var r struct {
	Path string `attr:"path"`
}
ann.Decode(&r)
```

### Builder

`builder` · [docs](builder/README.md)

An annotation says which door serves the handler. The builder writes the door,
the app modules and main.

| | |
|---|---|
| `@REST` | echo route |
| `@GRPC` | gRPC service |
| `@Function` | azure functions · cloudevents · kafka · kafka-go · lambda · nats · pubsub |
| `@MCP` | MCP tool |
| `@Job` | runs once and exits · CronJob |
| `@Middleware` | generated decorator |
| `@Default` | implementation picked by config |
| `@Server` | plugins of a server |

### Functions

`bootstrap/function` · [docs](bootstrap/README.md)

One event handler, six transports. Change the adapter in config; the handler
stays the same.

```go
func(ctx, in cloudevents.Event)
	(*cloudevents.Event, error)
```

### MCP

`factory/contrib/modelcontextprotocol/go-sdk/v1/{server,client}`

Your use cases become tools that agents call over
[MCP](https://github.com/modelcontextprotocol/go-sdk), and your service calls
other MCP servers. Server and client are factories like the rest: config keys,
plugins, and the server runs beside [echo](https://echo.labstack.com/) and
[gRPC](https://grpc.io/docs/languages/go/). Mark a handler with `@MCP` and the
builder registers it as a tool.

`server · HTTP` `server · STDIO` `client · HTTP` `client · COMMAND` `log` `otel`

```go
srv, _ := server.NewServer(ctx, mcplog.Register, otel.Register)
mcp.AddTool(srv.Instance(), &mcp.Tool{
	Name:        "reserve",
	Description: "holds items for an order",
}, reserve)
srv.Serve(ctx)
```

## From annotation to binary in two commands.

You write handlers, services and annotations. Everything that repeats between
services is generated, checked and rewritten on every run.

1. **Annotate.** Mark constructors with `@Provide`, handlers with a door,
   services with `@Middleware`, the default implementation with `@Default`,
   options with `@Config`, the plugins of a server with `@Server`.
2. **Generate.** Inject and builder run in one pass. A broken graph stops here
   and leaves your output untouched.

   ```bash
   go generate ./...
   ```

3. **Run.** Each app gets its own `main`, with the servers and plugins it
   declared.

   ```bash
   go run ./cmd/api
   ```

Generated output, rewritten on every run:

```text
cmd/api/main.go                // boost.Start(), then app api
gen/
├── builder/
│   ├── api/
│   │   ├── modules.go         // the doors of app api
│   │   └── rest_module.go     // @REST POST /stock/:sku/reserve
│   └── middleware/…/local/
│       └── new_middleware.go  // @Middleware(log, prometheus)
├── implementation/…/stock/
│   └── repository_implementation.go
│                              // @Default: memory | mongodb
├── internal/
│   ├── app/…                  // one module per @Provide,
│   └── pkg/…                  // @Config: same tree as yours
└── module.go                  // every inject module
```

## Upgrade boost. Keep your code.

A minor release never asks you to rewrite a service. That is a rule of the
project, enforced on every pull request.

### Public API is a contract

A component's public API does not break without a major version bump and a
written reason. Any observable change ships with a test first.

### SDK majors live side by side

Factory paths carry the upstream major. Moving to a new SDK is an import you
choose, not a forced migration.

```text
factory/contrib/go.mongodb.org/mongo-driver/v1   // keeps working
factory/contrib/go.mongodb.org/mongo-driver/v2   // new SDK major
```

### A gate that only moves forward

Every PR is compared to `develop`: format, lint, build and tests can stay the
same or get better, never worse.

```diff
 module github.com/acme/shop

 require (
-	github.com/xgodev/boost v1.1.1
+	github.com/xgodev/boost v1.6.2
 )
```

```console
$ go build ./... && go test ./...
ok  github.com/acme/shop/...
```

1 line changed in go.mod, **0 lines changed in your code**.

Releases: 2024-07 → `v1.1.0` → `v1.1.1` → **`v1.6.2`**
([all releases](https://github.com/xgodev/boost/releases)).

## There's a factory for that.

58 factories over the 46 libraries your services already run on. Each one
arrives configured, with plugins for traces, metrics and health checks, so your
team writes the business logic, not the plumbing.

### Web, RPC & MCP

| Library | Factory wraps |
|---|---|
| [Echo](https://echo.labstack.com/) | `echo/v4 · v5` |
| [chi](https://go-chi.io/) | `chi/v5` |
| [Fiber](https://gofiber.io/) | `fiber/v2` |
| [gRPC](https://grpc.io/docs/languages/go/) server · client | `grpc/v1` |
| [MCP](https://github.com/modelcontextprotocol/go-sdk) server · client | `go-sdk/v1` |
| HTTP/2 server | `x/net/v0` |
| GraphQL | `graphql-go/v0` |
| Resty | `resty/v2` |
| [fasthttp](https://github.com/valyala/fasthttp) client | `fasthttp/v1` |
| [drpc](https://github.com/storj/drpc) server · client | `drpc/v1` |

### Messaging

| Library | Factory wraps |
|---|---|
| Kafka (Confluent) | `confluent-kafka-go/v2` |
| Kafka (segmentio) | `kafka-go/v0` |
| NATS | `nats.go/v1` |
| Google Pub/Sub | `pubsub/v1 · v2` |
| Go CDK Pub/Sub | `gocloud.dev/v0` |
| CloudEvents | `sdk-go/v2` |
| Goka | `goka/v1` |
| SQS · SNS · Kinesis | `aws-sdk-go-v2` |

### Data

| Library | Factory wraps |
|---|---|
| [MongoDB](https://www.mongodb.com/docs/drivers/go/current/) | `mongo-driver/v1 · v2` |
| PostgreSQL | `pgx/v5` |
| Oracle | `godror/v0` |
| Cassandra | `gocql/v1 · v2` |
| Elasticsearch | `go-elasticsearch/v8` · `v9` |
| BigQuery · Bigtable · Firestore · API client | `cloud.google.com` |
| FTP | `ftp/v0` |
| [Azure Blob Storage](https://github.com/Azure/azure-sdk-for-go) | `azure-sdk-for-go/v1` |

### Cache

| Library | Factory wraps |
|---|---|
| Redis | `go-redis/v9` |
| BigCache | `bigcache/v3` |
| FreeCache | `freecache/v1` |
| BuntDB | `buntdb/v1` |
| memdb | `memdb/v1` |

### Observability

| Library | Factory wraps |
|---|---|
| OpenTelemetry | `otel/v1` |
| Datadog | `dd-trace-go/v1` |
| Prometheus | `client_golang/v1` |
| Hystrix | `hystrix-go/v0` |

### Logging & platform

| Library | Factory wraps |
|---|---|
| zap · zerolog · logrus | `wrapper/log` |
| fx | `fx/v1` |
| Cobra | `cobra/v1` |
| Kubernetes | `client-go/v0` |
| Vault | `vaultgo/v0` |
| ants worker pool | `ants/v2` |

### Why a factory, not just the SDK

Any SDK connects in five lines. Running it in a company takes the rest: config
per environment, auth, checks, logs, metrics, traces. That part is the same in
every service, so boost writes it once.

- **Configured, not coded.** Every option is a config key with a default, read
  from file or environment: `boost.factory.mongo.uri` is
  `BOOST_FACTORY_MONGO_URI`. One binary for every environment.
- **The plumbing comes built.** Defaults, auth, connect and ping, structured
  logs through `wrapper/log`. What a platform team rewrites for every SDK in
  every service is already there.
- **78 plugins to plug in.** OpenTelemetry, Prometheus, Datadog, health checks,
  retry, request id, CORS, gzip. Pass the plugin to the factory; switch it off
  with a key.

**By hand: MongoDB, in every service**

- Read URI, user and password from the environment, with defaults
- Build the client options and the authentication
- Connect and ping before taking traffic
- Trace every command with OpenTelemetry
- Count commands and errors in Prometheus
- Register a health check
- Log commands when you need to debug
- Then again for Kafka, Redis, gRPC, the HTTP client…

**With boost:** `factory/contrib/go.mongodb.org/mongo-driver/v2`

```go
conn, err := mongo.NewConn(ctx,
	mongootel.Register,   // traces
	mongoprom.Register,   // metrics
	mongohealth.Register, // health check
)
```

```console
$ export BOOST_FACTORY_MONGO_URI=mongodb://db:27017/shop
$ export BOOST_FACTORY_MONGO_AUTH_USERNAME=shop
$ export BOOST_FACTORY_MONGO_PLUGINS_PROMETHEUS_ENABLED=false
```

Same shape for every factory: `New…(ctx, plugins...)` and keys under
`boost.factory.<name>`.

## Coding agents that follow the rules.

The [`claude-plugin@xgodev`](https://github.com/xgodev/claude-plugin) plugin
teaches Claude Code boost's Iron Laws and the project layout, and runs the
comparative quality gate before every push.

```
/plugin marketplace add xgodev/claude-plugin
/plugin install claude-plugin@xgodev
```

1. **`boost.Start()` first.** Config and log exist only after it.
2. **Log only through `wrapper/log`.** `log.FromContext(ctx)`, never a
   third-party logger in app code.
3. **Config only through `wrapper/config`.** No `os.Getenv`; keys under
   `boost.factory.<x>.*`.
4. **One factory per component.** Never instantiate the upstream SDK directly.
5. **Backward compatible.** Public API breaks only with a major bump and a
   reason.

## Start with one line.

```bash
go get github.com/xgodev/boost
```

[See the examples](https://github.com/xgodev/boost/tree/main/examples) ·
[Star on GitHub](https://github.com/xgodev/boost)

---

boost · xgodev · [MIT](LICENSE)
