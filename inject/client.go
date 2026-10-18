package inject

import (
	"go/ast"
	"go/parser"
	"sort"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

// ClientFactory is a client type inject builds through its boost factory
// when an @Inject parameter asks for it.
type ClientFactory struct {
	Path    string                  // import path of the client type's package
	Type    string                  // name of the client type, always used as a pointer
	Build   string                  // factory import path + "." + the ...WithConfigPath constructor
	Root    string                  // config root of the factory
	Plugins map[string]ClientPlugin // by name: the last segment of the plugin's config root
}

// ClientPlugin is a plugin of a client factory. The value of its constructor
// registers itself through its Register method.
type ClientPlugin struct {
	Path string // import path of the plugin
	New  string // ...WithConfigPath constructor of the plugin
}

const boostFactories = "github.com/xgodev/boost/factory/contrib/"

// clientFactories are the clients inject builds, by Path + "." + Type.
var clientFactories = map[string]ClientFactory{}

func init() {
	resty := boostFactories + "go-resty/resty/v2"
	addClientFactory("github.com/go-resty/resty/v2", "Client", resty, "NewClientWithConfigPath", "boost.factory.resty", map[string]ClientPlugin{
		"log":        {Path: resty + "/plugins/local/wrapper/log", New: "NewLogWithConfigPath"},
		"health":     {Path: resty + "/plugins/local/extra/health", New: "NewHealthWithConfigPath"},
		"otel":       {Path: resty + "/plugins/contrib/dubonzi/otelresty/v1", New: "NewOtelrestyWithConfigPath"},
		"datadog":    {Path: resty + "/plugins/contrib/datadog/dd-trace-go/v1", New: "NewDatadogWithConfigPath"},
		"prometheus": {Path: resty + "/plugins/contrib/prometheus/client_golang/v1", New: "NewPrometheusWithConfigPath"},
		"retry":      {Path: resty + "/plugins/extra/retry", New: "NewRetryWithConfigPath"},
		"requestid":  {Path: resty + "/plugins/extra/requestid", New: "NewRequestIDWithConfigPath"},
	})

	redis := boostFactories + "redis/go-redis/v9"
	redisPlugins := func(health, others string) map[string]ClientPlugin {
		return map[string]ClientPlugin{
			"health":     {Path: redis + "/plugins/local/extra/health", New: health},
			"otel":       {Path: redis + "/plugins/native/otel/v9", New: others},
			"datadog":    {Path: redis + "/plugins/contrib/datadog/dd-trace-go/v1", New: others},
			"prometheus": {Path: redis + "/plugins/contrib/prometheus/client_golang/v1", New: others},
		}
	}
	addClientFactory("github.com/redis/go-redis/v9", "Client", redis, "NewClientWithConfigPath", "boost.factory.redis",
		redisPlugins("NewClientHealthWithConfigPath", "NewClientWithConfigPath"))
	addClientFactory("github.com/redis/go-redis/v9", "ClusterClient", redis, "NewClusterClientWithConfigPath", "boost.factory.redis",
		redisPlugins("NewClusterHealthWithConfigPath", "NewClusterWithConfigPath"))

	for _, version := range []string{"v1", "v2"} {
		mongo := boostFactories + "go.mongodb.org/mongo-driver/" + version
		addClientFactory(mongo, "Conn", mongo, "NewConnWithConfigPath", "boost.factory.mongo", map[string]ClientPlugin{
			"health":     {Path: mongo + "/plugins/local/health", New: "NewHealthWithConfigPath"},
			"otel":       {Path: mongo + "/plugins/contrib/go.opentelemetry.io/contrib/v0", New: "NewOtelMongoWithConfigPath"},
			"datadog":    {Path: mongo + "/plugins/contrib/datadog/dd-trace-go/v1", New: "NewDatadogWithConfigPath"},
			"prometheus": {Path: mongo + "/plugins/contrib/prometheus/client_golang/v1", New: "NewPrometheusWithConfigPath"},
		})
	}
}

func addClientFactory(path, tp, factory, build, root string, plugins map[string]ClientPlugin) {
	clientFactories[path+"."+tp] = ClientFactory{Path: path, Type: tp, Build: factory + "." + build, Root: root, Plugins: plugins}
}

// PluginNames returns the names of the factory's plugins, sorted.
func (f ClientFactory) PluginNames() []string {
	names := make([]string, 0, len(f.Plugins))
	for name := range f.Plugins {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// FactoryPath returns the import path of the factory.
func (f ClientFactory) FactoryPath() string {
	return f.Build[:strings.LastIndex(f.Build, ".")]
}

// Constructor returns the name of the factory's ...WithConfigPath constructor.
func (f ClientFactory) Constructor() string {
	return f.Build[strings.LastIndex(f.Build, ".")+1:]
}

// clientFactory returns the factory of the client type tp (*pkg.Type),
// written in the file of entry, when inject builds it.
func clientFactory(entry annotation.Entry, tp string) (ClientFactory, bool) {
	expr, err := parser.ParseExpr(tp)
	if err != nil {
		return ClientFactory{}, false
	}
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return ClientFactory{}, false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return ClientFactory{}, false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return ClientFactory{}, false
	}
	f, ok := clientFactories[entry.Imports[pkg.Name]+"."+sel.Sel.Name]
	return f, ok
}

// parsePlugins reads the plugins attribute of @Inject: "[a, b]" or a single
// name. It returns nil when raw is empty and a non-nil empty list for "[]".
func parsePlugins(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	open, closed := strings.HasPrefix(raw, "["), strings.HasSuffix(raw, "]")
	if open != closed {
		return nil, errors.NotValidf("plugins=%s: write the list as [a, b]", raw)
	}
	if open {
		raw = strings.TrimSpace(raw[1 : len(raw)-1])
	}

	names := []string{}
	if raw == "" {
		return names, nil
	}
	seen := map[string]bool{}
	for _, name := range strings.Split(raw, ",") {
		name = strings.TrimSpace(name)
		switch {
		case name == "":
			return nil, errors.NotValidf("plugins=[%s] has an empty name", raw)
		case seen[name]:
			return nil, errors.NotValidf("plugins=[%s] lists %s twice", raw, name)
		}
		seen[name] = true
		names = append(names, name)
	}
	return names, nil
}
