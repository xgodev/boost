package inject

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

type Generator struct {
	moduleName string
	graph      *graph.Graph[Component]
}

func NewGenerator(moduleName string, graph *graph.Graph[Component]) *Generator {
	return &Generator{
		moduleName: moduleName,
		graph:      graph,
	}
}

func (p *Generator) Generate(ctx context.Context) error {

	for _, vert := range p.graph.VerticesWithNoIncomingEdges() {
		err := p.module(ctx, vert)
		if err != nil {
			return errors.Wrap(err, errors.Internalf("error generating module file"))
		}
	}

	return nil
}

func (p *Generator) module(ctx context.Context, vertex *graph.Vertex[Component]) error {
	annoEntry := vertex.Value
	entry := annoEntry.Entry

	packageName := filepath.Base(entry.Path)
	funcName := entry.Func.Name

	data := ModuleData{
		PackageName:  packageName,
		FunctionName: funcName,
		ImportPath:   entry.Path,
		Alias:        alias(entry.Path),
		Entry:        entry,
		Type:         fetchType(entry.Annotations),
	}

	uniqueImports := make(map[string]struct{})

	for _, v := range vertex.Incoming() {

		entry := v.Value.Entry

		var a string
		if entry.Package != packageName {
			a = alias(entry.Path)
		}

		data.Modules = append(data.Modules, ImportData{Alias: a, Entry: entry})

		if entry.Package == packageName {
			continue
		}

		importPath := strings.ReplaceAll(entry.Path, "github.com/", "")
		fullImportPath := p.moduleName + "/gen/inject/" + importPath

		if _, exists := uniqueImports[a]; !exists {
			uniqueImports[a] = struct{}{}
			data.Imports = append(data.Imports, ImportData{Alias: a, Path: fullImportPath, Entry: entry})
		}
	}

	tmpl, err := NewTemplate()
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error creating template"))
	}

	repoPath := strings.ReplaceAll(entry.Path, "github.com/", "")
	fileName := fmt.Sprintf("%s_module.go", strings.ToLower(funcName))
	filePath := filepath.Join("gen", "inject", repoPath, fileName)

	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error creating directories"))
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error executing template"))
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Debugf(string(buf.Bytes()))
		return errors.Wrap(err, errors.Internalf("error formatting source"))
	}

	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error creating file"))
	}
	defer file.Close()

	_, err = file.Write(formatted)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error writing to file"))
	}

	for _, v := range vertex.Adjacent() {
		if err := p.module(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

func alias(packagePath string) string {
	hash := md5.Sum([]byte(packagePath))
	return "pkg_" + hex.EncodeToString(hash[:])
}

func fetchType(annons []annotation.Annotation) string {
	for _, ann := range annons {
		if strings.ToUpper(ann.Name) == AnnotationTypeINVOKE.String() {
			return AnnotationTypeINVOKE.String()
		}
	}

	return AnnotationTypePROVIDE.String()
}
