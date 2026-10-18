# boost website — mockup

Issue: [#74](https://github.com/xgodev/boost/issues/74). Mockup: [`mockup.html`](mockup.html),
with its styles in [`mockup.css`](mockup.css) and its script in [`mockup.js`](mockup.js).
Open the HTML in a browser; only Google Fonts load from the network.

Scope: the page's design and content. Hosting, domain and the production build
come later.

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
| hero | Write the handler. Boost writes the service. | `go get`, stats, code window with 3 tabs: annotated `service.go`, generated `new_module.go`, terminal |
| `#boot` | One boot. Everything ready after it. | `boost.Start()` → config → log → factory → fx → serve; one config key set by file, env or flag |
| `#packages` | Every piece has one job. | Bento: factories, plugins, wrappers, middleware onion, inject graph, annotation, builder doors (`@REST`, `@GRPC`, `@Function`, `@Middleware`, `@Default`), functions, MCP |
| `#builder` | From annotation to binary in two commands. | Annotate → `go generate ./...` → `go run ./cmd/api`, generated file tree |
| `#compat` | Upgrade boost. Keep your code. | Public API as a contract, SDK majors side by side (mongo-driver v1/v2), comparative quality gate, `go.mod` diff, release timeline |
| `#integrations` | 38 libraries. Same shape every time. | Upstream libraries by group, each with its version path |
| `#agents` | — | `claude-plugin@xgodev` install, the 5 Iron Laws |
| final CTA | Start with one line. | `go get` |

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

Counted from the repo on 2026-10-18. Recount before publishing.

| Number | What |
|---|---|
| 48 | factories: directories under `factory/contrib/` with a `New…`/`Register…` constructor in non-test, non-example Go files, outside `/plugins/` and `/internal/` |
| 58 | plugins: the same count, only `/plugins/` directories |
| 38 | upstream libraries: distinct `org/repo` under `factory/contrib/` |
| 219 | tags since 2024-07 (latest `v1.1.1`) |

The 48/58/38 include the MCP factory (server and client), which is on
`feature/70-mcp-factory` (#70); `develop` alone counts 46/54/37.

`annotation` and `inject` are on `develop` and `builder` is on feature branches
(#67, #72); none of them is in a tag yet. `@Default` (#75) is a design, not
code. The page shows all of them as released. The builder has no `@MCP` door:
#70 leaves it out of scope.

## Checks done

- rendered headless at 1440px and 390px
- every framework link answers 200
- no horizontal scroll at 390px
- code-window tabs work with the arrow keys; copy buttons fall back to
  selecting the text when the clipboard is blocked
