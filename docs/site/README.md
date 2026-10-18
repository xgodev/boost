# boost website — mockup

Issue: [#74](https://github.com/xgodev/boost/issues/74). Mockup: [`mockup.html`](mockup.html)
(single file, open it in a browser; only Google Fonts load from the network).

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
| `#packages` | Every piece has one job. | Bento: factories, plugins, wrappers, middleware onion, inject graph, annotation, builder doors, functions |
| `#builder` | From annotation to binary in two commands. | Annotate → `go generate ./...` → `go run ./cmd/api`, generated file tree |
| `#compat` | Upgrade boost. Keep your code. | Public API as a contract, SDK majors side by side (mongo-driver v1/v2), comparative quality gate, `go.mod` diff, release timeline |
| `#integrations` | 36 libraries. Same shape every time. | Upstream libraries by group, each with its version path |
| `#agents` | — | `claude-plugin@xgodev` install, the 5 Iron Laws |
| final CTA | Start with one line. | `go get` |

## Numbers used on the page

Counted from the repo on 2026-10-18. Recount before publishing.

| Number | What |
|---|---|
| 45 | factories under `factory/contrib/` |
| 36 | upstream libraries they wrap |
| 62 | factory plugins |
| 219 | tags since 2024-07 (latest `v1.1.1`) |

`annotation` and `inject` are on `develop` and `builder` is on feature branches
(#67, #72); none of them is in a tag yet. The page shows them as released.

## Checks done

- rendered headless at 1440px and 390px
- no horizontal scroll at 390px
- code-window tabs work with the arrow keys; copy buttons fall back to
  selecting the text when the clipboard is blocked
