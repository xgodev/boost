package koanf

import (
	"errors"
	"os"
)

// fileProvider is a koanf.Provider that hands a config file's bytes to a Parser.
type fileProvider string

// ReadBytes returns the file contents.
func (p fileProvider) ReadBytes() ([]byte, error) {
	return os.ReadFile(string(p))
}

// Read is not supported: a file needs a Parser.
func (p fileProvider) Read() (map[string]interface{}, error) {
	return nil, errors.New("file provider does not support Read")
}
