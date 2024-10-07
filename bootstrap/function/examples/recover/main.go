package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/bootstrap/function/middleware/recovery"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
)

type Message struct {
	Number *int
}

func main() {

	boost.Start()

	wrp := middleware.NewAnyErrorWrapper[string](context.Background(), "bootstrap", recovery.NewRecover[string]())
	fw := function.Wrapper[string](wrp, func(ctx context.Context, in cloudevents.Event) (string, error) {

		var msg Message
		if err := in.DataAs(&msg); err != nil {
			return "invalid", errors.Wrap(err, errors.NotValidf("the event data is"))
		}

		if msg.Number == nil {
			panic("panic")
		}

		return fmt.Sprintf("Hello, World! %v", msg.Number), nil
	})

	two := 2
	three := 3
	four := 4

	var numbers = []*int{nil, &two, &three, &four, nil}

	for _, number := range numbers {

		ev := cloudevents.NewEvent()
		err := ev.SetData("", &Message{Number: number})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		res, err := fw(context.Background(), ev)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(res)
	}
}
