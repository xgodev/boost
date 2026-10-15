package app

import (
	"context"
	"net/http"
	"net/url"

	"github.com/xgodev/boost/annotation/examples/simple/app/annotated"
)

// ExampleStruct Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyStructAnnotation(param=xpto)
type ExampleStruct struct {
}

// FooFunc Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @MyFuncAnnotation(param=xpto)
func FooFunc(ctx context.Context, value string) (h *annotated.Loren, err error) {
	return annotated.LorenMethod(ctx, &http.Request{URL: &url.URL{Path: value}})
}
