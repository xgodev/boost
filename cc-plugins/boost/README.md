# boost — Claude Code plugin

A single Claude Code plugin shipping all the skills for working with [github.com/xgodev/boost](https://github.com/xgodev/boost).

One install. 18 skills, each independently versioned in its own SKILL.md frontmatter.

## Install

```
/plugin marketplace add xgodev/boost
/plugin install boost@xgodev-boost
```

That's it. All 18 skills become discoverable in your Claude Code session.

## Update

```
/plugin update boost@xgodev-boost
```

## Skills shipped

| Skill | Subsystem | Status |
|---|---|---|
| `boost-core` | Iron Laws, `boost.Start`, `log.FromContext`, `config`, `model/errors` | mature (v0.3.0) |
| `boost-factory-echo` | HTTP APIs with Echo | stub |
| `boost-factory-resty` | Outbound HTTP clients | stub |
| `boost-factory-pubsub` | GCP Pub/Sub client factory | stub |
| `boost-bootstrap-function` | `function.New` / `fn.Run` plumbing | stub |
| `boost-bootstrap-adapter-pubsub` | Pub/Sub subscriber (incl. ctx-loss workaround) | stub |
| `boost-bootstrap-adapter-nats` | NATS subscriber | stub |
| `boost-bootstrap-adapter-kafka` | Kafka subscriber | stub |
| `boost-bootstrap-middleware` | recovery / logger / publisher / ignore_errors | stub |
| `boost-wrapper-publisher` | publisher drivers | stub |
| `boost-wrapper-cache` | cache drivers | stub |
| `boost-wrapper-log` | log backends | stub |
| `boost-wrapper-config` | koanf config wrapper | stub |
| `boost-fx-modules` | uber/fx modules | stub |
| `boost-extra-middleware` | `AnyErrorMiddleware` / `AnyErrorWrapper` | stub |
| `boost-extra-health` | health checkers | stub |
| `boost-extra-multiserver` | multi-server lifecycle | stub |
| `boost-maintainer` | contributor guide | stub |

Stubs are placeholders that route to `boost-core`. They get lapidated via TDD cycle when the corresponding feature is exercised — see [`../CONTRIBUTING.md`](../CONTRIBUTING.md).

## How discovery works

Claude Code reads each skill's `description` field at session start. When your context matches (e.g., you open a `.go` file that imports `github.com/xgodev/boost/factory/contrib/labstack/echo/v4`), the relevant skill activates automatically. You don't pick — Claude does, based on triggers in the skill description.

If you want to force-load a specific skill: invoke it via the Skill tool with the skill name.

## License

MIT — see [LICENSE](./LICENSE). Originally seeded from `samber/cc-skills-golang` (also MIT).
