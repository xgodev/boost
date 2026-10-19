package kafka

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/model/errors"
)

const defaultContentType = "application/json"

// decode turns a record into the event of the handler: CloudEvents binary
// mode when it has a ce_specversion header, structured mode when its
// content-type is application/cloudevents+json (or it has no content-type and
// the value is a CloudEvent), else the raw value as data.
func decode(msg kafka.Message) (cloudevents.Event, error) {

	contentType := ""
	binary := false
	for _, h := range msg.Headers {
		switch h.Key {
		case "content-type":
			contentType = string(h.Value)
		case "ce_specversion":
			binary = true
		}
	}

	switch {
	case binary:
		return decodeBinary(msg, contentType)
	case strings.HasPrefix(contentType, cloudevents.ApplicationCloudEventsJSON):
		return decodeStructured(msg)
	case contentType == "":
		if in, err := decodeStructured(msg); err == nil {
			return in, nil
		}
	}

	return decodeRaw(msg, contentType)
}

func decodeBinary(msg kafka.Message, contentType string) (cloudevents.Event, error) {

	in := cloudevents.NewEvent()

	for _, h := range msg.Headers {
		name, ok := strings.CutPrefix(h.Key, "ce_")
		if !ok {
			continue
		}
		value := string(h.Value)
		switch name {
		case "specversion":
			in.SetSpecVersion(value)
		case "id":
			in.SetID(value)
		case "source":
			in.SetSource(value)
		case "type":
			in.SetType(value)
		case "subject":
			in.SetSubject(value)
		case "dataschema":
			in.SetDataSchema(value)
		case "time":
			if t, ok := parseTime(value); ok {
				in.SetTime(t)
			}
		default:
			in.SetExtension(name, value)
		}
	}

	if err := setData(&in, contentType, msg.Value); err != nil {
		return in, err
	}

	if err := in.Validate(); err != nil {
		return in, errors.NotValidf("cloudevent in kafka headers: %s", err)
	}

	return in, nil
}

// parseTime reads RFC3339 and also the time.Time.String format the confluent
// publisher driver writes in ce_time.
func parseTime(value string) (time.Time, bool) {
	if t, err := types.ParseTime(value); err == nil {
		return t, true
	}
	if t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", value); err == nil {
		return t, true
	}
	return time.Time{}, false
}

func decodeStructured(msg kafka.Message) (cloudevents.Event, error) {
	var in cloudevents.Event
	if err := json.Unmarshal(msg.Value, &in); err != nil {
		return in, errors.NotValidf("structured cloudevent: %s", err)
	}
	if err := in.Validate(); err != nil {
		return in, errors.NotValidf("structured cloudevent: %s", err)
	}
	return in, nil
}

func decodeRaw(msg kafka.Message, contentType string) (cloudevents.Event, error) {

	in := cloudevents.NewEvent()
	in.SetID(fmt.Sprintf("%s-%d-%d", msg.Topic, msg.Partition, msg.Offset))
	in.SetSource(fmt.Sprintf("kafka://%s/%d", msg.Topic, msg.Partition))
	in.SetType("kafka.message")
	in.SetSubject(msg.Topic)
	if msg.Time.IsZero() {
		in.SetTime(time.Now())
	} else {
		in.SetTime(msg.Time)
	}

	return in, setData(&in, contentType, msg.Value)
}

func setData(in *cloudevents.Event, contentType string, value []byte) error {
	if contentType == "" {
		contentType = defaultContentType
	}
	if err := in.SetData(contentType, value); err != nil {
		return errors.NotValidf("kafka record data: %s", err)
	}
	return nil
}
