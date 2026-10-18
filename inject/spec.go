//go:generate go-enum -f=$GOFILE --marshal
package inject

import (
	"fmt"
	"strings"
)

// ENUM(MODULE,PROVIDE,INJECT,INVOKE,DEFAULT,CONFIG)
type AnnotationType int

// ENUM(MODULE,PATH,PACKAGE,FUNC)
type ModuleAttr int

// ENUM(DEFAULT,NAMED,GROUPED)
type AnnotationIDType int

// Annotation holds the arguments of a @Provide or @Inject annotation.
// Index points to the function result (@Provide) or parameter (@Inject);
// Name and Group map to fx `name:"..."` and `group:"..."` tags.
type Annotation struct {
	Index *int
	Name  string
	Group string
}

// IDType reports whether the annotation is a default, named or grouped value.
func (a *Annotation) IDType() AnnotationIDType {
	switch {
	case a.Name != "":
		return AnnotationIDTypeNAMED
	case a.Group != "":
		return AnnotationIDTypeGROUPED
	default:
		return AnnotationIDTypeDEFAULT
	}
}

// ID identifies the value kind of the annotation. Names and groups keep
// their case because fx matches tags case-sensitively.
func (a *Annotation) ID() string {
	switch a.IDType() {
	case AnnotationIDTypeNAMED:
		return strings.ToLower(AnnotationIDTypeNAMED.String()) + "_" + a.Name
	case AnnotationIDTypeGROUPED:
		return strings.ToLower(AnnotationIDTypeGROUPED.String()) + "_" + a.Group
	default:
		return strings.ToLower(AnnotationIDTypeDEFAULT.String())
	}
}

// Tag returns the fx struct tag for the annotation, or "" for a default value.
func (a *Annotation) Tag() string {
	switch a.IDType() {
	case AnnotationIDTypeNAMED:
		return fmt.Sprintf("name:%q", a.Name)
	case AnnotationIDTypeGROUPED:
		return fmt.Sprintf("group:%q", a.Group)
	default:
		return ""
	}
}
