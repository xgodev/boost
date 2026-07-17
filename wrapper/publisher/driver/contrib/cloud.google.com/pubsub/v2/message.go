package pubsub

import (
	"time"

	"cloud.google.com/go/pubsub/v2"
	v2 "github.com/cloudevents/sdk-go/v2"
)

// buildAttrs extracts CloudEvents attributes into a Pub/Sub attributes map.
func buildAttrs(ev *v2.Event) map[string]string {
	attrs := map[string]string{
		"ce_specversion": ev.SpecVersion(),
		"ce_id":          ev.ID(),
		"ce_source":      ev.Source(),
		"ce_type":        ev.Type(),
		"ce_time":        ev.Time().UTC().Format(time.RFC3339),
		"ce_path":        "/",
		"ce_subject":     ev.Subject(),
	}
	if ct := ev.DataContentType(); ct != "" {
		attrs["content-type"] = ct
	} else {
		attrs["content-type"] = "application/json"
	}
	return attrs
}

// getPartitionKey extracts the ordering key extension or uses the event ID.
func getPartitionKey(ev *v2.Event) (string, error) {
	if key, ok := ev.Extensions()["key"]; ok {
		return key.(string), nil
	}
	return ev.ID(), nil
}

func buildMessage(ev *v2.Event, orderingKey bool) *pubsub.Message {
	raw := ev.Data()
	attrs := buildAttrs(ev)
	msg := &pubsub.Message{ID: ev.ID(), Data: raw, Attributes: attrs}

	if orderingKey {
		if pk, err := getPartitionKey(ev); err == nil {
			msg.OrderingKey = pk
		}
	}

	return msg
}
