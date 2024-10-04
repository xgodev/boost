package inject

import (
	_ "embed"
	"fmt"
	"github.com/xgodev/boost/annotation"
	"text/template"
)

//go:embed templates/module.tmpl
var moduleTemplate string

type ModuleData struct {
	PackageName  string
	FunctionName string
	ImportPath   string
	Modules      []ImportData
	Imports      []ImportData
	Alias        string
	Entry        annotation.Entry
	Type         string
}

type ImportData struct {
	Alias string
	Path  string
	Entry annotation.Entry
}

func NewTemplate() (*template.Template, error) {
	tmpl, err := template.New("module").Parse(moduleTemplate)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}
	return tmpl, nil
}
