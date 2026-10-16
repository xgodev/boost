package app

import (
	"context"
)

// FooFunc Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @RestRouter(path=/, method=POST)
// @RestRouter(path=/foo, method=GET)
// @RestConsume(type=application/json)
// @RestConsume(type=application/yaml)
// @RestProduce(type=application/json)
// @RestQueryParam(name=foo, type=bool, required=true, description= tiam sed efficitur purus)
// @RestQueryParam(name=bar, type=string, required=true, description= tiam sed efficitur purus)
// @RestPathParam(name=foo, type=string, required=true, description=tiam sed efficitur purus)
// @RestPathParam(name=bar, type=int, required=true, description=tiam sed efficitur purus)
// @RestHeader(name=foo, type=string, required=true, description=tiam sed efficitur purus)
// @RestHeader(name=bar, type=string, required=true, description=tiam sed efficitur purus)
// @RestRequestBody(type=github.com/xgodev/boost/inject/examples/simple.Request)
// @RestResponse(code=201, type=github.com/xgodev/boost/inject/examples/simple.Response, description=tiam sed efficitur purus at lacinia magna)
// @IgnoredAnnotation(param=201)
func FooFunc(ctx context.Context, r string) {
}
