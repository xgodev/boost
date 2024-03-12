package simple

import (
	"context"
	"net/http"
)

// ExampleStruct Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @Package github.com/xgodev/boost/boost
// @RelativePackage examples/simple
// @App xpto
// @HandlerType HTTP
// @Type Interface
type ExampleStruct struct {
}

func (t *ExampleStruct) FooStructMethod(ctx context.Context, r *http.Request) (interface{}, error) {
	return Response{
		Message: "Hello world",
	}, nil
}

// FooMethod Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @Package github.com/xgodev/boost/boost
// @RelativePackage examples/simple
// @App xpto
// @HandlerType HTTP
// @Type Function
// @Path /foo
// @Path /
// @Method POST
// @Consume application/json
// @Consume application/yaml
// @Produce application/json
// @Param query foo bool true tiam sed efficitur purus
// @Param query bar string true tiam sed efficitur purus
// @Param path foo string tiam sed efficitur purus
// @Param path bar string tiam sed efficitur purus
// @Param header foo string true tiam sed efficitur purus
// @Param header bar string true tiam sed efficitur purus
// @Body github.com/xgodev/boost/boost/examples/simple.Request
// @Response 201 github.com/xgodev/boost/boost/examples/simple.Response tiam sed efficitur purus, at lacinia magna
func FooMethod(ctx context.Context, r *http.Request) (interface{}, error) {
	return Response{
		Message: "Hello world",
	}, nil
}

type Response struct {
	Message string
}
