package nats

import (
	"github.com/xgodev/boost/config"
)

const (
	root     = "faas.nats"
	subjects = root + ".subjects"
	queue    = root + ".queue"
)

func init() {
	config.Add(subjects, []string{"changeme"}, "nats listener subjects")
	config.Add(queue, "changeme", "nats listener queue")
}

// SubjectValue returns the subjects from the configuration via "faas.nats.subjects" key.
func SubjectsValue() []string {
	return config.Strings(subjects)
}

// QueueValue returns the queue from the configuration via "faas.nats.queue" key.
func QueueValue() string {
	return config.String(queue)
}
