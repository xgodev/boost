package sns

// Options of the sns publisher.
// @Config
type Options struct {
	Topic string `default:"arn:aws:sns:us-east-1:000000000000:example" desc:"arn of the topic the messages go to"`
}
