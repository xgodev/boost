# Elasticsearch v9 Integration for Boost

## Overview

Factory for the [go-elasticsearch](https://github.com/elastic/go-elasticsearch) v9 client and its bulk indexer. It sits next to the v8 factory (`factory/contrib/elastic/go-elasticsearch/v8`), reads the same config keys (`boost.factory.elasticsearch.*`) and keeps the same options, so moving a service from v8 to v9 is an import change.

## Features

- **Client from config**: `NewClient` reads `boost.factory.elasticsearch.*` and pings the node before returning.
- **Bulk indexer**: `NewBulkIndexer` builds an `esutil.BulkIndexer` from `boost.factory.elasticsearch.bulk.*`.
- **Plugins**: `Plugin` runs on the built client, in order; the first error is returned.
- **Health check**: `plugins/local/extra/health` registers a ping checker in `extra/health`.
- **Boost logging**: requests and bulk debug lines go to `wrapper/log`.

## Installation

```go
import (
    "github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v9"
)
```

## Usage

```go
boost.Start()

client, err := elasticsearch.NewClient(ctx)
if err != nil {
    return err
}

res, err := client.Info()
if err != nil {
    return err
}
defer res.Body.Close()
```

Other config roots and explicit options:

```go
client, err := elasticsearch.NewClientWithConfigPath(ctx, "app.search") // app.search.* over the defaults
client, err := elasticsearch.NewClientWithOptions(ctx, &elasticsearch.Options{Addresses: []string{"http://es:9200"}})
```

### Bulk indexer

```go
bi, err := elasticsearch.NewBulkIndexer(ctx, client)
if err != nil {
    return err
}
err = bi.Add(ctx, esutil.BulkIndexerItem{Action: "index", Body: strings.NewReader(`{"sku":"a"}`)})
// ...
err = bi.Close(ctx) // flushes what is left
```

`NewBulkIndexerWithConfigPath(ctx, "app.search", client)` reads `app.search.bulk.*`; what the path leaves out comes from `boost.factory.elasticsearch.bulk.*`.

### Health check

```go
client, err := elasticsearch.NewClient(ctx, health.Register)
```

A ping answered with an error status fails the check as `errors.ServiceUnavailable`.

## Configuration

| Key | Default | Description |
|---|---|---|
| `boost.factory.elasticsearch.addresses` | `["http://127.0.0.1:9200"]` | Elasticsearch nodes |
| `boost.factory.elasticsearch.username` | `""` | HTTP Basic Authentication user |
| `boost.factory.elasticsearch.password` | `""` | HTTP Basic Authentication password (hidden) |
| `boost.factory.elasticsearch.cloudID` | `""` | Elastic Cloud endpoint |
| `boost.factory.elasticsearch.APIKey` | `""` | base64 API key; overrides username and password (hidden) |
| `boost.factory.elasticsearch.CACert` | `""` | PEM-encoded certificate authorities |
| `boost.factory.elasticsearch.retryOnStatus` | `[502, 503, 504]` | statuses that are retried |
| `boost.factory.elasticsearch.disableRetry` | `false` | turn retries off |
| `boost.factory.elasticsearch.maxRetries` | `3` | maximum retries |
| `boost.factory.elasticsearch.retryBackoff` | `5ms` | backoff per attempt, multiplied by the attempt number |
| `boost.factory.elasticsearch.discoverNodesOnStart` | `false` | discover nodes when the client starts |
| `boost.factory.elasticsearch.discoverNodesInterval` | `0` | discover nodes periodically (0 = off) |
| `boost.factory.elasticsearch.enableMetrics` | `false` | collect client metrics |
| `boost.factory.elasticsearch.enableDebugLogger` | `false` | client debug logging |
| `boost.factory.elasticsearch.bulk.index` | `index` | default index of the bulk indexer |
| `boost.factory.elasticsearch.bulk.numWorkers` | number of CPUs | bulk workers |
| `boost.factory.elasticsearch.bulk.flushInterval` | `30s` | flush period |
| `boost.factory.elasticsearch.bulk.flushBytes` | `5242880` | flush threshold in bytes |
| `boost.factory.elasticsearch.bulk.pipeline` | `""` | ingest pipeline |
| `boost.factory.elasticsearch.bulk.timeout` | `2s` | bulk request timeout |
| `boost.factory.elasticsearch.plugins.health.name` | `elasticsearch` | health check name |
| `boost.factory.elasticsearch.plugins.health.description` | `default connection` | health check description |
| `boost.factory.elasticsearch.plugins.health.required` | `true` | the check is required |
| `boost.factory.elasticsearch.plugins.health.enabled` | `true` | register the check |

Environment variables follow the keys: `BOOST_FACTORY_ELASTICSEARCH_ADDRESSES`, `BOOST_FACTORY_ELASTICSEARCH_MAX__RETRIES` (`__` marks the camelCase boundary).

## fx

`fx/modules/factory/contrib/elastic/go-elasticsearch/v9` provides `*elasticsearch.Client` (`Module`) and `esutil.BulkIndexer` (`BulkModule`). Plugins come from an optional `[]elasticsearch.Plugin` in the graph.

## Differences from v8

- Option and plugin errors are returned; v8 logged them with `Fatalf` or panicked.
- A ping answered with an error status (401, 503, ...) fails with `errors.ServiceUnavailable`, in the client and in the health check; v8 accepted it.
- `retryBackoff` comes from the options, so `NewClientWithConfigPath` and `NewClientWithOptions` honour it; v8 read the global key.
- `enableRetryOnTimeout` is gone; it did nothing in v8 either.
- `NewBulkIndexerWithConfigPath` keeps the `boost.factory.elasticsearch.bulk.*` defaults the path leaves out.
- `APIKey` is hidden in the config listing, like `password`.
- `health.NewHealth` returns `(*Health, error)`.
