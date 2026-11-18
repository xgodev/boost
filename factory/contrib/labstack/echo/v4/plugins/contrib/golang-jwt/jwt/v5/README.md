# jwt — Bearer token authentication for echo

Rejects every request without a valid `Authorization: Bearer <token>` before
any handler runs. For a service outside a service mesh (Cloud Run, a VM):
under Istio the check is `RequestAuthentication` + `AuthorizationPolicy`, not
this plugin.

The token is a JWT ([`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt))
signed with a key of a JWKS ([`MicahParks/keyfunc/v3`](https://github.com/MicahParks/keyfunc)
fetches and caches it). Only asymmetric algorithms are accepted, `RS256` and
`ES256`: `alg: none` and `HS*` are refused, so the service holds no secret
able to mint a token.

## Config

Root `boost.factory.echo.plugins.jwt`:

| Key | Default | |
|---|---|---|
| `.enabled` | `false` | turns the plugin on |
| `.jwksUrl` | | where the public keys come from; required when enabled |
| `.issuer` | | `iss` every token must carry; required when enabled |
| `.audience` | | `aud` every token must carry; required when enabled |
| `.skipPaths` | `["/health", "/ready", "/live"]` | request paths served without a token (exact match on the URL path) |
| `.leeway` | `30s` | clock skew tolerated on `exp`, `nbf` and `iat` |

Env vars: `BOOST_FACTORY_ECHO_PLUGINS_JWT_ENABLED=true`,
`BOOST_FACTORY_ECHO_PLUGINS_JWT_JWKS__URL=https://login.example/.well-known/jwks.json`,
`BOOST_FACTORY_ECHO_PLUGINS_JWT_ISSUER=...`, `BOOST_FACTORY_ECHO_PLUGINS_JWT_AUDIENCE=...`.

Enabled without `jwksUrl`, `issuer` or `audience`, `Register` fails and the
server does not boot.

## Behavior

- No header, a scheme other than `Bearer`, a bad signature, an unknown
  `kid`, an expired or not-yet-valid token, a token without `exp`, a wrong
  `iss` or `aud`, an algorithm other than RS256/ES256 → `401` as
  `errors.Unauthorizedf` from `model/errors`; the `errorHandler` plugin
  renders the usual JSON envelope (`{"httpStatusCode":401,"errorCode":"401","message":"invalid token"}`).
  The reason is logged at debug level, never sent to the caller.
- A valid token → the claims go into the request context.
- The JWKS is read once at boot and refreshed every hour. An unknown `kid`
  refreshes it at once (key rotation), at most once every 5 minutes; inside
  that window an unknown `kid` is refused without waiting.
- A JWKS not reachable at boot (a service that publishes its own JWKS and
  is not listening yet) is not a boot error: it is logged and read again on
  the first token.

## Reading the caller in a handler

```go
import "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/contrib/golang-jwt/jwt/v5"

func (h *Handler) Handle(c echo.Context) error {
	claims, ok := jwt.ClaimsFromContext(c.Request().Context())
	if !ok {
		return errors.Unauthorizedf("no caller") // a skipped path, or the plugin off
	}
	account := claims.Subject() // the sub claim; other claims: claims["machine_id"]
	...
}
```

A handler test fakes a caller with `jwt.ContextWithClaims(ctx, jwt.Claims{"sub": "account-42"})`.

## Wiring

Builder app — list it in the `@Server` of the rest server:

```go
// @Server(type=rest, plugins=[recover, requestId, logger, errorHandler, health, jwt])
package rest
```

By hand — pass `jwt.Register` to the echo factory, after `errorHandler`:

```go
srv, err := echo.NewServer(ctx, errorhandler.Register, jwt.Register)
```

`NewJWT()`, `NewJWTWithConfigPath(path)` and `NewJWTWithOptions(o)` build
the plugin from the default root, another config path, or options in code.
