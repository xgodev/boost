package util

import (
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func payloadEvent(target string) *v2.Event {
	e := v2.NewEvent()
	e.SetID("123")
	e.SetSource("/orders")
	e.SetType("order.created")
	e.SetSubject("orders")
	if target != "" {
		e.SetExtension("target", target)
	}
	_ = e.SetData(v2.ApplicationJSON, map[string]string{"id": "123"})
	return &e
}

func TestPayload(t *testing.T) {
	testCases := []struct {
		name   string
		target string
		want   string
	}{
		{name: "no target: the whole event", want: `{"specversion":"1.0","id":"123","source":"/orders","type":"order.created","subject":"orders","datacontenttype":"application/json","data":{"id":"123"}}`},
		{name: "target data: only the data", target: "data", want: `{"id":"123"}`},
		{name: "another target: the whole event", target: "all", want: `{"specversion":"1.0","id":"123","source":"/orders","type":"order.created","subject":"orders","datacontenttype":"application/json","data":{"id":"123"},"target":"all"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Payload(payloadEvent(tc.target))
			require.NoError(t, err)
			assert.JSONEq(t, tc.want, string(got))
		})
	}
}

func TestPayload_DataThatIsNotJSON(t *testing.T) {
	e := payloadEvent("data")
	_ = e.SetData("text/plain", []byte("not json"))

	_, err := Payload(e)
	assert.Error(t, err)
}
