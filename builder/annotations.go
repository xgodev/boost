//go:generate go-enum -f=$GOFILE --marshal
package main

import (
	"strings"

	"github.com/xgodev/boost/model/errors"
	ustrings "github.com/xgodev/boost/utils/strings"

	"github.com/ettle/strcase"
	"github.com/xgodev/boost/wrapper/log"
)

// ENUM(APP,HANDLER_TYPE,PACKAGE,RELATIVE_PACKAGE,TYPE,PATH,METHOD,CONSUME,PRODUCE,PARAM,BODY,RESPONSE)
type AnnotationType int

type Annotation struct {
	AnnotationType AnnotationType
	Value          string
}

func (m *Annotation) SimpleValue() string {
	return strings.Split(m.Value, " ")[1]
}

type Annotations map[string][]string

func (m *Annotations) ToSpec() (spec Spec, err error) {

	var appsSpecs []*AppsSpec

	appsSpecsMap := make(map[string]*AppsSpec)

	for hf, as := range *m {
		app, htp, err := m.extractAppAndHandlerType(as)
		if err != nil {
			return spec, err
		}

		log.Infof("%s - %s", app.AnnotationType.String(), app.SimpleValue())
		log.Infof("%s - %s", htp.AnnotationType.String(), htp.SimpleValue())

		appsSpec, ok := appsSpecsMap[app.SimpleValue()]
		if !ok {
			appsSpec = &AppsSpec{}
			appsSpecsMap[app.SimpleValue()] = appsSpec
			appsSpecs = append(appsSpecs, appsSpec)
		}

		var ats []Annotation

		switch htp.SimpleValue() {
		case AppTypeHTTP.String():

			if appsSpec.Spec == nil {
				appsSpec.Spec = &HTTPSpec{
					Name: app.SimpleValue(),
				}
			}

			ats, err = m.extractParamsWithTypes(as,
				AnnotationTypePACKAGE,
				AnnotationTypeRELATIVEPACKAGE,
				AnnotationTypeTYPE,
				AnnotationTypePATH,
				AnnotationTypeMETHOD,
				AnnotationTypeCONSUME,
				AnnotationTypePRODUCE,
				AnnotationTypePARAM,
				AnnotationTypeBODY,
				AnnotationTypeRESPONSE)
			if err != nil {
				return spec, err
			}

		}
		appsSpec.Spec.SetFromAnnotations(hf, ats)
	}

	return Spec{
		Apps: appsSpecs,
	}, err
}

func (m *Annotations) extractAppAndHandlerType(as []string) (app Annotation, handlerType Annotation, err error) {
	a, err := m.extractParamsWithTypes(as, AnnotationTypeHANDLERTYPE, AnnotationTypeAPP)
	if err != nil {
		return Annotation{}, Annotation{}, err
	}
	if len(a) < 2 {
		return Annotation{}, Annotation{}, errors.NotFoundf("%s or %s not found in spec",
			AnnotationTypeHANDLERTYPE.String(), AnnotationTypeAPP.String())
	}
	for _, v := range a {
		switch v.AnnotationType {
		case AnnotationTypeHANDLERTYPE:
			handlerType = v
		case AnnotationTypeAPP:
			app = v
		}
	}

	if handlerType.Value == "" || app.Value == "" {
		return Annotation{}, Annotation{}, errors.NotFoundf("%s or %s not found in spec",
			AnnotationTypeHANDLERTYPE.String(), AnnotationTypeAPP.String())
	}

	return app, handlerType, err
}

func (m *Annotations) extractParamsWithTypes(a []string, tp ...AnnotationType) (as []Annotation, err error) {

	var astr []string
	for _, v := range tp {
		astr = append(astr, v.String())
	}

	for _, v := range a {
		an := m.extractAnn(v)
		if ustrings.SliceContains(astr, an) {
			at, err := ParseAnnotationType(an)
			if err != nil {
				log.Errorf("error on parse annotation %s. %s", an, err.Error())
				continue
			}
			as = append(as, Annotation{
				AnnotationType: at,
				Value:          v,
			})
		}
	}

	return as, nil
}

func (m *Annotations) extractAnn(a string) string {
	return strcase.ToCase(strings.Split(a, " ")[0], strcase.UpperCase, '_')
}
