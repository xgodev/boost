package kinesis

// Options of the kinesis publisher.
// @Config
type Options struct {
	Stream       string `default:"example" desc:"name of the stream the messages go to"`
	PartitionKey string `default:"messages" desc:"partition key of the records"`
}
