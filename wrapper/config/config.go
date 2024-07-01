package config

type Options struct {
	Hide bool
}

// Option is a func to set values in options.
type Option func(options *Options)

// WithHide sets hide option is true to config.
func WithHide() Option {
	return func(options *Options) {
		options.Hide = true
	}
}

// Config represents a flag configuration.
type Config struct {
	Key         string
	Value       interface{}
	Description string
	Options     *Options
}
