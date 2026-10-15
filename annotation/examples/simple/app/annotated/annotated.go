package annotated

import (
	"context"
	"net/http"
)

// Loren Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyStructAnnotation(param=loren)
type Loren struct {
	Path string
}

// LorenMethod Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyMethodAnnotation(param=loren)
func LorenMethod(ctx context.Context, r *http.Request) (*Loren, error) {
	return &Loren{Path: r.URL.Path}, nil
}
