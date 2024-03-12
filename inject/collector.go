package inject

import "github.com/xgodev/boost/annotation"

func CollectEntries(path string) ([]annotation.Entry, error) {
	collector, err := annotation.Collect(
		annotation.WithPath(path),
		annotation.WithFilters("Inject", "Provide", "Invoke"),
	)
	if err != nil {
		return []annotation.Entry{}, err
	}

	return collector.Entries(), nil
}
