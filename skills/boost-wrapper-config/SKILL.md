---
name: boost-wrapper-config
description: "Use when registering or reading configuration in a Go service that imports github.com/xgodev/boost/wrapper/config, or when migrating os.Getenv / viper / envconfig / koanf-direct calls in a boost service. Triggers on imports of github.com/xgodev/boost/wrapper/config, on questions about config.Add, config.String/Int/Bool/Duration, env var override semantics, or when a code review spots os.Getenv outside an init() registration."
license: MIT
metadata:
  author: jpfaria
  version: "0.3.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start` — the registry is empty until `boost.Start()` loads env+files.

## Register every tunable up front

```go
import (
    "time"

    "github.com/xgodev/boost/wrapper/config"
)

const root = "myapp.outbound"

func init() {
    // Pass the default in its NATIVE type — koanf infers the type from the
    // value. Env overrides still arrive as strings and are parsed on read,
    // so a Duration registered as 10*time.Second still accepts MYAPP_..="2h".
    config.Add(root+".subject", "default-topic", "downstream subject")       // string
    config.Add(root+".timeout", 10*time.Second, "publish timeout")           // time.Duration
    config.Add(root+".maxAttempts", 5, "retry budget")                       // int
    config.Add(root+".sampleRate", 0.01, "trace sampling rate")              // float64
    config.Add(root+".enabled", true, "feature toggle")                      // bool
    config.Add(root+".regions", []string{"sa-east-1"}, "allowed regions")    // []string (env: comma-separated)
    config.Add(root+".ports", []int{8080, 9090}, "listen ports")             // []int
    config.Add(root+".labels", map[string]string{"team": "cart"}, "tags")    // map[string]string
}

// later, after boost.Start ran — read with the matching typed getter:
subj := config.String(root + ".subject")
to   := config.Duration(root + ".timeout")
n    := config.Int(root + ".maxAttempts")
regs := config.Strings(root + ".regions")
```

Every `config.Add` call shows up in the boot banner and in `boost-config dump`. `os.Getenv` is invisible to both, so operators can't discover what's tunable.

**Register the default in its native type** (`10*time.Second`, `0.01`, `true`, `[]string{...}`) — not a stringified form (`"10s"`, `"0.01"`). The native literal self-documents the type and matches the getter; the env override still accepts a readable string because env values are always parsed on read.

## API surface

Pick the getter that matches the registered default's type.

| Getter | Returns |
|---|---|
| `config.String(key)` | string |
| `config.Int(key)` | int |
| `config.Int64(key)` | int64 |
| `config.Float64(key)` | float64 |
| `config.Bool(key)` | bool |
| `config.Duration(key)` | time.Duration |
| `config.Time(key, layout)` | time.Time |
| `config.Bytes(key)` | []byte |
| `config.Strings(key)` | []string |
| `config.Ints(key)` | []int |
| `config.Int64s(key)` | []int64 |
| `config.Float64s(key)` | []float64 |
| `config.Bools(key)` | []bool |
| `config.StringMap(key)` | map[string]string |
| `config.IntMap(key)` | map[string]int |
| `config.Int64Map(key)` | map[string]int64 |
| `config.Float64Map(key)` | map[string]float64 |
| `config.BoolMap(key)` | map[string]bool |
| `config.Unmarshal(&v)` / `config.UnmarshalWithPath(key, &v)` | decode a subtree into a struct |
| `config.Exists(key)` / `config.Get(key)` / `config.All()` | presence check / raw value / full dump |

The slice getter is `config.Strings` (plural) — there is no `StringSlice`. For a config subtree that maps onto a struct (nested objects, lists of objects), register the shape in a file/env and read it with `config.UnmarshalWithPath(key, &out)` instead of hand-parsing a JSON string.

## Env override is automatic

A key registered as `myapp.outbound.subject` is overridden at deploy time by env var `MYAPP_OUTBOUND_SUBJECT` (uppercased, dots → underscores). Framework-layer keys live under `boost.*` (e.g., `boost.factory.echo.port` ← `BOOST_FACTORY_ECHO_PORT`). Application-layer keys can use any namespace.

## See every config at boot, hide the secrets

Enable the startup config table (`key | default | resolved value`) — invaluable for spotting a bad override (a wrong `myapp.brand` shows up immediately instead of a silent boot fail-fast):

```
boost.print.config.enabled   = true   # env: BOOST_PRINT_CONFIG_ENABLED=true   (dev / .env only)
boost.print.config.maxLength = 25      # env: BOOST_PRINT_CONFIG_MAXLENGTH      (truncates long values)
```

Mark every secret with `config.WithHide()` so its value renders `****` in that table (the boot printer does `if entry.Options.Hide { v = "****" }`):

```go
config.Add("myapp.vtex.apptoken", "", "VTEX X-VTEX-API-AppToken", config.WithHide())
config.Add("myapp.dolphin.apikey", "", "Dolphin x-api-key",        config.WithHide())
```

`WithHide` only masks the *printed* value — `config.String(key)` still returns the real secret. Leave `boost.print.config.enabled` off (the default) in production.

## Red flags

| Red flag | Fix |
|---|---|
| `os.Getenv("FOO_BAR")` outside a `config.Add` registration | `config.Add("myapp.foo.bar", default, desc)` then `config.String(...)` |
| A secret key (token / apikey / password) registered without `config.WithHide()` | Add `config.WithHide()` — otherwise `boost.print.config.enabled=true` prints the secret into the boot log |
| Reading config inside `init()` (before `boost.Start` runs) | Read at call time, not init time |
| Mutating env vars in tests (`os.Setenv`) | Use `config.Add(...)` with the test default; if you must override at runtime, use a koanf provider in the test setup |
| Hard-coding what should be tunable (timeouts, URLs, retry budgets) | Register with `config.Add` and a sensible default |
| Stringified default for a typed value (`"10s"`, `"0.01"`, `"[]"`) | Pass the native type (`10*time.Second`, `0.01`, `[]string{}`) so it self-documents and matches the getter |
| Reaching for a non-existent `config.StringSlice` | Use `config.Strings` (plural); see the API table for the full set |
| Hand-parsing a JSON string for a struct/map config | `config.UnmarshalWithPath(key, &out)`, or a typed map getter (`config.StringMap`, …) |
| "Just one operator override, harmless" rationalization | Same fix — discoverability and test reproducibility don't allow exceptions |
