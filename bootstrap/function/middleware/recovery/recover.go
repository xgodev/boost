package recovery

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/wrapper/log"
)

type Recovery[T any] struct {
}

func NewRecovery[T any]() *Recovery[T] {
	return &Recovery[T]{}
}

func NewAnyErrorMiddleware[T any]() middleware.AnyErrorMiddleware[T] {
	return NewRecovery[T]()
}

func (c *Recovery[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (res T, err error) {
	defer func() {
		if r := recover(); r != nil {
			runtimeCallback := FancyHandleError(r)

			log.FromContext(ctx.GetContext()).WithTypeOf(*c).WithField("callers", runtimeCallback).Errorf("recovering: %v", r)
			err = fmt.Errorf(runtimeCallback)
		}
	}()

	res, err = ctx.Next(exec, fallbackFunc)

	//if r := recover(); r != nil {
	//	log.FromContext(ctx.GetContext()).WithTypeOf(*c).Errorf("recovering: %v", r)
	//	err = fmt.Errorf(FancyHandleError(r))
	//}
	return res, err
}

// this logs the function name as well.
func FancyHandleError(err any) string {
	// notice that we're using 1, so it will actually log the where
	// the error happened, 0 = this function, we don't want that.
	//pc, filename, line, _ := runtime.Caller(1)
	var pcs [10]uintptr
	n := runtime.Callers(1, pcs[:])
	iter := runtime.CallersFrames(pcs[:n])

	b := strings.Builder{}
	for {
		f, more := iter.Next()
		//fmt.Printf("  %s %s:%d\n", f.Function, f.File, f.Line)

		b.WriteString(fmt.Sprintf("%s %s:%d;", f.Function, f.File, f.Line))
		if !more {
			break
		}
	}

	return b.String()

	//log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
	//return fmt.Sprintf("[error] in %s[%s:%d] %v %v", runtime.FuncForPC(pc).Name(), filename, line, err, n)

}
