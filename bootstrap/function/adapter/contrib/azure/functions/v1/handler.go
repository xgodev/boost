package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/go-playground/validator/v10"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/sync/errgroup"
)

// Handler answers the Azure Functions host: it turns each invocation into
// CloudEvents and calls the function handler with them.
type Handler[T any] struct {
	fn      function.Handler[T]
	options *Options
}

// NewHandler returns a Handler that calls fn with the given options.
func NewHandler[T any](fn function.Handler[T], options *Options) *Handler[T] {
	return &Handler[T]{fn: fn, options: options}
}

// ServeHTTP answers the ping of the host on /, a forwarded http request on
// /{routePrefix}/..., and an invocation on /{function}.
func (h *Handler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := invocationID(r)
	logger := log.FromContext(r.Context()).WithTypeOf(*h).WithField("invocationId", id)
	ctx := logger.ToContext(r.Context())

	path := strings.Trim(r.URL.Path, "/")
	prefix := strings.Trim(h.options.RoutePrefix, "/")

	switch {
	case path == "":
		w.WriteHeader(http.StatusOK)
	case prefix != "" && strings.HasPrefix(path, prefix+"/"):
		h.forwarded(ctx, w, r, id, strings.TrimPrefix(path, prefix+"/"))
	case r.Method != http.MethodPost:
		http.Error(w, "an invocation is a POST", http.StatusMethodNotAllowed)
	default:
		h.invoke(ctx, w, r, id, path)
	}
}

// forwarded handles an http request the host forwarded as it came
// (enableForwardingHttpRequest / enableProxyingHttpRequest) and answers it.
func (h *Handler[T]) forwarded(ctx context.Context, w http.ResponseWriter, r *http.Request, id, function string) {
	in, err := fromRequest(r, id, function)
	if err != nil {
		writeHTTP(w, h.httpResult(ctx, nil, err))
		return
	}
	in.SetExtension(invocationIDExt, id)
	out, err := h.fn(ctx, in)
	writeHTTP(w, h.httpResult(ctx, out, err))
}

// invoke handles an invocation of the custom handler protocol.
func (h *Handler[T]) invoke(ctx context.Context, w http.ResponseWriter, r *http.Request, id, function string) {
	var req invokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeInvoke(w, http.StatusBadRequest, invokeResponse{Logs: []string{fmt.Sprintf("invalid invocation: %s", err)}})
		return
	}

	raw, err := h.trigger(req)
	if err != nil {
		writeInvoke(w, http.StatusBadRequest, invokeResponse{Logs: []string{err.Error()}})
		return
	}

	if isHTTPTrigger(raw) {
		h.invokeHTTP(ctx, w, raw, id, function)
		return
	}

	if mid := metadataID(req.Metadata); mid != "" {
		id = mid
	}

	var items []json.RawMessage
	batch := json.Unmarshal(raw, &items) == nil
	if !batch {
		items = []json.RawMessage{raw}
	}

	ins := make([]event.Event, len(items))
	for i, item := range items {
		itemID := id
		if batch {
			itemID = fmt.Sprintf("%s-%d", id, i)
		}
		in, err := fromItem(item, itemID, function)
		if err != nil {
			writeInvoke(w, http.StatusBadRequest, invokeResponse{Logs: []string{err.Error()}})
			return
		}
		in.SetExtension(invocationIDExt, invocationID(r))
		ins[i] = in
	}

	outs, logs := h.handle(ctx, ins)
	res := invokeResponse{Logs: logs}
	if len(logs) > 0 {
		// a failed invocation is retried by its trigger
		writeInvoke(w, http.StatusInternalServerError, res)
		return
	}

	var result any
	if batch {
		var results []any
		for _, out := range outs {
			if !isNil(out) {
				results = append(results, out)
			}
		}
		if len(results) > 0 {
			result = results
		}
	} else if !isNil(outs[0]) {
		result = outs[0]
	}

	if result != nil {
		b, err := json.Marshal(result)
		if err != nil {
			writeInvoke(w, http.StatusInternalServerError, invokeResponse{Logs: []string{fmt.Sprintf("encoding the result: %s", err)}})
			return
		}
		h.setOutput(&res, b)
	}
	writeInvoke(w, http.StatusOK, res)
}

// invokeHTTP handles an http trigger, answering its caller through the
// output: the invocation succeeds even when the handler fails.
func (h *Handler[T]) invokeHTTP(ctx context.Context, w http.ResponseWriter, raw json.RawMessage, id, function string) {
	var out T
	r, err := requestOf(raw)
	if err == nil {
		var in event.Event
		if in, err = fromRequest(r, id, function); err == nil {
			in.SetExtension(invocationIDExt, id)
			out, err = h.fn(ctx, in)
		}
	}

	res := invokeResponse{}
	hr := h.httpResult(ctx, out, err)
	if err != nil {
		res.Logs = []string{err.Error()}
	}
	b, merr := json.Marshal(hr)
	if merr != nil {
		writeInvoke(w, http.StatusInternalServerError, invokeResponse{Logs: []string{fmt.Sprintf("encoding the http response: %s", merr)}})
		return
	}
	h.setOutput(&res, b)
	writeInvoke(w, http.StatusOK, res)
}

// trigger is the trigger value of the invocation.
func (h *Handler[T]) trigger(req invokeRequest) (json.RawMessage, error) {
	if h.options.Trigger != "" {
		raw, ok := req.Data[h.options.Trigger]
		if !ok {
			return nil, errors.NotValidf("the invocation has no trigger %q", h.options.Trigger)
		}
		return raw, nil
	}
	if len(req.Data) != 1 {
		return nil, errors.NotValidf("the invocation has %d inputs: set boost.bootstrap.function.adapter.azure.trigger to the trigger binding name", len(req.Data))
	}
	for _, raw := range req.Data {
		return raw, nil
	}
	return nil, nil
}

// handle calls the function with every event, at most MaxWorkers at the same
// time, and returns their results in order and the errors that fail the
// invocation.
func (h *Handler[T]) handle(ctx context.Context, ins []event.Event) ([]T, []string) {
	logger := log.FromContext(ctx).WithTypeOf(*h)

	outs := make([]T, len(ins))
	var mu sync.Mutex
	var logs []string

	var g errgroup.Group
	g.SetLimit(max(1, h.options.MaxWorkers))
	for i, in := range ins {
		g.Go(func() error {
			out, err := h.fn(ctx, in)
			if err != nil && !ignored(err) {
				logger.Errorf("handling event %s: %v", in.ID(), err)
				mu.Lock()
				logs = append(logs, fmt.Sprintf("event %s: %s", in.ID(), err))
				mu.Unlock()
				return nil
			}
			outs[i] = out
			return nil
		})
	}
	_ = g.Wait() // intentional: the workers record their errors in logs and always return nil

	return outs, logs
}

func (h *Handler[T]) setOutput(res *invokeResponse, b []byte) {
	if h.options.Output == "" || h.options.Output == "$return" {
		res.ReturnValue = b
		return
	}
	res.Outputs = map[string]json.RawMessage{h.options.Output: b}
}

// httpResult is the http answer for the result of the handler: its JSON, 204
// when it is nil, or the status and body of the error.
func (h *Handler[T]) httpResult(ctx context.Context, out any, err error) httpResponse {
	jsonHeader := map[string]string{"Content-Type": "application/json"}

	if err != nil {
		if ignored(err) {
			return httpResponse{StatusCode: http.StatusNoContent}
		}
		log.FromContext(ctx).WithTypeOf(*h).Errorf("handling http request: %v", err)
		status := statusOf(err)
		b, _ := json.Marshal(response.Error{HttpStatusCode: status, ErrorCode: strconv.Itoa(status), Message: err.Error()}) // a struct of strings and an int always encodes
		return httpResponse{StatusCode: status, Headers: jsonHeader, Body: string(b)}
	}

	if isNil(out) {
		return httpResponse{StatusCode: http.StatusNoContent}
	}
	b, err := json.Marshal(out)
	if err != nil {
		return h.httpResult(ctx, nil, errors.Internalf("encoding the result: %s", err))
	}
	return httpResponse{StatusCode: http.StatusOK, Headers: jsonHeader, Body: string(b)}
}

func statusOf(err error) int {
	if _, ok := err.(validator.ValidationErrors); ok {
		return http.StatusUnprocessableEntity
	}
	return response.HTTPStatusFor(errors.Classify(err))
}

func writeHTTP(w http.ResponseWriter, hr httpResponse) {
	for k, v := range hr.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(hr.StatusCode)
	if hr.Body != "" {
		_, _ = w.Write([]byte(hr.Body)) // intentional: the caller is gone when the write fails
	}
}

func writeInvoke(w http.ResponseWriter, status int, res invokeResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(res) // intentional: the host is gone when the write fails
}

func ignored(err error) bool {
	pol, ok := errors.IgnoreOf(err)
	return ok && pol&errors.IgnoreAsSuccess != 0
}

// isNil reports whether the result of the handler is nil, also when it is a
// typed nil such as a nil *event.Event.
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
		return rv.IsNil()
	}
	return false
}
