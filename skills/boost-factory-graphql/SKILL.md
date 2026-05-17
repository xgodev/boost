---
name: boost-factory-graphql
description: "Use when exposing a GraphQL endpoint in a Go service via github.com/xgodev/boost/factory/contrib/graphql-go/graphql/v0. Covers NewHandler / NewHandlerWithConfig and how the resulting *handler.Handler attaches to an Echo route. Triggers on imports under factory/contrib/graphql-go/graphql/, on questions about GraphQL endpoints in a boost service, or on schema → handler wiring."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-factory-echo` (typical mount point).

```go
import (
    gqlfact "github.com/xgodev/boost/factory/contrib/graphql-go/graphql/v0"
    "github.com/graphql-go/graphql"
)

schema, _ := graphql.NewSchema(graphql.SchemaConfig{ /* ... */ })

h := gqlfact.NewHandler(&schema)

srv.POST("/graphql", echo.WrapHandler(h))
srv.GET("/graphql",  echo.WrapHandler(h))   // for GraphiQL
```

Use `NewHandlerWithConfig` when you need GraphiQL playground, query introspection toggles, or custom request parsers.

## Red flags

| Red flag | Fix |
|---|---|
| `handler.New(...)` from `graphql-go/handler` directly | `gqlfact.NewHandler(&schema)` |
| Schema constructed per request | Build once at startup |
| GraphQL errors leaking 500s through the Echo error_handler | Map them explicitly via the GraphQL formatter; don't return them as Go errors |
