package builder

import "github.com/xgodev/boost/annotation"

// CollectEntries returns the declarations under path, and in the packages
// they import, annotated with a door (@REST, @GRPC, @Function) or with the
// inject annotations (@Provide, @Inject, @Invoke) that build their receivers.
func CollectEntries(path string) ([]annotation.Entry, error) {
	collector, err := annotation.Collect(
		annotation.WithPath(path),
		annotation.WithFilters("Inject", "Provide", "Invoke", "REST", "GRPC", "Function"),
	)
	if err != nil {
		return []annotation.Entry{}, err
	}

	return collector.Entries(), nil
}
