package azure

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/google/uuid"
	"github.com/xgodev/boost/model/errors"
)

const (
	invocationIDHeader = "X-Azure-Functions-InvocationId"
	invocationIDExt    = "azureinvocationid"
	invokeType         = "azure.functions.invoke"
	httpType           = "azure.functions.http"
)

// invokeRequest is what the functions host posts to /{function}: the trigger
// and input bindings by name, and the trigger metadata.
type invokeRequest struct {
	Data     map[string]json.RawMessage
	Metadata map[string]json.RawMessage
}

// invokeResponse is the answer to the host: the output bindings by name, the
// logs of the invocation and the $return value.
type invokeResponse struct {
	Outputs     map[string]json.RawMessage `json:",omitempty"`
	Logs        []string
	ReturnValue json.RawMessage `json:",omitempty"`
}

// httpTrigger is the request of an http trigger inside an invocation.
type httpTrigger struct {
	Url     string
	Method  string
	Headers map[string][]string
	Body    json.RawMessage
}

// httpResponse is the value of an http output binding.
type httpResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

// eventGridEvent is an event of the Event Grid schema.
type eventGridEvent struct {
	ID        string          `json:"id"`
	Topic     string          `json:"topic"`
	Subject   string          `json:"subject"`
	EventType string          `json:"eventType"`
	EventTime time.Time       `json:"eventTime"`
	Data      json.RawMessage `json:"data"`
}

// hasKeys reports whether raw is a JSON object with all keys, matched with
// their case, unlike encoding/json.
func hasKeys(raw json.RawMessage, keys ...string) bool {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return false
	}
	for _, k := range keys {
		if _, ok := fields[k]; !ok {
			return false
		}
	}
	return true
}

func isHTTPTrigger(raw json.RawMessage) bool {
	return hasKeys(raw, "Url", "Method", "Headers")
}

// payload is the content of a binding value: the text of a JSON string (a
// queue or Service Bus message), else the JSON itself.
func payload(raw json.RawMessage) []byte {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return []byte(s)
	}
	return raw
}

// fromItem converts one trigger item: the CloudEvent it carries, an Event
// Grid event, or else the item as data, with id from the trigger metadata.
func fromItem(raw json.RawMessage, id, function string) (event.Event, error) {
	p := payload(raw)

	if in, ok := structured(p); ok {
		return in, nil
	}

	if hasKeys(p, "id", "eventType", "eventTime") {
		var eg eventGridEvent
		if err := json.Unmarshal(p, &eg); err == nil {
			in := event.New()
			in.SetID(eg.ID)
			in.SetSource(eg.Topic)
			in.SetSubject(eg.Subject)
			in.SetType(eg.EventType)
			in.SetTime(eg.EventTime)
			if err := in.SetData(event.ApplicationJSON, []byte(eg.Data)); err != nil {
				return in, errors.NotValidf("event grid data: %s", err)
			}
			return in, nil
		}
	}

	in, err := withData(p, "")
	in.SetID(id)
	in.SetSource("azure://functions/" + function)
	in.SetType(invokeType)
	in.SetSubject(function)
	return in, err
}

// fromRequest converts an http request, forwarded by the host or rebuilt
// from an http trigger: the CloudEvent in its headers or body, else the body
// as data.
func fromRequest(r *http.Request, id, function string) (event.Event, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return event.Event{}, errors.BadRequestf("reading the http body: %s", err)
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	if ev, err := cehttp.NewEventFromHTTPRequest(r); err == nil {
		return *ev, nil
	}

	// the query may carry the function key (code=...), so it stays out
	source := *r.URL
	source.RawQuery = ""

	in, err := withData(body, r.Header.Get("Content-Type"))
	in.SetID(id)
	in.SetSource(source.String())
	in.SetType(httpType)
	in.SetSubject(function)
	return in, err
}

// requestOf rebuilds the http request of an http trigger.
func requestOf(raw json.RawMessage) (*http.Request, error) {
	var t httpTrigger
	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, errors.BadRequestf("http trigger: %s", err)
	}
	u, err := url.Parse(t.Url)
	if err != nil {
		return nil, errors.BadRequestf("http trigger url: %s", err)
	}

	var body []byte
	if len(t.Body) > 0 && string(t.Body) != "null" {
		body = payload(t.Body)
	}

	r := &http.Request{
		Method: strings.ToUpper(t.Method),
		URL:    u,
		Header: http.Header{},
		Body:   io.NopCloser(bytes.NewReader(body)),
	}
	for k, v := range t.Headers {
		r.Header[http.CanonicalHeaderKey(k)] = v
	}
	return r, nil
}

// structured decodes p as a CloudEvent in structured mode.
func structured(p []byte) (event.Event, bool) {
	if !hasKeys(p, "specversion") {
		return event.Event{}, false
	}
	var in event.Event
	if err := json.Unmarshal(p, &in); err != nil || in.Validate() != nil {
		return event.Event{}, false
	}
	return in, true
}

// withData returns a new event carrying p as data in contentType, or, when
// it is empty, JSON when p is JSON and plain text otherwise.
func withData(p []byte, contentType string) (event.Event, error) {
	in := event.New()
	if contentType == "" {
		contentType = event.TextPlain
		if json.Valid(p) {
			contentType = event.ApplicationJSON
		}
	}
	if len(p) == 0 {
		in.SetDataContentType(contentType)
		return in, nil
	}
	if err := in.SetData(contentType, p); err != nil {
		return in, errors.NotValidf("invocation data: %s", err)
	}
	return in, nil
}

// invocationID is the id the host gave the invocation, or a new one.
func invocationID(r *http.Request) string {
	if id := r.Header.Get(invocationIDHeader); id != "" {
		return id
	}
	return uuid.NewString()
}

// metadataID is the id of the trigger message in the metadata, when it has one.
func metadataID(m map[string]json.RawMessage) string {
	for _, k := range []string{"Id", "MessageId"} {
		var id string
		if err := json.Unmarshal(m[k], &id); err == nil && id != "" {
			return id
		}
	}
	return ""
}
