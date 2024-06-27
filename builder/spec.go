//go:generate go-enum -f=$GOFILE --marshal
package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/xgodev/boost/wrapper/log"
)

// ENUM(HTTP,FUNCTION,GRPC)
type AppType int

// ENUM(INTERFACE,FUNCTION)
type Type int

type AppSpec interface {
	SetFromAnnotations(string, []Annotation)
}

type GRPCSpec struct {
	Name        string
	Description string
	Options     struct {
		GoPackage   string
		JavaPackage string
	}
	Methods []struct {
		Name string
		// Type    string (UNARY, STREAM)
		Handler         string
		Package         string
		RelativePackage string
		Type            Type
		Message         struct {
			Input  string
			Output string
		}
	}
}

func (s *GRPCSpec) SetFromAnnotations(handler string, an []Annotation) {}

type FunctionSpec struct {
	Name            string
	Description     string
	Handler         string
	Type            Type
	Package         string
	RelativePackage string
}

func (s *FunctionSpec) SetFromAnnotations(handler string, an []Annotation) {}

type HTTPSpec struct {
	Name        string
	Description string
	Endpoints   []EndpointHTTPSpec
}

func (s *HTTPSpec) SetFromAnnotations(handler string, an []Annotation) {
	es := EndpointHTTPSpec{
		Handler: handler,
		Method:  http.MethodGet,
	}

	for _, a := range an {
		switch a.AnnotationType {
		case AnnotationTypePACKAGE:
			es.Package = a.SimpleValue()
		case AnnotationTypeRELATIVEPACKAGE:
			es.RelativePackage = a.SimpleValue()
		case AnnotationTypeTYPE:
			v, err := ParseType(strings.ToUpper(a.SimpleValue()))
			if err != nil {
				log.Warnf("cannot parse http handler type.", err.Error())
			}
			es.Type = v
		case AnnotationTypePATH:
			es.Paths = append(es.Paths, strings.ToLower(a.SimpleValue()))
		case AnnotationTypeMETHOD:
			es.Method = strings.ToUpper(a.SimpleValue())
		case AnnotationTypeCONSUME:
			es.Consumes = append(es.Consumes, strings.ToLower(a.SimpleValue()))
		case AnnotationTypePRODUCE:
			es.Produces = append(es.Produces, strings.ToLower(a.SimpleValue()))
		case AnnotationTypePARAM:
			es.Parameters = append(es.Parameters, NewParameterHTTPSpecFromAnnotation(a))
		case AnnotationTypeBODY:
			es.Body = a.SimpleValue()
		case AnnotationTypeRESPONSE:
			es.Responses = append(es.Responses, NewResponseHTTPSpecFromAnnotation(a))
		}
	}

	s.Endpoints = append(s.Endpoints, es)
}

type EndpointHTTPSpec struct {
	Paths           []string
	Method          string
	Package         string
	RelativePackage string
	Description     string
	Handler         string
	Type            Type
	Consumes        []string
	Produces        []string
	Parameters      []ParameterHTTPSpec
	Responses       []ResponseHTTPSpec
	Body            string
}

type ResponseHTTPSpec struct {
	Description string
	Schema      string
	Code        int
}

func NewResponseHTTPSpecFromAnnotation(a Annotation) ResponseHTTPSpec {
	v := strings.Split(a.Value, " ")
	code, _ := strconv.Atoi(v[1])
	return ResponseHTTPSpec{
		Description: strings.Join(v[3:], " "),
		Code:        code,
		Schema:      v[2],
	}
}

type ParameterHTTPSpec struct {
	Name        string
	Description string
	Source      string
	Type        string
	Required    bool
	Validations struct{}
}

func NewParameterHTTPSpecFromAnnotation(a Annotation) ParameterHTTPSpec {
	v := strings.Split(a.Value, " ")
	req, _ := strconv.ParseBool(v[4])
	return ParameterHTTPSpec{
		Name:        v[2],
		Description: strings.Join(v[5:], " "),
		Source:      v[1],
		Type:        v[3],
		Required:    req,
		Validations: struct{}{},
	}
}

type Spec struct {
	Apps []*AppsSpec
}

type AppsSpec struct {
	Type AppType
	Spec AppSpec
}
