---
name: gqlgen-field-resolvers
description: "Use when building or refactoring a gqlgen (99designs/gqlgen) Go service — schema-first, Federation 2 — with DDD/clean architecture. Covers when a field is a model field vs a field resolver vs a DataLoader, binding a GraphQL type to a domain-carrying model without coupling application/domain to transport, and the gqlgen regen/federation pitfalls that break the build. NOT for github.com/graphql-go/graphql (code-first) — that is boost-factory-graphql, a different library."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

For the boost GraphQL **factory** (`graphql-go/graphql`, code-first) see
`boost-factory-graphql`. This skill is for **gqlgen** (schema-first + Federation).

## Model field vs field resolver vs DataLoader

Do NOT make every field a resolver — that is an anti-pattern (boilerplate). Pick per field:

| The field's value… | Use |
|---|---|
| comes straight from one already-fetched source (the aggregate/orderForm) | **model field / thin mapper** — gqlgen auto-resolves; cheapest |
| is derived/formatted from the parent, no I/O | **field resolver** (lazy: computed only when selected) — or a mapper, equivalent output |
| needs a SEPARATE external call **per item/entity** (N+1 fan-out) | **field resolver + DataLoader** (`graph-gophers/dataloader/v7`), injected per-request via middleware |

A single-source aggregate (one VTEX orderForm → whole Cart) needs **no DataLoader** — there is no N+1 to batch. DataLoader earns its place only on real fan-out (per-item SKU detail, Federation entity batching). Don't cargo-cult it.

## Bind the GraphQL type to a domain-carrying model (keeps DDD)

To get real per-field resolvers without leaking transport types into the core: bind the schema type to an **interface-layer** model that carries the domain entity, and resolve each field from it. Application/domain never import gqltypes.

```go
// internal/interface/graphql/model/cart_output.go
type CartOutput struct{ Source cart.Cart } // domain aggregate
func (CartOutput) IsEntity() {}             // REQUIRED for a Federation @key model
```
```yaml
# gqlgen.yml
models:
  CartOutput:
    model: <module>/internal/interface/graphql/model.CartOutput
    fields:
      items: { resolver: true }   # field resolver → mapper.ItemsToGQL(obj.Source)
      # ... one per derived/IO field
```
Root resolvers return `model.NewCartOutput(domainCart)`; field resolvers delegate to per-field mappers (`mapper.ItemsToGQL(obj.Source)` …). The mapper is the ONLY place gqltypes meet domain.

## Pitfalls that break the build (all hit in production)

| Symptom | Cause → Fix |
|---|---|
| `model.X does not satisfy fedruntime.Entity` | A `@key` type bound to a custom model needs `func (X) IsEntity() {}` |
| `duplicate method LoginType` in generated exec.go | Two schema fields camelize to the same resolver name (`login_type` + `loginType`). `fieldName` does NOT rename a resolver method — bind ONE to a **model method** (`fieldName: LoginTypeLegacy`, no `resolver: true`) and leave the other a resolver |
| `gqlgen generate` fails: `packages.Load … undefined` | The package must **compile/load BEFORE regen**. Don't delete the old mapper or break call sites first. Order: keep it loadable → regen → THEN fix resolver bodies. A half-run also corrupts resolver stubs — `git checkout` the generated + resolver files and retry from a clean, loadable state |
| field resolver reads `obj.SomeField` after rebinding to a domain model | The model has no gqltypes fields — read from `obj.Source` via a mapper |

## Verify

After any resolver refactor: `go build ./...`, `go test ./...`, and confirm the
emitted contract is byte-identical to the reference (introspection + a populated
record diffed field-by-field) — a field-resolver rewrite must not change output.
