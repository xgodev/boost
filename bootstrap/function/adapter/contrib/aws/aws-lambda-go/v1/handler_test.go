package lambda

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
)

var lc = &lambdacontext.LambdaContext{
	AwsRequestID:       "request-1",
	InvokedFunctionArn: "arn:aws:lambda:us-east-1:123456789012:function:fn",
}

func lambdaContext() context.Context {
	return lambdacontext.NewContext(context.Background(), lc)
}

func load(t *testing.T, name string) Event {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	var e Event
	if err := json.Unmarshal(b, &e); err != nil {
		t.Fatal(err)
	}
	return e
}

// recorder is a function handler that keeps every event it receives.
type recorder struct {
	mu  sync.Mutex
	ins []event.Event
	err error
}

func (r *recorder) handle(_ context.Context, in event.Event) (*event.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ins = append(r.ins, in)
	return nil, r.err
}

func handle(t *testing.T, name string) []event.Event {
	t.Helper()
	r := &recorder{}
	h := NewHandler[*event.Event](r.handle, &Options{MaxWorkers: 1})
	if err := h.Handle(lambdaContext(), load(t, name)); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	return r.ins
}

func only(t *testing.T, ins []event.Event) event.Event {
	t.Helper()
	if len(ins) != 1 {
		t.Fatalf("handler called %d times, want 1", len(ins))
	}
	in := ins[0]
	if got := in.Extensions()["awsrequestid"]; got != lc.AwsRequestID {
		t.Errorf("awsrequestid = %v, want %s", got, lc.AwsRequestID)
	}
	if got := in.Extensions()["invokedfunctionarn"]; got != lc.InvokedFunctionArn {
		t.Errorf("invokedfunctionarn = %v, want %s", got, lc.InvokedFunctionArn)
	}
	return in
}

func header(t *testing.T, in event.Event, id, typ, source string) {
	t.Helper()
	if in.ID() != id {
		t.Errorf("ID = %q, want %q", in.ID(), id)
	}
	if in.Type() != typ {
		t.Errorf("Type = %q, want %q", in.Type(), typ)
	}
	if in.Source() != source {
		t.Errorf("Source = %q, want %q", in.Source(), source)
	}
}

func data(t *testing.T, in event.Event) map[string]any {
	t.Helper()
	var m map[string]any
	if err := in.DataAs(&m); err != nil {
		t.Fatalf("DataAs() error = %v", err)
	}
	return m
}

func field(t *testing.T, m map[string]any, key string, want any) {
	t.Helper()
	if m[key] != want {
		t.Errorf("data[%q] = %v, want %v", key, m[key], want)
	}
}

func TestHandle_Kinesis(t *testing.T) {
	in := only(t, handle(t, "kinesis_success.json"))
	header(t, in, "6aee8d49-76fd-4900-a49b-eecac3552c5a", "receiver/internal/app/receiver/domain/model/item", "receiver-receiver")
	m := data(t, in)
	field(t, m, "id", "123")
	field(t, m, "source", "skyhub")
	field(t, m, "name", "Iphone 11 XS")
}

func TestHandle_S3(t *testing.T) {
	in := only(t, handle(t, "s3_success.json"))
	header(t, in, "b21b84d653bb07b05b1e6b33684dc11b", "ObjectCreated:Put", "aws:s3")
	var s3 events.S3Entity
	if err := in.DataAs(&s3); err != nil {
		t.Fatal(err)
	}
	if s3.Bucket.Name != "lambda-artifacts-deafc19498e3f2df" || s3.Object.Size != 1305107 || s3.Object.Sequencer != "0C0F6F405D6ED209E1" {
		t.Errorf("s3 entity = %+v", s3)
	}
}

func TestHandle_DynamoDB(t *testing.T) {
	in := only(t, handle(t, "dynamodb_success.json"))
	header(t, in, "7de3041dd709b024af6f29e4fa13d34c", "INSERT", "aws:dynamodb")
	if got := in.Extensions()["useridentitytype"]; got != "Service" {
		t.Errorf("useridentitytype = %v, want Service", got)
	}
	if got := in.Extensions()["useridentityprincipalid"]; got != "dynamodb.amazonaws.com" {
		t.Errorf("useridentityprincipalid = %v, want dynamodb.amazonaws.com", got)
	}
	var rec events.DynamoDBStreamRecord
	if err := in.DataAs(&rec); err != nil {
		t.Fatal(err)
	}
	if rec.SequenceNumber != "13021600000000001596893679" || rec.SizeBytes != 112 || rec.StreamViewType != "NEW_IMAGE" {
		t.Errorf("dynamodb record = %+v", rec)
	}
	if got := rec.Keys["Username"].String(); got != "John Doe" {
		t.Errorf("Keys[Username] = %q, want John Doe", got)
	}
}

func TestHandle_SNS(t *testing.T) {
	in := only(t, handle(t, "sns_success.json"))
	header(t, in, "95df01b4-ee98-5cb9-9903-4c221d41eb5e", "Notification", "aws:sns")
	field(t, data(t, in), "test", "123")
}

func TestHandle_SQSPlainBody(t *testing.T) {
	in := only(t, handle(t, "sqs_plain_success.json"))
	if in.Source() != "aws:sqs" {
		t.Errorf("Source = %q, want aws:sqs", in.Source())
	}
	field(t, data(t, in), "id", "abc")
}

func TestHandle_SQSCloudEventBody(t *testing.T) {
	in := only(t, handle(t, "sqs_success.json"))
	header(t, in, "123456", "myType", "mySource")
	field(t, data(t, in), "id", "abc")
}

func TestHandle_SQSWrapsSNS(t *testing.T) {
	in := only(t, handle(t, "sqs_wraps_sns_success.json"))
	if in.ID() != "2e1424d4-f796-459a-8184-9c92662be6da" {
		t.Errorf("ID = %q", in.ID())
	}
	m := data(t, in)
	field(t, m, "test", "123")
	field(t, m, "id", "id-123")
}

func TestHandle_SQSTextBodyAndEveryRecord(t *testing.T) {
	ins := handle(t, "sqs_text_success.json")
	if len(ins) != 2 {
		t.Fatalf("handler called %d times, want 2", len(ins))
	}
	if ins[0].DataContentType() != event.TextPlain || string(ins[0].Data()) != "hello world" {
		t.Errorf("first record = %s %q, want text/plain hello world", ins[0].DataContentType(), ins[0].Data())
	}
	field(t, data(t, ins[1]), "id", "second")
}

func TestHandle_EventBridge(t *testing.T) {
	in := only(t, handle(t, "eventbridge_schedule.json"))
	header(t, in, "53dc4d37-cffa-4f76-80c9-8b7d4a4d2eaa", "Scheduled Event", "aws.events")
	field(t, data(t, in), "job", "nightly")
	if in.Time().IsZero() {
		t.Error("Time is zero, want the event time")
	}
}

func TestHandle_Skip(t *testing.T) {
	r := &recorder{}
	h := NewHandler[*event.Event](r.handle, &Options{Skip: true, MaxWorkers: 1})
	if err := h.Handle(lambdaContext(), load(t, "sqs_success.json")); err != nil {
		t.Fatal(err)
	}
	if len(r.ins) != 0 {
		t.Errorf("handler called %d times, want 0", len(r.ins))
	}
}

func TestHandle_HandlerFails(t *testing.T) {
	r := &recorder{err: errors.New("boom")}
	h := NewHandler[*event.Event](r.handle, &Options{MaxWorkers: 4})
	err := h.Handle(lambdaContext(), load(t, "sqs_text_success.json"))
	if err == nil {
		t.Fatal("Handle() error = nil, want the handler error")
	}
	if !errors.Is(err, r.err) {
		t.Errorf("Handle() error = %v, want it to wrap %v", err, r.err)
	}
	if len(r.ins) != 2 {
		t.Errorf("handler called %d times, want 2: one failure must not stop the others", len(r.ins))
	}
}

func TestHandle_IgnoredError(t *testing.T) {
	defer errors.ResetRegistry()
	r := &recorder{err: errors.New("ignorable")}
	errors.Ignore(r.err, errors.IgnoreAsSuccess)

	h := NewHandler[*event.Event](r.handle, &Options{MaxWorkers: 1})
	if err := h.Handle(lambdaContext(), load(t, "sqs_success.json")); err != nil {
		t.Errorf("Handle() error = %v, want nil for an error ignored as success", err)
	}
}

func TestHandle_NoLambdaContext(t *testing.T) {
	h := NewHandler[*event.Event]((&recorder{}).handle, &Options{MaxWorkers: 1})
	if err := h.Handle(context.Background(), load(t, "sqs_success.json")); !errors.IsInternal(err) {
		t.Errorf("Handle() error = %v, want internal", err)
	}
}

func TestHandle_UnknownTrigger(t *testing.T) {
	r := &recorder{}
	h := NewHandler[*event.Event](r.handle, &Options{MaxWorkers: 1})
	err := h.Handle(lambdaContext(), Event{Records: []Record{{EventSource: "aws:unknown"}}})
	if !errors.IsNotImplemented(err) {
		t.Errorf("Handle() error = %v, want not implemented", err)
	}
	if len(r.ins) != 0 {
		t.Errorf("handler called %d times, want 0", len(r.ins))
	}
}

func TestHandle_NoRecords(t *testing.T) {
	r := &recorder{}
	h := NewHandler[*event.Event](r.handle, &Options{MaxWorkers: 1})
	if err := h.Handle(lambdaContext(), Event{}); err != nil {
		t.Fatal(err)
	}
	if len(r.ins) != 0 {
		t.Errorf("handler called %d times, want 0", len(r.ins))
	}
}

func TestNew_Command(t *testing.T) {
	cmd := New[*event.Event]()(func(context.Context, event.Event) (*event.Event, error) { return nil, nil })
	if cmd.Use != "lambda" {
		t.Errorf("Use = %q, want lambda", cmd.Use)
	}
}

func TestNewOptions_Defaults(t *testing.T) {
	config.Set(koanf.New())
	config.Load()

	o, err := NewOptions()
	if err != nil {
		t.Fatal(err)
	}
	if o.Skip || o.MaxWorkers != 10 {
		t.Errorf("options = %+v, want Skip false and MaxWorkers 10", o)
	}
}
