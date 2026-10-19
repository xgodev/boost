package background

import "time"

// Options of the local server.
// @Config
type Options struct {
	Duration time.Duration `default:"10s" desc:"how long the local server serves before it stops"`
}
