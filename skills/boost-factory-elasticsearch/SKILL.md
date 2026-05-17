---
name: boost-factory-elasticsearch
description: "Use when constructing an Elasticsearch client (or bulk indexer) in a Go service via github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8. Covers NewClient + NewBulkIndexer and the canonical connect/health shapes shipped under examples/{connect,health}/. Triggers on imports under factory/contrib/elastic/go-elasticsearch/, on questions about Elasticsearch in a boost service, or on bulk-indexer wiring."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`.

## Canonical examples (ship with boost)

- `factory/contrib/elastic/go-elasticsearch/v8/examples/connect/main.go`
- `factory/contrib/elastic/go-elasticsearch/v8/examples/health/main.go`

The `health` example is the reference for wiring this factory into `boost-extra-health` checkers.

## Client

```go
import esfact "github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8"

es, err := esfact.NewClient(ctx)
if err != nil { log.Fatalf("elasticsearch: %v", err) }
```

Configure addresses, auth, transport TLS under `boost.factory.elasticsearch.*` (override `BOOST_FACTORY_ELASTICSEARCH_*`).

## Bulk indexer (high-throughput ingest)

```go
bi, err := esfact.NewBulkIndexer(ctx, es)
if err != nil { log.Fatalf("es bulk: %v", err) }
defer bi.Close(ctx)

bi.Add(ctx, esutil.BulkIndexerItem{Action: "index", Index: "orders", Body: strings.NewReader(jsonBody)})
```

Tune flush size + interval via `boost.factory.elasticsearch.bulk.*`.

## Red flags

| Red flag | Fix |
|---|---|
| `elasticsearch.NewClient(esCfg)` directly | `esfact.NewClient(ctx)` |
| Per-document `es.Index` calls in a hot path | Use `BulkIndexer` — orders of magnitude faster |
| Forgetting `defer bi.Close(ctx)` | Add it — flushes pending batches |
