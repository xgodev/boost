# builder: generate every piece of boilerplate

Status (2026-10-18): piece 1 (`@Middleware`) implemented in #72 (PR #73,
stacked on #69), skill updated in claude-plugin 1.24.3. Next: piece 2.

## Intent

The builder already generates the inject modules, the doors and the `main` of
every app. The goal is that a service written on boost has **no boilerplate
left by hand**: the application writes handlers, services, datastores and
annotations; everything that repeats between services is generated. The
`claude-plugin` Go skill (`skills/dev/golang/`) changes with it, so the layout
it teaches stops asking for hand-written wiring.

## Boilerplate the claude-plugin layout asks for today

1. **Middleware around a call** (`extra/middleware`: log, cache, hystrix,
   datadog, prometheus, fallback). Today each method builds the wrapper and
   wraps its own body in `wrapper.Exec(ctx, id, func…)`. Direction: an
   annotation on the service; the builder generates a decorator of the
   service interface that runs the chain, wired with `fx.Decorate`, so the
   method body stays clean.
2. **Implementation switch** (`internal/pkg/fx/modules/<unit>/module.go`:
   `sync.Once`, `module(impl)`, `switch` on a config key, boot failure on an
   unknown value) when a contract has two or more implementations.
3. **`config.go`**: `config.Add(key, default, description)` plus accessors.
4. **Factory modules with the app's plugins**
   (`internal/pkg/fx/modules/lib/boost/factory/<factory>/module.go`).
5. **claude-plugin skill**: replace the hand-written pieces above with the
   annotations.

## Decided

- One piece at a time in the order above, the skill updated with each.
- Middleware: the list lives in the annotation (`@Middleware(log, ...)`),
  config only turns each plugin off. Design:
  `docs/superpowers/specs/2026-10-18-builder-middleware-design.md`.

## Open

- Annotation syntax of pieces 2-4.
