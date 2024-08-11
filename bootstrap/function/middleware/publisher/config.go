package publisher

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root       = middleware.Root + ".publisher"
	subject    = root + ".subject"
	deadletter = root + ".deadletter"
	dlenabled  = deadletter + ".enabled"
	dlsubject  = deadletter + ".subject"
	dlerrors   = deadletter + ".errors"
	retry      = root + ".retry"
	renabled   = retry + ".enabled"
	rbackoff   = retry + ".backoff"
)

func init() {
	config.Add(subject, "changeme", "defines output subject")
	config.Add(dlenabled, false, "enables dead letter")
	config.Add(dlsubject, "changeme", "defines dead letter output subject")
	config.Add(dlerrors, []error{}, "defines dead letter errors list")
	config.Add(renabled, false, "enables retry")
	config.Add(rbackoff, 3, "defines retry backoff")
}
