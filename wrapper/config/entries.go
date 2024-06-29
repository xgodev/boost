package config

var (
	entries []Config
)

// Add adds a flag configuration to Entries.
func Add(key string, value interface{}, description string, opts ...Option) {

	o := &Options{}

	for _, opt := range opts {
		opt(o)
	}

	entries = append(entries, Config{
		Key:         key,
		Value:       value,
		Description: description,
		Options:     o,
	})
}

// Entries returns the flag configuration list as an array.
func Entries() []Config {
	return entries
}

// SetEntries
func SetEntries(v []Config) {
	entries = v
}
