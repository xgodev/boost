# boost website

Issues: [#74](https://github.com/xgodev/boost/issues/74),
[#82](https://github.com/xgodev/boost/issues/82) (hero example),
[#88](https://github.com/xgodev/boost/issues/88) (publishing). Page:
[`index.html`](index.html), with its styles in [`site.css`](site.css), its script
in [`site.js`](site.js). Open the HTML in a browser; only Google Fonts load
from the network.

The root [`README.md`](../../README.md) mirrors this page: same sections, copy,
numbers, example and version. A change here changes it in the same PR.

## Publishing

Live at <https://xgodev.github.io/boost/>, published from `develop`.

- [`.github/workflows/site.yml`](../../.github/workflows/site.yml) deploys to
  GitHub Pages on every push to `develop` that touches `docs/site/**`. It can
  also be run by hand (`workflow_dispatch`).
- Everything in `docs/site` ships except `example/` and this README.
- Pages source is "GitHub Actions" (repo settings), not a branch. The
  `github-pages` environment allows deploys from `develop`.
- The repo's homepage (About) points to the site.

## Logo

Go gopher on a jetpack ([#81](https://github.com/xgodev/boost/issues/81)). The
gopher was designed by Renée French and is licensed
[CC BY 4.0](https://creativecommons.org/licenses/by/4.0/): the logo's tooltip
carries the credit, and so must every other place the logo ships.

| File | Use |
|---|---|
| `logo.svg` | the mark in the nav, next to the `boost` wordmark (Bricolage Grotesque 800, loaded only for those five letters) |
| `favicon.svg` | browser tab: head only, so it still reads at 16 px |
| `favicon-32.png` | tab icon for browsers without SVG favicons |
| `apple-touch-icon.png` | 180 px home-screen icon, full mark on navy |

## Library logos

`logos/` holds the 23 marks of the `#integrations` wall (#101). Each tile draws
its SVG as a CSS mask in `--fg-2` (`--fg` on hover), so every mark reads as one
monochrome set on the dark page whatever its original colours. A tile is a
whole-card link to the library's page, the same URL as in the full list.

| Source | Logos |
|---|---|
| [simple-icons](https://simpleicons.org/) 13.21.0 (CC0) | GraphQL, Kafka, NATS, Pub/Sub, AWS, MongoDB, PostgreSQL, Oracle, Cassandra, Elasticsearch, BigQuery, Firestore, Redis, OpenTelemetry, Prometheus, Datadog, Kubernetes, Vault, Google Cloud |
| [cncf/artwork](https://github.com/cncf/artwork) | gRPC, CloudEvents |
| [modelcontextprotocol](https://github.com/modelcontextprotocol/modelcontextprotocol) `docs/logo/light.svg` | MCP (the icon, without the wordmark) |
| `echo.labstack.com/favicon.svg` | Echo (background square removed) |

The marks are their owners' trademarks, shown only to name the library each
factory wraps. gRPC, MCP, Cassandra and Echo carry a cropped `viewBox`; gRPC
and Cassandra are wide, so their tiles set `--w`.

## Positioning

boost is a set of things for building enterprise Go services fast: boot,
config, logs, factories, middleware, health and code generation. The headline
and the lede speak about the whole set. Never present boost as only a code
generator or only a DI tool: "Write the handler. Boost writes the service." was
dropped for that reason.

## Goal

One page that sells what boost does well, using real code and real numbers:

- one boot (`boost.Start()`), then config and log are ready
- one factory per component, with plugins
- middleware around a call (log, prometheus, hystrix, …)
- `annotation`, `inject` and `builder`: annotate the handler, `go generate`
  writes the fx modules, the doors and `main`
- function adapters (cloudevents, kafka, nats, pubsub)
- MCP: use cases as tools agents call, and a client for other MCP servers
- backward compatibility: upgrade boost, keep your code

## Design system

Generated with the `ux-ui` skill from `claude-plugin@xgodev`
(`--design-system` for a Go framework developer site).

| Token | Value |
|---|---|
| Style | Dark Mode (OLED), dark only (`color-scheme: dark`) |
| Background | `#0F172A` (deep band `#0A1020`) |
| Surface | `#1B2336` |
| Accent ("run green") | `#22C55E` |
| Muted text | `#94A3B8` |
| Border | `#475569` |
| Display / code | JetBrains Mono |
| Body | IBM Plex Sans |
| Spacing | 4/8 scale |

The layout reads like one Go module: every section eyebrow is a real package
path (`// github.com/xgodev/boost/factory`, …).

## Sections

| Anchor | Headline | Shows |
|---|---|---|
| hero | Enterprise Go services. Built fast. | `go get`, stats, and a code window on the [shop example](#hero-example) with 3 tabs: the 9 files you write, the 12 files `go generate` writes, the run |
| `#boot` | One boot. Everything ready after it. | `boost.Start()` → config → log → factory → fx → serve; one config key set by file, env or flag |
| `#packages` | Every piece has one job. | Bento: factories, plugins, wrappers, middleware onion, inject graph, annotation, builder doors (`@REST`, `@GRPC`, `@Function`, `@MCP`, `@Middleware`, `@Default`), functions, MCP |
| `#builder` | From annotation to binary in two commands. | Annotate → `go generate ./...` → `go run ./cmd/api`, generated file tree |
| `#compat` | Upgrade boost. Keep your code. | Public API as a contract, SDK majors side by side (mongo-driver v1/v2), comparative quality gate, `go.mod` diff, release timeline |
| `#integrations` | There's a factory for that. | Logo wall of 23 libraries plus a `+18` tile; the other 18 by name; the full 41 by group with version paths in `<details id="all-libs">` (the `+18` tile opens it); why a factory: config keys, built-in plumbing, plugins; MongoDB by hand vs `mongo.NewConn` with plugins and env keys (#101) |
| `#agents` | — | `claude-plugin@xgodev` install, the 5 Iron Laws |
| final CTA | Start with one line. | `go get` |

## Hero example

The hero window is a small but complete app, [`example/shop`](example/shop)
(`github.com/acme/shop`): `POST /stock/:sku/reserve` reserves one unit and answers
409 when the item is sold out. It follows the xgodev layout from
`claude-plugin@xgodev` (`skills/dev/golang`):

| Layer | Path | Annotation |
|---|---|---|
| handler | `internal/app/api/rest/handler/stock/reserver/` | `@Provide`, `@REST` |
| service | `internal/pkg/domain/service/stock/reserver/` (interface) and `local/` | `@Provide`, `@Middleware(log, prometheus)` |
| repository | `internal/pkg/domain/repository/stock/` (interface) | |
| datastore | `internal/pkg/datastore/stock/memory/` and `mongodb/` | `@Default` on memory, `@Config` on the mongo options |

The service receives the repository; the mongo collection lives only in the
mongodb datastore, built by boost's mongo factory. The site never shows a DB
client inside a service.

Everything in the window is real: the sources, the files `go generate` wrote
(`gen/` and `cmd/api/main.go`) and the terminal output of a run. To change it:

1. edit the app in `example/shop` and run `go generate ./...` there
2. `go build ./...` and run it; paste new terminal output into `RUN` in
   [`example/build_hero.py`](example/build_hero.py)
3. from `docs/site`: `python3 example/build_hero.py example/shop index.html`,
   which rewrites the page between the `hero:start` and `hero:end` comments

`example/shop` is its own module with `replace github.com/xgodev/boost =>
../../../..`, so the repo's `go build ./...` skips it. It builds against
`develop`, but `go generate` needs the builder from the feature branches below.

## Framework links

Every third-party framework the page names links to its official page, so a
reader can check what boost builds on.

- **fx** always links to <https://uber-go.github.io/fx/>, everywhere, code
  included (`fx.New`, `fx.Module`, …).
- **Other frameworks** (echo, gRPC, MCP, kafka, zap, OpenTelemetry, …) link in
  prose, labels and lists.
- **Not linked:** code identifiers that are boost's own (`echo.NewServer` is
  boost's echo factory), boost package paths, config keys, annotation
  arguments, plugin chips, the middleware onion labels and terminal output.
- Links open in the same tab. They use `a.ref`: dotted underline that turns
  solid on hover, the accent focus ring on keyboard focus. The underline, not
  colour alone, marks a link.
- Every URL is checked to answer 200 before it goes on the page.

## Numbers used on the page

Counted from the repo on 2026-10-18 (#88, #131, #129, #126). Recount when a factory or library is added.

| Number | What |
|---|---|
| 53 | factories: directories under `factory/contrib/` with a `New…`/`Register…` constructor in non-test, non-example Go files, outside `/plugins/` and `/internal/` |
| 79 | plugins: the same count, only `/plugins/` directories |
| 41 | upstream libraries: the 43 `org/repo` directories under `factory/contrib/`, minus `cloud.google.com/grpc`, which wraps `google.golang.org/grpc` (already counted), and `apache/cassandra-gocql-driver`, the Apache home of `gocql/gocql` (#131). `cloud.google.com/api` is `google.golang.org/api` |
| 220 | tags since 2024-07 (latest `v1.2.0`) |

The 53/79/41 include the MCP factory (server and client, #70), echo v5 (#129) and Bigtable (#126).

`annotation`, `inject`, the MCP factory, `builder` and its doors (`@MCP`
since #79), `@Middleware`, `@Default`, `@Config` and the `gen/` layout that mirrors the project are on
`develop` and released in `v1.2.0`. `@MCP` puts a handler
`func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, Out, error)`
on the boost MCP server as a tool: the name is the handler in snake case and the
description is its doc comment (`builder/examples/simple/stocktool`).

## Checks done

- rendered headless at 1440px and 390px, every hero tab
- `example/shop` builds; `build_hero.py` rewrites the page byte for byte
- every framework link answers 200
- no horizontal scroll at 390px
- code-window tabs work with the arrow keys; copy buttons fall back to
  selecting the text when the clipboard is blocked
