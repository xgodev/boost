# Repository Guidelines

## Project Structure & Module Organization
- `bootstrap/` wires application startup and serverless entrypoints; use it as the glue layer when exposing services.
- `factory/` and `wrapper/` provide pluggable integrations (loggers, caches, messaging); extend these packages when adding new providers.
- `fx/` handles dependency injection and lifecycle; modules should export constructors compatible with this package.
- `extra/` hosts cross-cutting utilities (health checks, middleware, multiserver) while `model/` centralizes shared contracts.
- `examples/` demonstrates canonical setups, and `vendor/` locks dependencies after `make v`.

## Build, Test, and Development Commands
- `go test ./...` runs package tests directly; prefer `make test` to replicate CI’s `go test all` invocation.
- `go run ./start.go` boots the default sample service for local smoke-testing.
- `make v` tidies modules and refreshes the vendor tree; always run it when dependencies change.
- `make upgrade-deps` (or `make check-upgrade-deps`) audits module updates before adopting new versions.

## Coding Style & Naming Conventions
- Format Go code with `gofmt` (implied by `go fmt ./...`); commit only formatted sources.
- Keep packages focused; exported constructors follow the `NewThing` pattern, while interfaces live under `model/`.
- Use CamelCase for exported identifiers and snake_case for environment variables (e.g., `BOOST_FACTORY_ZAP_CONSOLE_LEVEL`).
- Document public types with concise GoDoc comments when behavior is non-trivial.

## Testing Guidelines
- Place tests alongside code using `*_test.go`; name cases `TestPackage_FeatureScenario` to mirror the component under test.
- Favor table-driven tests and cover both happy-path and failure handling, especially around wrappers and factories.
- When adding integrations, include example tests under `examples/` to show configuration and lifecycle usage.

## Commit & Pull Request Guidelines
- Follow the existing short, imperative subject style (e.g., “Adiciona timeout para envio”); group related changes per commit.
- Reference issues or tickets in the PR description, summarize behavior changes, and highlight any new env vars or CLI flags.
- Ensure PRs include reproduction or verification steps (commands, configs) so reviewers can validate factories and wrappers quickly.
- Attach screenshots or logs when altering observability outputs, and request review from maintainers owning the touched module.

## Configuration Tips
- Centralize default values in `config.go` and surface overrides via environment variables for consistency across modules.
- Use sample env blocks from `examples/` to validate new providers before wiring them into `bootstrap/`.
