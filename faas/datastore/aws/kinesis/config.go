package kinesis

import (
	"github.com/xgodev/boost/config"
)

const (
	root               = "faas.provider.kinesis"
	randomPartitionKey = root + ".randomPartitionKey"
)

func init() {
	config.Add(randomPartitionKey, false, "ramdomize partition key")
}

// RandomPartitionKeyValue returns if random partition key that is enabled or not via the "faas.provider.kinesis.randomPartitionKey" key.
// If not configured, the default is false.
func RandomPartitionKeyValue() bool {
	return config.Bool(randomPartitionKey)
}
