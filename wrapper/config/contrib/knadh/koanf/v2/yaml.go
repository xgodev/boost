package koanf

import "go.yaml.in/yaml/v3"

// yamlParser is a koanf.Parser for YAML config files.
type yamlParser struct{}

// Unmarshal decodes a YAML mapping into a config map.
func (yamlParser) Unmarshal(b []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Marshal encodes a config map as YAML.
func (yamlParser) Marshal(o map[string]interface{}) ([]byte, error) {
	return yaml.Marshal(o)
}
