package simple

import (
	"fmt"

	"github.com/jpfaria/tests/annotated"
)

type Loren struct {
	Text string
}

// Foo title
// @Provide (name=A, index=0)
// @Inject (index=0)
func Foo(ex1 *annotated.Loren) *Loren {
	return &Loren{
		Text: "Hello World",
	}
}

// FooBar title
// @Provide (index=0)
func FooBar() *Loren {
	return &Loren{
		Text: "Hello World",
	}
}

// FooBaz title
// @Provide (name=A, index=0)
func FooBaz() *annotated.Loren {
	return nil
}

// Bar title
// @Inject (index=0)
// @Invoke
func Bar(ex *Loren) {
	fmt.Printf("invoked: %s", ex.Text)
}

// Foz title
// @Inject (name=A, index=0)
// @Inject (name=A, index=1)
// @Inject (index=2)
// @Inject (index=3)
// @MyAnnotation
// @Invoke
func Foz(ex1 *Loren, ex2 *annotated.Loren, ex3 *Loren, ex4 *annotated.Loren) {
	fmt.Printf("invoked: %s", ex1.Text)
}
