# cloudevents http

App `worker` runs a function on the cloudevents posted to port 8080: it
replies with a greeting, which the publisher middleware publishes on subject
`replies` (the noop driver only logs it).

```shell
CONF=configs/config.yaml go run ./cmd/worker
```

```shell
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/cloudevents+json" \
  -d '{
    "specversion": "1.0",
    "id": "214b0346-7386-11ea-b6ae-acde48001122",
    "source": "curl",
    "type": "message.sent",
    "datacontenttype": "application/json",
    "data": {"message": "hi"}
  }'
```
