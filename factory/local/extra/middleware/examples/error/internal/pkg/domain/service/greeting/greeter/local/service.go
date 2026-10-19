// Package local greets through the middleware wrapper the app provides, with
// a fallback for the languages it does not know.
package local

import (
	"context"

	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/factory/local/extra/middleware/examples/error/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/local/extra/middleware/examples/error/internal/pkg/domain/service/greeting/greeter"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

var greetings = map[string]string{"en": "hello", "pt": "olá", "es": "hola"}

type service struct {
	wrapper *middleware.AnyErrorWrapper[*entity.Greeting]
}

// New provides the greeter on the wrapper its calls run through.
// @Provide(index=0)
func New(wrapper *middleware.AnyErrorWrapper[*entity.Greeting]) greeter.Service {
	return &service{wrapper: wrapper}
}

// Greet fails for a language it does not know, and the hystrix circuit
// breaker of the wrapper runs the fallback, which greets in English.
func (s *service) Greet(ctx context.Context, name, lang string) (*entity.Greeting, error) {
	return s.wrapper.Exec(ctx, name,
		func(ctx context.Context) (*entity.Greeting, error) {
			greeting, ok := greetings[lang]
			if !ok {
				return nil, errors.NotFoundf("greeting in language %q", lang)
			}
			return &entity.Greeting{Text: greeting + ", " + name}, nil
		},
		func(ctx context.Context, _ *entity.Greeting, err error) (*entity.Greeting, error) {
			log.FromContext(ctx).Infof("greeting in English: %v", err)
			return &entity.Greeting{Text: greetings["en"] + ", " + name}, nil
		})
}
