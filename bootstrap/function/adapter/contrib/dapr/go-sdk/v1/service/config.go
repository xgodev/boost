package service

import (
	"github.com/dapr/go-sdk/service/common"
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root = adapter.Root + ".dapr.service"
	subs = root + ".subscriptions"
)

func init() {
	config.Add(subs, []common.Subscription{}, "dapr event subscriptions")
}
