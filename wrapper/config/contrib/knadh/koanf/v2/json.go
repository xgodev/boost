package koanf

import "encoding/json"

// jsonParser is a koanf.Parser for JSON config files.
type jsonParser struct{}

// Unmarshal decodes a JSON object into a config map.
func (jsonParser) Unmarshal(b []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Marshal encodes a config map as JSON.
func (jsonParser) Marshal(o map[string]interface{}) ([]byte, error) {
	return json.Marshal(o)
}
