// Package fetcher fetches a web page.
package fetcher

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/google/internal/pkg/domain/model/entity"
)

// Service fetches the page.
type Service interface {
	Fetch(ctx context.Context) (*entity.Page, error)
}
