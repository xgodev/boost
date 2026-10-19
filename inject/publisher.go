package inject

import (
	"path"
	"slices"
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

	driverContrib = "github.com/xgodev/boost/wrapper/publisher/driver/contrib/"
)

// PublisherDriver is a driver of the publisher that the generated module can
// wire.
type PublisherDriver struct {
	Name    string // value of PublisherDriverKey
	Path    string // import path of the driver package, whose New builds it
	Package string // package name of the driver
	Client  string // "<import path>.<Type>" of the client New takes, "" when none
}

// publisherDrivers are the drivers PublisherDriverKey names, sorted by name.
// Each client is one inject builds (clientFactories).
var publisherDrivers = []PublisherDriver{
	{Name: "kafka", Path: driverContrib + "confluentinc/confluent-kafka-go/v2", Package: "confluent",
		Client: "github.com/confluentinc/confluent-kafka-go/v2/kafka.Producer"},
	{Name: "nats", Path: driverContrib + "nats-io/nats.go/v1", Package: "nats",
		Client: "github.com/nats-io/nats.go.Conn"},
	{Name: "noop", Path: "github.com/xgodev/boost/wrapper/publisher/driver/extra/noop", Package: "noop"},
	{Name: "pubsub", Path: driverContrib + "cloud.google.com/pubsub/v2", Package: "pubsub",
		Client: "cloud.google.com/go/pubsub/v2.Client"},
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

// Publisher is the *publisher.Publisher inject builds: one for the
// application, with the driver PublisherDriverKey names.
type Publisher struct {
	Consumers []string // keys of the components injecting it, sorted
	// Drivers are the names of the drivers generated, in the order of
	// publisherDrivers: the ones the consumers list, plus DefaultPublisherDriver;
	// all of them when a consumer lists none.
	Drivers []string
}

// PublisherOf returns the publisher inject builds: its consumers are the
// components with an unnamed, ungrouped @Inject of *publisher.Publisher that
// no annotated function provides. It has no consumers when inject builds none.
func PublisherOf(g *graph.Graph[Component]) Publisher {
	provided := providedValues(g)

	var p Publisher
	listed := map[string]bool{DefaultPublisherDriver: true}
	all := false
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for _, idx := range sortedIndexes(comp.Injects) {
			ann := comp.Injects[idx]
			tp := comp.Entry.Func.Parameters[idx].Type
			if ann.Name != "" || ann.Group != "" || !isPublisher(comp.Entry, tp) {
				continue
			}
			if _, ok := provided[injectID(comp.Entry, tp, ann)]; ok {
				continue
			}
			if !slices.Contains(p.Consumers, key) {
				p.Consumers = append(p.Consumers, key)
			}
			// validateDrivers checked the names when the graph was built
			drivers, _ := parseDrivers(ann.Drivers)
			all = all || ann.Drivers == ""
			for _, name := range drivers {
				listed[name] = true
			}
		}
	}

	if len(p.Consumers) == 0 {
		return p
	}
	for _, d := range publisherDrivers {
		if all || listed[d.Name] {
			p.Drivers = append(p.Drivers, d.Name)
		}
	}
	return p
}

// isPublisher tells whether tp, written in the file of entry, is
// *publisher.Publisher.
func isPublisher(entry annotation.Entry, tp string) bool {
	id, ok := pointerTo(entry, tp)
	return ok && id == publisherWrapper+".Publisher"
}

// File returns the path, relative to OutputDir, of the publisher's file.
func (p Publisher) File() string {
	return path.Join("client", publisherWrapper, "publisher.go")
}

// Import returns the import of the generated package holding the publisher,
// given the import path of the output dir (BaseImport).
func (p Publisher) Import(base string) ImportData {
	pth := path.Join(base, "client", publisherWrapper)
	return ImportData{Alias: Alias("client", "publisher", pth), Path: pth}
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
// default clients of its drivers that clients lacks: the generated package
// of a driver's client holds the module the publisher wires. Only the drivers
// of p are generated, so an app links no other driver nor its client. A
// driver whose client builds only under a constraint goes in a file of its
// own, with a fallback that fails the boot when the driver is chosen.
func renderPublisher(tmpl *template.Template, base string, p Publisher, clients []Client, files map[string][]byte) ([]Client, error) {
	data := PublisherData{Key: PublisherDriverKey, Default: DefaultPublisherDriver, Names: strings.Join(p.Drivers, ", ")}
	var missing []Client
	for _, name := range p.Drivers {
		d, _ := publisherDriver(name)
		drv := ImportData{Alias: Alias("driver", d.Package, d.Path), Path: d.Path}
		imports := []ImportData{drv}
		var options []string

		constraint := ""
		if d.Client != "" {
			c, ok := defaultClient(clients, d.Client)
			if !ok {
				c = Client{Factory: clientFactories[d.Client]}
				missing = append(missing, c)
			}
			imp, call := c.Module(base)
			imports = append(imports, imp)
			options = append(options, call)
			constraint = c.Factory.Constraint
		}
		options = append(options, "fx.Provide(fx.Private, "+drv.Alias+".New)")

		if constraint == "" {
			data.Imports = append(data.Imports, imports...)
			data.Cases = append(data.Cases, PublisherCase{Name: d.Name, Options: strings.Join(options, ", ")})
			continue
		}

		data.Cases = append(data.Cases, PublisherCase{Name: d.Name, Options: d.Name + "Driver()..."})
		driver := PublisherDriverData{Name: d.Name, Constraint: constraint, Options: strings.Join(options, ", "), Imports: imports}
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

// defaultClient returns the default client of type id ("<import path>.<Type>")
// among clients.
func defaultClient(clients []Client, id string) (Client, bool) {
	for _, c := range clients {
		if c.Name == "" && c.Factory.Path+"."+c.Factory.Type == id {
			return c, true
		}
	}
	return Client{}, false
}
