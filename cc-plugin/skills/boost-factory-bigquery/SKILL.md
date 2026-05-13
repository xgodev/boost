---
name: boost-factory-bigquery
description: "Use when constructing a Google BigQuery client in a Go service via github.com/xgodev/boost/factory/contrib/cloud.google.com/bigquery/v1. Covers NewClient + variants and the per-instance config root (boost.factory.gcp.bigquery.*) with nested apiOptions/grpcOptions. Triggers on imports under factory/contrib/cloud.google.com/bigquery/, on questions about BigQuery in a boost service, or NewClient construction."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`, `boost-factory-gcp-api`, `boost-factory-gcp-grpc`.

```go
import bq "github.com/xgodev/boost/factory/contrib/cloud.google.com/bigquery/v1"

client, err := bq.NewClient(ctx)
if err != nil { log.Fatalf("bigquery: %v", err) }
defer client.Close()
```

Configure under `boost.factory.gcp.bigquery.*`. The factory composes nested `boost.factory.gcp.bigquery.apiOptions.*` (GCP credentials, endpoint) and `boost.factory.gcp.bigquery.grpcOptions.*` (keepalive, retries, message size). Override individually:

```
BOOST_FACTORY_GCP_BIGQUERY_APIOPTIONS_PROJECTID=my-project
BOOST_FACTORY_GCP_BIGQUERY_GRPCOPTIONS_KEEPALIVE_TIME=30s
```

`plugins ...clientgrpc.Plugin` accepts gRPC interceptors (tracing, metrics).

## Red flags

| Red flag | Fix |
|---|---|
| `bigquery.NewClient(ctx, projectID)` from upstream SDK directly | `bq.NewClient(ctx)` |
| `GOOGLE_CLOUD_PROJECT` via `os.Getenv` | `BOOST_FACTORY_GCP_BIGQUERY_APIOPTIONS_PROJECTID` |
| Forgetting `defer client.Close()` | Add it — drains gRPC |
