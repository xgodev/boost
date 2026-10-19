# bootstrap/function/adapter/contrib/azure/functions/v1

Runs a [`function.Handler[T]`](../../../../../) as an Azure Functions
[custom handler](https://learn.microsoft.com/azure/azure-functions/functions-custom-handlers)
(the `azure` command): an http server the functions host calls on every
invocation.

```go
fxfact.NewApp(
	fx.Provide(func() function.Handler[*event.Event] { return handle }),
	azurefx.Module[*event.Event](),
	functionfx.Module[*event.Event](middlewares),
).Run()
```

Set `boost.bootstrap.function.default` to `azure`, or run `app azure`.

## Function app

`host.json` starts the binary as the custom handler; the host gives it the
port in `FUNCTIONS_CUSTOMHANDLER_PORT`, the default of `port`:

```json
{
  "version": "2.0",
  "customHandler": {
    "description": { "defaultExecutablePath": "app", "arguments": ["azure"] }
  }
}
```

Each function is a folder with a `function.json`; the host posts its
invocations to `/<folder>`:

```json
{
  "bindings": [
    { "type": "queueTrigger", "direction": "in", "name": "item", "queueName": "orders", "connection": "AzureWebJobsStorage" },
    { "type": "queue", "direction": "out", "name": "$return", "queueName": "orders-done", "connection": "AzureWebJobsStorage" }
  ]
}
```

## Triggers

The trigger binding (`trigger`, or the only input) becomes the CloudEvents:

| Trigger value | CloudEvent |
|---|---|
| a CloudEvent (JSON, or JSON in a message string) | that event |
| an Event Grid event | id, topic as source, subject, event type as type, event time, data |
| an http request (http trigger) | the CloudEvent of its headers or body; otherwise the body as data, the url (without query) as source, type `azure.functions.http` |
| anything else (queue, Service Bus, Event Hubs, timer, blob…) | the value as data, the message id of the metadata (else the invocation id) as id, source `azure://functions/<function>`, type `azure.functions.invoke`, the function as subject |
| an array (batch, `cardinality: many`) | one event per item, ids `<id>-<index>` |

A value that is not JSON becomes `text/plain` data. Every event carries the
extension `azureinvocationid`.

## Results

The result goes to the output binding `output`: `$return`, the default, is
the return value of the function; another name is that output binding. A
batch returns the array of the non-nil results.

An http trigger answers its caller through an `http` output binding (set
`output` to its name, e.g. `res`): the result as JSON with 200, 204 when it
is nil, or the error as a [`restresponse.Error`](../../../../../../../model/restresponse)
with its status. With `enableForwardingHttpRequest` the host forwards the
request itself to `/<routePrefix>/<route>`, and the answer is written to it the
same way.

## Errors

The items of a batch are handled up to `maxWorkers` at the same time. When
any of them fails, the invocation fails (500), so the trigger retries it; the
others are still handled. An error registered with
`errors.Ignore(err, errors.IgnoreAsSuccess)` does not fail it. The errors go
to the invocation logs.

## Config

| Key | Default | Meaning |
|---|---|---|
| `boost.bootstrap.function.adapter.azure.port` | `FUNCTIONS_CUSTOMHANDLER_PORT`, else `8080` | port the functions host calls |
| `boost.bootstrap.function.adapter.azure.trigger` | `""` | binding name of the trigger; empty takes the only input |
| `boost.bootstrap.function.adapter.azure.output` | `$return` | output binding that receives the result |
| `boost.bootstrap.function.adapter.azure.routePrefix` | `api` | route prefix of the http functions (`extensions.http.routePrefix` in `host.json`) |
| `boost.bootstrap.function.adapter.azure.maxWorkers` | `10` | items of a batch handled at the same time; `1` keeps their order |
| `boost.bootstrap.function.adapter.azure.enabled` | `true` | adds the adapter in the fx module |
