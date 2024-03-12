package errors

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorSuite struct {
	suite.Suite
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (s *ErrorSuite) TestIsForbidden() {

	tt := []struct {
		name string
		got  func() error
		want string
	}{
		{
			name: "uncomparable errors",
			got: func() error {
				err := Annotatef(newNonComparableError("uncomparable"), "annotation")
				return Annotatef(err, "another")
			},
			want: "another: annotation: uncomparable",
		},
		{
			name: "Errorf",
			got: func() error {
				return Errorf("first error")
			},
			want: "first error",
		},
		{
			name: "annotated error",
			got: func() error {
				err := Errorf("first error")
				return Annotatef(err, "annotation")
			},
			want: "annotation: first error",
		},
		{
			name: "test annotation format",
			got: func() error {
				err := Errorf("first %s", "error")
				return Annotatef(err, "%s", "annotation")
			},
			want: "annotation: first error",
		},
		{
			name: "wrapped error",
			got: func() error {
				err := newError("first error")
				return Wrap(err, newError("detailed error"))
			},
			want: "detailed error",
		},
		{
			name: "wrapped annotated error",
			got: func() error {
				err := Errorf("first error")
				err = Annotatef(err, "annotated")
				return Wrap(err, fmt.Errorf("detailed error"))
			},
			want: "detailed error",
		},
		{
			name: "annotated wrapped error",
			got: func() error {
				err := Errorf("first error")
				err = Wrap(err, fmt.Errorf("detailed error"))
				return Annotatef(err, "annotated")
			},
			want: "annotated: detailed error",
		},
		{
			name: "traced, and annotated",
			got: func() error {
				err := New("first error")
				err = Trace(err)
				err = Annotate(err, "some context")
				err = Trace(err)
				err = Annotate(err, "more context")
				return Trace(err)
			},
			want: "more context: some context: first error",
		},
		{
			name: "traced, and annotated, masked and annotated",
			got: func() error {
				err := New("first error")
				err = Trace(err)
				err = Annotate(err, "some context")
				err = Maskf(err, "masked")
				err = Annotate(err, "more context")
				return Trace(err)
			},
			want: "more context: masked: some context: first error",
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := t.got()
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

// This is an uncomparable error type, as it is a struct that supports the
// error interface (as opposed to a pointer type).
type error_ struct {
	info  string
	slice []string
}

// Create a non-comparable error
func newNonComparableError(message string) error {
	return error_{info: message}
}

func (e error_) Error() string {
	return e.info
}

func newError(message string) error {
	return testError{message}
}

// The testError is a value type error for ease of seeing results
// when the test fails.
type testError struct {
	message string
}

func (e testError) Error() string {
	return e.message
}
