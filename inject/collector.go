package inject

import "github.com/xgodev/boost/annotation"

// CollectEntries returns the functions under path, and in the packages they
// import, annotated with @Provide, @Inject, @Invoke or @Default.
func CollectEntries(path string) ([]annotation.Entry, error) {
	collector, err := annotation.Collect(
		annotation.WithPath(path),
		annotation.WithFilters("Inject", "Provide", "Invoke", "Default"),
	)
	if err != nil {
		return []annotation.Entry{}, err
	}

	return collector.Entries(), nil
}
