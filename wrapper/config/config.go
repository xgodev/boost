package config

// Config represents a flag configuration.
type Config struct {
	Key         string
	Value       interface{}
	Description string
	Options     *Options
}
