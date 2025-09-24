package pubsub

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root                  = "boost.wrapper.publisher.driver.pubsub.v2"
	logRoot               = ".log"
	orderingKey           = ".orderingKey"
	level                 = logRoot + ".level"
	publishTimeout        = ".publishTimeout"
	settings              = ".settings"
	delayThreshold        = settings + ".delayThreshold"
	countThreshold        = settings + ".countThreshold"
	timeout               = settings + ".timeout"
	flowControlSettings   = settings + ".flowControlSettings"
	limitExceededBehavior = flowControlSettings + ".limitExceededBehavior"
	enableCompression     = settings + ".enableCompression"
)

func init() {
	ConfigAdd(root)
}

//func getMaxProcs() int {
//	val := os.Getenv("GOMAXPROCS")
//	if val == "" {
//		return runtime.NumCPU() // fallback padrão
//	}
//	n, err := strconv.Atoi(val)
//	if err != nil || n <= 0 {
//		return runtime.NumCPU() // fallback se inválido
//	}
//	return n
//}

func ConfigAdd(path string) {
	config.Add(path+level, "DEBUG", "defines log level")
	config.Add(path+orderingKey, false, "defines ordering key")
	config.Add(path+delayThreshold, 10*time.Millisecond, "the maximum duration to wait before sending a batch of messages")
	config.Add(path+countThreshold, 100, "the maximum number of messages to include in a batch")
	config.Add(path+timeout, 60*time.Second, "the maximum duration to block Publish calls")
	config.Add(path+limitExceededBehavior, 1, "behavior when flow control limits are exceeded: Block or Ignore")
	config.Add(path+enableCompression, false, "whether to compress messages before sending")
	config.Add(path+publishTimeout, 60*time.Second, "the maximum duration to wait for a publish to complete")
}
