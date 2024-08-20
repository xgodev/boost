package ignore_errors

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = middleware.Root + ".ignore_errors"
	ierrors = root + ".errors"
)

func init() {
	config.Add(ierrors, []string{"internal"}, "defines dead letter errors list")
}
