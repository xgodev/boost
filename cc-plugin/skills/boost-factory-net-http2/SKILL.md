---
name: boost-factory-net-http2
description: "Use when configuring HTTP/2 server settings (max concurrent streams, idle timeout, frame size) in a Go HTTP service that imports github.com/xgodev/boost/factory/contrib/golang.org/x/net/v0. The factory composes its options into the Echo server (boost-factory-echo) automatically; reach for this skill when you need to tune HTTP/2-specific knobs that the default settings don't cover. Triggers on imports under factory/contrib/golang.org/x/net/, on questions about HTTP/2 tuning or H2 frame settings in a boost HTTP service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. Typically composed by `boost-factory-echo`.

The factory exposes HTTP/2 server-side tunables under `boost.factory.golang.x.net.http2.server.*` (override `BOOST_FACTORY_GOLANG_X_NET_HTTP2_SERVER_*`):

| Knob | What |
|---|---|
| `maxConcurrentStreams` | Per-connection stream cap |
| `idleTimeout` | When to close an idle connection |
| `maxReadFrameSize` | Largest HTTP/2 frame accepted |
| `maxUploadBufferPerConnection` | Per-conn write buffer |
| `maxUploadBufferPerStream` | Per-stream write buffer |

The actual subdirs and exported helpers vary by boost version — consult `factory/contrib/golang.org/x/net/v0/` in your boost dependency.

## When to reach here

Most services don't need this — Echo's default HTTP/2 wiring is fine. Reach for it when:

- Hosting a streaming gRPC service on the same port as HTTP/1 — bump stream caps.
- High-throughput ingest endpoints showing HTTP/2 backpressure in dashboards.
- Long-lived connections (SSE, websockets-over-h2) needing tuned idle timeouts.

## Red flags

| Red flag | Fix |
|---|---|
| Patching `http2.Server{}` fields after server start | Configure before via `BOOST_FACTORY_GOLANG_X_NET_HTTP2_SERVER_*` |
| Disabling HTTP/2 entirely to avoid tuning | Tune the knobs that hurt; HTTP/2 is the default for good reasons (multiplexing, header compression) |
