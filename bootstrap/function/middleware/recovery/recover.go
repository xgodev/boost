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
			runtimeCallback := GenerateStackTrace(r)

			log.FromContext(ctx.GetContext()).WithTypeOf(*c).WithField("callers", runtimeCallback).Errorf("recovering: %v", r)
			err = fmt.Errorf(runtimeCallback)
		}
	}()

	res, err = ctx.Next(exec, fallbackFunc)

	return res, err
}

func GenerateStackTrace(err any) string {
	var pcs [10]uintptr
	n := runtime.Callers(1, pcs[:])
	iter := runtime.CallersFrames(pcs[:n])

	b := strings.Builder{}
	for {
		f, more := iter.Next()
		b.WriteString(fmt.Sprintf("%s %s:%d;", f.Function, f.File, f.Line))
		if !more {
			break
		}
	}

	return b.String()
}
