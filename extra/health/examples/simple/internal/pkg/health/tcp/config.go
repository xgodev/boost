package tcp

import "github.com/xgodev/boost/wrapper/config"

const addr = "app.dependency.addr"

func init() {
	config.Add(addr, "localhost:6379", "host:port of the dependency the tcp checker dials")
}
