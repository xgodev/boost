---
name: boost-factory-kafka
description: "Use when constructing a raw Kafka producer or consumer (Confluent client) via github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2. Covers NewProducer / NewConsumer + variants and the canonical shapes shipped under examples/{producer,consumer}/. Use this skill for the FACTORY layer (raw *kafka.Producer / *kafka.Consumer); use boost-bootstrap-adapter-kafka for event-handler subscriber wiring with the middleware chain."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. For event-handler integration (CloudEvents + middleware chain) → `boost-bootstrap-adapter-kafka`.

## Canonical examples (ship with boost)

- `factory/contrib/confluentinc/confluent-kafka-go/v2/examples/producer/main.go`
- `factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/main.go`

## Construction

```go
import kafkafact "github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2"

producer, err := kafkafact.NewProducer(ctx)
consumer, err := kafkafact.NewConsumer(ctx)
```

Configure brokers, group id, security under `boost.factory.kafka.*` (override `BOOST_FACTORY_KAFKA_*`).

## Factory vs adapter

| Use case | Reach for |
|---|---|
| Custom Kafka pipeline (low-level Poll, manual offset commit) | `boost-factory-kafka` (raw client) |
| Event handler over CloudEvents semantics with middleware chain | `boost-bootstrap-adapter-kafka` |

## Red flags

| Red flag | Fix |
|---|---|
| `kafka.NewProducer(&kafka.ConfigMap{...})` directly | `kafkafact.NewProducer(ctx)` |
| Brokers/groupID via `os.Getenv` | `BOOST_FACTORY_KAFKA_*` |
| Forgetting `producer.Close()` / `consumer.Close()` on shutdown | Add it |
