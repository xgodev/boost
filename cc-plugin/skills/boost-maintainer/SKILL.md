---
name: boost-maintainer
description: "Use when contributing code to github.com/xgodev/boost itself — adding a new factory contrib, wrapper driver (publisher / cache / log / config), bootstrap adapter or middleware, Echo plugin, or fx module. Covers the strict layout convention (vendor/lib/v<major>/), the constructor trio (New / NewWithOptions / NewWithConfigPath), config registration in init() with ConfigAdd export for multi-instance, the multi-service SDK rule for AWS / Azure / GCP, and the // TODO(maintainer-review): convention for honest extrapolation. Triggers on file paths under boost's bootstrap/, factory/, wrapper/, extra/, fx/ trees, on questions about contrib layout, constructor patterns, ConfigAdd, or PR conventions for boost."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-wrapper-config` — config.Add semantics; new code registers its keys.
- `boost-wrapper-log` — logger access pattern in driver code.
- `boost-model-errors` — error verbs to mirror.

## Layout convention

| Adding | Path |
|---|---|
| Publisher driver | `wrapper/publisher/driver/contrib/<vendor>/<lib>/v<major>/` |
| Cache driver | `wrapper/cache/driver/contrib/<vendor>/<lib>/v<major>/` |
| Logger contrib | `wrapper/log/contrib/<vendor>/<lib>/v<major>/` |
| Function adapter | `bootstrap/function/adapter/contrib/<vendor>/<lib>/v<major>/` |
| Function middleware | `bootstrap/function/middleware/<name>/` |
| Factory contrib | `factory/contrib/<vendor>/<lib>/v<major>/` |
| Echo plugin | `factory/contrib/labstack/echo/v4/plugins/{native,extra,local}/<name>/` |
| Fx module | `fx/modules/<area>/<component>/` |

**Package name = the short library name** (`pubsub`, `nats`, `confluent`, `goka`, `sns`, `redis`), **not** the version directory leaf. Alias the upstream SDK at the import site if a name clash would occur:

```go
import asns "github.com/aws/aws-sdk-go-v2/service/sns"
```

## Multi-service SDKs (AWS SDK v2, Azure SDK, GCP SDK)

For `wrapper/`, `bootstrap/`, and `extra/` — split per service:

```
wrapper/publisher/driver/contrib/aws/sns/v1/
wrapper/publisher/driver/contrib/aws/sqs/v1/
wrapper/cache/driver/contrib/aws/dynamodb/v1/
```

Do **not** nest under the SDK module dir (`aws/aws-sdk-go-v2/v1/`). The umbrella SDK layout (`factory/contrib/aws/aws-sdk-go-v2/v1/client/<service>/`) is **exclusive to `factory/contrib/`** because factories ship clients that legitimately share an SDK version pin. Drivers don't.

## Constructor trio

```go
func New(ctx context.Context, c *upstream.Client) (publisher.Driver, error)
func NewWithOptions(ctx context.Context, c *upstream.Client, opts *Options) publisher.Driver
func NewWithConfigPath(ctx context.Context, c *upstream.Client, path string) (publisher.Driver, error)
```

Every driver/adapter exposes the same trio so call sites are interchangeable.

## Config registration

```go
// config.go
package <pkg>

import "github.com/xgodev/boost/wrapper/config"

const root = "boost.wrapper.publisher.driver.<name>"

func init() {
    config.Add(root+".log.level", "INFO", "log level")
    config.Add(root+".publishTimeout", "10s", "per-event publish timeout")
}

// ConfigAdd — exported so multi-instance consumers can register at a non-default path.
func ConfigAdd(path string) {
    config.Add(path+".log.level", "INFO", "log level")
    // mirror the rest with `path` instead of `root`
}
```

Multi-instance pattern: one binary publishes to multiple SNS topics / Pub/Sub projects / Kafka clusters by calling `<driver>.ConfigAdd("boost.wrapper.publisher.driver.<name>.<instance>")` per instance.

## Errors and logging

```go
import (
    "github.com/xgodev/boost/model/errors"
    "github.com/xgodev/boost/wrapper/log"
)

logger := log.FromContext(ctx).WithTypeOf(*p)
return nil, errors.Wrap(err, errors.Internalf("publish failed"))
```

Match the verbs in the closest existing driver. boost is consistent enough that diverging stands out in review.

## Honest extrapolation — `// TODO(maintainer-review):`

When the new driver has no direct precedent, mark every guess inline:

```go
// TODO(maintainer-review): SNS requires a full ARN. Falling back to options.TopicArn
// when ev.Subject() is not an ARN. Verify this is the wanted behavior, or require
// Subject to always be an ARN.
return p.options.TopicArn
```

Then call out the marked decisions in the PR description so reviewers know where their judgment is needed.

## What NOT to add

- ❌ Code under `bootstrap/` for a wrapper-layer concern (or vice versa). Layering is enforced by directory.
- ❌ A new top-level interface alongside an existing one (`Driver2` is a better `Driver`). Extend the existing one or open an RFC issue first.
- ❌ A direct dependency on a third-party DI / config / log library. Use `wrapper/config`, `wrapper/log`, and `fx` modules.
- ❌ Tests that `os.Setenv` to inject config — use `config.Add` with explicit defaults.
- ❌ A driver that exposes its concrete type as the public return value of `New(...)`. Return the interface.
- ❌ A `Close()` method on the `Driver` interface itself unless ALL existing implementations need it. Optional close = separate optional interface + feature-detect.

## When to escalate to a human

- New top-level layer (peer to `bootstrap/`, `factory/`, `wrapper/`) — architectural change, not a contrib.
- Boost-internal bug that blocks the canonical path (like the ctx-loss in adapter helpers, see `boost-bootstrap-adapter-pubsub`) — open an issue and propose a PR; don't unilaterally redesign the public API.
- Consumer asks for direct dependency on a non-boost framework (gin, fiber, zap-direct) — suggest the boost-equivalent first; only deviate with explicit user sign-off.

## Self-test before opening the PR

```bash
go build ./...
go vet ./...
go test ./...
make v 2>/dev/null || go mod vendor   # if you added a dep
```

- [ ] Files at `<area>/<kind>/contrib/<vendor>/<lib>/v<major>/` (or per-service for multi-service SDKs).
- [ ] Package name is the short library name, not the version leaf.
- [ ] Constructor trio present.
- [ ] `init()` calls `config.Add` for every tunable; `ConfigAdd(path)` exported for multi-instance.
- [ ] Errors via `model/errors`; logger via `wrapper/log.FromContext`.
- [ ] Every extrapolation marked with `// TODO(maintainer-review):` and called out in the PR.
- [ ] Public API change: NONE on existing interfaces.
- [ ] Mirror style matches the closest existing contrib.
