package config

var (
	configs []Config
)

// Add adds a flag configuration to Entries. A key added again replaces
// its entry in place, so Entries never holds a key twice: flag parsers
// refuse a duplicate, and two packages of one app, like the v1 and v2
// pubsub adapters, can register the same key.
func Add(key string, value interface{}, description string, opts ...Option) {

	o := &Options{}

	for _, opt := range opts {
		opt(o)
	}

	entry := Config{
		Key:         key,
		Value:       value,
		Description: description,
		Options:     o,
	}

	for i := range configs {
		if configs[i].Key == key {
			configs[i] = entry
			return
		}
	}

	configs = append(configs, entry)
}

// Entries returns the flag configuration list as an array.
func Entries() []Config {
	return configs
}

// SetEntries
func SetEntries(v []Config) {
	configs = v
}
