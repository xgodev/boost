# recovery

App `worker` echoes the number of the cloudevents posted to port 8080,
behind the recovery middleware. An event without number panics: recovery
turns the panic into an error, the adapter answers 500 and the worker keeps
running.

```shell
CONF=configs/config.yaml go run ./cmd/worker
```

A number is echoed; an event without number answers 500:

```shell
curl -i -X POST http://localhost:8080 \
  -H "Content-Type: application/cloudevents+json" \
  -d '{"specversion":"1.0","id":"1","source":"curl","type":"number.sent","datacontenttype":"application/json","data":{"number":2}}'
curl -i -X POST http://localhost:8080 \
  -H "Content-Type: application/cloudevents+json" \
  -d '{"specversion":"1.0","id":"2","source":"curl","type":"number.sent","datacontenttype":"application/json","data":{}}'
```
