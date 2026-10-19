# Wrapper

The Wrapper package provides adapters and abstractions for external libraries, allowing standardized and flexible integration with the Boost framework. This package is essential for ensuring interoperability and consistency when using third-party libraries.

## Overview

Wrapper acts as an abstraction layer between application code and external libraries, offering consistent interfaces and facilitating the replacement of underlying implementations without affecting client code. This promotes code maintainability and testability.

## Main Components

### Cache

The `cache` subpackage provides abstractions for caching systems, including:

- **Codec**: Serializers and deserializers for different formats (binary, gob, json, string)
- **Driver**: Implementations for different cache providers (Redis, BigCache, FreeCache)
- **Plugins**: Extensions for additional cache functionalities

This component allows efficient storage and retrieval of cached data, with support for multiple formats and backends.

### Config

The `config` subpackage offers a unified interface for application configuration, with support for:

- Multiple configuration sources (files, environment variables, flags)
- Different formats (YAML, JSON, ENV)
- Value validation and transformation

This facilitates configuration management across different environments and deployment scenarios.

Drivers: [`contrib/knadh/koanf/v1`](config/contrib/knadh/koanf/v1) (wired by `boost.Start()`) and
[`contrib/knadh/koanf/v2`](config/contrib/knadh/koanf/v2) (set with `config.Set(koanf.New())`).

### Log

The `log` subpackage provides abstractions for logging systems, with support for:

- Multiple log levels (debug, info, warn, error)
- Different output formats (text, JSON, CloudWatch)
- Integration with popular libraries (Zap, Zerolog, Logrus)

This component allows consistent and configurable recording of application events and information.

### Publisher

The `publisher` subpackage offers abstractions for message publishing/subscription systems, including:

- Drivers for different providers (Google Cloud Pub/Sub, Kafka, NATS, AWS SNS, SQS and Kinesis)
- Middleware for message interception and modification
- Utilities for monitoring and metrics

This facilitates asynchronous communication between components and services.

`Publisher.Publish` returns an error when the driver fails or when any event
fails. A per-event failure comes back as one `Internal` error that joins every
failed event's error (tagged with its event ID), so `errors.Is`/`errors.As`
still reach each cause. Middlewares keep receiving the per-event
`[]PublishOutput` (with each event's `Error`) and a nil chain error; the
aggregate error is built after the middleware chain.

The Kafka driver (`driver/contrib/confluentinc/confluent-kafka-go/v2`) returns
only after every produced message has its delivery report, bounded by the
`ctx` passed to `Publish`: a failed delivery (`TopicPartition.Error`, e.g.
message timeout or unknown partition) or a `ctx` that ends before the report
arrives becomes that event's `Error`. Reports go to a delivery channel per
`Publish` call, so they no longer reach the producer's `Events()` channel; the
driver logs each delivery at its own `log.level`.

The kafka-go driver (`driver/contrib/segmentio/kafka-go/v0`, on the writer of
`factory/contrib/segmentio/kafka-go/v0`) writes all events in one
`WriteMessages`, in CloudEvents binary mode: the data is the value, the
attributes and extensions are `ce_*` headers (`ce_time` in RFC3339) plus
`content-type`. The topic is the event subject; the key is the `key`
extension, else `partitionkey`, else the event ID. An event without subject,
or one kafka-go reports in `kafka.WriteErrors`, gets its own `Error`; any
other write error fails every event of the call. fx:
`fx/modules/wrapper/publisher/driver/contrib/segmentio/kafka-go/v0`, switched by
`boost.wrapper.publisher.driver.kafkago.enabled`.

The AWS drivers (`driver/contrib/aws/{sns,sqs,kinesis}/v1`, on the clients of
`factory/contrib/aws/aws-sdk-go-v2/v1/client/...`) send the event subject as
the destination: the SNS topic ARN, the SQS queue name (URL resolved and
cached) or the Kinesis stream name. The payload is the whole CloudEvent as
JSON, or only its data when the `target` extension is `data`
(`util.Payload`). SQS takes `MessageGroupId` from the `group` extension;
Kinesis takes the partition key from the `partitionkey` extension (the event
ID without it) and puts the events of one stream in one `PutRecords`, so a
failed stream fails all its events. Retries are the AWS SDK retryer's
(`boost.factory.aws.retryer.*`). fx: `fx/modules/wrapper/publisher/driver/contrib/aws/<service>/v1`,
switched by `boost.wrapper.publisher.driver.<service>.enabled`.

## How to Use

To use the Wrapper package in your application:

1. Import the specific subpackage containing the desired abstraction
2. Configure the specific implementation you want to use
3. Use the abstract interface in your code, independent of the underlying implementation

## Integration with Other Packages

Wrapper integrates with other components of the Boost framework:

- Bootstrap: For component initialization and configuration
- Factory: For creating instances of specific implementations
- Extra: For additional functionalities like middleware

## Extensibility

The Wrapper package was designed to be easily extensible. To add support for new libraries:

1. Create a new subpackage in `contrib` with the name of the library
2. Implement the necessary abstract interfaces
3. Add appropriate tests and documentation

## Contribution

Contributions to the Wrapper package are welcome. When contributing, make sure to follow the code standards and add appropriate tests for new functionalities.
