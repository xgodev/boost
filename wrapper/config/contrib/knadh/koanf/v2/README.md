# koanf v2 config driver

`config.Provider` for [`wrapper/config`](../../../../) built on
[koanf v2](https://github.com/knadh/koanf). It sits next to the
[v1 driver](../v1) and loads the same sources in the same order, so the keys,
flags and env vars of an app do not change.

`boost.Start()` still wires the v1 driver; choosing the default is tracked in
[#132](https://github.com/xgodev/boost/issues/132). To use v2 today, set it
before loading:

```go
import (
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v2"
)

func main() {
	config.Set(koanf.New())
	config.Load()

	host := config.String("app.server.host")
	_ = host
}
```

## Sources

`Load` merges, each one overriding the previous:

1. map defaults declared with `config.Add` (`map[string]string`,
   `map[string]int`, `map[string]int64`), loaded first so a later source that
   sets one entry keeps the others;
2. config files from the `CONF` env var (comma separated) or, when it is
   empty, the `--conf` flag;
3. environment variables;
4. command line flags. A flag left unset only fills a key no earlier source
   set, so its default never hides a file or env value.

Every key passed to `config.Add` becomes a flag typed by its default value.

### Environment names

`_` splits the key path and `__` joins camel case words:

| Env var | Key |
|---|---|
| `APP_SERVER_HOST` | `app.server.host` |
| `APP_MAX__IDLE__CONNS` | `app.maxIdleConns` |

### File formats

| Extension | Parser |
|---|---|
| `.json` | `encoding/json` |
| `.yaml`, `.yml` | `go.yaml.in/yaml/v3` |
| `.toml` | `github.com/knadh/koanf/parsers/toml/v2` |

Any other extension, or a source that cannot be read or parsed, makes `Load`
panic, as in v1.

## Why some adapters are local

The core (`github.com/knadh/koanf/v2`), the env provider
(`providers/env/v2`) and the TOML parser (`parsers/toml/v2`) are koanf's own
modules. koanf's flag (`posflag`), file, JSON and YAML modules use the same
import paths as directories inside the `github.com/knadh/koanf` v1 module,
which boost still requires for the v1 driver; Go rejects both in one build
("ambiguous import"). This package carries small equivalents instead:
`flags.go` (same types as `posflag`), `file.go`, `json.go` and `yaml.go`. When
v1 is removed they can be swapped for koanf's modules.

## Differences from v1

- `map[string]int` and `map[string]int64` read back through `config.IntMap`
  and `config.Int64Map`; v1 returns them empty.
- TOML is parsed by go-toml v2 (v1 uses go-toml v1).
