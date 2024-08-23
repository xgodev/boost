package ignore_errors

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"testing"
)

func TestExec(t *testing.T) {
	tests := []struct {
		name         string
		wantError    bool
		ignoreErrors []string
		exec         middleware.AnyErrorExecFunc[*cloudevents.Event]
	}{
		{
			name:      "when handler return success then should not return error",
			wantError: false,
			exec: func(ctx context.Context) (*cloudevents.Event, error) {
				ev := cloudevents.NewEvent()
				return &ev, nil
			},
		},
		{
			name:         "when handler return error and error is in ignore list then should ignore error",
			wantError:    false,
			ignoreErrors: []string{"internal"},
			exec: func(ctx context.Context) (*cloudevents.Event, error) {
				ev := cloudevents.NewEvent()
				return &ev, errors.Internalf("error")
			},
		},
		{
			name:         "when handler return error and error is not in ignore list then return error",
			wantError:    true,
			ignoreErrors: []string{"internal"},
			exec: func(ctx context.Context) (*cloudevents.Event, error) {
				ev := cloudevents.NewEvent()
				return &ev, errors.NotValidf("error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			p := NewAnyErrorMiddlewareWithOptions[*cloudevents.Event](&Options{
				Errors: tt.ignoreErrors,
			})

			ctx := middleware.NewAnyErrorContext[*cloudevents.Event]("test", "xpto")

			_, err := p.Exec(ctx, tt.exec, nil)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
