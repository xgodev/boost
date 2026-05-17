---
name: boost-factory-firestore
description: "Use when constructing a Google Firestore client in a Go service via github.com/xgodev/boost/factory/contrib/cloud.google.com/firestore/v1. Covers NewClient + variants and the per-instance config root (boost.factory.gcp.firestore.*) with nested apiOptions/grpcOptions. Triggers on imports under factory/contrib/cloud.google.com/firestore/, on questions about Firestore in a boost service, or NewClient construction."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`, `boost-factory-gcp-api`, `boost-factory-gcp-grpc`.

```go
import fs "github.com/xgodev/boost/factory/contrib/cloud.google.com/firestore/v1"

client, err := fs.NewClient(ctx)
if err != nil { log.Fatalf("firestore: %v", err) }
defer client.Close()
```

Configure under `boost.factory.gcp.firestore.*`. Composes `apiOptions.*` (GCP credentials, endpoint) + `grpcOptions.*` (keepalive, retries). Override at deploy via `BOOST_FACTORY_GCP_FIRESTORE_*`. `plugins ...clientgrpc.Plugin` accepts gRPC interceptors.

## Red flags

| Red flag | Fix |
|---|---|
| `firestore.NewClient(ctx, projectID)` from upstream SDK directly | `fs.NewClient(ctx)` |
| `GOOGLE_CLOUD_PROJECT` via `os.Getenv` | `BOOST_FACTORY_GCP_FIRESTORE_APIOPTIONS_PROJECTID` |
| Forgetting `defer client.Close()` | Add it |
