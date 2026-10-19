// Package boost publishes price received events. It shares its name with the
// package of the other event, as one package per event does.
package boost

// Publisher sends price received events to its subject.
type Publisher struct {
	Subject string
}

// NewPublisher builds the publisher from this package's Options.
// @Provide(index=0)
func NewPublisher(options *Options) *Publisher {
	return &Publisher{Subject: options.Subject}
}
