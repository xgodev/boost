package pubsub

import (
	"context"
	"sync"

	p "github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
	"gocloud.dev/pubsub"
)

var topicOnce sync.Once

// TopicModule fx module for gocloud pubsub topic. The topic is shut down on app stop.
func TopicModule() fx.Option {
	options := fx.Options()

	topicOnce.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle) (*pubsub.Topic, error) {
					t, err := p.NewTopic(ctx)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(t.Shutdown))
					return t, nil
				},
			),
		)
	})

	return options
}

var subscriptionOnce sync.Once

// SubscriptionModule fx module for gocloud pubsub subscription. The subscription is shut down on app stop.
func SubscriptionModule() fx.Option {
	options := fx.Options()

	subscriptionOnce.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle) (*pubsub.Subscription, error) {
					s, err := p.NewSubscription(ctx)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(s.Shutdown))
					return s, nil
				},
			),
		)
	})

	return options
}
