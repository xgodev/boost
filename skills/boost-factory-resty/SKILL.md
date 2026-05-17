---
name: boost-factory-resty
description: "Use when constructing outbound HTTP clients in a Go service that imports github.com/xgodev/boost/factory/contrib/go-resty/resty/v2. Covers per-target config namespacing via resty.ConfigAdd, NewClientWithConfigPath construction, and how each upstream HTTP dependency lives under its own boost.factory.resty.<target> root so timeouts, base URLs, retry policy, and auth headers are tunable per environment. Triggers on imports under factory/contrib/go-resty/, on questions about resty.ConfigAdd, NewClientWithConfigPath, or wiring multiple HTTP clients in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-start` — `boost.Start()` first.
- `boost-wrapper-config` — config namespacing semantics.

## Construction — one client per upstream

```go
import "github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"

func init() {
    resty.ConfigAdd("boost.factory.resty.customer")
    resty.ConfigAdd("boost.factory.resty.janis")
}

func main() {
    boost.Start()
    ctx := context.Background()

    customerHTTP, err := resty.NewClientWithConfigPath(ctx, "boost.factory.resty.customer")
    if err != nil { log.Fatalf("customer http: %v", err) }

    janisHTTP, err := resty.NewClientWithConfigPath(ctx, "boost.factory.resty.janis")
    if err != nil { log.Fatalf("janis http: %v", err) }

    // ... pass into your domain services ...
}
```

Each upstream gets its own config root, registered in `init()` via `resty.ConfigAdd("boost.factory.resty.<target>")`. The boot banner then enumerates timeouts, retries, base URL, etc., per target. Operators override via `BOOST_FACTORY_RESTY_CUSTOMER_*` / `BOOST_FACTORY_RESTY_JANIS_*`.

## Per-target tunables (under each `boost.factory.resty.<target>` root)

Standard knobs registered automatically:

| Key suffix | What |
|---|---|
| `.baseURL` | Base URL for every request |
| `.timeout` | Per-request timeout (`time.Duration`) |
| `.retryCount` | Retry budget |
| `.retryWaitTime` | Base backoff between retries |

Application-specific extras (auth header, client ID, etc.) typically live in your own service config, not `boost.factory.resty.*` — pass them explicitly to your client constructor.

## Red flags

| Red flag | Fix |
|---|---|
| `resty.New()` directly from the upstream SDK | Use `resty.NewClientWithConfigPath(ctx, "boost.factory.resty.<target>")` |
| Single `resty.NewClient` shared across multiple upstreams | One client per upstream, each with its own config root |
| Config root not registered via `resty.ConfigAdd` in `init()` | Register so the boot banner discovers the target |
| Hard-coded base URL or timeout in the client constructor | Tune via `boost.factory.resty.<target>.baseURL` / `.timeout` |
| Reading API keys / auth headers via `os.Getenv` | Register them as `myapp.<target>.apiKey` via `config.Add` (see `boost-wrapper-config`) |
