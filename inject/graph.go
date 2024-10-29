package inject

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event/datacodec/json"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	ustrings "github.com/xgodev/boost/utils/strings"
	"github.com/xgodev/boost/wrapper/log"
	"strings"
)

// NewGraphFromEntries processes a list of annotation entries, creating a dependency graph
// based on `Provide` and `Inject` annotations. It validates annotations and constructs
// a directed graph where each provider links to injectors that depend on it.
//
// Arguments:
//   - ctx: context for logging and error handling.
//   - entries: a slice of `annotation.Entry`, each representing a function or component
//     with associated annotations.
//
// Returns:
// - A pointer to a `graph.Graph[Component]` representing dependencies between components.
// - An error if the graph generation fails due to invalid annotations or missing providers.
func NewGraphFromEntries(ctx context.Context, entries []annotation.Entry) (*graph.Graph[Component], error) {
	// Debug: Encode entries to JSON and log
	if err := logJSON(ctx, entries); err != nil {
		return nil, err
	}

	out, in := map[string]Component{}, map[string][]Component{}

	// Parse each entry to populate `out` and `in` maps for providers and injectors
	for _, entry := range entries {
		// Skip entry if it is not a valid function or has invalid annotation combinations
		if !isValidEntry(entry) {
			continue
		}

		for _, ann := range entry.Annotations {
			// Validate each annotation name
			if !isValidAnnotation(ann.Name) {
				log.Warnf("Invalid annotation: %s", ann.Name)
				continue
			}

			annType, _ := ParseAnnotationType(strings.ToUpper(ann.Name))
			parsedAnnotation := Annotation{}
			if err := ann.Decode(&parsedAnnotation); err != nil {
				return nil, fmt.Errorf("decode error on %s in %s.%s: %w", ann.Name, entry.Path, entry.Func.Name, err)
			}

			// Process based on annotation type (PROVIDE or INJECT)
			switch annType {
			case AnnotationTypePROVIDE:
				if err := populateOutMap(entry, parsedAnnotation, out); err != nil {
					return nil, err
				}
			case AnnotationTypeINJECT:
				if err := populateInMap(entry, parsedAnnotation, in); err != nil {
					return nil, err
				}
			}
		}
	}

	// Build and return the graph from populated maps
	graph, err := buildGraphFromDependencies(ctx, out, in)
	if err != nil {
		return nil, err
	}

	return graph, nil
}

// populateOutMap adds a PROVIDE annotation entry to the `out` map.
//
// This function only adds entries where `Index` in the annotation matches the index
// of the function result in the entry.
//
// Arguments:
// - entry: `annotation.Entry` containing function metadata.
// - ann: parsed `Annotation` structure from the `Provide` annotation.
// - out: map of component IDs to their corresponding provider components.
//
// Returns:
// - Error if the annotation is missing a required `Index` field.
func populateOutMap(entry annotation.Entry, ann Annotation, out map[string]Component) error {
	if ann.Index == nil {
		return fmt.Errorf("index parameter is required in PROVIDE annotation for %s.%s", entry.Path, entry.Func.Name)
	}

	for i, res := range entry.Func.Results {
		if i != *ann.Index {
			continue
		}
		id := xid(entry.Package, res.Type, ann)
		if _, exists := out[id]; !exists {
			out[id] = Component{Entry: entry, An: ann}
		}
	}
	return nil
}

// populateInMap adds an INJECT annotation entry to the `in` map.
//
// This function only adds entries where `Index` in the annotation matches the index
// of the function parameter in the entry.
//
// Arguments:
// - entry: `annotation.Entry` containing function metadata.
// - ann: parsed `Annotation` structure from the `Inject` annotation.
// - in: map of component IDs to a slice of injector components.
//
// Returns:
// - Error if the annotation is missing a required `Index` field.
func populateInMap(entry annotation.Entry, ann Annotation, in map[string][]Component) error {
	if ann.Index == nil {
		return fmt.Errorf("index parameter is required in INJECT annotation for %s.%s", entry.Path, entry.Func.Name)
	}

	for i, param := range entry.Func.Parameters {
		if i != *ann.Index {
			continue
		}
		id := xid(entry.Package, param.Type, ann)
		in[id] = append(in[id], Component{Entry: entry, An: ann})
	}
	return nil
}

// buildGraphFromDependencies constructs a dependency graph from provider and injector mappings.
//
// Arguments:
// - ctx: context for logging and error handling.
// - out: map where keys are provider IDs and values are provider components.
// - in: map where keys are injector IDs and values are lists of injector components.
//
// Returns:
// - A pointer to `graph.Graph[Component]` with dependency edges.
// - Error if a dependency edge cannot be formed due to missing provider.
func buildGraphFromDependencies(ctx context.Context, out map[string]Component, in map[string][]Component) (*graph.Graph[Component], error) {
	g := graph.NewGraph[Component]()
	for _, provider := range out {
		providerID := gid(provider.Entry)
		g.AddVertex(providerID, provider)
	}

	for id, injectors := range in {
		provider, exists := out[id]
		if !exists {
			return nil, fmt.Errorf("provider not found for %s", id)
		}
		for _, injector := range injectors {
			injectorID := gid(injector.Entry)
			g.AddVertex(injectorID, injector)
			g.AddEdge(gid(provider.Entry), injectorID)
		}
	}

	return g, logJSON(ctx, g) // Log final graph
}

// gid generates a unique identifier for a graph vertex based on the entry's path and function name.
//
// Arguments:
// - entry: `annotation.Entry` struct for which a unique ID is required.
//
// Returns:
// - A string representing a unique identifier for the vertex.
func gid(entry annotation.Entry) string {
	return fmt.Sprintf("%s_%s", entry.Path, entry.Func.Name)
}

// xid generates a unique identifier for dependencies based on package, type, and annotation.
//
// This function ensures identifiers are unique even for references (`*Type`) and different
// annotations with the same type.
//
// Arguments:
// - pkg: the package path of the component.
// - tp: the type of the dependency (parameter or return type).
// - ann: the annotation applied to the component.
//
// Returns:
// - A string representing a unique identifier for the dependency.
func xid(pkg, tp string, ann Annotation) string {
	ref := strings.Contains(tp, "*")
	tp = strings.ReplaceAll(tp, "*", "")
	if !strings.Contains(tp, ".") {
		tp = fmt.Sprintf("%s.%s", pkg, tp)
	}
	if ref {
		tp = "*" + tp
	}
	return fmt.Sprintf("%s_%s", tp, ann.ID())
}

// isValidEntry checks if an entry is a valid function and has compatible annotations.
//
// This function ensures that the entry is a function and its annotations do not contain
// conflicting `Provide` and `Inject` annotations.
//
// Arguments:
// - entry: `annotation.Entry` to validate.
//
// Returns:
// - Boolean indicating whether the entry is valid.
func isValidEntry(entry annotation.Entry) bool {
	if !entry.IsFunc() || !isValidCombinedAnnotations(entry.Annotations) {
		log.Debugf("Invalid entry or annotations for %s", entry.Path)
		return false
	}
	return true
}

// isValidCombinedAnnotations checks for conflicting combinations of `Provide` and `Inject` annotations.
//
// Arguments:
// - annotations: a slice of `annotation.Annotation` to check.
//
// Returns:
// - Boolean indicating whether the annotations are valid.
func isValidCombinedAnnotations(annotations []annotation.Annotation) bool {
	var all []string
	for _, ann := range annotations {
		all = append(all, strings.ToUpper(ann.Name))
	}
	return !ustrings.SliceContainsAll(all, []string{AnnotationTypePROVIDE.String(), AnnotationTypeINJECT.String()})
}

// isValidAnnotation verifies if an annotation name is supported.
//
// Arguments:
// - name: the name of the annotation to validate.
//
// Returns:
// - Boolean indicating whether the annotation name is valid.
func isValidAnnotation(name string) bool {
	return ustrings.SliceContains(
		[]string{AnnotationTypeMODULE.String(), AnnotationTypePROVIDE.String(), AnnotationTypeINJECT.String(), AnnotationTypeINVOKE.String()},
		strings.ToUpper(name))
}

// logJSON logs the JSON-encoded representation of the data for debugging.
//
// Arguments:
// - ctx: context for error handling.
// - data: the data to be JSON-encoded and logged.
//
// Returns:
// - Error if encoding fails, nil otherwise.
func logJSON(ctx context.Context, data interface{}) error {
	bytes, err := json.Encode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}
	fmt.Println(string(bytes))
	return nil
}
