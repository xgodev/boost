# cloudevents nats

App `worker` runs a function on the events published to nats subject
`messages`: it replies with a greeting, which the publisher middleware
publishes on subject `replies` (the noop driver only logs it). The cloudevents nats
adapter reads structured cloudevents.

```shell
docker run -d --name nats -p 4222:4222 nats
CONF=configs/config.yaml go run ./cmd/worker
```

```shell
nats pub messages '{"specversion":"1.0","id":"1","source":"cli","type":"message.sent","datacontenttype":"application/json","data":{"message":"hi"}}'
```
