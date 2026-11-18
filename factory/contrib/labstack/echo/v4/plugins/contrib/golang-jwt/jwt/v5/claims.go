package jwt

import "context"

// Claims are the claims of the token that authenticated the request.
type Claims map[string]any

// Subject returns the sub claim, "" when the token has none.
func (c Claims) Subject() string {
	s, _ := c["sub"].(string) // intentional: a missing or non-string sub reads as no subject
	return s
}

type claimsKey struct{}

// ContextWithClaims returns a copy of ctx carrying claims, as the plugin
// leaves them for the handler; a handler test uses it to fake a caller.
func ContextWithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// ClaimsFromContext returns the claims the plugin validated for the request
// whose context is ctx, and false when the request carried no token (a
// skipped path, or the plugin off).
func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsKey{}).(Claims)
	return c, ok
}
