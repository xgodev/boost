---
name: boost-factory-grpc
description: "Use when constructing a gRPC client or server in a Go service via github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1 (with subdirs client/ and server/). Covers the canonical shapes shipped under client/examples/examplesvc and server/examples/{examplesvc,examplesvcautotls}, including TLS-enabled server wiring. Triggers on imports under factory/contrib/google.golang.org/grpc/, on questions about gRPC client dial options, server interceptors, or autoTLS in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.2.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`.

## Canonical examples (ship with boost)

- `factory/contrib/google.golang.org/grpc/v1/client/examples/examplesvc/` — minimal client wiring
- `factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvc/` — minimal server wiring
- `factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvcautotls/` — server with auto-TLS

Read `examplesvcautotls` before enabling TLS on a new service — it shows the certificate-manager wiring boost expects.

## Two halves

| Path | When |
|---|---|
| `factory/contrib/google.golang.org/grpc/v1/client/` | Outbound gRPC calls — dial options, interceptors |
| `factory/contrib/google.golang.org/grpc/v1/server/` | Inbound gRPC service — listener, interceptors, TLS |

Configure under `boost.factory.grpc.client.*` and `boost.factory.grpc.server.*` (override with the matching `BOOST_FACTORY_GRPC_*` envs).

## GCP-tuned variant

For talking to GCP gRPC APIs (Pub/Sub, BigQuery, Firestore), the cloud-google factories compose `factory/contrib/cloud.google.com/grpc/v1` internally. You normally don't import it directly — you configure its keys at the per-service factory's `apiOptions` / `grpcOptions` namespace.

## Error → gRPC code (and custom errors)

The server converts errors to gRPC status via `server.Error(err)`, which resolves
the code through `model/errors.Classify` (`NotFound`→`codes.NotFound`,
`NotValid`/`BadRequest`→`InvalidArgument`, `Conflict`/`AlreadyExists`→
`AlreadyExists`, `Unauthorized`→`Unauthenticated`, `Forbidden`→`PermissionDenied`,
`ServiceUnavailable`→`Unavailable`, `NotImplemented`→`Unimplemented`,
`TooManyRequests`→`ResourceExhausted`, `Timeout`→`DeadlineExceeded`, else
`Internal`). To map an app-specific error to a code (or ignore it → returns
`nil`/OK), register it at boot — see `boost-model-errors`
(`Register`/`RegisterMatch`/`Ignore`). Don't edit `server.Error` by hand.

## Red flags

| Red flag | Fix |
|---|---|
| `grpc.Dial(...)` with hand-built dial options | Use the client factory so config + interceptors are wired |
| `grpc.NewServer()` without going through the server factory | Use the server factory so default interceptors (recovery, logging, tracing) are installed |
| TLS config hand-rolled instead of mirroring `examplesvcautotls` | Mirror the example shape — cert lifecycle is easy to get wrong |
| Forgetting `defer conn.Close()` (client) or `srv.GracefulStop()` (server) | Add them |
