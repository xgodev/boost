package annotation

import (
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/tools/go/packages"
)

const annotationPrefix = "// @"

type filterFunc func(Entry) bool

type Collector struct {
	filters  []string
	pkgs     []string
	basePath string
	entries  []Entry
	// names maps the import path of every package seen while walking the
	// import graph to its package name.
	names map[string]string
}

func (c *Collector) Entries() []Entry {
	return c.entries
}

func (c *Collector) filterEntries(filter filterFunc) (entries []Entry) {
	for _, entry := range c.Entries() {
		if filter(entry) {
			entries = append(entries, entry)
		}
	}
	return
}

func (c *Collector) EntriesWithArgType(annotationName string, argumentType string) []Entry {
	return c.filterEntries(func(entry Entry) bool {
		return !entry.IsStruct() &&
			containsParamOfType(entry.Func.Parameters, argumentType) &&
			containsAnnotation(entry.Annotations, annotationName)
	})
}

func (c *Collector) EntriesWithResultType(annotation string, result string) []Entry {
	return c.filterEntries(func(entry Entry) bool {
		return !entry.IsStruct() &&
			containsAnnotation(entry.Annotations, annotation) &&
			containsParamOfType(entry.Func.Results, result)
	})
}

func (c *Collector) EntriesWith(annotation string) []Entry {
	return c.filterEntries(func(entry Entry) bool {
		return containsAnnotation(entry.Annotations, annotation)
	})
}

func (c *Collector) EntriesWithPrefix(prefix string) []Entry {
	return c.filterEntries(func(entry Entry) bool {
		for _, ann := range entry.Annotations {
			if strings.HasPrefix(ann.Name, prefix) {
				return true
			}
		}
		return false
	})
}

func containsParamOfType(parameters []EntryFuncType, paramType string) bool {
	for _, parameter := range parameters {
		if parameter.Type == paramType {
			return true
		}
	}

	return false
}

func containsAnnotation(annotations []Annotation, annoName string) bool {
	for _, annotation := range annotations {
		if annotation.Name == annoName {
			return true
		}
	}

	return false
}

// declarations returns, in source order, the struct types, methods and
// functions declared in the file. A method belongs to its receiver type
// wherever that type is declared.
func declarations(file *ast.File) (entries []Entry) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				entries = append(entries, structDeclarations(d)...)
			}
		case *ast.FuncDecl:
			entries = append(entries, funcDeclaration(d))
		}
	}
	return entries
}

func structDeclarations(genDecl *ast.GenDecl) (entries []Entry) {
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}
		doc := typeSpec.Doc
		if doc == nil && !genDecl.Lparen.IsValid() {
			doc = genDecl.Doc
		}
		entries = append(entries, Entry{
			Struct:   typeSpec.Name.Name,
			Fields:   structFields(structType),
			Comments: getComments(doc),
		})
	}
	return entries
}

// structFields returns one field per name, in order. An embedded field has
// no name.
func structFields(structType *ast.StructType) (fields []EntryField) {
	for _, field := range structType.Fields.List {
		tp := types.ExprString(field.Type)
		var tag string
		if field.Tag != nil {
			tag, _ = strconv.Unquote(field.Tag.Value)
		}
		if len(field.Names) == 0 {
			fields = append(fields, EntryField{Type: tp, Tag: tag})
			continue
		}
		for _, name := range field.Names {
			fields = append(fields, EntryField{Name: name.Name, Type: tp, Tag: tag})
		}
	}
	return fields
}

func funcDeclaration(funcDecl *ast.FuncDecl) Entry {
	entry := Entry{
		Comments: getComments(funcDecl.Doc),
		Func: EntryFunc{
			Name:       funcDecl.Name.Name,
			Parameters: getFuncParams(funcDecl.Type.Params),
			Results:    getFuncParams(funcDecl.Type.Results),
		},
	}
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		entry.Struct = receiverTypeName(funcDecl.Recv.List[0].Type)
	}
	return entry
}

// receiverTypeName returns T for the receivers T, *T, T[K] and *T[K, V].
func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverTypeName(t.X)
	case *ast.ParenExpr:
		return receiverTypeName(t.X)
	case *ast.IndexExpr:
		return receiverTypeName(t.X)
	case *ast.IndexListExpr:
		return receiverTypeName(t.X)
	}
	return ""
}

func getFuncParams(fieldList *ast.FieldList) []EntryFuncType {
	var params []EntryFuncType

	if fieldList != nil {
		for _, field := range fieldList.List {
			paramType := types.ExprString(field.Type)
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					params = append(params, EntryFuncType{
						Name: name.Name,
						Type: paramType,
					})
				}
			} else {
				params = append(params, EntryFuncType{Type: paramType})
			}
		}
	}

	return params
}

func getComments(group *ast.CommentGroup) []string {
	var comments []string

	if group != nil {
		for _, comment := range group.List {
			comments = append(comments, comment.Text)
		}
	}

	return comments
}

func (c *Collector) filterFiles(p *packages.Package) (entries []Entry) {
	var modName string
	if p.Module != nil {
		modName = p.Module.Path
	}

	for _, file := range p.Syntax {
		fileName := filepath.Base(p.Fset.File(file.Pos()).Name())
		imports := c.fileImports(file)

		for _, ens := range declarations(file) {
			ans, ok := c.getAnnotations(ens.Comments)
			if !ok {
				continue
			}
			name := ens.Func.Name
			if name == "" {
				name = ens.Struct
			}
			ens.File = fileName
			ens.Header = parseHeader(name, ens.Comments)
			ens.Path = p.PkgPath
			ens.Module = modName
			ens.Package = p.Name
			ens.Annotations = ans
			ens.Imports = imports
			entries = append(entries, ens)
		}
	}

	return entries
}

// fileImports maps the name each import has in the file to its path.
func (c *Collector) fileImports(file *ast.File) map[string]string {
	imports := map[string]string{}
	for _, spec := range file.Imports {
		// the parser only accepts well-formed string literals here
		importPath, _ := strconv.Unquote(spec.Path.Value)
		name := c.packageName(importPath)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		if name == "_" || name == "." {
			continue
		}
		imports[name] = importPath
	}
	return imports
}

// packageName returns the name of the package at importPath, falling back to
// the last path element that is not a major version (v2, v3, ...).
func (c *Collector) packageName(importPath string) string {
	if name, ok := c.names[importPath]; ok {
		return name
	}
	name := path.Base(importPath)
	if dir := path.Dir(importPath); dir != "." && majorVersion.MatchString(name) {
		name = path.Base(dir)
	}
	return name
}

var majorVersion = regexp.MustCompile(`^v[0-9]+$`)

func (c *Collector) getAnnotations(cmts []string) (ans []Annotation, ok bool) {
	for _, cmt := range cmts {
		if an, ok := c.extractAnnotation(cmt); ok {
			ans = append(ans, an)
		}
	}
	return ans, len(ans) > 0
}

// extractAnnotation parses a `// @Name` or `// @Name(key=value, ...)` comment
// whose name passes the filters. A malformed annotation that passes the
// filters is reported and ignored.
func (c *Collector) extractAnnotation(cmt string) (Annotation, bool) {
	if !strings.HasPrefix(cmt, annotationPrefix) {
		return Annotation{}, false
	}

	name, value, valid := parseAnnotation(cmt)
	if !c.isAllowedAnnotation(name) {
		log.Debugf("the annotation is not included in the filters and will be ignored. %s", cmt)
		return Annotation{}, false
	}
	if !valid {
		log.Warnf("The annotation does not follow the format and will be ignored. %s", cmt)
		return Annotation{}, false
	}

	log.Debugf("discovered annotation @%s with values (%s)", name, value)
	return NewAnnotation(name, value), true
}

func (c *Collector) isValidAnnotation(input string) bool {
	_, _, valid := parseAnnotation(input)
	return valid
}

func (c *Collector) isAllowedAnnotation(name string) bool {
	if len(c.filters) == 0 {
		return true
	}
	for _, filter := range c.filters {
		if strings.HasPrefix(name, filter) {
			return true
		}
	}
	return false
}

// parseAnnotation splits `// @Name(values)` into its name and the text between
// the parentheses. valid is false when the comment does not follow that format:
// every value must be a key=value pair or a bare name, and quotes and brackets
// must balance.
func parseAnnotation(cmt string) (name string, value string, valid bool) {
	if !strings.HasPrefix(cmt, annotationPrefix) {
		return "", "", false
	}
	body := strings.TrimSpace(cmt[len(annotationPrefix):])
	end := strings.IndexFunc(body, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '.'
	})
	if end < 0 {
		end = len(body)
	}
	name = body[:end]
	rest := strings.TrimSpace(body[end:])
	if name == "" {
		return "", "", false
	}
	if rest == "" {
		return name, "", true
	}
	if !strings.HasPrefix(rest, "(") || !strings.HasSuffix(rest, ")") {
		return name, "", false
	}
	value = strings.TrimSpace(rest[1 : len(rest)-1])
	items, balanced := splitTopLevel(value)
	if !balanced {
		return name, "", false
	}
	for _, item := range items {
		if !strings.Contains(item, "=") && !bareName.MatchString(strings.TrimSpace(item)) {
			return name, "", false
		}
	}
	return name, value, true
}

// bareName is a value without a key, such as the plugin names of
// @Middleware(log, hystrix). It is kept in Value and left out of Map.
var bareName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.-]*$`)

// parseHeader takes the title from the first line of the doc comment, without
// the leading declaration name, and the description from the lines after it.
// Annotation lines are not part of the header.
func parseHeader(name string, cmts []string) EntryHeader {
	var lines []string
	for _, cmt := range cmts {
		if strings.HasPrefix(cmt, annotationPrefix) {
			continue
		}
		lines = append(lines, strings.TrimSpace(strings.TrimPrefix(cmt, "//")))
	}
	if len(lines) == 0 {
		return EntryHeader{}
	}

	title := lines[0]
	if first, rest, found := strings.Cut(title, " "); found && first == name {
		title = rest
	} else if title == name {
		title = ""
	}

	return EntryHeader{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(strings.Join(lines[1:], "\n")),
	}
}
