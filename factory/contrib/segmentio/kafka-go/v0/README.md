# kafka-go

Factory for [segmentio/kafka-go](https://github.com/segmentio/kafka-go), Kafka
in pure Go (no cgo, unlike `confluentinc/confluent-kafka-go/v2`).

```go
w, err := kafka.NewWriter(ctx)   // *kafka.Writer, no Topic: each message names it
r, err := kafka.NewReader(ctx)   // *kafka.Reader of reader.topic
c, err := kafka.NewConn(ctx)     // *kafka.Conn to the first broker
```

Every constructor has `WithConfigPath` and `WithOptions` variants and returns
an error instead of panicking; an invalid reader config, which
`kafka.NewReader` would panic on, is an error too. The writer and the reader
connect on first use.

Plugins: `WriterPlugin` (`factory.Plugin[*kafka.Writer]`, e.g. a `Transport`
with SASL/TLS), `ReaderPlugin` (`factory.Plugin[*kafka.ReaderConfig]`, runs
before the reader is built, e.g. a `Dialer`) and `DialerPlugin`
(`factory.Plugin[*kafka.Dialer]`). kafka-go logs go to the boost logger: info
at trace, errors at error.

## Config (`boost.factory.kafkago`)

| Key | Default | |
|---|---|---|
| `brokers` | `[localhost:9092]` | broker addresses |
| `writer.balancer` | `crc32` | `roundrobin`, `leastbytes`, `hash`, `crc32`, `murmur2`; `crc32` is librdkafka's partitioner, so a key lands on the same partition as with confluent-kafka-go |
| `writer.requiredAcks` | `all` | `none`, `one`, `all` |
| `writer.compression` | empty | `gzip`, `snappy`, `lz4`, `zstd` |
| `writer.maxAttempts` | `10` | |
| `writer.batchSize` | `100` | |
| `writer.batchBytes` | `1048576` | |
| `writer.batchTimeout` | `10ms` | kafka-go's 1s would delay every synchronous publish |
| `writer.readTimeout`, `writer.writeTimeout` | `10s` | |
| `writer.async` | `false` | |
| `writer.allowAutoTopicCreation` | `false` | |
| `reader.groupId` | empty | empty reads `reader.partition` alone |
| `reader.topic` | empty | required by the reader |
| `reader.partition` | `0` | |
| `reader.queueCapacity` | `100` | |
| `reader.minBytes` / `reader.maxBytes` | `1` / `10485760` | |
| `reader.maxWait`, `reader.readBatchTimeout` | `2s` | |
| `reader.commitInterval` | `0` | 0 commits synchronously |
| `reader.startOffset` | `last` | `first`, `last` (group with no commit) |
| `conn.type` | `leader` | `leader` (of `conn.topic`/`conn.partition`), `partition` (on the first broker, no lookup), `server` |
| `conn.network` | `tcp` | |
| `conn.topic`, `conn.partition` | empty, `0` | |
| `conn.timeout` | `10s` | dial timeout |

Env vars join camelCase words with `__`, e.g.
`BOOST_FACTORY_KAFKAGO_WRITER_REQUIRED__ACKS=one`.

## fx

`fx/modules/factory/contrib/segmentio/kafka-go/v0`: `WriterModule` and
`ReaderModule` provide the writer and the reader and close them on app stop
(which flushes pending messages). The publisher driver is
`wrapper/publisher/driver/contrib/segmentio/kafka-go/v0`.
