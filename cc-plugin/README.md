# golang-boost — Claude Code plugin

A single Claude Code skill that teaches AI agents how to write and review Go code on top of [github.com/xgodev/boost](https://github.com/xgodev/boost).

## Install

```
/plugin marketplace add xgodev/boost
/plugin install golang-boost@xgodev
```

## Update

```
/plugin update golang-boost@xgodev
```

## What it covers

- The three Iron Laws every boost service obeys (`boost.Start`, handler typing, config layer)
- Structured logging via `wrapper/log.FromContext`
- Configuration via `wrapper/config` (and why `os.Getenv` is a violation)
- The `model/errors` type system + how Echo's `error_handler` and the function `publisher` deadletter middleware route on it
- Canonical `main.go` shapes for HTTP APIs (Echo) and Pub/Sub functions
- The documented production workaround for the ctx-loss in adapter helpers (with the required `// TODO(boost-upstream):` annotation)
- Maintainer-side layout convention for new factory contribs, wrapper drivers, bootstrap adapters, and Echo plugins
- Red flags + self-test checklist

## Discovery

Claude Code activates the skill automatically when your session has a Go file that imports `github.com/xgodev/boost` or you ask about its APIs. You don't need to invoke it manually.

## License

MIT — see [LICENSE](./LICENSE).
