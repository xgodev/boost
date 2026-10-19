package inject

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// NewGraphFromEntries builds the dependency graph of the annotated functions.
//
// Every function with @Provide or @Invoke becomes a vertex. An edge goes from
// a provider to each function that consumes what it provides: parameters with
// @Inject must have a provider in the graph, while parameters without it are
// linked when a default provider of the same type exists and are otherwise
// left for the application to provide. An @Inject client parameter
// (*resty.Client, ...) without a provider is a client inject builds (see
// Clients).
//
// Two or more providers of the same non-group value form an implementation
// set when one of them is marked @Default (see Implementations).
//
// A struct annotated @Config is a vertex too: it provides *Struct, filled
// from config.
//
// It returns a NotValid error for malformed annotations, a dependency cycle,
// two @Config in a package, an invalid implementation set or invalid client
// plugins, AlreadyExists when two functions provide the same non-group value
// and none is @Default or two clients share a file, and NotFound when an
// @Inject has no provider.
func NewGraphFromEntries(ctx context.Context, entries []annotation.Entry) (*graph.Graph[Component], error) {
	g := graph.NewGraph[Component]()
	entries = withImports(entries)

	providers, err := addComponents(g, entries)
	if err != nil {
		return nil, err
	}
	if err := checkConfigs(g); err != nil {
		return nil, err
	}

	if err := linkComponents(ctx, g, providers); err != nil {
		return nil, err
	}

	if _, remaining, err := g.TopologicalSort(); err != nil {
		return nil, errors.NotValidf("dependency cycle among %s", strings.Join(remaining, ", "))
	}

	if _, err := Implementations(g); err != nil {
		return nil, err
	}
	if _, err := Clients(g); err != nil {
		return nil, err
	}

	log.FromContext(ctx).Debugf("inject graph built with %d components", len(g.Vertices))

	return g, nil
}

// addComponents adds a vertex per annotated function and returns the keys of
// the components providing each value id.
func addComponents(g *graph.Graph[Component], entries []annotation.Entry) (map[string][]string, error) {
	providers := map[string][]string{}
	groups := map[string]bool{}

	for _, entry := range entries {
		comp, ok, err := newComponent(entry)
		if err != nil {
			return nil, err
		}

		key := comp.Key()
		if _, exists := g.Vertices[key]; !ok || exists {
			// exists: the collector may load the same package by dir and by import path
			continue
		}
		g.AddVertex(key, comp)

		for _, idx := range sortedIndexes(comp.Provides) {
			ann := comp.Provides[idx]
			id := xid(entry, comp.Result(idx), ann)
			providers[id] = append(providers[id], key)
			groups[id] = ann.Group != ""
		}
	}

	for _, id := range sortedKeys(providers) {
		if keys := providers[id]; len(keys) > 1 && !groups[id] && !hasDefault(g, keys) {
			return nil, errors.AlreadyExistsf("%s is provided by more than one function: %s; mark one @Default to choose the implementation by config",
				id, strings.Join(keys, ", "))
		}
	}

	return providers, nil
}

// checkConfigs fails when a package declares more than one @Config struct:
// the keys of a package share its root key.
func checkConfigs(g *graph.Graph[Component]) error {
	byPath := map[string]string{}
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		if !comp.IsConfig() {
			continue
		}
		if other, dup := byPath[comp.Entry.Path]; dup {
			return errors.NotValidf("%s and %s are both @Config in %s: the keys of a package share its root key; move one to its own package",
				other, key, comp.Entry.Path)
		}
		byPath[comp.Entry.Path] = key
	}
	return nil
}

// hasDefault reports whether one of the components is marked @Default.
func hasDefault(g *graph.Graph[Component], keys []string) bool {
	for _, key := range keys {
		if g.Vertices[key].Value.Default {
			return true
		}
	}
	return false
}

// linkComponents adds an edge from each provider to the components consuming
// what it provides.
func linkComponents(ctx context.Context, g *graph.Graph[Component], providers map[string][]string) error {
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for i, param := range comp.Entry.Func.Parameters {
			ann, injected := comp.Injects[i]
			id := injectID(comp.Entry, param.Type, ann)

			keys := providers[id]
			if _, client := clientFactory(comp.Entry, param.Type); len(keys) == 0 && injected && client && ann.Group == "" {
				log.FromContext(ctx).Debugf("parameter %d of %s (%s) is a client inject builds", i, key, param.Type)
				continue
			}
			if len(keys) == 0 && injected && ann.Group == "" && isPublisher(comp.Entry, param.Type) {
				if ann.Name != "" && !envSegment.MatchString(ann.Name) {
					return errors.NotValidf("@Inject on %s parameter %d: the name %s of the %s inject builds must match %s so an environment variable reaches its keys",
						key, i, ann.Name, param.Type, envSegment)
				}
				log.FromContext(ctx).Debugf("parameter %d of %s (%s) is the publisher inject builds", i, key, param.Type)
				continue
			}
			if len(keys) == 0 && injected {
				return errors.NotFoundf("provider not found for %s, required by parameter %d of %s", id, i, key)
			}
			if len(keys) == 0 {
				log.FromContext(ctx).Debugf("parameter %d of %s (%s) has no annotated provider; it must be provided by the application", i, key, param.Type)
				continue
			}

			for _, provider := range keys {
				if provider == key {
					return errors.NotValidf("%s depends on the value it provides (%s)", key, id)
				}
				g.AddEdge(provider, key)
			}
		}
	}
	return nil
}

// Provider returns the component that provides the default value (no name,
// no group) of type tp, written as in package pkg ("*Handler" in package
// orders and "*orders.Handler" elsewhere are the same type), and the index of
// that result.
//
// Deprecated: Provider matches types by package name, so it cannot tell two
// packages with the same name apart; use ProviderFor.
func Provider(g *graph.Graph[Component], pkg, tp string) (Component, int, bool) {
	none := Annotation{}
	want := qualifyName(pkg, tp) + "_" + none.ID()
	return findProvider(g, func(comp Component, idx int, ann Annotation) bool {
		return qualifyName(comp.Entry.Package, comp.Result(idx))+"_"+ann.ID() == want
	})
}

// ProviderFor returns the component that provides the default value (no
// name, no group) of type tp, written in the file of entry, and the index of
// that result. Types are matched by import path.
func ProviderFor(g *graph.Graph[Component], entry annotation.Entry, tp string) (Component, int, bool) {
	want := xid(entry, tp, Annotation{})
	return findProvider(g, func(comp Component, idx int, ann Annotation) bool {
		return xid(comp.Entry, comp.Result(idx), ann) == want
	})
}

func findProvider(g *graph.Graph[Component], match func(Component, int, Annotation) bool) (Component, int, bool) {
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for _, idx := range sortedIndexes(comp.Provides) {
			if match(comp, idx, comp.Provides[idx]) {
				return comp, idx, true
			}
		}
	}
	return Component{}, 0, false
}

// parsedAnnotation is an inject annotation decoded from an entry.
type parsedAnnotation struct {
	kind AnnotationType
	ann  Annotation
}

// newComponent parses the inject annotations of an entry. It returns false
// when the entry does not take part in injection.
func newComponent(entry annotation.Entry) (Component, bool, error) {
	parsed, err := parseAnnotations(entry)
	if err != nil || len(parsed) == 0 {
		return Component{}, false, err
	}

	for _, p := range parsed {
		if p.kind == AnnotationTypeCONFIG {
			return newConfig(entry, parsed)
		}
	}

	if !entry.IsFunc() {
		log.Warnf("inject annotations on %s.%s are ignored: only package-level functions can be injected", entry.Path, entry.Struct)
		return Component{}, false, nil
	}

	comp := Component{
		Entry:    entry,
		Provides: map[int]Annotation{},
		Injects:  map[int]Annotation{},
	}

	if comp.Type, err = componentType(entry, parsed); err != nil {
		return Component{}, false, err
	}

	for _, p := range parsed {
		switch p.kind {
		case AnnotationTypeDEFAULT:
			comp.Default = true
		case AnnotationTypePROVIDE:
			err = addIndexed(comp.Provides, p.ann, len(entry.Func.Results), "@Provide", "result", entry)
		case AnnotationTypeINJECT:
			err = addIndexed(comp.Injects, p.ann, len(entry.Func.Parameters), "@Inject", "parameter", entry)
		}
		if err != nil {
			return Component{}, false, err
		}
	}

	if err := validateProvides(entry, comp.Provides); err != nil {
		return Component{}, false, err
	}
	if err := validateDefault(comp); err != nil {
		return Component{}, false, err
	}
	if err := validateInjects(entry, comp.Injects); err != nil {
		return Component{}, false, err
	}

	return comp, true, nil
}

// parseAnnotations decodes the @Provide, @Inject, @Invoke, @Default and
// @Config annotations of an entry, ignoring the others.
func parseAnnotations(entry annotation.Entry) ([]parsedAnnotation, error) {
	var parsed []parsedAnnotation
	for _, ann := range entry.Annotations {
		annType, err := ParseAnnotationType(strings.ToUpper(ann.Name))
		if err != nil || annType == AnnotationTypeMODULE {
			continue
		}

		var a Annotation
		if hint, ok := noArguments[annType]; ok {
			if strings.TrimSpace(ann.Value) != "" {
				return nil, errors.NotValidf("@%s on %s takes no arguments; %s", ann.Name, gid(entry), hint)
			}
			parsed = append(parsed, parsedAnnotation{kind: annType})
			continue
		}
		if err := ann.Decode(&a); err != nil {
			return nil, errors.NotValidf("decode error on @%s in %s: %v", ann.Name, gid(entry), err)
		}
		parsed = append(parsed, parsedAnnotation{kind: annType, ann: a})
	}
	return parsed, nil
}

// noArguments are the annotations without arguments, with what replaces them.
var noArguments = map[AnnotationType]string{
	AnnotationTypeDEFAULT: "the implementation is named after its package",
	AnnotationTypeCONFIG:  "the keys are named after the package and the fields",
}

// componentType tells whether the annotations make the function a provider
// or an invoker.
func componentType(entry annotation.Entry, parsed []parsedAnnotation) (AnnotationType, error) {
	kinds := map[AnnotationType]bool{}
	for _, p := range parsed {
		kinds[p.kind] = true
	}

	switch {
	case kinds[AnnotationTypePROVIDE] && kinds[AnnotationTypeINVOKE]:
		return 0, errors.NotValidf("%s cannot have both @Provide and @Invoke", gid(entry))
	case kinds[AnnotationTypeDEFAULT] && !kinds[AnnotationTypePROVIDE]:
		return 0, errors.NotValidf("@Default on %s needs a @Provide: it marks the default implementation of the provided value", gid(entry))
	case kinds[AnnotationTypePROVIDE]:
		return AnnotationTypePROVIDE, nil
	case kinds[AnnotationTypeINVOKE]:
		return AnnotationTypeINVOKE, nil
	default:
		return 0, errors.NotValidf("%s has @Inject but neither @Provide nor @Invoke", gid(entry))
	}
}

func validateProvides(entry annotation.Entry, provides map[int]Annotation) error {
	for idx, ann := range provides {
		if entry.Func.Results[idx].Type == "error" {
			return errors.NotValidf("@Provide on %s points to the error result %d", gid(entry), idx)
		}
		if ann.Name != "" && ann.Group != "" {
			return errors.NotValidf("@Provide on %s result %d cannot have both name and group", gid(entry), idx)
		}
		if ann.Plugins != "" {
			return errors.NotValidf("@Provide on %s result %d cannot have plugins: plugins go on the @Inject of a client parameter", gid(entry), idx)
		}
		if ann.Drivers != "" {
			return errors.NotValidf("@Provide on %s result %d cannot have drivers: drivers go on the @Inject of a *publisher.Publisher parameter", gid(entry), idx)
		}
	}
	return nil
}

// validateDefault checks that a @Default function provides one non-group value.
func validateDefault(comp Component) error {
	if !comp.Default {
		return nil
	}
	if len(comp.Provides) != 1 {
		return errors.NotValidf("@Default on %s needs exactly one @Provide", gid(comp.Entry))
	}
	for idx, ann := range comp.Provides {
		if ann.Group != "" {
			return errors.NotValidf("@Default on %s result %d: a group has no default implementation, every provider joins it", gid(comp.Entry), idx)
		}
	}
	return nil
}

func validateInjects(entry annotation.Entry, injects map[int]Annotation) error {
	for idx, ann := range injects {
		if ann.Name != "" && ann.Group != "" {
			return errors.NotValidf("@Inject on %s parameter %d cannot have both name and group", gid(entry), idx)
		}
		if _, ok := sliceElem(entry.Func.Parameters[idx].Type); ann.Group != "" && !ok {
			return errors.NotValidf("@Inject(group=%s) on %s parameter %d must be a slice", ann.Group, gid(entry), idx)
		}
		if err := validatePlugins(entry, idx, ann); err != nil {
			return err
		}
		if err := validateDrivers(entry, idx, ann); err != nil {
			return err
		}
	}
	return nil
}

// validateDrivers checks the drivers of an @Inject: the parameter is a
// *publisher.Publisher and every driver is one inject generates.
func validateDrivers(entry annotation.Entry, idx int, ann Annotation) error {
	if ann.Drivers == "" {
		return nil
	}

	tp := entry.Func.Parameters[idx].Type
	if !isPublisher(entry, tp) {
		return errors.NotValidf("@Inject on %s parameter %d has drivers, but %s is not *publisher.Publisher", gid(entry), idx, tp)
	}

	drivers, err := parseDrivers(ann.Drivers)
	if err != nil {
		return errors.NewNotValid(err, fmt.Sprintf("@Inject on %s parameter %d", gid(entry), idx))
	}
	for _, name := range drivers {
		if _, ok := publisherDriver(name); !ok {
			return errors.NotValidf("@Inject on %s parameter %d: the publisher has no driver %s; its drivers are %s",
				gid(entry), idx, name, publisherDriverNames())
		}
	}
	return nil
}

// validatePlugins checks the plugins of an @Inject: the parameter is a client
// inject builds and every plugin is one of its factory.
func validatePlugins(entry annotation.Entry, idx int, ann Annotation) error {
	if ann.Plugins == "" {
		return nil
	}

	tp := entry.Func.Parameters[idx].Type
	f, ok := clientFactory(entry, tp)
	switch {
	case ann.Group != "":
		return errors.NotValidf("@Inject on %s parameter %d cannot have both group and plugins: plugins configure one client", gid(entry), idx)
	case !ok:
		return errors.NotValidf("@Inject on %s parameter %d has plugins, but %s is not a client inject builds (%s)",
			gid(entry), idx, tp, strings.Join(supportedClients(), ", "))
	}

	plugins, err := parsePlugins(ann.Plugins)
	if err != nil {
		return errors.NewNotValid(err, fmt.Sprintf("@Inject on %s parameter %d", gid(entry), idx))
	}
	for _, plugin := range plugins {
		if _, ok := f.Plugins[plugin]; !ok {
			return errors.NotValidf("@Inject on %s parameter %d: %s has no plugin %s; its plugins are %s",
				gid(entry), idx, tp, plugin, strings.Join(f.PluginNames(), ", "))
		}
	}
	return nil
}

// addIndexed validates the index of a @Provide/@Inject and stores it.
func addIndexed(dst map[int]Annotation, ann Annotation, size int, kind, what string, entry annotation.Entry) error {
	if ann.Index == nil {
		return errors.NotValidf("index parameter is required in %s annotation for %s", kind, gid(entry))
	}
	idx := *ann.Index
	if idx < 0 || idx >= size {
		return errors.NotValidf("%s index %d on %s is out of range: the function has %d %s(s)", kind, idx, gid(entry), size, what)
	}
	if _, exists := dst[idx]; exists {
		return errors.NotValidf("%s index %d is declared twice on %s", kind, idx, gid(entry))
	}
	dst[idx] = ann
	return nil
}

// gid generates a unique identifier for a graph vertex based on the entry's
// import path and declaration name.
func gid(entry annotation.Entry) string {
	return fmt.Sprintf("%s.%s", entry.Path, declName(entry))
}

// declName is the name of the annotated function, or of the struct for a
// struct entry.
func declName(entry annotation.Entry) string {
	if entry.Func.Name == "" {
		return entry.Struct
	}
	return entry.Func.Name
}

// xid identifies a provided or injected value by its type, qualified by
// import path in the file of entry, and its annotation (default, named or
// grouped).
func xid(entry annotation.Entry, tp string, ann Annotation) string {
	return fmt.Sprintf("%s_%s", qualifyType(entry, tp), ann.ID())
}

// injectID returns the id a parameter depends on. A group parameter `[]T`
// depends on the providers of T in that group.
func injectID(entry annotation.Entry, tp string, ann Annotation) string {
	if ann.Group != "" {
		if elem, ok := sliceElem(tp); ok {
			tp = elem
		}
	}
	return xid(entry, tp, ann)
}

func sortedIndexes(m map[int]Annotation) []int {
	idx := make([]int, 0, len(m))
	for i := range m {
		idx = append(idx, i)
	}
	sort.Ints(idx)
	return idx
}

func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
