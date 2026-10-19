package builder

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// run is one generation of the application under src: the inject modules and
// doors go to out and the mains to cmd, all absolute.
type run struct {
	src, out, cmd string
	modules       int // inject modules
	apps          int
	inject        *inject.Generator
	doors         *Generator
}

func newRun(ctx context.Context, src, out, cmd string) (*run, error) {
	r := &run{}
	for _, p := range []struct {
		dst  *string
		path string
	}{{&r.src, src}, {&r.out, out}, {&r.cmd, cmd}} {
		if p.path == "" {
			continue
		}
		abs, err := filepath.Abs(p.path)
		if err != nil {
			return nil, errors.Wrap(err, errors.Internalf("error resolving %s", p.path))
		}
		*p.dst = abs
	}

	modulePath, moduleDir, err := inject.FindModule(r.src)
	if err != nil {
		return nil, err
	}
	log.FromContext(ctx).Debugf("scanning %s (module %s)", r.src, modulePath)

	entries, err := CollectEntries(r.src)
	if err != nil {
		return nil, err
	}
	g, err := inject.NewGraphFromEntries(ctx, entries)
	if err != nil {
		return nil, err
	}
	spec, err := NewSpec(entries, g, LoadPackages(r.src))
	if err != nil {
		return nil, err
	}

	opts := Options{
		GeneratorOptions: inject.GeneratorOptions{ModulePath: modulePath, ModuleDir: moduleDir, OutputDir: r.out,
			DefaultClients: spec.DefaultClients(), NamedClients: spec.NamedClients()},
		CmdDir: r.cmd,
	}
	r.modules, r.apps = len(g.Vertices), len(spec.Apps)
	r.inject = inject.NewGenerator(g, opts.GeneratorOptions)
	r.doors = NewGenerator(spec, opts)
	return r, nil
}

// Generate writes the inject modules and the doors of the annotations under
// src to out, and the main of every app to cmd (none when cmd is empty): what
// cmd/builder does. It writes nothing when an annotation is not valid or the
// main of an app was written by hand.
func Generate(ctx context.Context, src, out, cmd string) error {
	r, err := newRun(ctx, src, out, cmd)
	if err != nil {
		return err
	}

	// render the doors and mains before writing anything, so an error leaves
	// out and cmd untouched
	if _, err := r.doors.Render(ctx); err != nil {
		return err
	}
	if _, err := r.doors.RenderCmd(ctx); err != nil {
		return err
	}
	if err := r.inject.Generate(ctx); err != nil {
		return err
	}
	if err := r.doors.Generate(ctx); err != nil {
		return err
	}

	log.FromContext(ctx).Infof("generated %d inject modules and %d apps in %s, mains in %s", r.modules, r.apps, r.out, r.cmd)
	return nil
}

// Stale returns the files Generate would write, change or remove under out
// and cmd, as slash-separated paths relative to src, sorted; none when the
// generated code is up to date. It writes nothing, so a test of the
// application can fail when gen/ or a main is stale:
//
//	stale, err := builder.Stale(ctx, ".", "gen", "cmd")
func Stale(ctx context.Context, src, out, cmd string) ([]string, error) {
	r, err := newRun(ctx, src, out, cmd)
	if err != nil {
		return nil, err
	}

	files, err := r.inject.Render(ctx)
	if err != nil {
		return nil, err
	}
	doors, err := r.doors.Render(ctx)
	if err != nil {
		return nil, err
	}
	for name, content := range doors {
		files[name] = content
	}
	mains, err := r.doors.RenderCmd(ctx)
	if err != nil {
		return nil, err
	}

	var stale []string
	report := func(p string) error {
		rel, err := filepath.Rel(r.src, p)
		if err != nil {
			return errors.Wrap(err, errors.Internalf("error resolving %s", p))
		}
		stale = append(stale, filepath.ToSlash(rel))
		return nil
	}
	if err := diff(r.out, files, report, inject.GeneratedHeader, generatedHeader); err != nil {
		return nil, err
	}
	if r.cmd != "" {
		if err := diff(r.cmd, mains, report, generatedHeader); err != nil {
			return nil, err
		}
	}

	sort.Strings(stale)
	return stale, nil
}

// diff reports the files of want, keyed by their slash-separated path
// relative to dir, that are missing or differ on disk, and the .go files
// under dir that start with one of headers and are not in want: the ones the
// generators would remove.
func diff(dir string, want map[string][]byte, report func(path string) error, headers ...string) error {
	for name, content := range want {
		p := filepath.Join(dir, filepath.FromSlash(name))
		got, err := os.ReadFile(p)
		if err != nil && !os.IsNotExist(err) {
			return errors.Wrap(err, errors.Internalf("error reading %s", p))
		}
		if err != nil || !bytes.Equal(got, content) {
			if err := report(p); err != nil {
				return err
			}
		}
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(p) != ".go" {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		if _, ok := want[filepath.ToSlash(rel)]; ok {
			return nil
		}
		content, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for _, header := range headers {
			if bytes.HasPrefix(content, []byte(header)) {
				return report(p)
			}
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error comparing the generated files in %s", dir))
	}
	return nil
}
