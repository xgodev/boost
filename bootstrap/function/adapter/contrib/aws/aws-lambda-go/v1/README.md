# bootstrap/function/adapter/contrib/aws/aws-lambda-go/v1

Runs a [`function.Handler[T]`](../../../../../) as an AWS Lambda function
(the `lambda` command), on [aws-lambda-go](https://github.com/aws/aws-lambda-go).

```go
fxfact.NewApp(
	fx.Provide(func() function.Handler[*event.Event] { return handle }),
	lambdafx.Module[*event.Event](),
	functionfx.Module[*event.Event](middlewares),
).Run()
```

Set `boost.bootstrap.function.default` to `lambda`, or run `app lambda`.

## Triggers

Every record of the invocation becomes one CloudEvent:

| Trigger | CloudEvent |
|---|---|
| SQS | the body, when it is a CloudEvent; an SNS message delivered to SQS is unwrapped first; otherwise the body is the data and the message id the id |
| SNS | the message, when it is a CloudEvent; otherwise the message is the data, the message id the id and the SNS type the type |
| Kinesis | the record data, when it is a CloudEvent; otherwise the data |
| S3 | the S3 entity as data, the object key as id |
| DynamoDB Streams | the stream record as data; the user identity in `useridentitytype` and `useridentityprincipalid` |
| EventBridge (schedules, CloudWatch Events) | the detail as data, the detail type as type |

A payload that is not JSON becomes `text/plain` data. When the payload is not a
CloudEvent, id, type and source fall back to the record's event id, event name
and event source. Every event carries the extensions `awsrequestid` and
`invokedfunctionarn`.

A trigger it does not know fails the invocation with a not-implemented error.

## Errors

The records of one invocation are handled up to `maxWorkers` at the same
time. When any of them fails, the invocation fails, so Lambda delivers the
batch again; the others are still handled. An error registered with
`errors.Ignore(err, errors.IgnoreAsSuccess)` does not fail it.

## Config

| Key | Default | Meaning |
|---|---|---|
| `boost.bootstrap.function.adapter.lambda.skip` | `false` | acknowledges every invocation without calling the handler |
| `boost.bootstrap.function.adapter.lambda.maxWorkers` | `10` | records of one invocation handled at the same time; `1` keeps their order |
