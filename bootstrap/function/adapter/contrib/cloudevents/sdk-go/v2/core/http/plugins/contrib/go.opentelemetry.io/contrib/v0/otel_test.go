package contrib

import (
	"context"
	"testing"

	ce "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/stretchr/testify/assert"
)

func TestRegister_DisabledKeepsTheOptions(t *testing.T) {
	opts := []cehttp.Option{ce.WithPort(8081), ce.WithPath("/events")}

	got := NewOtelWithOptions(&Options{Enabled: false}).Register(context.Background(), opts)

	assert.Len(t, got, 2, "a disabled plugin leaves the port and path of the adapter as they are")
}

func TestRegister_EnabledAddsTheMiddleware(t *testing.T) {
	opts := []cehttp.Option{ce.WithPort(8081), ce.WithPath("/events")}

	got := NewOtelWithOptions(&Options{Enabled: true}).Register(context.Background(), opts)

	assert.Len(t, got, 4, "the round tripper and the middleware after the adapter's options")
}
