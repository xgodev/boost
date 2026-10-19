package kafka

// Options of the kafka publisher.
// @Config
type Options struct {
	Topic string `default:"messages" desc:"topic the messages go to"`
}
