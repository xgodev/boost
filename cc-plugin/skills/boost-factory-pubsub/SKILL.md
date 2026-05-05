---
name: boost-factory-pubsub
description: "Use when constructing a Google Cloud Pub/Sub client in a Go service that imports github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1. Covers fpubsub.NewClient construction, lifecycle (defer Close), and how the resulting client feeds both the publisher driver (boost-wrapper-publisher) and the subscriber adapter (boost-bootstrap-adapter-pubsub). Triggers on imports under factory/contrib/cloud.google.com/pubsub/, on questions about NewClient or pubsub project ID configuration in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-start` — `boost.Start()` first.
- `boost-wrapper-config` — project ID is read from config, not `os.Getenv`.

## Construction

```go
import fpubsub "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"

pb, err := fpubsub.NewClient(ctx)
if err != nil { log.Fatalf("pubsub client: %v", err) }
defer pb.Close()
```

The client reads its project ID and connection options from boost config (`boost.factory.gcp.pubsub.apiOptions.projectId` and friends). Override at deploy via `BOOST_FACTORY_GCP_PUBSUB_APIOPTIONS_PROJECTID=...`.

## Lifecycle

The same `*pubsub.Client` typically feeds both:
- A publisher driver (see `boost-wrapper-publisher`)
- A subscriber adapter (see `boost-bootstrap-adapter-pubsub`)

Construct once, share, and `defer pb.Close()` so the gRPC connection drains on shutdown.

## Red flags

| Red flag | Fix |
|---|---|
| `pubsub.NewClient(ctx, projectID)` directly from the upstream SDK | Use `fpubsub.NewClient(ctx)` so config + observability instrumentation are wired |
| Reading `GOOGLE_CLOUD_PROJECT` via `os.Getenv` | Use `BOOST_FACTORY_GCP_PUBSUB_APIOPTIONS_PROJECTID` (or override the koanf key) |
| Forgetting `defer pb.Close()` | Add it — graceful gRPC shutdown |
| Constructing two clients (one for publish, one for subscribe) | Construct one, share it |
