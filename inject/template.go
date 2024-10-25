package inject

import (
	_ "embed"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
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
		return nil, errors.Wrap(err, errors.Internalf("error parsing template"))
	}
	return tmpl, nil
}
