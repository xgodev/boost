package builder

import (
	"strings"

	"github.com/xgodev/boost/inject"
)

// jobDoor provides the @Job handler as the job.Handler the job module runs
// once: the app exits when it returns.
func jobDoor(base, app string, job Job) (doorData, error) {
	d := door(app, jobFile, "runs the @Job handler once; the app exits when it returns.")
	imps := newImports()

	var params []string
	var modules []string
	include := func(provider inject.Component) {
		imp, mod := inject.ModuleOf(base, provider.Entry)
		imps.add(imp.Alias, imp.Path)
		modules = append(modules, mod)
	}
	call, err := handlerFunc(imps, job.Handler, map[string]string{}, &params, include)
	if err != nil {
		return doorData{}, err
	}

	handlerType := imps.add("bsjob", jobPath) + ".Handler"
	d.Options = append(d.Options, modules...)
	d.Options = append(d.Options,
		"fx.Provide(func("+strings.Join(params, ", ")+") "+handlerType+" {\nreturn "+call+"\n})",
		imps.add("jobfx", jobFxPath)+".Module()",
	)
	d.Imports = imps.sorted()
	return d, nil
}
