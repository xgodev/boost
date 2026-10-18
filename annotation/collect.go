package annotation

import (
	"errors"
	"sort"
	"strings"

	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/tools/go/packages"
)

const (
	graphMode  = packages.NeedName | packages.NeedModule | packages.NeedImports | packages.NeedDeps
	syntaxMode = packages.NeedName | packages.NeedModule | packages.NeedFiles |
		packages.NeedCompiledGoFiles | packages.NeedSyntax
)

// Collect loads every package under the base path (as `go list ./...` sees it
// from there: hidden, underscore, testdata and vendor directories are skipped)
// plus the packages they import, and returns the annotated declarations.
// Standard library packages are never scanned; WithPackages restricts the scan
// to the given import paths and the packages below them.
func Collect(options ...Option) (*Collector, error) {
	c := &Collector{}
	for _, opt := range options {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	if c.basePath == "" {
		return nil, errors.New("basePath are required")
	}
	log.Tracef("starting to collect annotations. filters: %v packages: %v basePath: %s", c.filters, c.pkgs, c.basePath)

	paths, err := c.packagesToScan()
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return c, nil
	}

	pkgs, err := packages.Load(&packages.Config{Mode: syntaxMode, Dir: c.basePath}, paths...)
	if err != nil {
		return nil, err
	}
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].PkgPath < pkgs[j].PkgPath })
	for _, p := range pkgs {
		c.entries = append(c.entries, c.filterFiles(p)...)
	}

	return c, nil
}

// packagesToScan walks the import graph of the packages under the base path and
// returns, sorted, the import paths whose annotations must be collected. The
// imports of a package that is not scanned are not followed.
func (c *Collector) packagesToScan() ([]string, error) {
	roots, err := packages.Load(&packages.Config{Mode: graphMode, Dir: c.basePath}, "./...")
	if err != nil {
		return nil, err
	}

	selected := make(map[string]bool)
	c.names = make(map[string]string)
	packages.Visit(roots, func(p *packages.Package) bool {
		c.names[p.PkgPath] = p.Name
		if p.Module == nil || !c.isAllowedPackage(p.PkgPath) {
			log.Debugf("the package %s will be ignored", p.PkgPath)
			return false
		}
		selected[p.PkgPath] = true
		return true
	}, nil)

	paths := make([]string, 0, len(selected))
	for path := range selected {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

// isAllowedPackage reports whether the import path is one of the packages
// given to WithPackages or lives below one of them.
func (c *Collector) isAllowedPackage(pkgPath string) bool {
	if len(c.pkgs) == 0 {
		return true
	}
	for _, n := range c.pkgs {
		if pkgPath == n || strings.HasPrefix(pkgPath, n+"/") {
			return true
		}
	}
	return false
}

type Option func(*Collector) error

// WithFilters keeps only the annotations whose name starts with one of the
// filters: "Rest" keeps @Rest, @RestRouter and @RestResponse.
func WithFilters(filters ...string) Option {
	return func(c *Collector) error {
		if filters == nil {
			return errors.New("no types informed")
		}
		c.filters = filters
		return nil
	}
}

// WithPackages restricts the scan to the given import paths and the packages
// below them, both under the base path and among their imports.
func WithPackages(pkgs ...string) Option {
	return func(c *Collector) error {
		if pkgs == nil {
			return errors.New("no packages informed")
		}
		c.pkgs = append(c.pkgs, pkgs...)
		return nil
	}
}

// WithPath sets the directory to scan. It must be inside a Go module.
func WithPath(path string) Option {
	return func(c *Collector) error {
		if path == "" {
			return errors.New("no basePath informed")
		}
		c.basePath = path
		return nil
	}
}
