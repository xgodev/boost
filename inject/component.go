package inject

import "github.com/xgodev/boost/annotation"

// Component is an annotated function that takes part in the injection graph.
type Component struct {
	Entry annotation.Entry
	// Type is AnnotationTypePROVIDE for constructors and AnnotationTypeINVOKE
	// for functions run when the application starts.
	Type AnnotationType
	// Provides maps a result index to its @Provide annotation.
	Provides map[int]Annotation
	// Injects maps a parameter index to its @Inject annotation. Parameters
	// without @Inject are resolved by type.
	Injects map[int]Annotation
}

// Key uniquely identifies the component in the graph.
func (c Component) Key() string {
	return gid(c.Entry)
}
