package sqs

// Options of the sqs publisher.
// @Config
type Options struct {
	Queue string `default:"example" desc:"name of the queue the messages go to"`
}
