package publisher

import (
	stderrors "errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/model/errors"
)

type PublishOutput struct {
	Event *cloudevents.Event
	Error error
}

// outputsError joins the error of every failed output, each tagged with its
// event ID, into one Internal error; nil when every event was published.
func outputsError(outputs []PublishOutput) error {
	var errs []error
	for _, out := range outputs {
		if out.Error != nil {
			errs = append(errs, fmt.Errorf("event %s: %w", out.Event.ID(), out.Error))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.NewInternal(stderrors.Join(errs...),
		fmt.Sprintf("unable to publish %d of %d events", len(errs), len(outputs)))
}
