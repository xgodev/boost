package builder

import (
	"regexp"
	"sort"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

const (
	fxAdapterRoot = "github.com/xgodev/boost/fx/modules/bootstrap/function/adapter/contrib"

	defaultApp        = "api"
	defaultRESTServer = "rest"
	defaultGRPCServer = "grpc"
	functionFile      = "function" // the function door of an app is function_module.go
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

// functionArgs are the arguments of @Function(adapter=nats, app=api).
type functionArgs struct {
	Adapter string `attr:"adapter"`
	App     string `attr:"app"`
}

var (
	restMethods = map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true, "HEAD": true, "OPTIONS": true,
	}
	// functionAdapters maps the adapter of @Function to the import paths of
	// the fx modules that run it.
	functionAdapters = map[string][]string{
		"cloudevents": {fxAdapterRoot + "/cloudevents/sdk-go/v2/core/http"},
		"kafka":       {fxAdapterRoot + "/confluentinc/confluent-kafka-go/v2"},
		"nats":        {fxAdapterRoot + "/nats-io/nats.go/v1"},
		// the version (v1 or v2) is picked by config, so both are included
		"pubsub": {fxAdapterRoot + "/cloud.google.com/pubsub/v1", fxAdapterRoot + "/cloud.google.com/pubsub/v2"},
	}

	restKeys     = []string{"method", "path", "server", "app"}
	grpcKeys     = []string{"register", "server", "app"}
	functionKeys = []string{"adapter", "app"}

	// identName is the form of app and server names: they become a package,
	// a file and a function name and a config path element.
	identName = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
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

func validateNames(ann, where, app, server string) error {
	if !identName.MatchString(app) {
		return errors.NotValidf("@%s on %s: app name %q must match %s", ann, where, app, identName)
	}
	if server != "" && (!identName.MatchString(server) || server == functionFile) {
		return errors.NotValidf("@%s on %s: server name %q must match %s and not be %q", ann, where, server, identName, functionFile)
	}
	return nil
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
