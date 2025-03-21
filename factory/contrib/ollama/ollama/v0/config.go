package ollama

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root     = "boost.factory.ollama"
	endpoint = ".endpoint"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+endpoint, "http://localhost:11434", "define ollama server url")
}
