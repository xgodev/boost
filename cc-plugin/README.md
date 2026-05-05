# golang-boost — Claude Code plugin

Single Claude Code plugin shipping one skill per `github.com/xgodev/boost` component (20 skills total). Skills are component-aligned and cross-reference each other via `REQUIRED BACKGROUND` markers — Claude loads only the ones relevant to your current task.

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

### Foundations

| Skill | Component (under `github.com/xgodev/boost/`) |
|---|---|
| `boost-start` | the `boost.Start()` boot sequence |
| `boost-wrapper-log` | `wrapper/log` — `log.FromContext`, structured logging |
| `boost-wrapper-log-backends` | `factory/contrib/{zap,zerolog,logrus}/v1` — picking the backend |
| `boost-wrapper-config` | `wrapper/config` — `config.Add`, typed accessors, env override |
| `boost-wrapper-publisher` | `wrapper/publisher` — driver-agnostic publishing |
| `boost-wrapper-cache` | `wrapper/cache` — Manager[T] + drivers (Redis, allegro, stretchr) + codecs |
| `boost-model-errors` | `model/errors` — typed error catalog + matchers |

### Factories

| Skill | Component |
|---|---|
| `boost-factory-echo` | `factory/contrib/labstack/echo/v4` — HTTP API factory |
| `boost-factory-resty` | `factory/contrib/go-resty/resty/v2` — outbound HTTP clients |
| `boost-factory-pubsub` | `factory/contrib/cloud.google.com/pubsub/v1` — Pub/Sub client |

### Bootstrap (event-driven functions)

| Skill | Component |
|---|---|
| `boost-bootstrap-function` | `bootstrap/function` — generic plumbing + `Handler[T]` rule |
| `boost-bootstrap-adapter-pubsub` | `bootstrap/.../pubsub/v1` — incl. ctx-loss workaround |
| `boost-bootstrap-adapter-nats` | `bootstrap/.../nats-io/nats.go/v1` — same workaround pattern |
| `boost-bootstrap-adapter-kafka` | `bootstrap/.../confluent-kafka-go/v2` — same workaround pattern |
| `boost-bootstrap-middleware` | recovery → logger → publisher chain order + deadletter wrapping |

### Extra & DI

| Skill | Component |
|---|---|
| `boost-extra-middleware` | `extra/middleware` — `NewAnyErrorWrapper` for the workaround pattern |
| `boost-extra-health` | `extra/health` — checkers, liveness vs readiness |
| `boost-extra-multiserver` | `extra/multiserver` — coordinated multi-listener lifecycle |
| `boost-fx-modules` | `fx/modules` — uber/fx wiring at scale |

### Contributing

| Skill | Component |
|---|---|
| `boost-maintainer` | guide for adding a new contrib (driver / adapter / plugin / module) |

## Discovery

Claude reads each skill's `description` field at session start. When your context matches (Go file importing the corresponding boost subpackage, or a question about its APIs), the relevant skill activates. Cross-references via `REQUIRED BACKGROUND` keep load minimal — a dev working only on HTTP API typically loads `boost-start`, `boost-wrapper-log`, `boost-wrapper-config`, `boost-model-errors`, and `boost-factory-echo`, not the function/Pub/Sub/Kafka ones.

## License

MIT — see [LICENSE](./LICENSE).
