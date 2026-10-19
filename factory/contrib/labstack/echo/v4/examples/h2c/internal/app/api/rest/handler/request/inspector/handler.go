// Package inspector serves GET /: what the server sees of the request.
package inspector

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

const page = `<code>
  Protocol: %s<br>
  Host: %s<br>
  Remote Address: %s<br>
  Method: %s<br>
  Path: %s<br>
</code>
`

// Handler answers with the request as the server received it.
type Handler struct{}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New() *Handler {
	return &Handler{}
}

// Handle answers with the protocol, host, remote address, method and path
// of the request.
// @REST(method=GET, path=/)
func (h *Handler) Handle(c echo.Context) error {
	req := c.Request()
	return c.HTML(http.StatusOK, fmt.Sprintf(page, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
}
