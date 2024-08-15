package service

import (
	"github.com/dapr/go-sdk/service/common"
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root          = adapter.Root + ".dapr.service"
	subs          = root + ".subscriptions"
	srvInv        = root + ".service.invocation"
	srvInvEnabled = srvInv + ".enabled"
	srvInvName    = srvInv + ".name"
)

func init() {
	config.Add(subs, []common.Subscription{}, "dapr event subscriptions")
	config.Add(srvInvEnabled, false, "dapr service invocation enabled")
	config.Add(srvInvName, "/events", "dapr service invocation name")
}
