package builder

import (
	"go/token"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

const (
	fxAdapterRoot    = "github.com/xgodev/boost/fx/modules/bootstrap/function/adapter/contrib"
	fxMiddlewareRoot = "github.com/xgodev/boost/fx/modules/bootstrap/function/middleware"
	middlewareRoot   = "github.com/xgodev/boost/bootstrap/function/middleware"

	defaultApp        = "api"
	defaultRESTServer = "rest"
	defaultGRPCServer = "grpc"
	defaultMCPServer  = "mcp"
	functionFile      = "function"   // the function door of an app is function_module.go
	middlewareDir     = "middleware" // gen/builder/middleware holds the decorators, so no app takes the name
)

// restArgs are the arguments of @REST(method=GET, path=/x, server=rest, app=api).
type restArgs struct {
	Method string `attr:"method"`
	Path   string `attr:"path"`
	Server string `attr:"server"`
	App    string `attr:"app"`
}

// grpcArgs are the arguments of @GRPC(register=pb.RegisterXServer, server=grpc, app=api).
type grpcArgs struct {
	Register string `attr:"register"`
	Server   string `attr:"server"`
	App      string `attr:"app"`
}

// mcpArgs are the arguments of @MCP(name=check_stock, server=mcp, app=api).
type mcpArgs struct {
	Name   string `attr:"name"`
	Server string `attr:"server"`
	App    string `attr:"app"`
}

// functionArgs are the arguments of @Function(adapter=nats, app=api, middlewares=[recovery, logger]).
type functionArgs struct {
	Adapter     string `attr:"adapter"`
	App         string `attr:"app"`
	Middlewares string `attr:"middlewares"`
}

// functionMiddleware is a middleware of @Function: the fx module that
// provides it and the type it provides, *Type[T] of package Pkg.
type functionMiddleware struct {
	Module inject.ImportData
	Pkg    inject.ImportData
	Type   string
}

var (
	restMethods = map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true, "HEAD": true, "OPTIONS": true,
	}
	// functionAdapters maps the adapter of @Function to the fx modules that
	// run it, imported under fixed aliases by the generated function door.
	functionAdapters = map[string][]inject.ImportData{
		"cenats":      {{Alias: "cenatsfx", Path: fxAdapterRoot + "/cloudevents/sdk-go/v2/contrib/nats-io/nats.go/v1"}},
		"cloudevents": {{Alias: "cloudeventsfx", Path: fxAdapterRoot + "/cloudevents/sdk-go/v2/core/http"}},
		"kafka":       {{Alias: "kafkafx", Path: fxAdapterRoot + "/confluentinc/confluent-kafka-go/v2"}},
		"lambda":      {{Alias: "lambdafx", Path: fxAdapterRoot + "/aws/aws-lambda-go/v1"}},
		"nats":        {{Alias: "natsfx", Path: fxAdapterRoot + "/nats-io/nats.go/v1"}},
		// each version is enabled by its own config, so both are included
		"pubsub": {
			{Alias: "pubsubv1fx", Path: fxAdapterRoot + "/cloud.google.com/pubsub/v1"},
			{Alias: "pubsubv2fx", Path: fxAdapterRoot + "/cloud.google.com/pubsub/v2"},
		},
	}

	// functionMiddlewares maps the names in @Function(middlewares=[...]) to
	// the bootstrap/function middlewares.
	functionMiddlewares = map[string]functionMiddleware{
		"ignore_errors": {
			Module: inject.ImportData{Alias: "ignoreerrorsfx", Path: fxMiddlewareRoot + "/ignore_errors"},
			Pkg:    inject.ImportData{Alias: "ignoreerrorsmw", Path: middlewareRoot + "/ignore_errors"},
			Type:   "IgnoreErrors",
		},
		"logger": {
			Module: inject.ImportData{Alias: "loggerfx", Path: fxMiddlewareRoot + "/logger"},
			Pkg:    inject.ImportData{Alias: "loggermw", Path: middlewareRoot + "/logger"},
			Type:   "Logger",
		},
		"opentelemetry": {
			Module: inject.ImportData{Alias: "opentelemetryfx", Path: fxMiddlewareRoot + "/opentelemetry"},
			Pkg:    inject.ImportData{Alias: "opentelemetrymw", Path: middlewareRoot + "/go.opentelemetry.io/otel/v1"},
			Type:   "OpenTelemetry",
		},
		"prometheus": {
			Module: inject.ImportData{Alias: "prometheusfx", Path: fxMiddlewareRoot + "/prometheus"},
			Pkg:    inject.ImportData{Alias: "prometheusmw", Path: middlewareRoot + "/prometheus"},
			Type:   "Prometheus",
		},
		"publisher": {
			Module: inject.ImportData{Alias: "publisherfx", Path: fxMiddlewareRoot + "/publisher"},
			Pkg:    inject.ImportData{Alias: "publishermw", Path: middlewareRoot + "/publisher"},
			Type:   "Publisher",
		},
		"recovery": {
			Module: inject.ImportData{Alias: "recoveryfx", Path: fxMiddlewareRoot + "/recovery"},
			Pkg:    inject.ImportData{Alias: "recoverymw", Path: middlewareRoot + "/recovery"},
			Type:   "Recovery",
		},
	}

	restKeys     = []string{"method", "path", "server", "app"}
	grpcKeys     = []string{"register", "server", "app"}
	mcpKeys      = []string{"name", "server", "app"}
	functionKeys = []string{"adapter", "app", "middlewares"}

	// identName is the form of app and server names: they become a package,
	// a file and a function name and a config path element.
	identName = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
	// toolName is the form of an MCP tool name.
	toolName = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`)
)

// decode checks that the annotation only uses the given keys and copies its
// values into args.
func decode(ann annotation.Annotation, where string, keys []string, args any) error {
	allowed := map[string]bool{}
	for _, k := range keys {
		allowed[k] = true
	}
	for k := range ann.Map {
		if !allowed[k] {
			return errors.NotValidf("@%s on %s has unknown key %q; valid keys: %s", ann.Name, where, k, strings.Join(keys, ", "))
		}
	}
	if err := ann.Decode(args); err != nil {
		return errors.NotValidf("@%s on %s: %v", ann.Name, where, err)
	}
	return nil
}

// parseList reads a list value written as [a, b]: names not empty, none twice.
func parseList(raw, where, key string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	inner, open := strings.CutPrefix(raw, "[")
	inner, closed := strings.CutSuffix(inner, "]")
	if !open || !closed {
		return nil, errors.NotValidf("%s: %s=%s: write the list as [a, b]", where, key, raw)
	}
	if strings.TrimSpace(inner) == "" {
		return nil, nil
	}

	var names []string
	seen := map[string]bool{}
	for _, name := range strings.Split(inner, ",") {
		name = strings.TrimSpace(name)
		switch {
		case name == "":
			return nil, errors.NotValidf("%s: %s=%s has an empty name", where, key, raw)
		case seen[name]:
			return nil, errors.NotValidf("%s: %s=%s lists %s twice", where, key, raw, name)
		}
		seen[name] = true
		names = append(names, name)
	}
	return names, nil
}

func validateNames(ann, where, app, server string) error {
	// the app name is the package name of gen/builder/<app>
	if !identName.MatchString(app) || token.IsKeyword(app) || app == "main" || app == middlewareDir {
		return errors.NotValidf("@%s on %s: app name %q must match %s and not be a Go keyword, main or %s", ann, where, app, identName, middlewareDir)
	}
	if server != "" && (!identName.MatchString(server) || server == functionFile) {
		return errors.NotValidf("@%s on %s: server name %q must match %s and not be %q", ann, where, server, identName, functionFile)
	}
	return nil
}

// snakeCase turns a Go name into a tool name: CheckStock is check_stock and
// ListHTTPOrders is list_http_orders.
func snakeCase(name string) string {
	runes := []rune(name)
	var b strings.Builder
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && nextLower) {
				b.WriteByte('_')
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// toolDescription is the doc comment of the handler: its title, capitalized
// once the leading name is gone, and the lines after it.
func toolDescription(header annotation.EntryHeader) string {
	desc := header.Title
	if desc != "" {
		r, size := utf8.DecodeRuneInString(desc)
		desc = string(unicode.ToUpper(r)) + desc[size:]
	}
	if header.Description != "" {
		if desc != "" {
			desc += "\n"
		}
		desc += header.Description
	}
	return desc
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func sortedNames[V any](m map[string]V) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
