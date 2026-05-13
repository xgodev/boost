# golang-boost — Claude Code plugin

Single Claude Code plugin shipping one skill per `github.com/xgodev/boost` component (53 skills total). Skills are component-aligned, cross-reference each other via `REQUIRED BACKGROUND` markers, and — for factories with framework-shipped examples — anchor their guidance on the boost source's own `examples/` dir.

## Install

```
/plugin marketplace add xgodev/boost
/plugin install golang-boost@xgodev
```

## Update

```
/plugin update golang-boost@xgodev
```

## Skills shipped (53)

### Foundations (7)

`boost-start`, `boost-wrapper-log`, `boost-wrapper-log-backends`, `boost-wrapper-config`, `boost-wrapper-publisher`, `boost-wrapper-cache`, `boost-model-errors`.

### Factories (36)

Every contrib under `factory/contrib/` is covered. Skills marked ✦ are anchored on framework-shipped `examples/` dirs.

**HTTP / API**
- `boost-factory-echo` ✦ — Echo HTTP server + plugin set
- `boost-factory-resty` ✦ — Resty outbound HTTP client
- `boost-factory-grpc` ✦ — google.golang.org/grpc client + server (incl. autoTLS)
- `boost-factory-graphql` — graphql-go handler
- `boost-factory-cloudevents` — CloudEvents HTTP receiver
- `boost-factory-net-http2` — HTTP/2 server tunables

**Messaging**
- `boost-factory-pubsub` ✦ — GCP Pub/Sub client
- `boost-factory-kafka` ✦ — Confluent Kafka raw producer/consumer
- `boost-factory-nats` — raw NATS connection
- `boost-factory-goka` — Goka emitter
- `boost-factory-gocloud-pubsub` ✦ — provider-portable URL-driven pub/sub

**Databases**
- `boost-factory-mongo` ✦ — mongo-driver v1 + v2
- `boost-factory-cassandra` ✦ — gocql session
- `boost-factory-pgx` — jackc/pgx PostgreSQL
- `boost-factory-godror` — godror Oracle
- `boost-factory-bigquery` — Google BigQuery
- `boost-factory-firestore` — Google Firestore
- `boost-factory-elasticsearch` ✦ — go-elasticsearch client + bulk indexer

**Caches & embedded stores**
- `boost-factory-redis` ✦ — go-redis client + cluster
- `boost-factory-bigcache` — allegro/bigcache
- `boost-factory-freecache` — coocood/freecache
- `boost-factory-buntdb` — embedded KV
- `boost-factory-memdb` — hashicorp in-memory indexed store

**Cloud SDKs / composition**
- `boost-factory-aws` ✦ — AWS SDK v2 umbrella + per-service examples
- `boost-factory-gcp-api` — GCP API options composed by every cloud.google.com factory
- `boost-factory-gcp-grpc` — GCP gRPC dial options composed by every cloud.google.com factory
- `boost-factory-k8s` — Kubernetes clientset

**Observability**
- `boost-factory-otel` — OpenTelemetry exporters + readers
- `boost-factory-datadog` — Datadog APM bridge
- `boost-factory-prometheus` — Prometheus metrics endpoint

**Resilience & utility**
- `boost-factory-hystrix` — circuit breaker per upstream
- `boost-factory-ants` — bounded goroutine pool
- `boost-factory-vault` — Vault secret retrieval
- `boost-factory-cobra` — multi-command CLI
- `boost-factory-fx` — fx app builder
- `boost-factory-ftp` — legacy FTP partner connections

### Bootstrap (event-driven functions) (5)

`boost-bootstrap-function`, `boost-bootstrap-adapter-pubsub` (incl. ctx-loss workaround), `boost-bootstrap-adapter-nats`, `boost-bootstrap-adapter-kafka`, `boost-bootstrap-middleware`.

### Extra & DI (4)

`boost-extra-middleware`, `boost-extra-health`, `boost-extra-multiserver`, `boost-fx-modules`.

### Contributing (1)

`boost-maintainer` — layout convention, constructor trio, `// TODO(maintainer-review):` discipline, multi-service SDK rule.

## Discovery

Claude reads each skill's `description` field at session start. When your context matches (Go file importing the corresponding boost subpackage, or a question about its APIs), the relevant skill activates. Cross-references via `REQUIRED BACKGROUND` keep load minimal — a dev working only on HTTP API typically loads `boost-start`, `boost-wrapper-log`, `boost-wrapper-config`, `boost-model-errors`, and `boost-factory-echo`, not the function/Pub/Sub/Kafka ones.

## Validation note

The 7 foundations + 5 bootstrap + 11 factories with framework-shipped examples (✦) were extracted from real boost source patterns or validated through TDD cycles against the `examples/` dirs.

The remaining 25 factory skills (databases, caches, cloud SDKs, observability, utilities) were inferred from each factory's constructor trio + `ConfigAdd` shape — they apply the same `<area>/<vendor>/<lib>/v<major>/` convention codified in `boost-maintainer`. They cover canonical construction and known pitfalls but should be re-validated against real-world failures when first exercised on a non-trivial task.

## Adding a new skill

Add `boost-<area>-<name>` when:
- A boost component ships an `examples/` dir, OR
- A boost component has non-obvious wiring (multi-step, plugin order, ctx-loss workaround), OR
- An AI agent demonstrably fails on the component in a real task.

Skip placeholder skills. Pure trivial wrappers don't need a dedicated skill — `boost-maintainer` already documents the universal layout/trio convention.

## License

MIT — see [LICENSE](./LICENSE).
