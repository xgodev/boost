package multiserver

import (
	"context"
	"testing"
	"time"
)

const wait = 5 * time.Second

// blockingServer serves until released, like a real listener.
type blockingServer struct {
	started chan struct{}
	release chan struct{}
}

func newBlockingServer() *blockingServer {
	return &blockingServer{started: make(chan struct{}), release: make(chan struct{})}
}

func (b *blockingServer) Serve(context.Context) {
	close(b.started)
	<-b.release
}

func (b *blockingServer) Shutdown(context.Context) {}

func waitClosed(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(wait):
		t.Fatal(msg)
	}
}

// Check runs on request goroutines (the echo health plugin) while Serve's
// goroutines record servers going down: it must be race-free and report
// the server that stopped.
func TestCheckReportsStoppedServerWhileServing(t *testing.T) {
	ctx := context.Background()
	up, down := newBlockingServer(), newBlockingServer()

	served := make(chan struct{})
	go func() {
		defer close(served)
		Serve(ctx, up, down)
	}()
	waitClosed(t, up.started, "first server did not start")
	waitClosed(t, down.started, "second server did not start")

	stopChecking := make(chan struct{})
	checked := make(chan error, 1)
	go func() {
		for {
			select {
			case <-stopChecking:
				checked <- nil
				return
			default:
				if err := Check(ctx); err != nil {
					checked <- err
					return
				}
			}
		}
	}()

	close(down.release)

	select {
	case err := <-checked:
		if err == nil {
			t.Fatal("Check stopped without reporting the stopped server")
		}
	case <-time.After(wait):
		close(stopChecking)
		t.Fatal("Check kept reporting healthy after a server stopped")
	}

	close(up.release)
	waitClosed(t, served, "Serve did not return after every server stopped")
}
