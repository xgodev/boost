package annotation

import (
	"errors"
	"golang.org/x/tools/go/packages"
	"path/filepath"
)

func Collect(options ...Option) (*Collector, error) {
	c := &Collector{
		pkgProcessed: make(map[string]bool),
		pkgConfig: &packages.Config{
			Mode: packages.NeedName | packages.NeedTypesInfo | packages.NeedSyntax |
				packages.NeedFiles | packages.NeedImports | packages.NeedDeps | packages.NeedTypes |
				packages.NeedEmbedFiles | packages.NeedExportFile | packages.NeedModule | packages.NeedCompiledGoFiles,
		},
	}
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

	err := filepath.Walk(c.basePath, c.visitPath)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type Option func(*Collector) error

func WithFilters(filters ...string) Option {
	return func(c *Collector) error {
		if filters == nil {
			return errors.New("no types informed")
		}
		c.filters = filters
		return nil
	}
}

func WithPackages(pkgs ...string) Option {
	return func(c *Collector) error {
		if pkgs == nil {
			return errors.New("no packages informed")
		}
		for _, pkg := range pkgs {
			c.pkgs = append(c.pkgs, pkg)
		}
		return nil
	}
}

func WithPath(path string) Option {
	return func(c *Collector) error {
		if path == "" {
			return errors.New("no basePath informed")
		}
		c.basePath = path
		return nil
	}
}
