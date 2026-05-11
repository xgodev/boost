# golang-boost — Claude Code plugin

Single Claude Code plugin shipping one skill per `github.com/xgodev/boost` component (28 skills total). Skills are component-aligned, cross-reference each other via `REQUIRED BACKGROUND` markers, and — for factories — anchor their guidance on the framework's own `examples/` dir so devs can verify the canonical shape against runnable code.

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
| `boost-wrapper-cache` | `wrapper/cache` — `Manager[T]` + drivers (Redis, allegro, stretchr) + codecs |
| `boost-model-errors` | `model/errors` — typed error catalog + matchers |

### Factories (each backed by the component's `examples/` dir)

| Skill | Component |
|---|---|
| `boost-factory-echo` | `factory/contrib/labstack/echo/v4` — HTTP API + plugin set |
| `boost-factory-resty` | `factory/contrib/go-resty/resty/v2` — outbound HTTP |
| `boost-factory-pubsub` | `factory/contrib/cloud.google.com/pubsub/v1` |
| `boost-factory-mongo` | `factory/contrib/go.mongodb.org/mongo-driver/v1` and `v2` |
| `boost-factory-cassandra` | `factory/contrib/gocql/gocql/v1` |
| `boost-factory-redis` | `factory/contrib/redis/go-redis/v9` (single + cluster) |
| `boost-factory-elasticsearch` | `factory/contrib/elastic/go-elasticsearch/v8` (client + bulk indexer) |
| `boost-factory-kafka` | `factory/contrib/confluentinc/confluent-kafka-go/v2` (raw producer/consumer) |
| `boost-factory-aws` | `factory/contrib/aws/aws-sdk-go-v2/v1` umbrella + S3/SNS/SQS/Kinesis examples |
| `boost-factory-grpc` | `factory/contrib/google.golang.org/grpc/v1` (client, server, autoTLS server) |
| `boost-factory-gocloud-pubsub` | `factory/contrib/gocloud.dev/pubsub/v0` (provider-portable) |

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

## Adding a new skill

Add a new `boost-<area>-<name>` skill when:
- A boost component ships an `examples/` dir AND
- The component has a non-trivial canonical shape (multi-step wiring, gotchas, plugin order) that an AI agent could miss.

Pure thin wrappers (one constructor, no surprises) don't need a dedicated skill — `boost-maintainer` covers the universal layout/trio convention.

Skip placeholder skills. If you don't have content from real usage or framework examples, don't ship a stub.

## License

MIT — see [LICENSE](./LICENSE).
