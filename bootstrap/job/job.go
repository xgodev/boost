// Package job is a function an app runs once and then exits, the shape of a
// Kubernetes CronJob. fx/modules/bootstrap/job runs it; builder generates
// the app from @Job.
package job

import "context"

// Handler is the work of a job. It returns when the work is done, with the
// error that makes the app exit non-zero, or when ctx is cancelled because
// the app is stopping (SIGTERM).
type Handler func(ctx context.Context) error
