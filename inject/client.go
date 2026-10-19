package inject

import (
	"fmt"
	"go/ast"
	"go/parser"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
)

// ClientFactory is a client type inject builds through its boost factory
// when an @Inject parameter asks for it.
type ClientFactory struct {
	Path    string                  // import path of the client type's package
	Type    string                  // name of the client type, always used as a pointer
	Package string                  // name of the client type's package ("resty", "redis", "mongo")
	Build   string                  // factory import path + "." + the ...WithConfigPath constructor
	Root    string                  // config root of the factory
	Plugins map[string]ClientPlugin // by name: the last segment of the plugin's config root
	Close   bool                    // the generated module closes the client when the app stops
}

// ClientPlugin is a plugin of a client factory. The value of its constructor
// registers itself through its Register method, or the method named by
// Register when the factory takes a plugin of another shape.
type ClientPlugin struct {
	Path     string // import path of the plugin
	New      string // ...WithConfigPath constructor of the plugin
	Register string // method passed to the factory, "" for Register
}

const boostFactories = "github.com/xgodev/boost/factory/contrib/"

// clientFactories are the clients inject builds, by Path + "." + Type.
var clientFactories = map[string]ClientFactory{}

func init() {
	resty := boostFactories + "go-resty/resty/v2"
	addClientFactory("github.com/go-resty/resty/v2", "Client", "resty", resty, "NewClientWithConfigPath", "boost.factory.resty", map[string]ClientPlugin{
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
	addClientFactory("github.com/redis/go-redis/v9", "Client", "redis", redis, "NewClientWithConfigPath", "boost.factory.redis",
		redisPlugins("NewClientHealthWithConfigPath", "NewClientWithConfigPath"))
	addClientFactory("github.com/redis/go-redis/v9", "ClusterClient", "redis", redis, "NewClusterClientWithConfigPath", "boost.factory.redis",
		redisPlugins("NewClusterHealthWithConfigPath", "NewClusterWithConfigPath"))

	for _, version := range []string{"v1", "v2"} {
		mongo := boostFactories + "go.mongodb.org/mongo-driver/" + version
		addClientFactory(mongo, "Conn", "mongo", mongo, "NewConnWithConfigPath", "boost.factory.mongo", map[string]ClientPlugin{
			"health":     {Path: mongo + "/plugins/local/health", New: "NewHealthWithConfigPath"},
			"otel":       {Path: mongo + "/plugins/contrib/go.opentelemetry.io/contrib/v0", New: "NewOtelMongoWithConfigPath"},
			"datadog":    {Path: mongo + "/plugins/contrib/datadog/dd-trace-go/v1", New: "NewDatadogWithConfigPath"},
			"prometheus": {Path: mongo + "/plugins/contrib/prometheus/client_golang/v1", New: "NewPrometheusWithConfigPath"},
		})
	}

	nats := boostFactories + "nats-io/nats.go/v1"
	addClosingClientFactory("github.com/nats-io/nats.go", "Conn", "nats", nats, "NewConnWithConfigPath", "boost.factory.nats", map[string]ClientPlugin{
		"health": {Path: nats + "/plugins/local/health", New: "NewHealthWithConfigPath", Register: "Plugin"},
		"log":    {Path: nats + "/plugins/local/log", New: "NewLogWithConfigPath"},
	})

	grpc := boostFactories + "google.golang.org/grpc/v1/client/plugins"
	addClosingClientFactory("cloud.google.com/go/pubsub/v2", "Client", "pubsub", boostFactories+"cloud.google.com/pubsub/v2", "NewClientWithConfigPath", "boost.factory.gcp.pubsub", map[string]ClientPlugin{
		"log":        {Path: grpc + "/local/wrapper/log", New: "NewLogWithConfigPath"},
		"datadog":    {Path: grpc + "/contrib/datadog/dd-trace-go/v1", New: "NewDatadogWithConfigPath"},
		"prometheus": {Path: grpc + "/contrib/prometheus/client_golang/v1", New: "NewPrometheusWithConfigPath"},
	})

	for path, factory := range map[string]string{
		"github.com/gocql/gocql":                      boostFactories + "gocql/gocql/v1",
		"github.com/apache/cassandra-gocql-driver/v2": boostFactories + "apache/cassandra-gocql-driver/v2",
	} {
		addClosingClientFactory(path, "Session", "gocql", factory, "NewSessionWithConfigPath", "boost.factory.gocql", map[string]ClientPlugin{
			"health": {Path: factory + "/plugins/local/extra/health", New: "NewHealthWithConfigPath"},
		})
	}
}

func addClientFactory(path, tp, pkg, factory, build, root string, plugins map[string]ClientPlugin) {
	clientFactories[path+"."+tp] = ClientFactory{Path: path, Type: tp, Package: pkg, Build: factory + "." + build, Root: root, Plugins: plugins}
}

// addClosingClientFactory adds a client whose module closes it when the app
// stops: its Close method releases a connection the factory opened.
func addClosingClientFactory(path, tp, pkg, factory, build, root string, plugins map[string]ClientPlugin) {
	addClientFactory(path, tp, pkg, factory, build, root, plugins)
	f := clientFactories[path+"."+tp]
	f.Close = true
	clientFactories[path+"."+tp] = f
}

// supportedClients returns the client types inject builds, sorted.
func supportedClients() []string {
	types := make([]string, 0, len(clientFactories))
	for key := range clientFactories {
		types = append(types, "*"+key)
	}
	sort.Strings(types)
	return types
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

// Client is a client inject builds through its factory: one per client type
// and name asked by @Inject parameters that no annotated function provides.
type Client struct {
	Factory   ClientFactory
	Name      string   // fx name of the client, "" for the default client
	Plugins   []string // in the declared order, nil when no consumer lists them
	Consumers []string // keys of the components injecting the client, sorted
}

// Clients returns the clients inject builds, sorted by type and name. A
// consumer without plugins takes the client as another consumer declares it.
// It is NotValid when plugins are listed for a client the application
// provides, a name is not reachable from an environment variable or two
// consumers of a client list different plugins, and AlreadyExists when two
// clients share a file (a client named "default" and the unnamed one).
func Clients(g *graph.Graph[Component]) ([]Client, error) {
	provided := map[string]string{} // value id -> key of its first provider
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for _, idx := range sortedIndexes(comp.Provides) {
			id := xid(comp.Entry.Package, comp.Result(idx), comp.Provides[idx])
			if _, ok := provided[id]; !ok {
				provided[id] = key
			}
		}
	}

	clients := map[string]*Client{}
	declaredBy := map[string]string{} // client id -> key of the consumer listing its plugins
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for _, idx := range sortedIndexes(comp.Injects) {
			ann := comp.Injects[idx]
			tp := comp.Entry.Func.Parameters[idx].Type
			f, ok := clientFactory(comp.Entry, tp)
			if !ok || ann.Group != "" {
				continue
			}

			if provider, ok := provided[injectID(comp.Entry.Package, tp, ann)]; ok {
				if ann.Plugins != "" {
					return nil, errors.NotValidf("@Inject on %s parameter %d has plugins, but %s is provided by %s: plugins apply only to a client inject builds",
						key, idx, tp, provider)
				}
				continue
			}
			if ann.Name != "" && !envSegment.MatchString(ann.Name) {
				return nil, errors.NotValidf("@Inject on %s parameter %d: client name %s must match %s so an environment variable reaches its keys",
					key, idx, ann.Name, envSegment)
			}

			plugins, err := parsePlugins(ann.Plugins)
			if err != nil {
				return nil, err
			}

			id := f.Path + "." + f.Type + "_" + ann.Name
			c, ok := clients[id]
			if !ok {
				c = &Client{Factory: f, Name: ann.Name}
				clients[id] = c
			}
			if plugins != nil {
				if c.Plugins != nil && !slices.Equal(c.Plugins, plugins) {
					return nil, errors.NotValidf("%s asks for the %s with plugins %s, but %s asks for it with %s: the consumers of a client list the same plugins, or only one lists them",
						declaredBy[id], c, pluginList(c.Plugins), key, pluginList(plugins))
				}
				c.Plugins = plugins
				declaredBy[id] = key
			}
			if !slices.Contains(c.Consumers, key) {
				c.Consumers = append(c.Consumers, key)
			}
		}
	}

	list := make([]Client, 0, len(clients))
	for _, c := range clients {
		list = append(list, *c)
	}
	sort.Slice(list, func(i, j int) bool {
		a, b := list[i].Factory, list[j].Factory
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		return list[i].Name < list[j].Name
	})

	files := map[string]Client{}
	for _, c := range list {
		if other, dup := files[c.File()]; dup {
			return nil, errors.AlreadyExistsf("%s, asked by %s, and %s, asked by %s, would both be written to %s; rename the client",
				other, strings.Join(other.Consumers, ", "), c, strings.Join(c.Consumers, ", "), c.File())
		}
		files[c.File()] = c
	}
	return list, nil
}

func pluginList(plugins []string) string {
	return "[" + strings.Join(plugins, ", ") + "]"
}

func (c Client) nameOrDefault() string {
	if c.Name == "" {
		return "default"
	}
	return c.Name
}

// String describes the client in messages: "quote client (*github.com/go-resty/resty/v2.Client)".
func (c Client) String() string {
	return c.nameOrDefault() + " client (*" + c.Factory.Path + "." + c.Factory.Type + ")"
}

// Root returns the config root of the client: the factory's root for the
// default client, <factory root>.clients.<name> for a named one.
func (c Client) Root() string {
	if c.Name == "" {
		return c.Factory.Root
	}
	return c.Factory.Root + ".clients." + c.Name
}

// PluginRoot returns the config root of a plugin of the client.
func (c Client) PluginRoot(plugin string) string {
	return c.Root() + ".plugins." + plugin
}

// Ident is the exported identifier of the client in generated code: its
// name plus its type ("QuoteClient", "DefaultClusterClient").
func (c Client) Ident() string {
	return upperFirst(c.nameOrDefault()) + c.Factory.Type
}

// File returns the path, relative to OutputDir, of the client's file.
func (c Client) File() string {
	return path.Join("client", c.Factory.Path, strings.ToLower(c.nameOrDefault())+"_"+strings.ToLower(c.Factory.Type)+".go")
}

// Import returns the import of the generated package holding the client,
// given the import path of the output dir (BaseImport).
func (c Client) Import(base string) ImportData {
	p := path.Join(base, "client", c.Factory.Path)
	return ImportData{Alias: Alias("client", c.Factory.Package, p), Path: p}
}

// Module returns the import of the client's package and the call of the
// module that provides the client.
func (c Client) Module(base string) (ImportData, string) {
	imp := c.Import(base)
	return imp, imp.Alias + "." + c.Ident() + "Module()"
}

// clientData collects what the client template needs for one client.
func clientData(c Client) ClientData {
	ident := c.Ident()
	data := ClientData{
		Package: packageIdent(c.Factory.Package),
		Ident:   ident,
		Var:     strings.ToLower(ident[:1]) + ident[1:],
		Desc:    c.String(),
		Root:    c.Root(),
		Env:     ConfigEnv(c.Root(), "*"),
		Module:  c.Factory.Path + "." + ident,
		Target:  "new" + ident,
		Type:    "*factory." + c.Factory.Type,
		Build:   c.Factory.Constructor(),
		Names:   strings.Join(c.Plugins, ", "),
		Close:   c.Factory.Close,
		Imports: []ImportData{
			{Alias: "factory", Path: c.Factory.FactoryPath()},
			{Alias: "contextfx", Path: "github.com/xgodev/boost/fx/modules/core/context"},
		},
	}
	if c.Factory.Path != c.Factory.FactoryPath() {
		data.Type = "*client." + c.Factory.Type
		data.Imports = append(data.Imports, ImportData{Alias: "client", Path: c.Factory.Path})
	}
	if c.Name != "" {
		data.Target = "fx.Annotate(" + data.Target + ", fx.ResultTags(" + goString(fmt.Sprintf("name:%q", c.Name)) + "))"
	}
	for _, name := range c.Plugins {
		plugin := c.Factory.Plugins[name]
		alias := name + "plugin"
		register := plugin.Register
		if register == "" {
			register = "Register"
		}
		data.Plugins = append(data.Plugins, ClientPluginData{Name: name, Alias: alias, New: plugin.New, Register: register})
		data.Imports = append(data.Imports, ImportData{Alias: alias, Path: plugin.Path})
	}
	return data
}

// clientsByConsumer maps the key of each consumer to the clients it injects.
func clientsByConsumer(clients []Client) map[string][]Client {
	byConsumer := map[string][]Client{}
	for _, c := range clients {
		for _, key := range c.Consumers {
			byConsumer[key] = append(byConsumer[key], c)
		}
	}
	return byConsumer
}
