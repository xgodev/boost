# inject: config registration from `@Config` structs -- design

Issue #80. Piece 3 of the builder boilerplate roadmap
(`docs/builder-boilerplate/README.md`), after `@Middleware` (#72) and the
implementation switch (#75).

## Intent

A package that needs settings writes a `config.go` by hand today: key
constants, `config.Add` per key in `init()`, an `Options` struct and a
`NewOptions()` that calls `config.NewOptionsWithPath` (the
`factory/contrib/*/config.go` pattern, claude-plugin Go LAW 25). The
maintainer wants that boilerplate generated: the package declares only the
struct, the generator writes the registration and provides the filled
struct. The package itself never imports `wrapper/config`.

## Decided (by the maintainer)

```go
// package resty

// @Config
type Options struct {
    BaseURL string        `default:"http://localhost" desc:"quote API base URL"`
    Timeout time.Duration `default:"2s" desc:"call timeout"`
}

// @Provide(index=0)
func New(c *resty.Client, o *Options) quote.Service { ... }
```

- `@Config` on a struct; each field carries `default` and `desc` tags.
- The generator writes the `config.Add` calls and provides `*Options` to any
  constructor that receives it.
- The generated file lives in `gen/`; the package does not import `config`.
- Tag defaults cover simple types only; anything else fails generation with
  a clear error (no silent skip).

## Behaviour

- **The struct.** `@Config` takes no arguments and goes on an exported,
  non-generic struct outside package `main`. At most one `@Config` per
  package. `@Config` on a function is NotValid.
- **Fields.** Every exported field is a key, unless tagged `config:"-"`.
  Supported types: `string`, `bool`, `int`, `int8`..`int64`, `uint`,
  `uint8`..`uint64`, `float32`, `float64`, `time.Duration`, `[]string`.
  Anything else (a struct, a map, a pointer, a named type) is NotValid,
  naming the field and the type. Embedded fields are NotValid. Unexported
  fields are ignored unless they carry `default`/`desc`/`hide`/`config`,
  which is NotValid.
- **Tags.**
  - `desc:"..."` required: every key is documented.
  - `default:"..."` optional; missing means the zero value. It is parsed at
    generation time for the field type: `strconv` rules for bool and
    numbers (with the bit size), `time.ParseDuration` for durations, a comma
    list for `[]string` (items trimmed, `""` is the empty list). A default
    that does not parse is NotValid.
  - `hide:"true"` registers the key with `config.WithHide()` (secrets).
  - `config:"name"` names the key segment; it must be lowerCamel
    (`^[a-z][a-zA-Z0-9]*$`) so an env var reaches it.
- **Keys.** Root: the package import path relative to the Go module, in
  dots, with the same segment rules as the implementation key (`_` splits,
  `a-b` -> `aB`, lowercased; full import path outside the module). Field
  segment: the field name split into words (`BaseURL` -> base, url;
  `SoldOut` -> sold, out) in lowerCamel (`baseUrl`, `soldOut`), which is what
  the env loader produces from `BASE__URL`. Example: `internal/pkg/client/quote/resty`
  + `BaseURL` -> `internal.pkg.client.quote.resty.baseUrl`, env
  `INTERNAL_PKG_CLIENT_QUOTE_RESTY_BASE__URL`.
- **Collisions.** Two keys equal (between configs, or a config key and an
  implementation key) is AlreadyExists. A key that is the parent of another
  (`a.b` and `a.b.c`) is AlreadyExists: koanf cannot hold both a value and a
  map there.
- **Generated file.** `gen/<import path>/<lower struct>_module.go`:
  `init()` with one `config.Add(root+".<segment>", <typed default>, desc[,
  config.WithHide()])` per key, each commented with its env var, and
  `<Struct>Module()` (once-guarded like every module) providing
  `func() (*src.Struct, error)` that returns
  `config.NewOptionsWithPath[src.Struct](root)`. Defaults are Go literals of
  the field type (`2 * time.Second`, `int8(5)`, `[]string{"a", "b"}`).
- **Wiring.** The struct is a provider component of `*Struct` in the inject
  graph: constructors receiving `*Struct` (same package or `*pkg.Struct`)
  get an edge, their modules include `<Struct>Module()`, the root
  `gen.Module()` includes it, builder apps reach it through their doors.
  A function that also provides `*Struct` is the usual AlreadyExists.
- Builder and inject collect `@Config` (filter `Config`).

## Inside

- `annotation.Entry` gains `Fields []EntryField{Name, Type, Tag}` for struct
  entries (one per name; embedded field: empty `Name`).
- `inject.Component` gains `Config []ConfigField`, set for a `@Config`
  struct; `Component.Result(idx)` returns the provided type (`*Struct` for a
  config, the function result otherwise) and replaces the direct
  `Entry.Func.Results[idx]` reads in inject and builder. `gid` and
  `ModuleOf` name a struct component by its struct name.
- A `config` template next to `module`; `Render` picks it for config
  components and checks the key collisions with the implementation keys.

## Example

`builder/examples/simple/stock/memory` gets `@Config Options{SoldOut
[]string}` (default `sold-out`), received by `memory.New`. The e2e test
boots app api in a child process with
`BUILDER_EXAMPLES_SIMPLE_STOCK_MEMORY_SOLD__OUT=book,pen`: `book` and `pen`
are refused (409) and `sold-out` is accepted (201). Without the env the current
behaviour stays.

## Non-goals

- Nested structs, maps, custom types and `TextUnmarshaler` fields.
- Several `@Config` per package; a config per implementation name.
- Validation of values beyond type (ranges, required non-empty).
- Rewriting the hand-written `factory/contrib` configs.

## Left to the builder

Template wording, error message wording, test names, the word-split
details beyond the examples above.
