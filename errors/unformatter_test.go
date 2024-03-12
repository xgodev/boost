package errors

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UnformatterSuite struct {
	suite.Suite
}

func TestUnformatterSuite(t *testing.T) {
	suite.Run(t, new(UnformatterSuite))
}

func (s *UnformatterSuite) TestErrorStack() {

	tt := []struct {
		name string
		got  func() error
	}{
		{
			name: "single error stack",
			got: func() error {
				return New("first error") //err single
			},
		},
		{
			name: "wrapped error",
			got: func() error {
				err := New("first error")                    //err wrapped-0
				return Wrap(err, newError("detailed error")) //err wrapped-1
			},
		},
		{
			name: "annotated error",
			got: func() error {
				err := New("first error")          //err annotated-0
				return Annotate(err, "annotation") //err annotated-1
			},
		},
		{
			name: "annotated wrapped error",
			got: func() error {
				err := Errorf("first error")                  //err ann-wrap-0
				err = Wrap(err, fmt.Errorf("detailed error")) //err ann-wrap-1
				return Annotatef(err, "annotated")            //err ann-wrap-2
			},
		},
		{
			name: "traced, and annotated",
			got: func() error {
				err := New("first error")           //err stack-0
				err = Trace(err)                    //err stack-1
				err = Annotate(err, "some context") //err stack-2
				err = Trace(err)                    //err stack-3
				err = Annotate(err, "more context") //err stack-4
				return Trace(err)                   //err stack-5
			},
		},
		{
			name: "uncomparable, wrapped with a value error",
			got: func() error {
				err := newNonComparableError("first error") //err mixed-0
				err = Trace(err)                            //err mixed-1
				err = Wrap(err, newError("value error"))    //err mixed-2
				err = Maskf(err, "masked")                  //err mixed-3
				err = Annotate(err, "more context")         //err mixed-4
				return Trace(err)                           //err mixed-5
			},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := t.got()
			stack := ErrorStack(got)
			tracer, _ := got.(tracer)
			stackTrace := tracer.StackTrace()
			s.Assert().True(reflect.DeepEqual(stackTrace, strings.Split(stack, "\n")), "got  %v\nwant %v", stackTrace, strings.Split(stack, "\n"))
		})
	}
}

type tracer interface {
	StackTrace() []string
}
