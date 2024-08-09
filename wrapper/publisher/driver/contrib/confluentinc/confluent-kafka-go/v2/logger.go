package confluent

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/xgodev/boost/wrapper/log"
)

type Logger struct {
	level string
	producer *kafka.Producer
}

func NewLogger(producer *kafka.Producer, level string) *Logger {
	return &Logger{level: level, producer: producer}
}

func (s *Logger) Start() {
	// Listen to all the events on the default events channel
	go func() {
		for e := range s.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				// The message delivery report, indicating success or
				// permanent failure after retries have been exhausted.
				// Application level retries won't help since the client
				// is already configured to do that.
				m := ev
				if m.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v", m.TopicPartition.Error)
				} else {
					s.log("Delivered message to topic %s [%d] at offset %v",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				// Generic client instance-level errors, such as
				// broker connection failures, authentication issues, etc.
				//
				// These errors should generally be considered informational
				// as the underlying client will automatically try to
				// recover from any errors encountered, the application
				// does not need to take action on them.
				log.Errorf("Error: %v\n", ev)
			default:
				log.Warnf("Ignored event: %s\n", ev)
			}
		}
	}()
}

func (s *Logger) log(format string, args ...interface{}) {
	switch s.level {
	case "INFO":
		log.Infof(format, args...)
	case "TRACE":
		log.Tracef(format, args...)
	default:
		log.Debugf(format, args...)
	}
}