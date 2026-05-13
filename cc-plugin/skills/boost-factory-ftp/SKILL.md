---
name: boost-factory-ftp
description: "Use when establishing an FTP server connection in a Go service via github.com/xgodev/boost/factory/contrib/jlaffaye/ftp/v0. Covers NewServerConn + variants and the typical use case (legacy file-drop integrations). Triggers on imports under factory/contrib/jlaffaye/ftp/, on questions about FTP in a boost service, or on connecting to legacy partners."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`.

```go
import ftpfact "github.com/xgodev/boost/factory/contrib/jlaffaye/ftp/v0"

conn, err := ftpfact.NewServerConn(ctx)
if err != nil { log.Fatalf("ftp: %v", err) }
defer conn.Quit()
```

Configure host, port, user, password, timeout under `boost.factory.ftp.*` (override `BOOST_FACTORY_FTP_*`).

Multi-target pattern: `ftpfact.ConfigAdd("boost.factory.ftp.<partner>")` per partner FTP, then `NewServerConnWithConfigPath(ctx, "boost.factory.ftp.<partner>")`.

## Red flags

| Red flag | Fix |
|---|---|
| `ftp.Dial(addr)` directly | `ftpfact.NewServerConn(ctx)` |
| Credentials via `os.Getenv` | `BOOST_FACTORY_FTP_*` |
| Forgetting `defer conn.Quit()` | Add it |
| Sharing a single connection across goroutines | jlaffaye/ftp connections are NOT goroutine-safe — one per worker |
