# middlewares

App `worker` echoes the number of the cloudevents posted to port 8080,
behind the prometheus, logger and recovery middlewares. An event without
number panics: recovery turns the panic into an error, the adapter answers
500, the logger logs the error and prometheus counts it in
`boost_function_messages_processed_total{status="error"}`.

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
