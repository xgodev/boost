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

// Generate iterates over the graph and generates all necessary modules recursively
func (p *Generator) Generate(ctx context.Context) error {
	for _, vert := range p.graph.VerticesWithNoIncomingEdges() {
		// Start the recursive generation of modules
		err := p.generateModuleRecursive(ctx, vert)
		if err != nil {
			return errors.Wrap(err, errors.Internalf("error generating module file"))
		}
	}
	return nil
}

// This function handles the recursive module generation
func (p *Generator) generateModuleRecursive(ctx context.Context, vertex *graph.Vertex[Component]) error {
	// Step 1: Collect the metadata for the current module
	moduleData, err := p.collectMetadata(ctx, vertex)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error collecting metadata"))
	}

	// Step 2: Generate the file using the collected metadata
	err = p.generateFile(ctx, moduleData)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error generating module file"))
	}

	// Step 3: Recursively process adjacent vertices (dependencies)
	for _, v := range vertex.Adjacent() {
		err := p.generateModuleRecursive(ctx, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// Collects the metadata needed to generate a module
func (p *Generator) collectMetadata(ctx context.Context, vertex *graph.Vertex[Component]) (ModuleData, error) {
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
	return data, nil
}

// Responsible for generating the Go file based on the metadata
func (p *Generator) generateFile(ctx context.Context, data ModuleData) error {
	tmpl, err := NewTemplate()
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error creating template"))
	}

	repoPath := strings.ReplaceAll(data.ImportPath, "github.com/", "")
	fileName := fmt.Sprintf("%s_module.go", strings.ToLower(data.FunctionName))
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
