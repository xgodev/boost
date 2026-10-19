package inject

import (
	"fmt"
	"path"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
)

const (
	// publisherWrapper is the import path of the publisher inject builds.
	publisherWrapper = "github.com/xgodev/boost/wrapper/publisher"
	// PublisherDriverKey chooses the driver of the publisher inject builds.
	// It is not boost.wrapper.publisher.driver: that key holds the options of
	// the drivers.
	PublisherDriverKey = "boost.wrapper.publisher.driver.name"
	// DefaultPublisherDriver is the driver when PublisherDriverKey is unset.
	DefaultPublisherDriver = "noop"
	// publisherRoot is the config root of the publishers inject builds; a
	// named publisher keeps its driver key under publisherRoot.clients.<name>.
	publisherRoot = "boost.wrapper.publisher"

	driverContrib = "github.com/xgodev/boost/wrapper/publisher/driver/contrib/"
)

// PublisherDriver is a driver of the publisher that the generated module can
// wire.
type PublisherDriver struct {
	Name    string // value of PublisherDriverKey
	Path    string // import path of the driver package, whose New builds it
	Package string // package name of the driver
	Client  string // "<import path>.<Type>" of the client New takes, "" when none
	// ClientParam is the index of the client among the parameters of New.
	ClientParam int
}

// publisherDrivers are the drivers PublisherDriverKey names, sorted by name.
// Each client is one inject builds (clientFactories).
var publisherDrivers = []PublisherDriver{
	{Name: "kafka", Path: driverContrib + "confluentinc/confluent-kafka-go/v2", Package: "confluent",
		Client: "github.com/confluentinc/confluent-kafka-go/v2/kafka.Producer", ClientParam: 1},
	{Name: "nats", Path: driverContrib + "nats-io/nats.go/v1", Package: "nats",
		Client: "github.com/nats-io/nats.go.Conn"},
	{Name: "noop", Path: "github.com/xgodev/boost/wrapper/publisher/driver/extra/noop", Package: "noop"},
	{Name: "pubsub", Path: driverContrib + "cloud.google.com/pubsub/v2", Package: "pubsub",
		Client: "cloud.google.com/go/pubsub/v2.Client", ClientParam: 1},
}

// publisherDriverNames returns the driver names, comma separated.
func publisherDriverNames() string {
	names := make([]string, 0, len(publisherDrivers))
	for _, d := range publisherDrivers {
		names = append(names, d.Name)
	}
	return strings.Join(names, ", ")
}

// publisherDriver returns the driver named name.
func publisherDriver(name string) (PublisherDriver, bool) {
	for _, d := range publisherDrivers {
		if d.Name == name {
			return d, true
		}
	}
	return PublisherDriver{}, false
}

// parseDrivers reads the drivers attribute of @Inject: "[a, b]" or a single
// name. A list without names is NotValid: the publisher needs a driver.
func parseDrivers(raw string) ([]string, error) {
	names, err := parseList("drivers", raw)
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, errors.NotValidf("drivers=[] lists no driver: omit it to generate all of them (%s)", publisherDriverNames())
	}
	return names, nil
}

// Publisher is a *publisher.Publisher inject builds: the unnamed one, with
// the driver PublisherDriverKey names and the default clients, or a named
// one, with the driver its Key names and the clients of its name.
type Publisher struct {
	Name      string   // fx name of the publisher, "" for the unnamed one
	Consumers []string // keys of the components injecting it, sorted
	// Drivers are the names of the drivers generated, in the order of
	// publisherDrivers: the ones the consumers list, plus DefaultPublisherDriver;
	// all of them when a consumer lists none.
	Drivers []string
}

// PublisherOf returns the unnamed publisher inject builds: its consumers are
// the components with an unnamed, ungrouped @Inject of *publisher.Publisher
// that no annotated function provides. It has no consumers when inject builds
// none.
func PublisherOf(g *graph.Graph[Component]) Publisher {
	for _, p := range Publishers(g) {
		if p.Name == "" {
			return p
		}
	}
	return Publisher{}
}

// Publishers returns the publishers inject builds, the unnamed one first and
// the named ones by name: one per name of the ungrouped @Inject of
// *publisher.Publisher that no annotated function provides. Each has the
// drivers its consumers list.
func Publishers(g *graph.Graph[Component]) []Publisher {
	provided := providedValues(g)

	byName := map[string]*Publisher{}
	listed := map[string]map[string]bool{}
	all := map[string]bool{}
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for _, idx := range sortedIndexes(comp.Injects) {
			ann := comp.Injects[idx]
			tp := comp.Entry.Func.Parameters[idx].Type
			if ann.Group != "" || !isPublisher(comp.Entry, tp) {
				continue
			}
			if _, ok := provided[injectID(comp.Entry, tp, ann)]; ok {
				continue
			}
			p, ok := byName[ann.Name]
			if !ok {
				p = &Publisher{Name: ann.Name}
				byName[ann.Name] = p
				listed[ann.Name] = map[string]bool{DefaultPublisherDriver: true}
			}
			if !slices.Contains(p.Consumers, key) {
				p.Consumers = append(p.Consumers, key)
			}
			// validateDrivers checked the names when the graph was built
			drivers, _ := parseDrivers(ann.Drivers)
			all[ann.Name] = all[ann.Name] || ann.Drivers == ""
			for _, name := range drivers {
				listed[ann.Name][name] = true
			}
		}
	}

	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	list := make([]Publisher, 0, len(names))
	for _, name := range names {
		p := byName[name]
		for _, d := range publisherDrivers {
			if all[name] || listed[name][d.Name] {
				p.Drivers = append(p.Drivers, d.Name)
			}
		}
		list = append(list, *p)
	}
	return list
}

// isPublisher tells whether tp, written in the file of entry, is
// *publisher.Publisher.
func isPublisher(entry annotation.Entry, tp string) bool {
	id, ok := pointerTo(entry, tp)
	return ok && id == publisherWrapper+".Publisher"
}

// File returns the path, relative to OutputDir, of the publisher's file. A
// named publisher has a package of its own, under the unnamed one's.
func (p Publisher) File() string {
	return path.Join("client", publisherWrapper, p.Name, "publisher.go")
}

// Key returns the config key choosing the driver of the publisher:
// PublisherDriverKey for the unnamed one,
// boost.wrapper.publisher.clients.<name>.driver.name for a named one.
func (p Publisher) Key() string {
	if p.Name == "" {
		return PublisherDriverKey
	}
	return publisherRoot + ".clients." + p.Name + ".driver.name"
}

// String describes the publisher in messages and generated docs.
func (p Publisher) String() string {
	if p.Name == "" {
		return "the publisher"
	}
	return "the " + p.Name + " publisher"
}

// pkg returns the name of the package of the publisher's file.
func (p Publisher) pkg() string {
	if p.Name == "" {
		return "publisher"
	}
	return packageIdent(p.Name)
}

// Import returns the import of the generated package holding the publisher,
// given the import path of the output dir (BaseImport).
func (p Publisher) Import(base string) ImportData {
	pth := path.Join(base, "client", publisherWrapper, p.Name)
	return ImportData{Alias: Alias("client", p.pkg(), pth), Path: pth}
}

// Module returns the import of the publisher's package and the call of the
// module that provides it.
func (p Publisher) Module(base string) (ImportData, string) {
	imp := p.Import(base)
	return imp, imp.Alias + ".PublisherModule()"
}

// constrainedFile returns the path of the file holding the driver d, whose
// client builds only under constraint, and of its fallback.
func (p Publisher) constrainedFile(d PublisherDriver, constraint string) (string, string) {
	dir := path.Dir(p.File())
	return path.Join(dir, "publisher_"+d.Name+".go"), path.Join(dir, "publisher_"+d.Name+"_no"+constraint+".go")
}

// renderPublisher adds the files of the publisher to files and returns the
// clients of its drivers that clients lacks, the default ones or those of
// the publisher's name: the generated package of a driver's client holds the
// module the publisher wires. Only the drivers
// of p are generated, so an app links no other driver nor its client. A
// driver whose client builds only under a constraint goes in a file of its
// own, with a fallback that fails the boot when the driver is chosen.
func renderPublisher(tmpl *template.Template, base string, p Publisher, clients []Client, files map[string][]byte) ([]Client, error) {
	data := PublisherData{Key: p.Key(), Default: DefaultPublisherDriver, Names: strings.Join(p.Drivers, ", "),
		Package: p.pkg(), Name: p.Name, Desc: p.String(), Provide: "wrapper.New", Module: publisherWrapper + ".Publisher"}
	if p.Name != "" {
		data.Provide = "fx.Annotate(wrapper.New, fx.ResultTags(" + goString(fmt.Sprintf("name:%q", p.Name)) + "))"
		data.Module += "_" + p.Name
	}
	var missing []Client
	for _, name := range p.Drivers {
		d, _ := publisherDriver(name)
		drv := ImportData{Alias: Alias("driver", d.Package, d.Path), Path: d.Path}
		imports := []ImportData{drv}
		var options []string

		constraint := ""
		if d.Client != "" {
			c, ok := clientOf(clients, d.Client, p.Name)
			if !ok {
				c = Client{Factory: clientFactories[d.Client], Name: p.Name}
				missing = append(missing, c)
			}
			imp, call := c.Module(base)
			imports = append(imports, imp)
			options = append(options, call)
			constraint = c.Factory.Constraint
		}
		options = append(options, "fx.Provide(fx.Private, "+driverNew(drv.Alias, d, p.Name)+")")

		if constraint == "" {
			data.Imports = append(data.Imports, imports...)
			data.Cases = append(data.Cases, PublisherCase{Name: d.Name, Options: strings.Join(options, ", ")})
			continue
		}

		data.Cases = append(data.Cases, PublisherCase{Name: d.Name, Options: d.Name + "Driver()..."})
		driver := PublisherDriverData{Package: data.Package, Name: d.Name, Constraint: constraint, Options: strings.Join(options, ", "), Imports: imports}
		file, fallback := p.constrainedFile(d, constraint)
		for name, tmplName := range map[string]string{file: "publisher_driver", fallback: "publisher_unsupported"} {
			src, err := ExecuteTemplate(tmpl, tmplName, driver)
			if err != nil {
				return nil, errors.Wrap(err, errors.Internalf("error generating the %s driver of the publisher", d.Name))
			}
			files[name] = src
		}
	}

	src, err := ExecuteTemplate(tmpl, "publisher", data)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error generating the publisher"))
	}
	files[p.File()] = src
	return missing, nil
}

// driverNew returns the constructor of the driver d given to fx.Provide: a
// named publisher's driver takes the client of that name.
func driverNew(alias string, d PublisherDriver, name string) string {
	if name == "" || d.Client == "" {
		return alias + ".New"
	}
	tags := make([]string, d.ClientParam+1)
	for i := range tags {
		tags[i] = `""`
	}
	tags[d.ClientParam] = goString(fmt.Sprintf("name:%q", name))
	return "fx.Annotate(" + alias + ".New, fx.ParamTags(" + strings.Join(tags, ", ") + "))"
}

// defaultClient returns the default client of type id ("<import path>.<Type>")
// among clients.
func defaultClient(clients []Client, id string) (Client, bool) {
	return clientOf(clients, id, "")
}

// clientOf returns the client of type id ("<import path>.<Type>") and name
// among clients, the default one for "".
func clientOf(clients []Client, id, name string) (Client, bool) {
	for _, c := range clients {
		if c.Name == name && c.Factory.Path+"."+c.Factory.Type == id {
			return c, true
		}
	}
	return Client{}, false
}
