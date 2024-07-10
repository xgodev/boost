package publisher

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = middleware.Root + ".publisher"
	subject = root + ".subject"
)

func init() {
	config.Add(subject, "changemw", "defines output subject")
}
