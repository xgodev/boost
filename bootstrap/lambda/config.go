package lambda

import (
	"github.com/xgodev/boost/config"
)

const (
	root = "faas.lambda"
	skip = root + ".skip"
)

func init() {
	config.Add(skip, false, "skip all triggers")
}
