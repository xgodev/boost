package mockhttp

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMockServer(t *testing.T) {
	t.Run("start mock server at specified port", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))
		ms.Start(t)
		defer ms.Teardown()

		addr := "127.0.0.1:60000"
		require.Eventually(t, func() bool {
			_, err := net.Dial("tcp", addr)
			return err == nil
		}, 2*time.Second, 200*time.Millisecond)
	})

	t.Run("start mock server at any port available", func(t *testing.T) {
		ms := NewMockServer()
		ms.Start(t)
		defer ms.Teardown()

		addr := fmt.Sprintf("127.0.0.1:%d", ms.Port())
		require.Eventually(t, func() bool {
			_, err := net.Dial("tcp", addr)
			return err == nil
		}, 2*time.Second, 200*time.Millisecond)
	})

	t.Run("fail if unmapped route is called", func(t *testing.T) {
		mockT := new(testing.T)

		ms := NewMockServer(WithPort(60000))

		ms.Get("/get").Return(StatusCode(http.StatusNoContent))

		ms.Start(mockT)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		_, _ = http.Post(fmt.Sprintf("%s/foo", url), "", nil)

		require.True(t, mockT.Failed())
	})

	t.Run("mock request with return builder", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))

		ms.Get("/get").Return(StatusCode(http.StatusNoContent))

		ms.Start(t)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.Get(fmt.Sprintf("%s/get", url))
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("mock request with body and status return builder in wrong order", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))

		ms.Post("/post").Return(
			JSONBody(`{"result": true}`),
			StatusCode(http.StatusCreated),
		)

		ms.Start(t)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.Post(fmt.Sprintf("%s/post", url), "text/html", nil)
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusCreated, response.StatusCode)

		body, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		require.JSONEq(t, `{"result": true}`, string(body))
	})

	t.Run("mock request with header return builder", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))

		headers := http.Header{
			"X-Foo": []string{"bar"},
		}
		ms.Get("/get").Return(
			StatusCode(http.StatusNoContent),
			Headers(headers),
		)

		ms.Start(t)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.Get(fmt.Sprintf("%s/get", url))
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusNoContent, response.StatusCode)
		for k, v := range headers {
			require.Equal(t, v, response.Header[k])
		}
	})

	t.Run("mock request with json file return builder", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))

		ms.Get("/get").Return(
			StatusCode(http.StatusOK),
			JSONFileBody(t, "./fixtures/body.json"),
		)

		ms.Start(t)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.Get(fmt.Sprintf("%s/get", url))
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusOK, response.StatusCode)

		body, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		require.JSONEq(t, `{"result": true}`, string(body))
	})

	t.Run("mock request with query param matcher", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))

		ms.Get(
			"/get",
			MatchQueryParams(url.Values{"foo": []string{"bar"}}),
		).Return(
			StatusCode(http.StatusNoContent),
		)

		ms.Start(t)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.Get(fmt.Sprintf("%s/get?foo=bar", url))
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("mock request with header matcher", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))

		ms.Get(
			"/get",
			MatchHeader(http.Header{"X-App": []string{"foo"}}),
		).Return(
			StatusCode(http.StatusNoContent),
		)

		ms.Start(t)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/get", url), nil)
		require.NoError(t, err)

		request.Header.Set("X-App", "foo")

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.DefaultClient.Do(request)
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("mock request with json body matcher", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))

		jsonBody := `{"result": true}`
		ms.Post(
			"/post",
			MatchJSONBody(jsonBody),
		).Return(
			StatusCode(http.StatusCreated),
		)

		ms.Start(t)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		bodyReader := strings.NewReader(jsonBody)
		request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/post", url), bodyReader)
		require.NoError(t, err)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.DefaultClient.Do(request)
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("mock request with no return builder", func(t *testing.T) {
		ms := NewMockServer(WithPort(60000))

		ms.Get("/get")

		ms.Start(t)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.Get(fmt.Sprintf("%s/get", url))
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusOK, response.StatusCode)

		body, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		require.Empty(t, body)
	})

	t.Run("verifies that all mocked endpoint where called", func(t *testing.T) {
		mockT := new(testing.T)

		ms := NewMockServer(WithPort(60000))

		ms.Get("/get").Return(StatusCode(http.StatusNoContent))
		ms.Post("/post").Return(StatusCode(http.StatusOK))

		ms.Start(mockT)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.Get(fmt.Sprintf("%s/get", url))
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusNoContent, response.StatusCode)

		ms.AssertExpectations()
		require.True(t, mockT.Failed())
	})

	t.Run("verifies that an endpoint was not called", func(t *testing.T) {
		mockT := new(testing.T)

		ms := NewMockServer(WithPort(60000))

		getEndpoint := ms.Get("/get").Return(StatusCode(http.StatusNoContent)).Endpoint()
		postEndpoint := ms.Post("/post").Return(StatusCode(http.StatusOK)).Endpoint()

		ms.Start(mockT)
		defer ms.Teardown()

		url := ms.URL()
		require.Equal(t, "http://127.0.0.1:60000", url)

		var response *http.Response
		require.Eventually(t, func() bool {
			r, err := http.Get(fmt.Sprintf("%s/get", url))
			if err != nil {
				return false
			}
			response = r
			return true
		}, 2*time.Second, 200*time.Millisecond)

		require.Equal(t, http.StatusNoContent, response.StatusCode)

		ms.AssertNotCalled(postEndpoint)
		require.False(t, mockT.Failed())

		ms.AssertNotCalled(getEndpoint)
		require.True(t, mockT.Failed())
	})
}
