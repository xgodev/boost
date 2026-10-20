# bootstrap/job

A job is a function an app runs once and then exits, the shape of a
Kubernetes CronJob:

```go
// Handler is the work of a job.
type Handler func(ctx context.Context) error
```

It returns when the work is done, with the error that makes the app exit
non-zero, or when `ctx` is cancelled because the app is stopping (SIGTERM).
[`fx/modules/bootstrap/job`](../../fx/modules/bootstrap/job) runs it;
[builder](../../builder#job-an-app-that-runs-once-and-exits) generates the app
from `@Job`:

```go
// Run checks the stock of the configured items.
// @Job(app=stock-audit)
func (a *Auditor) Run(ctx context.Context) error { ... }
```
