// Package fetcher serves GET /google and GET /bing.
package fetcher

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/multi/internal/pkg/domain/service/page/fetcher"
)

// Response is the body of GET /google and GET /bing.
type Response struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

// Handler fetches a page and tells how it answered.
type Handler struct {
	google fetcher.Service
	bing   fetcher.Service
}

// New provides the handler to the generated rest door with both fetchers.
// @Provide(index=0)
// @Inject(name=google, index=0)
// @Inject(name=bing, index=1)
func New(google, bing fetcher.Service) *Handler {
	return &Handler{google: google, bing: bing}
}

// Google fetches Google.
// @REST(method=GET, path=/google)
func (h *Handler) Google(c echo.Context) error {
	return fetch(c, h.google)
}

// Bing fetches Bing.
// @REST(method=GET, path=/bing)
func (h *Handler) Bing(c echo.Context) error {
	return fetch(c, h.bing)
}

func fetch(c echo.Context, svc fetcher.Service) error {
	page, err := svc.Fetch(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{URL: page.URL, Status: page.Status})
}
