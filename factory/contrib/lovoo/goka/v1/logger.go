package goka

import (
	"github.com/xgodev/boost/wrapper/log"
)

type Logger struct {
}

func (s *Logger) Print(msgs ...interface{}) {
	log.Printf("%v", msgs)
}

func (s *Logger) Println(msgs ...interface{}) {
	log.Printf("%v", msgs)
}

func (s *Logger) Printf(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}

func (s *Logger) Debugf(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}
