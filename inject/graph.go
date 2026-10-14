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
// left for the application to provide.
//
// It returns a NotValid error for malformed annotations or a dependency
// cycle, AlreadyExists when two functions provide the same non-group value
// and NotFound when an @Inject has no provider.
func NewGraphFromEntries(ctx context.Context, entries []annotation.Entry) (*graph.Graph[Component], error) {
	g := graph.NewGraph[Component]()

	providers, err := addComponents(g, entries)
	if err != nil {
		return nil, err
	}

	if err := linkComponents(ctx, g, providers); err != nil {
		return nil, err
	}

	if _, remaining, err := g.TopologicalSort(); err != nil {
		return nil, errors.NotValidf("dependency cycle among %s", strings.Join(remaining, ", "))
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
			id := xid(entry.Package, entry.Func.Results[idx].Type, ann)
			providers[id] = append(providers[id], key)
			groups[id] = ann.Group != ""
		}
	}

	for _, id := range sortedKeys(providers) {
		if keys := providers[id]; len(keys) > 1 && !groups[id] {
			return nil, errors.AlreadyExistsf("%s is provided by more than one function: %s", id, strings.Join(keys, ", "))
		}
	}

	return providers, nil
}

// linkComponents adds an edge from each provider to the components consuming
// what it provides.
func linkComponents(ctx context.Context, g *graph.Graph[Component], providers map[string][]string) error {
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for i, param := range comp.Entry.Func.Parameters {
			ann, injected := comp.Injects[i]
			id := injectID(comp.Entry.Package, param.Type, ann)

			keys := providers[id]
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
	if err := validateInjects(entry, comp.Injects); err != nil {
		return Component{}, false, err
	}

	return comp, true, nil
}

// parseAnnotations decodes the @Provide, @Inject and @Invoke annotations of
// an entry, ignoring the others.
func parseAnnotations(entry annotation.Entry) ([]parsedAnnotation, error) {
	var parsed []parsedAnnotation
	for _, ann := range entry.Annotations {
		annType, err := ParseAnnotationType(strings.ToUpper(ann.Name))
		if err != nil || annType == AnnotationTypeMODULE {
			continue
		}

		var a Annotation
		if err := ann.Decode(&a); err != nil {
			return nil, errors.NotValidf("decode error on @%s in %s: %v", ann.Name, gid(entry), err)
		}
		parsed = append(parsed, parsedAnnotation{kind: annType, ann: a})
	}
	return parsed, nil
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
// import path and function name.
func gid(entry annotation.Entry) string {
	return fmt.Sprintf("%s.%s", entry.Path, entry.Func.Name)
}

// xid identifies a provided or injected value by its qualified type and
// annotation (default, named or grouped).
func xid(pkg, tp string, ann Annotation) string {
	return fmt.Sprintf("%s_%s", qualifyType(pkg, tp), ann.ID())
}

// injectID returns the id a parameter depends on. A group parameter `[]T`
// depends on the providers of T in that group.
func injectID(pkg, tp string, ann Annotation) string {
	if ann.Group != "" {
		if elem, ok := sliceElem(tp); ok {
			tp = elem
		}
	}
	return xid(pkg, tp, ann)
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
