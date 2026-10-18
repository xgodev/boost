package builder

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

const testModule = "example.com/app"

// fileImports is the import block every test entry is declared under.
var fileImports = map[string]string{
	"context":     "context",
	"echo":        "github.com/labstack/echo/v4",
	"cloudevents": "github.com/cloudevents/sdk-go/v2",
	"event":       "github.com/cloudevents/sdk-go/v2/event",
	"healthpb":    "google.golang.org/grpc/health/grpc_health_v1",
	"orders":      testModule + "/orders",
}

// decl builds the entry of a function (recv "") or method declared in
// package pkg of testModule.
func decl(pkg, recv, name string, params, results []string, anns ...annotation.Annotation) annotation.Entry {
	entry := annotation.Entry{
		Path:        testModule + "/" + pkg,
		Package:     path.Base(pkg),
		Struct:      recv,
		Func:        annotation.EntryFunc{Name: name},
		Annotations: anns,
		Imports:     fileImports,
	}
	for _, p := range params {
		entry.Func.Parameters = append(entry.Func.Parameters, annotation.EntryFuncType{Type: p})
	}
	for _, r := range results {
		entry.Func.Results = append(entry.Func.Results, annotation.EntryFuncType{Type: r})
	}
	return entry
}

func ann(name, value string) annotation.Annotation { return annotation.NewAnnotation(name, value) }

var (
	echoHandler  = []string{"echo.Context"}
	errResult    = []string{"error"}
	eventParams  = []string{"context.Context", "cloudevents.Event"}
	eventResults = []string{"*cloudevents.Event", "error"}
)

// validEntries declares two REST servers, a gRPC service and a function app.
func validEntries() []annotation.Entry {
	return []annotation.Entry{
		decl("orders", "", "New", nil, []string{"*Handler"}, ann("Provide", "index=0")),
		decl("orders", "Handler", "Create", echoHandler, errResult, ann("REST", "method=post, path=/orders")),
		decl("orders", "Handler", "Get", echoHandler, errResult, ann("REST", "method=GET, path=/orders/:id")),
		decl("probe", "", "Live", echoHandler, errResult,
			ann("REST", "method=GET, path=/live, server=admin"),
			ann("REST", "method=HEAD, path=/live, server=admin")),
		decl("health", "", "NewServer", nil, []string{"*Server"},
			ann("Provide", "index=0"), ann("GRPC", "register=healthpb.RegisterHealthServer")),
		decl("events", "", "New", nil, []string{"Handler"}, ann("Provide", "index=0")),
		decl("events", "Handler", "Handle", []string{"context.Context", "event.Event"}, eventResults,
			ann("Function", "adapter=nats, app=worker")),
	}
}

func newTestSpec(t *testing.T, entries []annotation.Entry) (*Spec, error) {
	t.Helper()
	g, err := inject.NewGraphFromEntries(context.Background(), entries)
	require.NoError(t, err)
	return NewSpec(entries, g, nil)
}

func routeNames(routes []Route) (out []string) {
	for _, r := range routes {
		name := r.Handler.Entry.Func.Name
		if r.Handler.Entry.Struct != "" {
			name = r.Handler.Entry.Struct + "." + name
		}
		out = append(out, r.Method+" "+r.Path+" "+r.Handler.Entry.Package+"."+name)
	}
	return out
}

func TestNewSpec(t *testing.T) {
	spec, err := newTestSpec(t, append(validEntries(), validEntries()...)) // duplicates are collapsed
	require.NoError(t, err)

	require.Len(t, spec.Apps, 2)
	api, worker := spec.Apps[0], spec.Apps[1]
	assert.Equal(t, "api", api.Name)
	assert.Equal(t, "worker", worker.Name)

	require.Len(t, api.Servers, 3)
	admin, grpcSrv, rest := api.Servers[0], api.Servers[1], api.Servers[2]

	assert.Equal(t, Server{App: "api", Name: "admin", Kind: KindREST, Routes: admin.Routes}, admin)
	assert.Equal(t, []string{"GET /live probe.Live", "HEAD /live probe.Live"}, routeNames(admin.Routes))
	assert.Nil(t, admin.Routes[0].Handler.Receiver, "a package function needs no provider")

	assert.Equal(t, "rest", rest.Name)
	assert.Equal(t, KindREST, rest.Kind)
	assert.Equal(t, []string{"POST /orders orders.Handler.Create", "GET /orders/:id orders.Handler.Get"}, routeNames(rest.Routes),
		"routes are sorted by path, then method; the method is upper-cased")
	recv := rest.Routes[0].Handler.Receiver
	require.NotNil(t, recv)
	assert.Equal(t, "*Handler", recv.Type)
	assert.Equal(t, testModule+"/orders.New", recv.Provider.Key())

	assert.Equal(t, "grpc", grpcSrv.Name)
	assert.Equal(t, KindGRPC, grpcSrv.Kind)
	require.Len(t, grpcSrv.Services, 1)
	assert.Equal(t, "healthpb.RegisterHealthServer", grpcSrv.Services[0].Register)
	assert.Equal(t, testModule+"/health.NewServer", grpcSrv.Services[0].Provider.Key())
	assert.Equal(t, 0, grpcSrv.Services[0].Index)

	assert.Empty(t, worker.Servers)
	require.NotNil(t, worker.Function)
	assert.Equal(t, "nats", worker.Function.Adapter)
	assert.Equal(t, "Handle", worker.Function.Handler.Entry.Func.Name)
	require.NotNil(t, worker.Function.Handler.Receiver)
	assert.Equal(t, "Handler", worker.Function.Handler.Receiver.Type, "a value provider serves a value receiver")
	assert.Nil(t, api.Function)
}

func TestNewSpec_ReceiverProvidedByAnotherPackage(t *testing.T) {
	entries := []annotation.Entry{
		decl("wiring", "", "NewOrders", nil, []string{"*orders.Handler"}, ann("Provide", "index=0")),
		decl("orders", "Handler", "Create", echoHandler, errResult, ann("REST", "method=POST, path=/orders")),
	}

	spec, err := newTestSpec(t, entries)
	require.NoError(t, err)
	recv := spec.Apps[0].Servers[0].Routes[0].Handler.Receiver
	require.NotNil(t, recv)
	assert.Equal(t, testModule+"/wiring.NewOrders", recv.Provider.Key())
	assert.Equal(t, "*Handler", recv.Type)
}

func TestNewSpec_Errors(t *testing.T) {
	live := func(anns ...annotation.Annotation) annotation.Entry {
		return decl("probe", "", "Live", echoHandler, errResult, anns...)
	}
	rest := func(v string) annotation.Annotation { return ann("REST", v) }
	grpcSvc := func(v string) annotation.Entry {
		return decl("health", "", "NewServer", nil, []string{"*Server"}, ann("Provide", "index=0"), ann("GRPC", v))
	}
	function := func(v string) annotation.Entry {
		return decl("events", "", "Handle", eventParams, eventResults, ann("Function", v))
	}

	testCases := []struct {
		name    string
		entries []annotation.Entry
		kind    func(error) bool
		msg     string
	}{
		{"unknown method", []annotation.Entry{live(rest("method=FETCH, path=/live"))}, errors.IsNotValid, "method FETCH"},
		{"missing method", []annotation.Entry{live(rest("path=/live"))}, errors.IsNotValid, "method"},
		{"missing path", []annotation.Entry{live(rest("method=GET"))}, errors.IsNotValid, "path"},
		{"unknown key", []annotation.Entry{live(rest("method=GET, path=/live, sever=admin"))}, errors.IsNotValid, "sever"},
		{"value of the wrong type", []annotation.Entry{live(rest("method=1, path=/live"))}, errors.IsNotValid, "@REST"},
		{"duplicate route", []annotation.Entry{
			live(rest("method=GET, path=/live")),
			decl("probe", "", "Ready", echoHandler, errResult, rest("method=GET, path=/live")),
		}, errors.IsAlreadyExists, "GET /live"},
		{"rest params", []annotation.Entry{decl("probe", "", "Live", []string{"context.Context"}, errResult, rest("method=GET, path=/live"))}, errors.IsNotValid, "func(echo.Context) error"},
		{"rest results", []annotation.Entry{decl("probe", "", "Live", echoHandler, []string{"string", "error"}, rest("method=GET, path=/live"))}, errors.IsNotValid, "func(echo.Context) error"},
		{"echo from another module", []annotation.Entry{func() annotation.Entry {
			e := live(rest("method=GET, path=/live"))
			e.Imports = map[string]string{"echo": "example.com/fake/echo"}
			return e
		}()}, errors.IsNotValid, "func(echo.Context) error"},
		{"receiver without provider", []annotation.Entry{decl("orders", "Handler", "Create", echoHandler, errResult, rest("method=POST, path=/orders"))}, errors.IsNotFound, "Handler"},
		{"on a struct", []annotation.Entry{decl("orders", "Handler", "", nil, nil, rest("method=POST, path=/orders"))}, errors.IsNotValid, "function or method"},
		{"unexported handler", []annotation.Entry{decl("probe", "", "live", echoHandler, errResult, rest("method=GET, path=/live"))}, errors.IsNotValid, "exported"},
		{"unexported receiver", []annotation.Entry{
			decl("orders", "", "New", nil, []string{"*handler"}, ann("Provide", "index=0")),
			decl("orders", "handler", "Create", echoHandler, errResult, rest("method=POST, path=/orders")),
		}, errors.IsNotValid, "exported"},
		{"package main", []annotation.Entry{func() annotation.Entry {
			e := live(rest("method=GET, path=/live"))
			e.Package = "main"
			return e
		}()}, errors.IsNotValid, "package main"},
		{"invalid server name", []annotation.Entry{live(rest("method=GET, path=/live, server=Admin"))}, errors.IsNotValid, "server name"},
		{"reserved server name", []annotation.Entry{live(rest("method=GET, path=/live, server=function"))}, errors.IsNotValid, "server name"},
		{"invalid app name", []annotation.Entry{live(rest("method=GET, path=/live, app=my-app"))}, errors.IsNotValid, "app name"},
		{"app name is a Go keyword", []annotation.Entry{live(rest("method=GET, path=/live, app=go"))}, errors.IsNotValid, "app name"},
		{"app name main", []annotation.Entry{live(rest("method=GET, path=/live, app=main"))}, errors.IsNotValid, "app name"},
		{"server used by two kinds", []annotation.Entry{live(rest("method=GET, path=/live, server=grpc")), grpcSvc("register=healthpb.RegisterHealthServer")}, errors.IsNotValid, "grpc"},
		{"grpc without provide", []annotation.Entry{decl("health", "", "NewServer", nil, []string{"*Server"}, ann("GRPC", "register=healthpb.RegisterHealthServer"))}, errors.IsNotValid, "@Provide"},
		{"grpc with named provide", []annotation.Entry{decl("health", "", "NewServer", nil, []string{"*Server"},
			ann("Provide", "index=0, name=health"), ann("GRPC", "register=healthpb.RegisterHealthServer"))}, errors.IsNotValid, "@Provide"},
		{"grpc missing register", []annotation.Entry{grpcSvc("server=grpc")}, errors.IsNotValid, "register"},
		{"grpc register not imported", []annotation.Entry{grpcSvc("register=missingpb.RegisterHealthServer")}, errors.IsNotValid, "missingpb"},
		{"grpc register not a function name", []annotation.Entry{grpcSvc(`register="healthpb.Register()"`)}, errors.IsNotValid, "register"},
		{"grpc unexported local register", []annotation.Entry{grpcSvc("register=register")}, errors.IsNotValid, "register"},
		{"grpc duplicate service", []annotation.Entry{func() annotation.Entry {
			e := grpcSvc("register=healthpb.RegisterHealthServer")
			e.Annotations = append(e.Annotations, ann("GRPC", "register=healthpb.RegisterHealthServer"))
			return e
		}()}, errors.IsAlreadyExists, "RegisterHealthServer"},
		{"two functions in an app", []annotation.Entry{function("adapter=nats"), decl("events", "", "Other", eventParams, eventResults, ann("Function", "adapter=kafka"))}, errors.IsAlreadyExists, "@Function"},
		{"unknown adapter", []annotation.Entry{function("adapter=sqs")}, errors.IsNotValid, "adapter sqs"},
		{"missing adapter", []annotation.Entry{function("app=worker")}, errors.IsNotValid, "adapter"},
		{"function params", []annotation.Entry{decl("events", "", "Handle", []string{"context.Context", "string"}, eventResults, ann("Function", "adapter=nats"))}, errors.IsNotValid, "func(context.Context, cloudevents.Event) (T, error)"},
		{"function results", []annotation.Entry{decl("events", "", "Handle", eventParams, []string{"error"}, ann("Function", "adapter=nats"))}, errors.IsNotValid, "func(context.Context, cloudevents.Event) (T, error)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newTestSpec(t, tc.entries)
			require.Error(t, err)
			assert.True(t, tc.kind(err), "unexpected error kind: %v", err)
			assert.Contains(t, err.Error(), tc.msg)
		})
	}
}

// pluginEntries provide the []Plugin of the servers of validEntries: a
// default value for rest and grpc, a value named after admin, and values
// no server takes.
func pluginEntries() []annotation.Entry {
	imports := map[string]string{
		"boostecho":  "github.com/xgodev/boost/factory/contrib/labstack/echo/v4",
		"grpcserver": "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server",
		"other":      "example.com/other",
	}
	entries := []annotation.Entry{
		decl("plugins", "", "Admin", nil, []string{"[]boostecho.Plugin"}, ann("Provide", "name=admin, index=0")),
		decl("plugins", "", "Rest", nil, []string{"[]boostecho.Plugin"}, ann("Provide", "index=0")),
		decl("plugins", "", "Grpc", nil, []string{"[]grpcserver.Plugin"}, ann("Provide", "index=0")),
		decl("plugins", "", "Other", nil, []string{"[]other.Plugin"}, ann("Provide", "name=admin, index=0")),
		decl("plugins", "", "Single", nil, []string{"boostecho.Plugin"}, ann("Provide", "name=admin, index=0")),
		decl("plugins", "", "Grouped", nil, []string{"[]boostecho.Plugin"}, ann("Provide", "group=admin, index=0")),
	}
	for i := range entries {
		entries[i].Imports = imports
	}
	return entries
}

func pluginNames(plugins []inject.Component) (out []string) {
	for _, p := range plugins {
		out = append(out, p.Entry.Func.Name)
	}
	return out
}

func TestNewSpec_ServerPlugins(t *testing.T) {
	spec, err := newTestSpec(t, append(validEntries(), pluginEntries()...))
	require.NoError(t, err)

	api := spec.Apps[0]
	require.Len(t, api.Servers, 3)
	admin, grpcSrv, rest := api.Servers[0], api.Servers[1], api.Servers[2]

	assert.Equal(t, []string{"Admin"}, pluginNames(admin.Plugins), "a named server takes the []Plugin named after it")
	assert.Equal(t, []string{"Rest"}, pluginNames(rest.Plugins), "rest takes the default []boostecho.Plugin")
	assert.Equal(t, []string{"Grpc"}, pluginNames(grpcSrv.Plugins), "grpc takes the default []grpcserver.Plugin")
}

func TestNewSpec_ImplementationSetErrors(t *testing.T) {
	provide, deflt := ann("Provide", "index=0"), ann("Default", "")
	// orders.Handler has two implementations, memory and unlimited
	handlers := func(door annotation.Entry) []annotation.Entry {
		return []annotation.Entry{
			decl("orders/memory", "", "New", nil, []string{"*orders.Handler"}, provide, deflt),
			decl("orders/unlimited", "", "New", nil, []string{"*orders.Handler"}, provide),
			door,
		}
	}

	testCases := []struct {
		name    string
		entries []annotation.Entry
		msg     string
	}{
		{"rest receiver", handlers(decl("orders", "Handler", "Create", echoHandler, errResult, ann("REST", "method=POST, path=/orders"))), "@REST"},
		{"function receiver", handlers(decl("orders", "Handler", "Handle", eventParams, eventResults, ann("Function", "adapter=nats"))), "@Function"},
		{"grpc constructor", []annotation.Entry{
			decl("health/memory", "", "NewServer", nil, []string{"healthpb.HealthServer"}, provide, deflt, ann("GRPC", "register=healthpb.RegisterHealthServer")),
			decl("health/unlimited", "", "NewServer", nil, []string{"healthpb.HealthServer"}, provide),
		}, "@GRPC"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newTestSpec(t, tc.entries)
			require.Error(t, err)
			assert.True(t, errors.IsNotValid(err), "unexpected error kind: %v", err)
			assert.Contains(t, err.Error(), tc.msg)
			assert.Contains(t, err.Error(), "memory, unlimited", "the message names the implementations")
		})
	}
}
