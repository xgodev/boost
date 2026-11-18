package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
	e "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	testIssuer   = "https://issuer.test"
	testAudience = "api.test"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// key is a signing key published in the test JWKS under kid.
type key struct {
	kid    string
	method jwtgo.SigningMethod
	signer any
	jwk    map[string]string
}

func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

func rsaKey(t *testing.T, kid string) key {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key{kid: kid, method: jwtgo.SigningMethodRS256, signer: k, jwk: map[string]string{
		"kty": "RSA", "kid": kid, "alg": "RS256", "use": "sig",
		"n": b64(k.N.Bytes()), "e": b64(big.NewInt(int64(k.E)).Bytes()),
	}}
}

func ecKey(t *testing.T, kid string) key {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	return key{kid: kid, method: jwtgo.SigningMethodES256, signer: k, jwk: map[string]string{
		"kty": "EC", "kid": kid, "alg": "ES256", "use": "sig", "crv": "P-256",
		"x": b64(k.X.FillBytes(make([]byte, 32))), "y": b64(k.Y.FillBytes(make([]byte, 32))),
	}}
}

// jwks is a JWKS served by httptest: the keys it publishes can change and
// it counts the times it is fetched.
type jwks struct {
	*httptest.Server
	mu      sync.Mutex
	keys    []key
	status  int
	fetches atomic.Int32
}

func newJWKS(t *testing.T, keys ...key) *jwks {
	t.Helper()
	s := &jwks{keys: keys, status: http.StatusOK}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		s.fetches.Add(1)
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.status != http.StatusOK {
			w.WriteHeader(s.status)
			return
		}
		set := struct {
			Keys []map[string]string `json:"keys"`
		}{}
		for _, k := range s.keys {
			set.Keys = append(set.Keys, k.jwk)
		}
		assert.NoError(t, json.NewEncoder(w).Encode(set))
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *jwks) publish(status int, keys ...key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status, s.keys = status, keys
}

func options(url string) *Options {
	return &Options{
		Enabled:   true,
		JwksUrl:   url,
		Issuer:    testIssuer,
		Audience:  testAudience,
		SkipPaths: []string{"/health"},
		Leeway:    30 * time.Second,
	}
}

// serve builds an echo server with plugin, the JSON error handler and two
// routes: /me answers the subject the plugin left in the context, /health
// whether the request carried claims.
func serve(t *testing.T, plugin echo.Plugin) *echo.Server {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	srv := echo.NewServerWithOptions(ctx, &echo.Options{}, plugin)
	srv.Instance().HTTPErrorHandler = echo.ErrorHandlerJSON
	srv.GET("/me", func(c e.Context) error {
		claims, ok := ClaimsFromContext(c.Request().Context())
		if !ok {
			return c.String(http.StatusOK, "anonymous")
		}
		return c.String(http.StatusOK, claims.Subject())
	})
	srv.GET("/health", func(c e.Context) error {
		_, ok := ClaimsFromContext(c.Request().Context())
		return c.String(http.StatusOK, map[bool]string{true: "claims", false: "no claims"}[ok])
	})
	return srv
}

func call(srv *echo.Server, path, authorization string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if authorization != "" {
		req.Header.Set(e.HeaderAuthorization, authorization)
	}
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, req)
	return rec
}

func sign(t *testing.T, k key, claims jwtgo.MapClaims) string {
	t.Helper()
	tok := jwtgo.NewWithClaims(k.method, claims)
	tok.Header["kid"] = k.kid
	s, err := tok.SignedString(k.signer)
	require.NoError(t, err)
	return s
}

func valid(mod func(jwtgo.MapClaims)) jwtgo.MapClaims {
	now := time.Now()
	c := jwtgo.MapClaims{
		"sub": "account-42",
		"iss": testIssuer,
		"aud": testAudience,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	}
	if mod != nil {
		mod(c)
	}
	return c
}

func assertUnauthorized(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
	var body response.Error
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, http.StatusUnauthorized, body.HttpStatusCode)
	assert.Equal(t, "401", body.ErrorCode)
	assert.NotEmpty(t, body.Message)
}

func TestMiddleware(t *testing.T) {
	rs, es := rsaKey(t, "rs-1"), ecKey(t, "es-1")
	ghost := rsaKey(t, "ghost")
	impostor := ghost
	impostor.kid = rs.kid
	set := newJWKS(t, rs, es)
	srv := serve(t, NewJWTWithOptions(options(set.URL)).Register)

	hs := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, valid(nil))
	hs.Header["kid"] = rs.kid
	hs256, err := hs.SignedString([]byte("shared secret"))
	require.NoError(t, err)

	none := jwtgo.NewWithClaims(jwtgo.SigningMethodNone, valid(nil))
	none.Header["kid"] = rs.kid
	algNone, err := none.SignedString(jwtgo.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	cases := []struct {
		name          string
		path          string
		authorization string
		want          string // body of the 200
		unauthorized  bool
	}{
		{name: "valid RS256", path: "/me", authorization: "Bearer " + sign(t, rs, valid(nil)), want: "account-42"},
		{name: "valid ES256", path: "/me", authorization: "Bearer " + sign(t, es, valid(nil)), want: "account-42"},
		{name: "scheme is case insensitive", path: "/me", authorization: "bearer " + sign(t, rs, valid(nil)), want: "account-42"},
		{name: "expired inside the leeway", path: "/me", authorization: "Bearer " + sign(t, rs, valid(func(c jwtgo.MapClaims) {
			c["exp"] = time.Now().Add(-10 * time.Second).Unix()
		})), want: "account-42"},
		{name: "non-string subject reads as empty", path: "/me", authorization: "Bearer " + sign(t, rs, valid(func(c jwtgo.MapClaims) {
			c["sub"] = 42
		})), want: ""},
		{name: "skip path needs no token", path: "/health", want: "no claims"},
		{name: "missing header", path: "/me", unauthorized: true},
		{name: "wrong scheme", path: "/me", authorization: "Basic dXNlcjpwYXNz", unauthorized: true},
		{name: "bearer without token", path: "/me", authorization: "Bearer ", unauthorized: true},
		{name: "malformed token", path: "/me", authorization: "Bearer not.a.jwt", unauthorized: true},
		{name: "expired", path: "/me", authorization: "Bearer " + sign(t, rs, valid(func(c jwtgo.MapClaims) {
			c["exp"] = time.Now().Add(-time.Hour).Unix()
		})), unauthorized: true},
		{name: "without exp", path: "/me", authorization: "Bearer " + sign(t, rs, valid(func(c jwtgo.MapClaims) {
			delete(c, "exp")
		})), unauthorized: true},
		{name: "not yet valid", path: "/me", authorization: "Bearer " + sign(t, rs, valid(func(c jwtgo.MapClaims) {
			c["nbf"] = time.Now().Add(time.Hour).Unix()
		})), unauthorized: true},
		{name: "wrong issuer", path: "/me", authorization: "Bearer " + sign(t, rs, valid(func(c jwtgo.MapClaims) {
			c["iss"] = "https://someone.else"
		})), unauthorized: true},
		{name: "wrong audience", path: "/me", authorization: "Bearer " + sign(t, rs, valid(func(c jwtgo.MapClaims) {
			c["aud"] = "other.api"
		})), unauthorized: true},
		{name: "unknown kid", path: "/me", authorization: "Bearer " + sign(t, ghost, valid(nil)), unauthorized: true},
		{name: "signed by a key not the kid's", path: "/me", authorization: "Bearer " + sign(t, impostor, valid(nil)), unauthorized: true},
		{name: "alg none", path: "/me", authorization: "Bearer " + algNone, unauthorized: true},
		{name: "HS256", path: "/me", authorization: "Bearer " + hs256, unauthorized: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := call(srv, tc.path, tc.authorization)
			if tc.unauthorized {
				assertUnauthorized(t, rec)
				return
			}
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tc.want, rec.Body.String())
		})
	}
}

func TestUnknownKidRefreshesTheJWKSOnceThenValidates(t *testing.T) {
	old, rotated, ghost := rsaKey(t, "old"), ecKey(t, "rotated"), rsaKey(t, "ghost")
	set := newJWKS(t, old)
	srv := serve(t, NewJWTWithOptions(options(set.URL)).Register)
	require.EqualValues(t, 1, set.fetches.Load(), "the JWKS is read at boot")

	set.publish(http.StatusOK, old, rotated)

	rec := call(srv, "/me", "Bearer "+sign(t, rotated, valid(nil)))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "account-42", rec.Body.String())
	assert.EqualValues(t, 2, set.fetches.Load(), "the unknown kid refreshed the JWKS once")

	rec = call(srv, "/me", "Bearer "+sign(t, rotated, valid(nil)))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.EqualValues(t, 2, set.fetches.Load(), "a known kid is read from the cache")

	assertUnauthorized(t, call(srv, "/me", "Bearer "+sign(t, ghost, valid(nil))))
	assert.EqualValues(t, 2, set.fetches.Load(), "a second unknown kid inside the rate limit fetches nothing")
}

func TestJWKSDownAtBootIsFetchedOnTheFirstToken(t *testing.T) {
	k := rsaKey(t, "k")
	set := newJWKS(t)
	set.publish(http.StatusServiceUnavailable)
	srv := serve(t, NewJWTWithOptions(options(set.URL)).Register)
	require.EqualValues(t, 1, set.fetches.Load())

	set.publish(http.StatusOK, k)

	rec := call(srv, "/me", "Bearer "+sign(t, k, valid(nil)))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "account-42", rec.Body.String())
}

func TestDisabledWiresNothing(t *testing.T) {
	set := newJWKS(t)
	o := options(set.URL)
	o.Enabled = false
	srv := serve(t, NewJWTWithOptions(o).Register)

	rec := call(srv, "/me", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "anonymous", rec.Body.String())
	assert.Zero(t, set.fetches.Load(), "off never reads the JWKS")
}

func TestRegisterRefusesAnIncompleteConfig(t *testing.T) {
	for name, mod := range map[string]func(*Options){
		"no jwksUrl":  func(o *Options) { o.JwksUrl = "" },
		"no issuer":   func(o *Options) { o.Issuer = "" },
		"no audience": func(o *Options) { o.Audience = "" },
		"bad jwksUrl": func(o *Options) { o.JwksUrl = "not a url" },
	} {
		t.Run(name, func(t *testing.T) {
			o := options("http://jwks.test/keys")
			mod(o)
			err := NewJWTWithOptions(o).Register(context.Background(), &echo.Server{})
			assert.Error(t, err)
		})
	}
}

func TestRegisterReadsTheConfig(t *testing.T) {
	k := rsaKey(t, "k")
	set := newJWKS(t, k)
	t.Cleanup(config.Load)
	t.Setenv("BOOST_FACTORY_ECHO_PLUGINS_JWT_ENABLED", "true")
	t.Setenv("BOOST_FACTORY_ECHO_PLUGINS_JWT_JWKS__URL", set.URL)
	t.Setenv("BOOST_FACTORY_ECHO_PLUGINS_JWT_ISSUER", testIssuer)
	t.Setenv("BOOST_FACTORY_ECHO_PLUGINS_JWT_AUDIENCE", testAudience)
	config.Load()

	srv := serve(t, Register)

	assertUnauthorized(t, call(srv, "/me", ""))
	assert.Equal(t, "no claims", call(srv, "/health", "").Body.String(), "/health is skipped by default")
	assert.Equal(t, "account-42", call(srv, "/me", "Bearer "+sign(t, k, valid(nil))).Body.String())
}

func TestDefaults(t *testing.T) {
	// Read from the registered entries: the loaded config keeps what other
	// tests set through the environment.
	got := map[string]any{}
	for _, c := range config.Entries() {
		if strings.HasPrefix(c.Key, root+".") {
			got[strings.TrimPrefix(c.Key, root)] = c.Value
		}
	}
	assert.Equal(t, map[string]any{
		enabled:   false,
		jwksURL:   "",
		issuer:    "",
		audience:  "",
		skipPaths: []string{"/health", "/ready", "/live"},
		leeway:    30 * time.Second,
	}, got)
}

func TestNewJWTWithConfigPath(t *testing.T) {
	ConfigAdd("app.auth")
	t.Cleanup(config.Load)
	t.Setenv("APP_AUTH_ISSUER", "https://app.issuer")
	config.Load()

	j, err := NewJWTWithConfigPath("app.auth")
	require.NoError(t, err)
	assert.Equal(t, "https://app.issuer", j.options.Issuer)
}

func TestConstructorsReturnTheConfigError(t *testing.T) {
	ConfigAdd("app.broken")
	t.Cleanup(config.Load)
	t.Setenv("BOOST_FACTORY_ECHO_PLUGINS_JWT_LEEWAY", "soon")
	t.Setenv("APP_BROKEN_LEEWAY", "soon")
	config.Load()

	_, err := NewJWT()
	assert.Error(t, err)
	_, err = NewJWTWithConfigPath("app.broken")
	assert.Error(t, err)
	assert.Error(t, Register(context.Background(), &echo.Server{}))
}

func TestContextWithClaims(t *testing.T) {
	_, ok := ClaimsFromContext(context.Background())
	assert.False(t, ok)

	got, ok := ClaimsFromContext(ContextWithClaims(context.Background(), Claims{"sub": "a"}))
	assert.True(t, ok)
	assert.Equal(t, "a", got.Subject())
}
