---
name: boost-factory-gcp-grpc
description: "Use when configuring shared GCP gRPC options (keepalive, retry policy, message size, dial-options) consumed by all cloud.google.com/* factories that talk to GCP gRPC services — bigquery, firestore, pubsub. Provided by github.com/xgodev/boost/factory/contrib/cloud.google.com/grpc/v1. Triggers on imports under factory/contrib/cloud.google.com/grpc/, on questions about GCP gRPC tunables, keepalive, or retry settings in a boost service. For generic (non-GCP) gRPC, see boost-factory-grpc."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. Composed by `boost-factory-bigquery`, `boost-factory-firestore`, `boost-factory-pubsub`. For non-GCP gRPC client/server → `boost-factory-grpc`.

## What it provides

GCP-tuned gRPC dial options registered under `boost.factory.gcp.grpc.*` and composed by every concrete GCP factory at `<root>.grpcOptions.*`. Sensible defaults for GCP API endpoints (keepalive that survives Google's idle close, retry policy that respects GCP's recommended backoff, message-size caps appropriate for streaming responses).

## Tunables (typical)

- `grpcOptions.keepalive.time` / `keepalive.timeout` / `keepalive.permitWithoutStream`
- `grpcOptions.retries.maxAttempts`, retry policy
- `grpcOptions.message.maxRecvSize` / `maxSendSize`

Override per-service: tighter retries on Firestore only:

```
BOOST_FACTORY_GCP_FIRESTORE_GRPCOPTIONS_RETRIES_MAXATTEMPTS=2
```

## Red flags

| Red flag | Fix |
|---|---|
| Building dial options by hand for a GCP service | Configure `<root>.grpcOptions.*` instead |
| Disabling keepalive thinking it saves resources | Google's intermediaries close idle gRPC streams; keepalive prevents costly reconnect storms |
| One global gRPC config across mixed GCP and non-GCP services | Use this skill ONLY for cloud.google.com factories; generic gRPC uses `boost-factory-grpc` |
