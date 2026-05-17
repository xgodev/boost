---
name: boost-wrapper-publisher
description: "Use when publishing CloudEvents from a Go service via github.com/xgodev/boost/wrapper/publisher, or when wiring a driver under wrapper/publisher/driver/contrib/. Covers the publisher.New(driver) entrypoint, the Driver interface contract, available drivers (Pub/Sub, Kafka, NATS, Goka), and how the function publisher middleware (boost-bootstrap-middleware) consumes a publisher for deadletter routing. Triggers on imports under wrapper/publisher/, on questions about publisher.New, Driver interface, or republishing in a function handler."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`. For function deadletter integration → `boost-bootstrap-middleware`. For adding a new driver to boost → `boost-maintainer`.

## Construct via the driver-agnostic wrapper

```go
import (
    "github.com/xgodev/boost/wrapper/publisher"
    drvpubsub "github.com/xgodev/boost/wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v1"
)

drv, err := drvpubsub.New(ctx, pb)   // pb is *pubsub.Client (see boost-factory-pubsub)
if err != nil { log.Fatalf("publisher driver: %v", err) }

pub := publisher.New(drv)
```

`publisher.New(driver)` returns a `Publisher` that consumers — including the function publisher middleware — depend on. Driver swaps don't require call-site changes.

## Available drivers (out of the box)

| Driver | Path |
|---|---|
| Google Cloud Pub/Sub | `wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v1` |
| Confluent Kafka | `wrapper/publisher/driver/contrib/confluentinc/confluent-kafka-go/v2` |
| NATS | `wrapper/publisher/driver/contrib/nats-io/nats.go/v1` |
| Lovoo Goka | `wrapper/publisher/driver/contrib/lovoo/goka/v1` |

## Republishing inside a function handler

Wire `pub` into the function publisher middleware (see `boost-bootstrap-middleware`):

```go
pmi, _ := pm.NewAnyErrorMiddleware[*cloudevents.Event](pub)
```

The middleware fires on a handler's successful return — the returned `*cloudevents.Event` gets published via the wrapped driver. The handler does NOT call `pub.Publish` manually; that bypasses middleware bookkeeping (deadletter, ordering, instrumentation).

## Topic selection

The driver typically uses the event's `Subject()` to route:

| Driver | What `Subject()` means |
|---|---|
| GCP Pub/Sub | Topic ID |
| Kafka | Topic |
| NATS | Subject |
| SNS (if added) | Full ARN — see `boost-maintainer` for the extrapolation discipline |

Set the destination via `out.SetSubject("topic-name")` before returning from the handler.

## Red flags

| Red flag | Fix |
|---|---|
| Calling `pb.Topic("foo").Publish(ctx, ...)` directly inside a handler | Return `*cloudevents.Event` from the handler; let the publisher middleware drive `pub.Publish` |
| Building a publisher per request instead of once at startup | Construct `pub` in `main`, pass to handlers via deps |
| Using a third-party publisher abstraction in parallel | Use `wrapper/publisher` as the only abstraction so middleware sees everything |
