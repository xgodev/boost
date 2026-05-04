---
name: boost-extra-middleware
description: "Generic middleware abstractions at github.com/xgodev/boost/extra/middleware — AnyErrorMiddleware[T] interface, NewAnyErrorWrapper[T] composition. Used by every function middleware and by the production graceful-shutdown workaround. Apply when composing middleware chains outside the boost-bootstrap-middleware presets. STATUS: stub — not yet lapidated via the TDD cycle. Treat as a placeholder routing to boost-core until first concrete usage forces a RED-GREEN-REFACTOR pass."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents working in Go projects on github.com/xgodev/boost.
metadata:
  author: jpfaria
  version: "0.0.1"
  status: stub
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

> **STATUS: STUB.** This skill ships as a placeholder. Until the first feature lands that exercises this subsystem in production, follow `boost-core` (Iron Laws + cross-cutting principles) and read the existing boost source under the corresponding directory directly.
>
> **When this skill is lapidated:** the contributor adding the feature MUST run a TDD cycle (RED baseline → GREEN write skill → REFACTOR plug loopholes) per the `superpowers:writing-skills` discipline, and update version + status here. See `cc-plugins/CONTRIBUTING.md`.

## Scope (planned)

Generic middleware abstractions at github.com/xgodev/boost/extra/middleware — AnyErrorMiddleware[T] interface, NewAnyErrorWrapper[T] composition. Used by every function middleware and by the production graceful-shutdown workaround. Apply when composing middleware chains outside the boost-bootstrap-middleware presets.

## Triggers (planned)

- TBD (added during the GREEN phase based on observed baseline failures)

## Required background

- `boost-core` (Iron Laws + boot/log/config/errors)

## Reference paths in boost source

- TBD (add the canonical source files this skill should cite)

## Lapidation checklist

- [ ] Run RED baseline — at least one realistic scenario without this skill loaded
- [ ] Document baseline failures verbatim
- [ ] Write minimal SKILL.md addressing those specific failures
- [ ] Run GREEN — verify agents now comply
- [ ] REFACTOR — close loopholes from new rationalizations
- [ ] Bump `metadata.version` to `0.1.0` and `status` to `mature`
- [ ] Update this plugin's `plugin.json` version to match
- [ ] Add evals/ directory with `evals.json` if any pressure scenarios remain
