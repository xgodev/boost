package inject

import (
	"go/ast"
	"go/parser"
	"path"
	"sort"
	"strings"

	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
)

// Implementation is an implementation set: two or more functions providing
// the same non-group value, one of them marked @Default. The generator
// writes a switch module that includes only the member a config key names.
type Implementation struct {
	Type    string               // name of the provided type ("Service")
	Path    string               // import path of the package declaring the type
	Name    string               // fx name of the value, "" for the default value
	Default string               // package name of the @Default member
	Members map[string]Component // by package name: the implementation's name
}

// Implementations returns the implementation sets of the graph, sorted by
// type and name. It is NotValid when a set has more than one @Default, a
// @Default value has no other provider, a member does not have exactly one
// @Provide, two members share a package name, the members do not provide the
// same named type or a consumer of the set lives in a member's package.
func Implementations(g *graph.Graph[Component]) ([]Implementation, error) {
	providers := map[string][]Component{}
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for _, idx := range sortedIndexes(comp.Provides) {
			ann := comp.Provides[idx]
			if ann.Group != "" {
				continue
			}
			id := xid(comp.Entry.Package, comp.Entry.Func.Results[idx].Type, ann)
			providers[id] = append(providers[id], comp)
		}
	}

	ids := make([]string, 0, len(providers))
	for id := range providers {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	sets := map[string]Implementation{} // by value id
	for _, id := range ids {
		set, ok, err := newImplementation(id, providers[id])
		if err != nil {
			return nil, err
		}
		if ok {
			sets[id] = set
		}
	}

	if err := checkConsumers(g, sets); err != nil {
		return nil, err
	}

	list := make([]Implementation, 0, len(sets))
	for _, set := range sets {
		list = append(list, set)
	}
	sort.Slice(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		return a.Name < b.Name
	})
	return list, nil
}

// ImplementationOf returns the implementation set comp is a member of.
func ImplementationOf(g *graph.Graph[Component], comp Component) (Implementation, bool) {
	sets, err := Implementations(g)
	if err != nil {
		return Implementation{}, false
	}
	for _, set := range sets {
		if member, ok := set.Members[comp.Entry.Package]; ok && member.Key() == comp.Key() {
			return set, true
		}
	}
	return Implementation{}, false
}

// newImplementation returns the implementation set of the providers of the
// value id, or false when none of them is @Default.
func newImplementation(id string, comps []Component) (Implementation, bool, error) {
	var defaults []string
	for _, comp := range comps {
		if comp.Default {
			defaults = append(defaults, comp.Key())
		}
	}
	switch {
	case len(defaults) == 0:
		return Implementation{}, false, nil
	case len(defaults) > 1:
		return Implementation{}, false, errors.NotValidf("%s has more than one @Default: %s; mark only the implementation wired when the config key is unset",
			id, strings.Join(defaults, ", "))
	case len(comps) == 1:
		return Implementation{}, false, errors.NotValidf("@Default on %s: %s has no other implementation; @Default picks among two or more @Provide of the same value",
			defaults[0], id)
	}

	set := Implementation{Members: map[string]Component{}}
	for _, comp := range comps {
		if len(comp.Provides) != 1 {
			return Implementation{}, false, errors.NotValidf("%s implements %s and must have exactly one @Provide: only the chosen implementation's module is wired", comp.Key(), id)
		}
		var idx int
		for i, ann := range comp.Provides {
			idx, set.Name = i, ann.Name
		}

		tp := comp.Entry.Func.Results[idx].Type
		typePath, typeName, ok := namedType(comp, tp)
		if !ok {
			return Implementation{}, false, errors.NotValidf("%s implements %s with %s; an implementation set needs a named type (T, *T, pkg.T or *pkg.T) whose package can be resolved",
				comp.Key(), id, tp)
		}
		if set.Path != "" && (set.Path != typePath || set.Type != typeName) {
			return Implementation{}, false, errors.NotValidf("the implementations of %s provide different types: %s.%s and %s.%s",
				id, set.Path, set.Type, typePath, typeName)
		}
		set.Path, set.Type = typePath, typeName

		name := comp.Entry.Package
		if other, dup := set.Members[name]; dup {
			return Implementation{}, false, errors.NotValidf("%s and %s implement %s in packages named %s; the implementation is named after its package, rename one",
				other.Key(), comp.Key(), id, name)
		}
		set.Members[name] = comp
		if comp.Default {
			set.Default = name
		}
	}
	return set, true, nil
}

// namedType resolves the import path and name of a result type written in
// the file of comp: T, *T, pkg.T or *pkg.T.
func namedType(comp Component, tp string) (string, string, bool) {
	expr, err := parser.ParseExpr(tp)
	if err != nil {
		return "", "", false
	}
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}

	switch e := expr.(type) {
	case *ast.Ident:
		if isPredeclaredType(e.Name) {
			return "", "", false
		}
		return comp.Entry.Path, e.Name, true
	case *ast.SelectorExpr:
		pkg, ok := e.X.(*ast.Ident)
		if !ok {
			return "", "", false
		}
		p := comp.Entry.Imports[pkg.Name]
		return p, e.Sel.Name, p != ""
	}
	return "", "", false
}

// checkConsumers rejects a consumer of an implementation set declared in the
// package of one of its members: the generated switch imports the members'
// modules, so that package's module would import its own switch.
func checkConsumers(g *graph.Graph[Component], sets map[string]Implementation) error {
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for i, param := range comp.Entry.Func.Parameters {
			set, ok := sets[injectID(comp.Entry.Package, param.Type, comp.Injects[i])]
			if !ok {
				continue
			}
			for _, name := range set.Names() {
				if member := set.Members[name]; member.Entry.Path == comp.Entry.Path && member.Key() != key {
					return errors.NotValidf("%s consumes %s.%s, implemented by %s in the same package; the generated switch would import itself, move %s to another package",
						key, set.Path, set.Type, member.Key(), comp.Entry.Func.Name)
				}
			}
		}
	}
	return nil
}

// Names returns the names of the implementations, sorted.
func (i Implementation) Names() []string {
	names := make([]string, 0, len(i.Members))
	for name := range i.Members {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Ident is the exported identifier of the set in generated code: the type
// name, plus the value name when the value is named ("Service",
// "ServiceFormal").
func (i Implementation) Ident() string {
	ident := upperFirst(i.Type)
	if i.Name != "" {
		ident += upperFirst(packageIdent(i.Name))
	}
	return ident
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Key returns the config key that chooses the implementation, given the Go
// module path of the application: the type's package path relative to the
// module (its full import path outside the module), plus the value name,
// plus "implementation", in dots. Every segment is lowercased, "_" splits
// segments and "a-b" becomes "aB", so the key is reachable from an
// environment variable (FOO_BAR__BAZ -> foo.barBaz).
func (i Implementation) Key(modulePath string) string {
	rel := i.Path
	switch {
	case i.Path == modulePath:
		rel = path.Base(modulePath)
	case strings.HasPrefix(i.Path, modulePath+"/"):
		rel = strings.TrimPrefix(i.Path, modulePath+"/")
	}

	segments := keySegments(rel)
	if i.Name != "" {
		segments = append(segments, keySegments(i.Name)...)
	}
	return strings.Join(append(segments, "implementation"), ".")
}

// keySegments splits s into config key segments.
func keySegments(s string) []string {
	var segments []string
	for _, seg := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return r == '/' || r == '.' || r == '_' }) {
		words := strings.Split(seg, "-")
		for j := 1; j < len(words); j++ {
			if words[j] != "" {
				words[j] = strings.ToUpper(words[j][:1]) + words[j][1:]
			}
		}
		segments = append(segments, strings.Join(words, ""))
	}
	return segments
}

// File returns the path, relative to OutputDir, of the switch file.
func (i Implementation) File() string {
	return path.Join("implementation", i.Path, strings.ToLower(i.Ident())+"_implementation.go")
}

// Import returns the import of the generated package holding the switch,
// given the import path of the output dir (BaseImport).
func (i Implementation) Import(base string) ImportData {
	p := path.Join(base, "implementation", i.Path)
	return ImportData{Alias: Alias("impl", path.Base(i.Path), p), Path: p}
}

// Module returns the import of the switch package and the call of the
// module that includes the chosen implementation.
func (i Implementation) Module(base string) (ImportData, string) {
	imp := i.Import(base)
	return imp, imp.Alias + "." + i.Ident() + "ImplementationModule()"
}

// Chosen returns the import of the switch package and the call returning
// the name of the chosen implementation.
func (i Implementation) Chosen(base string) (ImportData, string) {
	imp := i.Import(base)
	return imp, imp.Alias + "." + i.Ident() + "Implementation()"
}
