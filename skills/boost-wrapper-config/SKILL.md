---
name: boost-wrapper-config
description: "Use when registering or reading configuration in a Go service that imports github.com/xgodev/boost/wrapper/config, or when migrating os.Getenv / viper / envconfig / koanf-direct calls in a boost service. Triggers on imports of github.com/xgodev/boost/wrapper/config, on questions about config.Add, config.String/Int/Bool/Duration, env var override semantics, or when a code review spots os.Getenv outside an init() registration."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start` — the registry is empty until `boost.Start()` loads env+files.

## Register every tunable up front

```go
import "github.com/xgodev/boost/wrapper/config"

const root = "myapp.outbound"

func init() {
    config.Add(root+".subject", "default-topic", "downstream subject")
    config.Add(root+".timeout", "10s", "publish timeout")
    config.Add(root+".maxAttempts", 5, "retry budget")
}

// later, after boost.Start ran
subj := config.String(root + ".subject")
to   := config.Duration(root + ".timeout")
n    := config.Int(root + ".maxAttempts")
```

Every `config.Add` call shows up in the boot banner and in `boost-config dump`. `os.Getenv` is invisible to both, so operators can't discover what's tunable.

## API surface

| API | Returns |
|---|---|
| `config.String(key)` | string |
| `config.Int(key)` | int |
| `config.Bool(key)` | bool |
| `config.Duration(key)` | time.Duration |
| `config.Float64(key)` | float64 |
| `config.StringSlice(key)` | []string |

## Env override is automatic

A key registered as `myapp.outbound.subject` is overridden at deploy time by env var `MYAPP_OUTBOUND_SUBJECT` (uppercased, dots → underscores). Framework-layer keys live under `boost.*` (e.g., `boost.factory.echo.port` ← `BOOST_FACTORY_ECHO_PORT`). Application-layer keys can use any namespace.

## Red flags

| Red flag | Fix |
|---|---|
| `os.Getenv("FOO_BAR")` outside a `config.Add` registration | `config.Add("myapp.foo.bar", default, desc)` then `config.String(...)` |
| Reading config inside `init()` (before `boost.Start` runs) | Read at call time, not init time |
| Mutating env vars in tests (`os.Setenv`) | Use `config.Add(...)` with the test default; if you must override at runtime, use a koanf provider in the test setup |
| Hard-coding what should be tunable (timeouts, URLs, retry budgets) | Register with `config.Add` and a sensible default |
| "Just one operator override, harmless" rationalization | Same fix — discoverability and test reproducibility don't allow exceptions |
