package azure

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/model/errors"
)

type recorder struct {
	mu   sync.Mutex
	ins  []cloudevents.Event
	out  func(in cloudevents.Event) (*event.Event, error)
	errs map[string]error
}

func (r *recorder) handle(_ context.Context, in cloudevents.Event) (*event.Event, error) {
	r.mu.Lock()
	r.ins = append(r.ins, in)
	r.mu.Unlock()
	if err, ok := r.errs[in.ID()]; ok {
		return nil, err
	}
	if r.out != nil {
		return r.out(in)
	}
	return nil, nil
}

func reply(id string) func(cloudevents.Event) (*event.Event, error) {
	return func(cloudevents.Event) (*event.Event, error) {
		out := event.New()
		out.SetID(id)
		out.SetSource("orders")
		out.SetType("order.created")
		_ = out.SetData(event.ApplicationJSON, map[string]int{"id": 1})
		return &out, nil
	}
}

func options() *Options {
	return &Options{Output: "$return", RoutePrefix: "api", MaxWorkers: 1}
}

func invoke(t *testing.T, h http.Handler, function, payload string) (*httptest.ResponseRecorder, invokeResponse) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/"+function, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Azure-Functions-InvocationId", "inv-1")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	var res invokeResponse
	if w.Code == http.StatusOK {
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res), w.Body.String())
	}
	return w, res
}

func TestHandlerAnswersThePingOfTheHost(t *testing.T) {
	r := &recorder{}
	w := httptest.NewRecorder()
	NewHandler[*event.Event](r.handle, options()).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, r.ins)
}

func TestHandlerTurnsAQueueMessageIntoAnEvent(t *testing.T) {
	r := &recorder{out: reply("out-1")}
	w, res := invoke(t, NewHandler[*event.Event](r.handle, options()), "orders", `{
		"Data": {"item": "{\"id\":1005}"},
		"Metadata": {"Id": "msg-1", "DequeueCount": 1, "sys": {"MethodName": "orders"}}
	}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	require.Len(t, r.ins, 1)
	in := r.ins[0]
	assert.Equal(t, "msg-1", in.ID())
	assert.Equal(t, "azure://functions/orders", in.Source())
	assert.Equal(t, "azure.functions.invoke", in.Type())
	assert.Equal(t, "orders", in.Subject())
	assert.Equal(t, event.ApplicationJSON, in.DataContentType())
	assert.JSONEq(t, `{"id":1005}`, string(in.Data()))
	assert.Equal(t, "inv-1", in.Extensions()["azureinvocationid"])

	var out event.Event
	require.NoError(t, json.Unmarshal(res.ReturnValue, &out), "the result goes to $return")
	assert.Equal(t, "out-1", out.ID())
	assert.Empty(t, res.Outputs)
}

func TestHandlerKeepsAPlainTextMessageAsText(t *testing.T) {
	r := &recorder{}
	w, _ := invoke(t, NewHandler[*event.Event](r.handle, options()), "orders", `{"Data": {"item": "hello"}, "Metadata": {}}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	require.Len(t, r.ins, 1)
	assert.Equal(t, event.TextPlain, r.ins[0].DataContentType())
	assert.Equal(t, "hello", string(r.ins[0].Data()))
	assert.Equal(t, "inv-1", r.ins[0].ID(), "without a message id the invocation id is the id")
}

func TestHandlerUsesTheCloudEventInTheMessage(t *testing.T) {
	r := &recorder{}
	w, _ := invoke(t, NewHandler[*event.Event](r.handle, options()), "orders", `{
		"Data": {"item": "{\"specversion\":\"1.0\",\"id\":\"ce-1\",\"source\":\"shop\",\"type\":\"order.created\",\"data\":{\"id\":1}}"},
		"Metadata": {"Id": "msg-1"}
	}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	require.Len(t, r.ins, 1)
	assert.Equal(t, "ce-1", r.ins[0].ID())
	assert.Equal(t, "shop", r.ins[0].Source())
	assert.Equal(t, "order.created", r.ins[0].Type())
	assert.JSONEq(t, `{"id":1}`, string(r.ins[0].Data()))
	assert.Equal(t, "inv-1", r.ins[0].Extensions()["azureinvocationid"])
}

func TestHandlerMapsAnEventGridEvent(t *testing.T) {
	r := &recorder{}
	w, _ := invoke(t, NewHandler[*event.Event](r.handle, options()), "blobs", `{
		"Data": {"ev": {
			"id": "eg-1", "topic": "/subscriptions/x/storageAccounts/acc", "subject": "/blobServices/default/containers/c/blobs/a.txt",
			"eventType": "Microsoft.Storage.BlobCreated", "eventTime": "2026-10-18T10:00:00.5Z",
			"data": {"url": "https://acc.blob.core.windows.net/c/a.txt"}, "dataVersion": "1", "metadataVersion": "1"
		}},
		"Metadata": {}
	}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	require.Len(t, r.ins, 1)
	in := r.ins[0]
	assert.Equal(t, "eg-1", in.ID())
	assert.Equal(t, "/subscriptions/x/storageAccounts/acc", in.Source())
	assert.Equal(t, "Microsoft.Storage.BlobCreated", in.Type())
	assert.Equal(t, "/blobServices/default/containers/c/blobs/a.txt", in.Subject())
	assert.Equal(t, time.Date(2026, 10, 18, 10, 0, 0, 500_000_000, time.UTC), in.Time().UTC())
	assert.JSONEq(t, `{"url": "https://acc.blob.core.windows.net/c/a.txt"}`, string(in.Data()))
}

func TestHandlerHandlesEveryItemOfABatch(t *testing.T) {
	r := &recorder{out: func(in cloudevents.Event) (*event.Event, error) {
		if string(in.Data()) == `{"n":2}` {
			return nil, nil
		}
		return reply("out-" + string(in.Data()))(in)
	}}
	o := options()
	o.MaxWorkers = 3
	w, res := invoke(t, NewHandler[*event.Event](r.handle, o), "hub", `{
		"Data": {"events": ["{\"n\":1}", "{\"n\":2}", "{\"n\":3}"]},
		"Metadata": {}
	}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	assert.Len(t, r.ins, 3)
	var outs []event.Event
	require.NoError(t, json.Unmarshal(res.ReturnValue, &outs))
	require.Len(t, outs, 2, "a nil result is left out")
	assert.Equal(t, `out-{"n":1}`, outs[0].ID(), "results keep the order of the batch")
	assert.Equal(t, `out-{"n":3}`, outs[1].ID())
}

func TestHandlerFailsTheInvocationWhenTheHandlerFails(t *testing.T) {
	r := &recorder{errs: map[string]error{"msg-1": errors.Internalf("db down")}}
	w, _ := invoke(t, NewHandler[*event.Event](r.handle, options()), "orders", `{"Data": {"item": "{}"}, "Metadata": {"Id": "msg-1"}}`)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "a non-2xx fails the invocation, so the trigger retries it")
	var res invokeResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	require.NotEmpty(t, res.Logs)
	assert.Contains(t, res.Logs[0], "db down")
}

func TestHandlerTakesAnIgnoredErrorAsSuccess(t *testing.T) {
	gone := stderrors.New("gone")
	errors.Ignore(gone, errors.IgnoreAsSuccess)
	r := &recorder{errs: map[string]error{"msg-1": gone}}
	w, _ := invoke(t, NewHandler[*event.Event](r.handle, options()), "orders", `{"Data": {"item": "{}"}, "Metadata": {"Id": "msg-1"}}`)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
}

func TestHandlerWritesTheResultToTheConfiguredOutput(t *testing.T) {
	r := &recorder{out: reply("out-1")}
	o := options()
	o.Output = "message"
	w, res := invoke(t, NewHandler[*event.Event](r.handle, o), "orders", `{"Data": {"item": "{}"}, "Metadata": {}}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	require.Contains(t, res.Outputs, "message")
	var out event.Event
	require.NoError(t, json.Unmarshal(res.Outputs["message"], &out))
	assert.Equal(t, "out-1", out.ID())
	assert.Empty(t, res.ReturnValue)
}

func TestHandlerPicksTheTriggerAmongTheInputs(t *testing.T) {
	r := &recorder{}
	payload := `{"Data": {"item": "{\"a\":1}", "blob": "{\"b\":2}"}, "Metadata": {}}`

	w, _ := invoke(t, NewHandler[*event.Event](r.handle, options()), "orders", payload)
	assert.Equal(t, http.StatusBadRequest, w.Code, "with two inputs the trigger must be configured")
	assert.Contains(t, w.Body.String(), "trigger")

	o := options()
	o.Trigger = "item"
	w, _ = invoke(t, NewHandler[*event.Event](r.handle, o), "orders", payload)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.Len(t, r.ins, 1)
	assert.JSONEq(t, `{"a":1}`, string(r.ins[0].Data()))
}

func TestHandlerRejectsAnInvalidInvocation(t *testing.T) {
	r := &recorder{}
	h := NewHandler[*event.Event](r.handle, options())

	w, _ := invoke(t, h, "orders", `not json`)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/orders", nil))
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Empty(t, r.ins)
}

func httpOutput(t *testing.T, res invokeResponse, name string) httpResponse {
	t.Helper()
	require.Contains(t, res.Outputs, name)
	var out httpResponse
	require.NoError(t, json.Unmarshal(res.Outputs[name], &out))
	return out
}

func TestHandlerAnswersAnHTTPTriggerThroughTheHTTPOutput(t *testing.T) {
	r := &recorder{out: reply("out-1")}
	o := options()
	o.Output = "res"
	w, res := invoke(t, NewHandler[*event.Event](r.handle, o), "orders", `{
		"Data": {"req": {
			"Url": "http://localhost:7071/api/orders", "Method": "POST", "Query": {}, "Params": {},
			"Headers": {"Content-Type": ["application/json"]},
			"Body": "{\"id\":1005}"
		}},
		"Metadata": {}
	}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	require.Len(t, r.ins, 1)
	in := r.ins[0]
	assert.Equal(t, "inv-1", in.ID())
	assert.Equal(t, "http://localhost:7071/api/orders", in.Source())
	assert.Equal(t, "azure.functions.http", in.Type())
	assert.Equal(t, "orders", in.Subject())
	assert.JSONEq(t, `{"id":1005}`, string(in.Data()))

	out := httpOutput(t, res, "res")
	assert.Equal(t, http.StatusOK, out.StatusCode)
	assert.Equal(t, "application/json", out.Headers["Content-Type"])
	var ev event.Event
	require.NoError(t, json.Unmarshal([]byte(out.Body), &ev))
	assert.Equal(t, "out-1", ev.ID())
}

func TestHandlerReadsACloudEventFromTheHTTPTriggerHeaders(t *testing.T) {
	r := &recorder{}
	o := options()
	o.Output = "res"
	w, res := invoke(t, NewHandler[*event.Event](r.handle, o), "orders", `{
		"Data": {"req": {
			"Url": "http://localhost:7071/api/orders", "Method": "POST",
			"Headers": {"Content-Type": ["application/json"], "Ce-Specversion": ["1.0"], "Ce-Id": ["ce-1"], "Ce-Source": ["shop"], "Ce-Type": ["order.created"]},
			"Body": "{\"id\":1}"
		}},
		"Metadata": {}
	}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	require.Len(t, r.ins, 1)
	assert.Equal(t, "ce-1", r.ins[0].ID())
	assert.Equal(t, "shop", r.ins[0].Source())
	assert.Equal(t, "order.created", r.ins[0].Type())
	assert.Equal(t, http.StatusNoContent, httpOutput(t, res, "res").StatusCode, "a nil result answers 204")
}

func TestHandlerAnswersTheHTTPTriggerWithTheStatusOfTheError(t *testing.T) {
	r := &recorder{errs: map[string]error{"inv-1": errors.NotFoundf("order 1005")}}
	o := options()
	o.Output = "res"
	w, res := invoke(t, NewHandler[*event.Event](r.handle, o), "orders", `{
		"Data": {"req": {"Url": "http://localhost:7071/api/orders", "Method": "POST", "Headers": {}, "Body": "{}"}},
		"Metadata": {}
	}`)
	require.Equal(t, http.StatusOK, w.Code, "the invocation succeeds; the caller gets the error")

	out := httpOutput(t, res, "res")
	assert.Equal(t, http.StatusNotFound, out.StatusCode)
	assert.Contains(t, out.Body, "order 1005")
	require.NotEmpty(t, res.Logs)
}

func TestHandlerServesAForwardedHTTPRequest(t *testing.T) {
	r := &recorder{out: reply("out-1")}
	h := NewHandler[*event.Event](r.handle, options())

	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ce-Specversion", "1.0")
	req.Header.Set("Ce-Id", "ce-1")
	req.Header.Set("Ce-Source", "shop")
	req.Header.Set("Ce-Type", "order.created")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	var out event.Event
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	assert.Equal(t, "out-1", out.ID())

	require.Len(t, r.ins, 1)
	assert.Equal(t, "ce-1", r.ins[0].ID())
	assert.JSONEq(t, `{"id":1}`, string(r.ins[0].Data()))
}

func TestHandlerServesAForwardedRawRequest(t *testing.T) {
	r := &recorder{errs: map[string]error{"inv-9": errors.NotValidf("bad order")}}
	h := NewHandler[*event.Event](r.handle, options())

	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Azure-Functions-InvocationId", "inv-9")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "bad order")
	require.Len(t, r.ins, 1)
	assert.Equal(t, "/api/orders", r.ins[0].Source())
	assert.Equal(t, "orders", r.ins[0].Subject())
	assert.Equal(t, "azure.functions.http", r.ins[0].Type())
}
