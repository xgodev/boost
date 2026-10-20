# Annotation

The Annotation package reads annotations from Go doc comments and returns the
annotated declarations, so code generators (such as `inject`) can build on them.

## Format

An annotation is a doc comment line that starts with `// @`:

```go
// FooFunc handles the request.
// It is called by the router.
// @RestRouter(path=/foo, method=GET)
// @RestQueryParam(name=tags, type=string, description="comma, separated")
// @Invoke
func FooFunc(ctx context.Context) error { ... }
```

- `// @Name` or `// @Name(key=value, ...)`. A space between the name and the
  parentheses is allowed; `//@Name` (no space after `//`) is not an annotation.
- Every value is a `key=value` pair or a bare name (a letter or `_`, then
  letters, digits, `_`, `.` or `-`), as in `// @Middleware(log, hystrix)`. Bare
  names stay in `Value`/`RawValue()` and are left out of `Map`. Words after the
  name without parentheses (`// @Param query foo bool`) or with spaces inside a
  value without `=` are not a valid annotation and are ignored with a warning.
- Values are split on top-level commas. Commas inside double quotes or brackets
  belong to the value, and only the first `=` separates the key
  (`filter=a=b` is `filter: "a=b"`).
- Unquoted values become `int`, `float64` or `bool` (`true`/`false`) when they
  parse as such; quoted values are always strings.

## Usage

```go
collector, err := annotation.Collect(
	annotation.WithPath("."),                          // required: a directory inside a Go module
	annotation.WithFilters("Rest", "Inject"),          // optional: annotation name prefixes
	annotation.WithPackages("github.com/acme/service"), // optional: import paths to scan
)
if err != nil {
	return err
}

for _, entry := range collector.EntriesWith("RestRouter") {
	for _, ann := range entry.Annotations {
		var router struct {
			Path   string `attr:"path"`
			Method string `attr:"method"`
		}
		if err := ann.Decode(&router); err != nil {
			return err
		}
	}
}
```

### Options

| Option | Behavior |
|---|---|
| `WithPath(dir)` | Scans the packages under `dir` as `go list ./...` sees them from there: hidden (`.`), `_`, `testdata` and `vendor` directories are skipped. The packages they import are scanned too. |
| `WithFilters(names...)` | Keeps only annotations whose name starts with one of the filters: `Rest` keeps `@Rest`, `@RestRouter` and `@RestResponse`. Without it, every annotation is kept. |
| `WithPackageDocs()` | Also collects the annotations of package doc comments: one entry per file whose package doc has one. Without it, only declarations are collected. |
| `WithPackages(paths...)` | Scans only these import paths and the packages below them, both under the path and among their imports. The imports of a package that is not scanned are not followed. Without it, every non-standard-library package reached is scanned. |

Standard library packages are never scanned.

### Entries

Each `Entry` is an annotated struct type, method, function or package doc:

- `IsStruct()`: a struct type. Interfaces and other non-struct types are not reported.
  `Fields` lists its fields in order, one per name (`A, B string` is two), with
  the type as written and the tag without the backquotes; an embedded field has
  an empty `Name`.
- `IsMethod()`: a method, with `Struct` set to its receiver type (`T`, `*T` or a
  generic `T[K]`), wherever in the package the type is declared.
- `IsFunc()`: a package-level function.
- `IsPackage()`: the package doc comment of a file (`// Package x ...` above
  `package x`), only with `WithPackageDocs()`. `Struct` and `Func` are empty and
  `Header.Title` drops the leading `Package x`.

`Path` is the import path, `Package` the package name, `Module` the module path
and `File` the base name of the file that declares it. `Header.Title` is the first
doc line without the leading declaration name and `Header.Description` the doc
lines after it; annotation lines are not part of the header. `Annotation.Value`
(`RawValue()`) keeps the text between the parentheses as written.
`Imports` maps the name the declaring file uses for each import (its alias, or
the imported package's name, so `github.com/labstack/echo/v4` is `echo`) to the
import path, which resolves a type such as `echo.Context` written in the file;
blank and dot imports are left out.

Entries are returned per package (sorted by import path), in source order.

## Examples

Run from this directory:

```bash
go run ./examples/func
```

- `examples/func`: a function with several `@Rest*` annotations.
- `examples/method`: struct types and methods.
- `examples/decode`: decoding an annotation into a struct with `attr` tags.
- `examples/simple`: annotations found through an imported package.
