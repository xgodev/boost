package mockhttp

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"testing"
)

type Matcher func(t *testing.T, r *http.Request)

func MatchQueryParams(qp url.Values) Matcher {
	return func(t *testing.T, r *http.Request) {
		t.Helper()
		assert.Equal(t, qp, r.URL.Query())
	}
}

func MatchHeader(headers http.Header) Matcher {
	return func(t *testing.T, r *http.Request) {
		t.Helper()
		for k, v := range headers {
			assert.Equal(t, v, r.Header[k])
		}
	}
}

func MatchJSONBody(jsonBody string) Matcher {
	return func(t *testing.T, r *http.Request) {
		t.Helper()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
			return
		}
		assert.JSONEq(t, jsonBody, string(body))
	}
}
