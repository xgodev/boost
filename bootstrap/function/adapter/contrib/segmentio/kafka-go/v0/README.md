# bootstrap/function/adapter/contrib/segmentio/kafka-go/v0

Runs a [`function.Handler[T]`](../../../../../) on the messages of a kafka
consumer group (the `kafkago` command), on
[segmentio/kafka-go](https://github.com/segmentio/kafka-go): pure Go, no cgo.
The reader comes from the [kafka-go factory](../../../../../../../factory/contrib/segmentio/kafka-go/v0).

```go
fxfact.NewApp(
	fx.Provide(func() function.Handler[*event.Event] { return handle }),
	kafkagofx.Module[*event.Event](),
	functionfx.Module[*event.Event](middlewares),
).Run()
```

Set `boost.bootstrap.function.default` to `kafkago`, or run `app kafkago`.
The consumer group is `boost.factory.kafkago.reader.groupId`, required;
brokers, offsets and fetch sizes are the other `boost.factory.kafkago.reader.*`
keys.

## Messages

Every message becomes one CloudEvent:

| Message | CloudEvent |
|---|---|
| `ce_specversion` header (binary mode) | the `ce_*` headers as attributes, other `ce_<name>` as extensions, the value as data in `content-type` |
| `content-type: application/cloudevents+json` (structured mode) | the value |
| no `content-type` and the value is a CloudEvent | the value |
| anything else | the value as data; id `<topic>-<partition>-<offset>`, source `kafka://<topic>/<partition>`, type `kafka.message`, subject the topic, time the message time |

This is the format the kafka-go and confluent publisher drivers write, so
their events arrive unchanged.

## Order, retries and commits

The messages of a partition are handled in order, one at a time; partitions
run in parallel. A message is committed after its handler succeeds. A failed
one is retried, waiting `backoffBase`, doubled on each retry up to
`maxBackoff`; after `retryLimit` retries it is committed and skipped. A
message that cannot be decoded is committed and skipped.

On shutdown the subscriber stops fetching, waits for the messages being
handled and leaves the group. A message interrupted in its backoff stays
uncommitted, so the group delivers it again.

## Config

| Key | Default | Meaning |
|---|---|---|
| `boost.bootstrap.function.adapter.kafkago.topics` | `[changeme]` | topics read by the consumer group |
| `boost.bootstrap.function.adapter.kafkago.buffer` | `100` | messages queued per partition while its handler is busy |
| `boost.bootstrap.function.adapter.kafkago.backoff` | `true` | waits between retries of a failed message |
| `boost.bootstrap.function.adapter.kafkago.backoffBase` | `1s` | wait before the first retry |
| `boost.bootstrap.function.adapter.kafkago.maxBackoff` | `20s` | max wait between retries |
| `boost.bootstrap.function.adapter.kafkago.retryLimit` | `-1` | retries before the message is committed and skipped; `-1` retries forever |
| `boost.bootstrap.function.adapter.kafkago.enabled` | `true` | adds the adapter in the fx module |
