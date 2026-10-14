package inject

import (
	"os"
	"path/filepath"

	"github.com/xgodev/boost/model/errors"
	"golang.org/x/mod/modfile"
)

// FindModule walks up from dir to the nearest go.mod and returns the module
// path and its root directory.
func FindModule(dir string) (modulePath string, moduleDir string, err error) {
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", "", errors.Wrap(err, errors.Internalf("error resolving %s", dir))
	}

	for current := dir; ; current = filepath.Dir(current) {
		data, err := os.ReadFile(filepath.Join(current, "go.mod"))
		if err == nil {
			modulePath = modfile.ModulePath(data)
			if modulePath == "" {
				return "", "", errors.NotValidf("go.mod in %s has no module directive", current)
			}
			return modulePath, current, nil
		}
		if !os.IsNotExist(err) {
			return "", "", errors.Wrap(err, errors.Internalf("error reading go.mod in %s", current))
		}
		if parent := filepath.Dir(current); parent == current {
			return "", "", errors.NotFoundf("no go.mod found in %s or its parents", dir)
		}
	}
}
