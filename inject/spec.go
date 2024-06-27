//go:generate go-enum -f=$GOFILE --marshal
package inject

import (
	"strings"
)

// ENUM(MODULE,PROVIDE,INJECT,INVOKE)
type AnnotationType int

// ENUM(MODULE,PATH,PACKAGE,FUNC)
type ModuleAttr int

// ENUM(DEFAULT,NAMED,GROUPED)
type AnnotationIDType int

type Annotation struct {
	Index *int
	Name  string
	Group string
}

func (a *Annotation) ID() string {
	var fields []string
	if a.Name == "" && a.Group == "" {
		fields = append(fields, AnnotationIDTypeDEFAULT.String())
	} else if a.Name != "" {
		fields = append(fields, AnnotationIDTypeNAMED.String(), a.Name)
	} else if a.Group != "" {
		fields = append(fields, AnnotationIDTypeGROUPED.String(), a.Group)
	}

	return strings.ToLower(
		strings.Join(fields, "_"))
}
