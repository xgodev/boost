package kafka

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// One app with both modules: each has its own once, so using the writer does
// not drop the reader. The modules (and contextfx) provide only on the first
// call, so a second app in this process would get nothing.
func TestModules_ProvideWriterAndReaderAndCloseThemOnStop(t *testing.T) {
	key := "BOOST_FACTORY_KAFKAGO_READER_TOPIC"
	t.Setenv(key, "orders")
	boost.Start()
	t.Cleanup(func() {
		os.Setenv(key, "")
		boost.Start()
	})

	var (
		w *kafka.Writer
		r *kafka.Reader
	)
	app := fxtest.New(t, WriterModule(), ReaderModule(), fx.Populate(&w, &r))
	app.RequireStart()
	require.NotNil(t, w)
	require.Equal(t, "orders", r.Config().Topic)
	app.RequireStop()

	require.ErrorIs(t, w.WriteMessages(context.Background(), kafka.Message{Topic: "t"}), io.ErrClosedPipe)
	_, err := r.ReadMessage(context.Background())
	require.ErrorIs(t, err, io.EOF)
}
