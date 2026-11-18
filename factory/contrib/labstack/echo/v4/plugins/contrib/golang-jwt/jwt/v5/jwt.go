package jwt

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	jwtgo "github.com/golang-jwt/jwt/v5"
	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// algorithms are the only signing methods accepted: asymmetric, so the
// service holds no secret able to mint a token.
var algorithms = []string{"RS256", "ES256"}

// jwksTimeout bounds every request to the JWKS.
const jwksTimeout = 10 * time.Second

// Register registers a new jwt plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewJWT()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// JWT represents the jwt authentication plugin for echo server.
type JWT struct {
	options *Options
}

// NewJWTWithOptions returns a new jwt plugin with options.
func NewJWTWithOptions(options *Options) *JWT {
	return &JWT{options: options}
}

// NewJWTWithConfigPath returns a new jwt plugin with options from config path.
func NewJWTWithConfigPath(path string) (*JWT, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewJWTWithOptions(o), nil
}

// NewJWT returns a new jwt plugin with default options.
func NewJWT() (*JWT, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewJWTWithOptions(o), nil
}

// Register registers this jwt plugin for echo server: every request outside
// the skip paths needs a valid Bearer token. The JWKS is fetched in the
// background until ctx ends; a JWKS not reachable at boot is not an error,
// it is fetched again on the first unknown kid.
func (i *JWT) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	if i.options.JwksUrl == "" || i.options.Issuer == "" || i.options.Audience == "" {
		return errors.NotValidf("jwt plugin needs %s, %s and %s", root+jwksURL, root+issuer, root+audience)
	}

	logger.Trace("enabling jwt middleware in echo")

	kf, err := keyfunc.NewDefaultOverrideCtx(ctx, []string{i.options.JwksUrl}, keyfunc.Override{
		Client:      &http.Client{Timeout: jwksTimeout},
		HTTPTimeout: jwksTimeout,
		// An unknown kid never holds a request waiting for the refresh rate
		// limiter: it is refused while the limiter is closed.
		RateLimitWaitMax: time.Millisecond,
		RefreshErrorHandlerFunc: func(u string) func(context.Context, error) {
			return func(ctx context.Context, err error) {
				log.FromContext(ctx).Errorf("failed to refresh the JWKS %s: %v", u, err)
			}
		},
	})
	if err != nil {
		return errors.Wrap(err, errors.Internalf("jwt plugin could not read the JWKS %s", i.options.JwksUrl))
	}

	server.Use(i.middleware(kf.KeyfuncCtx))

	logger.Debug("jwt middleware successfully enabled in echo")

	return nil
}

func (i *JWT) middleware(keys func(context.Context) jwtgo.Keyfunc) e.MiddlewareFunc {
	skip := make(map[string]struct{}, len(i.options.SkipPaths))
	for _, p := range i.options.SkipPaths {
		skip[p] = struct{}{}
	}
	parser := jwtgo.NewParser(
		jwtgo.WithValidMethods(algorithms),
		jwtgo.WithIssuer(i.options.Issuer),
		jwtgo.WithAudience(i.options.Audience),
		jwtgo.WithLeeway(i.options.Leeway),
		jwtgo.WithExpirationRequired(),
	)

	return func(next e.HandlerFunc) e.HandlerFunc {
		return func(c e.Context) error {
			req := c.Request()
			if _, ok := skip[req.URL.Path]; ok {
				return next(c)
			}

			scheme, raw, found := strings.Cut(req.Header.Get(e.HeaderAuthorization), " ")
			if !found || !strings.EqualFold(scheme, "Bearer") || raw == "" {
				return errors.Unauthorizedf("missing bearer token")
			}

			claims := jwtgo.MapClaims{}
			if _, err := parser.ParseWithClaims(raw, claims, keys(req.Context())); err != nil {
				log.FromContext(req.Context()).Debugf("jwt rejected: %v", err)
				return errors.Unauthorizedf("invalid token")
			}

			c.SetRequest(req.WithContext(ContextWithClaims(req.Context(), Claims(claims))))
			return next(c)
		}
	}
}
