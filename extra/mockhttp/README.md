# mockhttp

HTTP mock server for end-to-end tests. Declare the endpoints the code under
test calls, start the server, point the app at `ms.URL()` and assert every
endpoint was hit.

```go
func TestCreateBook(t *testing.T) {
	ms := mockhttp.NewMockServer() // or mockhttp.WithPort(12000)
	ms.Get("/isbn", mockhttp.MatchQueryParams(url.Values{"title": {"Foundation"}})).
		Return(mockhttp.StatusCode(http.StatusOK), mockhttp.JSONBody(`{"isbn":"9780345317988"}`))

	ms.Start(t) // define every endpoint before Start
	defer ms.Teardown()

	// call the app, configured with ms.URL() as the upstream address

	ms.AssertExpectations()
}
```

| API | What it does |
|---|---|
| `Get`, `Post`, `Put`, `Patch`, `Delete`, `Head` | register an endpoint (chi pattern) with optional matchers |
| `MatchQueryParams`, `MatchHeader`, `MatchJSONBody` | fail the test when the request does not match |
| `StatusCode`, `Headers`, `JSONBody`, `JSONFileBody` | build the response |
| `AssertExpectations`, `AssertNotCalled` | check which endpoints were called |
| `Router`, `Server` | the underlying `chi.Router` and `httptest.Server` |

Unmapped routes fail the test with 404/405.
