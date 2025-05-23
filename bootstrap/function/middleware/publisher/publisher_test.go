package publisher

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/mocks"
	"testing"
)

func TestExec(t *testing.T) {
	tests := []struct {
		name             string
		subject          string
		deadLetterErrors []string
		mock             func(driver *mocks.Driver)
		exec             middleware.AnyErrorExecFunc[*cloudevents.Event]
	}{
		{
			name:    "when handler return success then should publish to test",
			subject: "test",
			mock: func(driver *mocks.Driver) {
				driver.On("Publish", mock.Anything, mock.Anything).Times(1).Return(nil, nil)
			},
			exec: func(ctx context.Context) (*cloudevents.Event, error) {
				ev := cloudevents.NewEvent()
				return &ev, nil
			},
		},
		{
			name:             "when handler return error and error is in deadletter list then should publish to deadletter",
			subject:          "deadletter",
			deadLetterErrors: []string{"internal"},
			mock: func(driver *mocks.Driver) {
				driver.On("Publish", mock.Anything, mock.Anything).Times(1).Return(nil, nil)
			},
			exec: func(ctx context.Context) (*cloudevents.Event, error) {
				ev := cloudevents.NewEvent()
				return &ev, errors.Internalf("error")
			},
		},
		{
			name:             "when handler return error and error is not in deadletter list then not should publish",
			subject:          "",
			deadLetterErrors: []string{"internal"},
			mock: func(driver *mocks.Driver) {
				driver.On("Publish", mock.Anything, mock.Anything).Maybe().Times(0).Return(nil)
			},
			exec: func(ctx context.Context) (*cloudevents.Event, error) {
				ev := cloudevents.NewEvent()
				return &ev, errors.NotValidf("error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDrv := &mocks.Driver{}

			p := NewAnyErrorMiddlewareWithOptions[*cloudevents.Event](publisher.New(mockDrv), &Options{
				Subject: "test",
				Deadletter: DeadletterOptions{
					Enabled: true,
					Errors:  tt.deadLetterErrors,
					Subject: "deadletter",
				},
			})

			tt.mock(mockDrv)

			ctx := middleware.NewAnyErrorContext[*cloudevents.Event]("test", "xpto")

			got, _ := p.Exec(ctx, tt.exec, nil)

			assert.Equal(t, tt.subject, got.Subject())
		})
	}
}
