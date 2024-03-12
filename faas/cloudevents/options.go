package cloudevents

import "strings"

// HandlerWrapperOptions represents options for handler wrapped in middleware.
type HandlerWrapperOptions struct {
	IDsToDiscard []string
}

// DefaultHandlerWrapperOptions returns handler wrapped in middleware options.
func DefaultHandlerWrapperOptions() (*HandlerWrapperOptions, error) {
	o := new(HandlerWrapperOptions)

	rawIDs := HandleDiscardEventsIDValue()
	ids := strings.Split(rawIDs, ",")
	if nil != ids && len(ids) >= 1 && ids[0] == "" {
		return o, nil
	}

	o.IDsToDiscard = ids

	return o, nil
}
