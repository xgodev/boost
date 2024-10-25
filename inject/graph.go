package inject

import (
	"context"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
	ustrings "github.com/xgodev/boost/utils/strings"
	"github.com/xgodev/boost/wrapper/log"
	"strings"
)

func NewGraphFromEntries(ctx context.Context, entries []annotation.Entry) (*graph.Graph[Component], error) {

	out := make(map[string]Component)
	in := make(map[string][]Component)

	for _, entry := range entries {
		if !entry.IsFunc() {
			continue
		}

		if !isValidCombinedAnnotations(entry.Annotations) {
			log.Warnf("the combined annotations is invalid")
			continue
		}

		for _, ann := range entry.Annotations {
			if !isValidAnnotation(ann.Name) {
				log.Warnf("the annotation %s is invalid", ann.Name)
				continue
			}

			annType, _ := ParseAnnotationType(strings.ToUpper(ann.Name))

			a := Annotation{}

			err := ann.Decode(&a)
			if err != nil {
				return nil, errors.Wrap(err, errors.NotValidf("failed to decode annotation %s in the entry %s.%s", ann.Name, entry.Path, entry.Func.Name))
			}

			switch annType {
			case AnnotationTypePROVIDE:

				if a.Index == nil {
					return nil, errors.NotValidf("the index parameter is required on the annotation %s in the entry %s.%s", ann.Name, entry.Path, entry.Func.Name)
				}

				index := *a.Index

				for i, res := range entry.Func.Results {
					if index != i {
						continue
					}

					id := xid(entry.Package, res.Type, a)
					if _, ok := out[id]; !ok {
						out[id] = Component{
							Entry: entry,
							An:    a,
						}
					}
				}

			case AnnotationTypeINJECT:

				if a.Index == nil {
					return nil, errors.NotValidf("the index parameter is required on the annotation %s in the entry %s.%s", ann.Name, entry.Path, entry.Func.Name)
				}

				index := *a.Index

				for i, param := range entry.Func.Parameters {
					if index != i {
						continue
					}

					id := xid(entry.Package, param.Type, a)

					if _, ok := out[id]; !ok {
						in[id] = make([]Component, 0)
					}
					in[id] = append(in[id], Component{
						Entry: entry,
						An:    a,
					})
				}

			case AnnotationTypeINVOKE:
			case AnnotationTypeMODULE:
			}
		}

	}

	graph := graph.NewGraph[Component]()
	for id, ae := range out {
		graph.AddVertex(id, ae)
	}

	for id, aes := range in {

		if outAnnoEntry, ok := out[id]; ok {
			for _, inb := range aes {

				graph.AddVertex(gid(inb.Entry), inb)
				graph.AddEdge(gid(outAnnoEntry.Entry), gid(inb.Entry))

			}
		} else {
			return nil, errors.NotFoundf("provider not found for %s", id)
		}

	}

	return graph, nil
}

func gid(entry annotation.Entry) string {
	return strings.Join([]string{entry.Path, entry.Func.Name}, "_")
}

func xid(pkg string, tp string, ann Annotation) string {

	ref := false
	if strings.Contains(tp, "*") {
		ref = true
	}

	tp = strings.ReplaceAll(tp, "*", "")
	if !strings.Contains(tp, ".") {
		tp = strings.Join([]string{pkg, ".", tp}, "")
	}

	if ref {
		tp = strings.Join([]string{"*", tp}, "")
	}

	return strings.Join([]string{tp, ann.ID()}, "_")
}

func isValidCombinedAnnotations(annons []annotation.Annotation) bool {
	var all []string
	for _, ann := range annons {
		all = append(all, strings.ToUpper(ann.Name))
	}

	return !ustrings.SliceContainsAll(all, []string{AnnotationTypePROVIDE.String(), AnnotationTypeINJECT.String()})
}

func isValidAnnotation(value string) bool {
	if ustrings.SliceContains([]string{
		AnnotationTypeMODULE.String(),
		AnnotationTypePROVIDE.String(),
		AnnotationTypeINJECT.String(),
		AnnotationTypeINVOKE.String()},
		strings.ToUpper(value)) {
		return true
	}
	return false
}
