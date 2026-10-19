# azure

Factories for the [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go):
a token credential and, under `client/azblob`, a Blob Storage client.

## Credential

```go
cred, err := azure.NewCredential(ctx)
```

With client id, client secret and tenant id all set it returns a client
secret credential; otherwise the default credential chain (environment,
workload identity, managed identity, az cli).

| Key | Default | Description |
|---|---|---|
| `boost.factory.azure.clientId` | `""` | azure client id (hidden) |
| `boost.factory.azure.clientSecret` | `""` | azure client secret (hidden) |
| `boost.factory.azure.tenantId` | `""` | azure tenant id |

## azblob

```go
blobs, err := azblob.NewClient(ctx, cred)
```

With a connection string the credential is ignored (handy for Azurite).
Otherwise the client authenticates with the credential against `url`, or
`https://<accountName>.blob.core.windows.net/` when `url` is empty.

| Key | Default | Description |
|---|---|---|
| `boost.factory.azure.azblob.accountName` | `unknown` | storage account name (hidden) |
| `boost.factory.azure.azblob.url` | `""` | service url |
| `boost.factory.azure.azblob.connectionString` | `""` | connection string (hidden) |

`Plugin = factory.Plugin[*azblob.ClientOptions]` customizes the client
options (retry, transport, telemetry) before the client is built.

## fx

`fx/modules/factory/contrib/azure/azure-sdk-for-go/v1` provides
`azcore.TokenCredential`; `.../v1/client/azblob` provides `*azblob.Client`.
