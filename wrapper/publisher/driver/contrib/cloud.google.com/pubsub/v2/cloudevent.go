package pubsub

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub/v2"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/model/errors"
)

func generateCloudEvent(event *v2.Event) (*pubsub.Message, error) {

	// Convert event data
	var data map[string]interface{}
	if err := event.DataAs(&data); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("failed to convert event data"))
	}

	// Serialize to JSON
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("failed to marshal event data"))
	}

	// Build attributes
	attrs := map[string]string{
		"ce_specversion": event.SpecVersion(),
		"ce_id":          event.ID(),
		"ce_source":      event.Source(),
		"ce_type":        event.Type(),
		"ce_time":        event.Time().UTC().Format(time.RFC3339),
		"ce_path":        "/",
		"ce_subject":     event.Subject(),
	}

	if ct := event.DataContentType(); ct != "" {
		attrs["content-type"] = ct
	} else {
		attrs["content-type"] = "application/json"
	}

	msg := &pubsub.Message{ID: event.ID(), Data: raw, Attributes: attrs, PublishTime: time.Now()}
	/*	if p.options.OrderingKey {
		if pk, err := p.getPartitionKey(ev); err == nil {
			msg.OrderingKey = pk
		}
	}*/

	return msg, nil
}
