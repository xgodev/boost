# vault

Factory for the [mittwald/vaultgo](https://github.com/mittwald/vaultgo) Vault
client, plus managers that keep other clients configured with the
credentials a Vault secret holds.

## Client

```go
client, err := vault.NewClient(ctx)
```

| Key | Default | Description |
|---|---|---|
| `boost.factory.vault.addr` | `""` | vault address |
| `boost.factory.vault.caPath` | `""` | ca path |
| `boost.factory.vault.type` | `TOKEN` | auth type, `TOKEN` or `K8S` (anything else is an error) |
| `boost.factory.vault.token` | `XPTOTOKEN` | vault token |
| `boost.factory.vault.k8s.role` | `default` | k8s role |
| `boost.factory.vault.k8s.jwt.file` | `""` | jwt file |
| `boost.factory.vault.k8s.jwt.content` | `""` | jwt content, used when file is empty |

## Managers

A `Manager` receives the data of a secret. `ManagerPool` reads each
manager's secret, renames its keys with `keys` and calls `Configure`. With
the watcher on, it renews the lease and, when the lease ends, closes the
manager and configures it again with a fresh secret.

```go
redisManager, err := redismanager.NewManager(redisOptions)
redisManager.OnChange(func(c *goredis.Client) { repo.Use(c) })
vault.ManageAll(ctx, redisManager)
```

`ClientManager[C]` is the generic manager behind the contrib ones: it builds
a client from the secret, keeps the latest (`Client()`) and passes each new
one to `OnChange`. A failed build keeps the previous client.

| Contrib manager | Client | Config root | Reads |
|---|---|---|---|
| `managers/contrib/redis/go-redis/v9` | `*redis.Client`, `*redis.ClusterClient` | `boost.factory.vault.managers.redis` | `password` |
| `managers/contrib/elastic/go-elasticsearch/v8` | `*elasticsearch.Client` | `boost.factory.vault.managers.elasticsearch` | `username`, `password` |
| `managers/contrib/go.mongodb.org/mongo-driver/v1` | `*mongo.Conn` | `boost.factory.vault.managers.mongo` | `username`, `password` |

Each manager root takes:

| Key | Default | Description |
|---|---|---|
| `<root>.secretPath` | `""` | full secret path, e.g. `v1/database/creds/orders` |
| `<root>.watcher.enabled` | `true` | renew the lease and reconfigure when it ends |
| `<root>.watcher.increment` | `120` | lease increment, in seconds |
| `<root>.keys` | `{username: username, password: password}` | secret key → manager key |
