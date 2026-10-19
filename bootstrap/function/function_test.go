package function_test

import (
	"context"
	"os"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	co "github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func handle(context.Context, cloudevents.Event) (*cloudevents.Event, error) { return nil, nil }

// command is an adapter command named name that records it ran.
func command(name string, ran *[]string) function.CmdFunc[*cloudevents.Event] {
	return func(function.Handler[*cloudevents.Event]) *co.Command {
		return &co.Command{Use: name, RunE: func(*co.Command, []string) error {
			*ran = append(*ran, name)
			return nil
		}}
	}
}

// withArgs runs the test with the process arguments a binary gets.
func withArgs(t *testing.T, args ...string) {
	t.Helper()
	saved := os.Args
	os.Args = append([]string{"app"}, args...)
	t.Cleanup(func() { os.Args = saved })
}

func newFunction(t *testing.T) *function.Function[*cloudevents.Event] {
	t.Helper()
	f, err := function.New[*cloudevents.Event]()
	require.NoError(t, err)
	return f
}

// An app wires the one adapter config chose (#203): with no argument and no
// boost.bootstrap.function.default, Run runs it.
func TestRun_TheOnlyAdapterRunsWithoutArguments(t *testing.T) {
	withArgs(t)
	var ran []string

	require.NoError(t, newFunction(t).Run(context.Background(), handle, command("nats", &ran)))

	assert.Equal(t, []string{"nats"}, ran)
}

func TestRun_AnArgumentChoosesTheAdapter(t *testing.T) {
	withArgs(t, "gcp_pubsub")
	var ran []string

	require.NoError(t, newFunction(t).Run(context.Background(), handle, command("nats", &ran), command("gcp_pubsub", &ran)))

	assert.Equal(t, []string{"gcp_pubsub"}, ran)
}

// A function that ends up with nothing to run fails, so the app exits
// non-zero instead of printing the usage and stopping as if it had worked.
func TestRun_NothingToRunFails(t *testing.T) {
	t.Run("no adapter", func(t *testing.T) {
		withArgs(t)

		err := newFunction(t).Run(context.Background(), handle)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no adapter")
	})

	t.Run("two adapters and none chosen", func(t *testing.T) {
		withArgs(t)
		var ran []string

		err := newFunction(t).Run(context.Background(), handle, command("nats", &ran), command("gcp_pubsub", &ran))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "nats")
		assert.Contains(t, err.Error(), "gcp_pubsub")
		assert.Contains(t, err.Error(), function.Root+".default")
		assert.Empty(t, ran)
	})
}
