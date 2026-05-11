---
name: boost-factory-aws
description: "Use when constructing an AWS SDK v2 aws.Config in a Go service via github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1. Covers NewConfig + variants and the umbrella SDK pattern (one shared aws.Config, then per-service clients constructed under factory/.../client/<service>/). Service-specific examples ship under examples/{kinesis,s3,sns,sqs}/. Triggers on imports under factory/contrib/aws/aws-sdk-go-v2/, on questions about AWS credentials, regions, or sharing an aws.Config across multiple AWS services in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. For wrapper-side drivers targeting AWS services, see `boost-maintainer` (multi-service SDK rule: drivers go under `<area>/aws/<service>/v<major>/`, NOT under the SDK module dir).

## Canonical examples (ship with boost)

- `factory/contrib/aws/aws-sdk-go-v2/v1/examples/kinesis/`
- `factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/`
- `factory/contrib/aws/aws-sdk-go-v2/v1/examples/sns/`
- `factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/`

Each service example shows the canonical shape: build aws.Config once via the factory, then construct per-service clients with `<service>.NewFromConfig(cfg)`.

## Construction

```go
import awsfact "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"

cfg, err := awsfact.NewConfig(ctx)
if err != nil { log.Fatalf("aws config: %v", err) }

snsClient := sns.NewFromConfig(cfg)
sqsClient := sqs.NewFromConfig(cfg)
```

Configure region, credentials provider, retry policy under `boost.factory.aws.*` (override `BOOST_FACTORY_AWS_*`). Standard AWS env vars (`AWS_REGION`, `AWS_PROFILE`) are honored too.

## Umbrella SDK layout (factory side only)

`factory/contrib/aws/aws-sdk-go-v2/v1/client/<service>/` is where per-service convenience clients live. This grouping is **exclusive to `factory/contrib/`** because factories ship clients that share an SDK version pin. For wrapper drivers, the layout is **different** — split per service: `wrapper/publisher/driver/contrib/aws/sns/v1/`, NOT `wrapper/publisher/driver/contrib/aws/aws-sdk-go-v2/v1/sns/`. See `boost-maintainer`.

## Red flags

| Red flag | Fix |
|---|---|
| `config.LoadDefaultConfig(ctx)` from `aws/config` directly | `awsfact.NewConfig(ctx)` |
| One aws.Config per service client | Build once, share |
| AWS credentials in code or `os.Getenv` reads | Use the credentials chain (env, IAM role, profile) — boost picks it up |
| Wrapper driver placed under `wrapper/.../aws/aws-sdk-go-v2/v1/sns/` | Should be `wrapper/.../aws/sns/v1/` |
