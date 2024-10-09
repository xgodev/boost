package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/bootstrap/function/middleware/logger"
	"github.com/xgodev/boost/bootstrap/function/middleware/prometheus"
	"github.com/xgodev/boost/bootstrap/function/middleware/recovery"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
)

type Message struct {
	Number *int
}

func main() {

	boost.Start()

	l, _ := logger.NewLogger[*cloudevents.Event]()
	r := recovery.NewRecovery[*cloudevents.Event]()
	p, _ := prometheus.NewPrometheus[*cloudevents.Event]()

	var mids = []middleware.AnyErrorMiddleware[*cloudevents.Event]{r, l, p}

	wrp := middleware.NewAnyErrorWrapper[*cloudevents.Event](context.Background(), "bootstrap", mids...)
	fw := function.Wrapper[*cloudevents.Event](wrp, func(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {

		var msg Message
		if err := in.DataAs(&msg); err != nil {
			return nil, errors.Wrap(err, errors.NotValidf("the event data is"))
		}

		if msg.Number == nil {
			panic("panic")
		}

		out := cloudevents.NewEvent()
		if err := out.SetData(cloudevents.ApplicationJSON, &Message{Number: msg.Number}); err != nil {
			return nil, errors.Wrap(err, errors.Internalf("failed to set data"))
		}

		return &out, nil
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
