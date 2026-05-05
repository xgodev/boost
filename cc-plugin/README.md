# golang-boost — Claude Code plugin

Single Claude Code plugin shipping one skill per `github.com/xgodev/boost` component. Skills are component-aligned (one per boost subsystem) and cross-reference each other via `REQUIRED BACKGROUND` markers — Claude loads only the ones relevant to your current task.

## Install

```
/plugin marketplace add xgodev/boost
/plugin install golang-boost@xgodev
```

## Update

```
/plugin update golang-boost@xgodev
```

## Skills shipped

| Skill | Component (under `github.com/xgodev/boost/`) |
|---|---|
| `boost-start` | the `boost.Start()` boot sequence |
| `boost-wrapper-log` | `wrapper/log` — `log.FromContext`, structured logging |
| `boost-wrapper-config` | `wrapper/config` — `config.Add`, typed accessors, env override |
| `boost-wrapper-publisher` | `wrapper/publisher` — driver-agnostic publishing |
| `boost-model-errors` | `model/errors` — typed error catalog + matchers |
| `boost-factory-echo` | `factory/contrib/labstack/echo/v4` — HTTP API factory |
| `boost-factory-pubsub` | `factory/contrib/cloud.google.com/pubsub/v1` — Pub/Sub client |
| `boost-bootstrap-function` | `bootstrap/function` — generic function plumbing + `Handler[T]` rule |
| `boost-bootstrap-adapter-pubsub` | `bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v1` — incl. ctx-loss workaround |
| `boost-bootstrap-middleware` | `bootstrap/function/middleware` — recovery / logger / publisher / ignore_errors stack |
| `boost-extra-middleware` | `extra/middleware` — generic `NewAnyErrorWrapper` for the workaround pattern |
| `boost-maintainer` | guide for adding a new contrib (driver / adapter / plugin / module) |

Components NOT yet covered (no validated content): `wrapper/cache`, `wrapper/log/contrib/*`, NATS / Kafka adapters, `extra/health`, `extra/multiserver`, `fx/modules`, `factory/contrib/go-resty`. New skills land only when there's concrete evidence of need (real-world failure observed when Claude codes against that component).

## Discovery

Claude reads each skill's `description` field at session start. When your context matches (Go file importing the corresponding boost subpackage, or a question about its APIs), the relevant skill activates. You don't pick — Claude does.

## License

MIT — see [LICENSE](./LICENSE).
