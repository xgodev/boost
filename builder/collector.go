package builder

import "github.com/xgodev/boost/annotation"

// CollectEntries returns the declarations under path, and in the packages
// they import, annotated with a door (@REST, @GRPC, @MCP, @Function, @Job) or with the
// inject annotations (@Provide, @Inject, @Invoke, @Default, @Config) that
// build their receivers, the @Middleware of constructors and the package docs
// with the @Server of each server.
func CollectEntries(path string) ([]annotation.Entry, error) {
	collector, err := annotation.Collect(
		annotation.WithPath(path),
		annotation.WithFilters("Inject", "Provide", "Invoke", "Default", "Config", "REST", "GRPC", "MCP", "Function", "Job", "Middleware", "Server"),
		annotation.WithPackageDocs(),
	)
	if err != nil {
		return []annotation.Entry{}, err
	}

	return collector.Entries(), nil
}
