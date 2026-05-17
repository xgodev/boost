---
name: boost-bootstrap-adapter-kafka
description: "Use when writing or reviewing a Go event-driven service that consumes Kafka via github.com/xgodev/boost/bootstrap/function/adapter/contrib/confluentinc/confluent-kafka-go/v2. Covers consumer-group wiring, offset commit semantics, topic subscription, and the same ctx-loss issue documented for the Pub/Sub adapter (helper.go hard-codes context.Background, requiring the workaround pattern for graceful shutdown). Triggers on imports under bootstrap/function/adapter/contrib/confluentinc/, on questions about Kafka consumers in a boost function, or on signal handling for Kafka workers."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-bootstrap-function` — handler typing rule.
- `boost-bootstrap-middleware` — recovery/logger/publisher chain.
- `boost-extra-middleware` — `NewAnyErrorWrapper` for the workaround.
- `boost-bootstrap-adapter-pubsub` — same shape, full workaround pattern documented there.

## Canonical (prototype / dev)

```go
import (
    akafka "github.com/xgodev/boost/bootstrap/function/adapter/contrib/confluentinc/confluent-kafka-go/v2"
    "github.com/xgodev/boost/bootstrap/function"
)

fn, _ := function.New[*cloudevents.Event](rec, lmi, pmi)
fn.Run(ctx, handle, akafka.New[*cloudevents.Event](consumer))
```

Topics, consumer group, broker list, and offset reset behavior are configured via `boost.bootstrap.function.adapter.kafka.*` (override via `BOOST_BOOTSTRAP_FUNCTION_ADAPTER_KAFKA_*`).

## Production caveat — same ctx-loss as Pub/Sub

`bootstrap/function/adapter/contrib/confluentinc/confluent-kafka-go/v2/helper.go:41` hard-codes:

```go
err := subscriber.Subscribe(context.Background())
```

SIGTERM does not gracefully drain in-flight messages. Apply the **same workaround pattern** documented in `boost-bootstrap-adapter-pubsub`: bypass `fn.Run`, build the chain via `extra/middleware.NewAnyErrorWrapper`, drive `akafka.NewSubscriber` with a signal-aware ctx, and add the `// TODO(boost-upstream):` annotation naming the offending file.

## Consumer group semantics

Kafka delivers each message to exactly one member of a consumer group. The boost adapter respects the group config — set `boost.bootstrap.function.adapter.kafka.groupID` so multiple replicas of your service share the partition load.

Offsets commit on successful handler return (post-publisher middleware). A handler error propagates as a nack — the message replays per Kafka's redelivery semantics. Wrap errors via `bootsterrors.Wrap` (see `boost-model-errors`) so the deadletter middleware can route by type.

## Red flags

| Red flag | Fix |
|---|---|
| `kafka.Consumer.Poll(...)` / `ReadMessage(...)` loops directly | Use `akafka.NewSubscriber(...).Subscribe(ctx)` or `function.New + fn.Run` |
| Bypass of `fn.Run` without `// TODO(boost-upstream):` naming `helper.go:41` | Add the comment, OR accept ungraceful shutdown |
| Reading `KAFKA_BROKERS` / `KAFKA_GROUP_ID` via `os.Getenv` | Use `BOOST_BOOTSTRAP_FUNCTION_ADAPTER_KAFKA_*` overrides |
| Manual offset commit inside the handler | Let the publisher middleware drive commit on success (default) |
