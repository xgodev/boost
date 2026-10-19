package koanf

import (
	"fmt"
	"path/filepath"

	"github.com/knadh/koanf/parsers/toml/v2"
	k "github.com/knadh/koanf/v2"
)

// parsers maps a config file extension to the koanf.Parser that reads it.
var parsers = map[string]k.Parser{
	".json": jsonParser{},
	".yaml": yamlParser{},
	".yml":  yamlParser{},
	".toml": toml.Parser(),
}

// parserFor returns the parser of path's extension.
func parserFor(path string) (k.Parser, error) {
	parser, ok := parsers[filepath.Ext(path)]
	if !ok {
		return nil, fmt.Errorf("error on check extension of file %s", path)
	}
	return parser, nil
}
