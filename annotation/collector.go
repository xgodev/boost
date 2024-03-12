package annotation

import (
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

type filterFunc func(Entry) bool

type Collector struct {
	filters      []string
	pkgs         []string
	basePath     string
	pkgProcessed map[string]bool
	pkgConfig    *packages.Config
	entries      []Entry
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
			hasSpecificAnnotation(entry, annotation) &&
			hasSpecificResultType(entry, result)
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

func (c *Collector) visitPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() && !strings.Contains(path, "/.") && !strings.Contains(path, "/vendor") {
		log.Tracef("visiting basePath %s", path)
		return c.load(path)
	}
	return nil
}

func (c *Collector) isValidLoadValue(value string) bool {
	if isCorePackage(value) {
		log.Debugf("the package %s will be ignored as it is part of the golang core project.", value)
		return false
	}

	if !c.isAllowedPackage(value) {
		log.Debugf("the package %s is not in the allowed list", value)
		return false
	}

	if c.pkgProcessed[value] {
		log.Debugf("the package %s has already been processed", value)
		return false
	}

	return true
}

func (c *Collector) load(value string) (err error) {

	if !c.isValidLoadValue(value) {
		return nil
	}

	log.Tracef("loading... %s", value)

	pkgs, err := packages.Load(c.pkgConfig, value)
	if err != nil {
		return err
	}

	c.pkgProcessed[value] = true

	for _, p := range pkgs {

		if p.Module == nil {
			continue
		}

		log.Tracef("parsing package %s", p.String())

		for _, imp := range p.Imports {

			log.Debugf("parsing import %s", imp)

			err := c.load(imp.String())
			if err != nil {
				return err
			}
		}

		entries, err := c.filterFiles(p)
		if err != nil {
			return err
		}

		c.entries = append(c.entries, entries...)
	}

	return nil
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

func getStructInfos(file *ast.File) (entries []Entry) {

	for _, decl := range file.Decls {

		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					structInfo := Entry{
						Struct:   typeSpec.Name.Name,
						Comments: getComments(genDecl.Doc),
					}
					entries = append(entries, structInfo)
					entries = append(entries, getStructMethods(file, typeSpec.Name.Name)...)
				}
			}
		}
	}

	return entries
}

// getStructMethods returns all the methods associated with the provided struct.
func getStructMethods(file *ast.File, structName string) (entries []Entry) {

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == structName {
						entries = append(entries, Entry{
							Struct:   structName,
							Comments: getComments(funcDecl.Doc),
							Func: EntryFunc{
								Name:       funcDecl.Name.Name,
								Parameters: getFuncParams(funcDecl.Type.Params),
								Results:    getFuncParams(funcDecl.Type.Results),
							},
						})
					}
				}
			}
		}
	}

	return entries
}

func hasSpecificResultType(entry Entry, result string) bool {
	for _, res := range entry.Func.Results {
		if res.Type == result {
			return true
		}
	}
	return false
}

func hasSpecificAnnotation(entry Entry, annotation string) bool {
	for _, ann := range entry.Annotations {
		if ann.Name == annotation {
			return true
		}
	}
	return false
}

func getFuncInfos(file *ast.File) (entries []Entry) {

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {

			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name != "" {
						continue
					}
				}
			}

			funcInfo := Entry{
				Func: EntryFunc{
					Name:       funcDecl.Name.Name,
					Parameters: getFuncParams(funcDecl.Type.Params),
					Results:    getFuncParams(funcDecl.Type.Results),
				},
				Comments: getComments(funcDecl.Doc),
			}
			entries = append(entries, funcInfo)
		}
	}

	return entries
}

func getFuncParams(fieldList *ast.FieldList) []EntryFuncType {
	var params []EntryFuncType

	if fieldList != nil {
		for _, field := range fieldList.List {
			paramType := types.ExprString(field.Type)
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					//paramInfo := fmt.Sprintf("%s %s", name.Name, paramType)
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

func (c *Collector) filterFiles(p *packages.Package) (entries []Entry, err error) {

	for _, file := range p.Syntax {

		var modName string
		if p.Module != nil {
			modName = p.Module.Path
		}

		var cap []Entry
		cap = append(cap, getStructInfos(file)...)
		cap = append(cap, getFuncInfos(file)...)

		for _, ens := range cap {

			ens.File = file.Name.String()
			ens.Header = c.parseHeader(ens.Comments)
			ens.Path = p.PkgPath
			ens.Module = modName
			ens.Package = p.Name

			if ans, ok := c.getAnnotations(ens.Comments); ok {
				ens.Annotations = ans
				entries = append(entries, ens)
			}

		}

	}

	return entries, err
}

func (c *Collector) getAnnotations(cmts []string) (ans []Annotation, ok bool) {

	log.Tracef("get comments comments")

	var contains bool
	for _, cmt := range cmts {
		an, ok := c.extractAnnotation(cmt)
		if !ok {
			continue
		}
		contains = true
		ans = append(ans, an)
	}
	if !contains {
		log.Debugf("there is no annotation in the comment entry")
		return nil, false
	}
	return ans, true
}

func (c *Collector) extractAnnotation(cmt string) (Annotation, bool) {

	log.Tracef("extracting an annotation from the comment. %s", cmt)

	if !strings.HasPrefix(cmt, "// @") {
		log.Debugf("the comment is not an annotation. %s", cmt)
		return Annotation{}, false
	}

	if !c.isValidAnnotation(cmt) {
		log.Warnf("The annotation does not follow the format and will be ignored. %s", cmt)
		return Annotation{}, false
	}

	var allowed bool
	for _, filter := range c.filters {
		a := strings.Join([]string{"@", filter}, "")
		if strings.Contains(cmt, a) {
			allowed = true
			break
		}
	}

	if len(c.filters) == 0 {
		allowed = true
	}

	if !allowed {
		log.Warnf("The annotation is valid but will be ignored as it is not included in the filters. %s", cmt)
		return Annotation{}, false
	}

	name, value := c.splitNameValue(cmt)
	log.Infof("discovered annotation @%s with values (%s)", name, value)
	return NewAnnotation(name, value), true
}

func (c *Collector) splitNameValue(cmt string) (name string, value string) {
	fields := strings.Split(cmt, "(")
	names := strings.Split(fields[0], "@")
	if len(fields) > 1 {
		value = strings.ReplaceAll(strings.TrimSpace(fields[1]), ")", "")
	}
	return strings.TrimSpace(names[1]), value
}

func (c *Collector) isValidAnnotation(input string) bool {

	log.Tracef("checking if it is a valid annotation. %s", input)

	split := strings.SplitN(input, "@", 2)

	// Checking if string starts with "// @"
	if !(len(split) == 2 && split[0] == "// ") {
		return false
	}

	// Checking if there are parentheses and if they have valid key-value pairs
	parenSplit := strings.SplitN(split[1], "(", 2)
	if len(parenSplit) == 1 {
		spcSplit := strings.Split(parenSplit[0], " ")
		if len(spcSplit) > 1 {
			return false
		}
		log.Debugf("there is a valid annotation in the comment")
		return true
	}
	if len(parenSplit) == 2 {
		if !strings.Contains(parenSplit[1], ")") {
			return false
		}
		parenContent := strings.Trim(parenSplit[1], " )")
		pairs := strings.Split(parenContent, ",")
		for _, pair := range pairs {
			if !strings.Contains(pair, "=") {
				return false
			}
		}
		log.Debugf("there is a valid annotation in the comment")
		return true
	}

	return false
}

func (c *Collector) isAllowedPackage(pkgPath string) bool {

	if c.pkgs == nil {
		return true
	}

	for _, n := range c.pkgs {
		if strings.Contains(pkgPath, n) {
			return true
		}
	}
	return false
}

func (c *Collector) parseHeader(cmts []string) EntryHeader {

	log.Tracef("parsing header on the comment group")

	if len(cmts) == 0 {
		return EntryHeader{}
	}

	w := strings.Split(strings.ReplaceAll(cmts[0], "//", ""), " ")
	var title string
	if len(w) > 2 {
		title = strings.Join(w[2:], " ")
	}

	return EntryHeader{
		Title:       title,
		Description: "// TODO",
	}
}
