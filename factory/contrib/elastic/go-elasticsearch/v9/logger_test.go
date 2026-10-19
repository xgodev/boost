package elasticsearch

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/model/errors"
)

func TestLoggerLogsTheRoundTripAtTheLevelOfItsOutcome(t *testing.T) {
	cases := []struct {
		name   string
		status int
		level  string
	}{
		{"success", http.StatusOK, "debug"},
		{"client error", http.StatusNotFound, "warn"},
		{"server error", http.StatusInternalServerError, "error"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := record(t)
			req := httptest.NewRequest(http.MethodPost, "http://node:9200/idx/_doc", strings.NewReader("abc"))
			res := &http.Response{StatusCode: c.status, Body: http.NoBody}

			err := (&Logger{}).LogRoundTrip(req, res, nil, time.Now(), time.Second)

			require.NoError(t, err)
			lines := rec.lines(c.level)
			require.Len(t, lines, 1)
			assert.Equal(t, map[string]interface{}{
				"method":      http.MethodPost,
				"status_code": c.status,
				"duration":    time.Second,
				"req_bytes":   int64(3),
				"res_bytes":   int64(0),
				"url":         "http://node:9200/idx/_doc",
			}, lines[0].fields)
		})
	}
}

func TestLoggerLogsTheTransportErrorWithoutAResponse(t *testing.T) {
	rec := record(t)
	req := httptest.NewRequest(http.MethodGet, "http://node:9200/", nil)

	err := (&Logger{}).LogRoundTrip(req, nil, errors.New("connection refused"), time.Now(), time.Millisecond)

	require.NoError(t, err)
	lines := rec.lines("error")
	require.Len(t, lines, 1)
	assert.Equal(t, "connection refused", lines[0].msg)
	assert.NotContains(t, lines[0].fields, "status_code", "there is no response to read a status from")
}

func TestLoggerCountsTheResponseBody(t *testing.T) {
	rec := record(t)
	req := httptest.NewRequest(http.MethodGet, "http://node:9200/", nil)
	res := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("hello"))}

	require.NoError(t, (&Logger{}).LogRoundTrip(req, res, nil, time.Now(), 0))

	lines := rec.lines("debug")
	require.Len(t, lines, 1)
	assert.Equal(t, int64(5), lines[0].fields["res_bytes"])
}

func TestLoggerAsksForBothBodies(t *testing.T) {
	assert.True(t, (&Logger{}).RequestBodyEnabled())
	assert.True(t, (&Logger{}).ResponseBodyEnabled())
}
