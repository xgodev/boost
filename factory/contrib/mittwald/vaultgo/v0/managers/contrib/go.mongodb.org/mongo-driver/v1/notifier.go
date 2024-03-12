package mongo

import "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"

// Observer
type Observer interface {
	OnNotify(*mongo.Conn)
}

type Notifier interface {
	Register(Observer)
	Unregister(Observer)
	Notify(*mongo.Conn)
}
