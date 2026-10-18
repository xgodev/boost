package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/inject"
)

// TestDecorator_InternalPackage checks that the decorator of a constructor
// under internal/ is generated where the modules of app api may import it.
func TestDecorator_InternalPackage(t *testing.T) {
	entry := annotation.Entry{
		Path:    testModule + "/internal/pkg/domain/service/reservestock/local",
		Package: "local",
		Func:    annotation.EntryFunc{Name: "New"},
	}
	m := Middleware{Provider: inject.Component{Entry: entry}}

	dir := "builder/middleware/" + testModule + "/internal_/pkg/domain/service/reservestock/local"
	assert.Equal(t, dir+"/new_middleware.go", decoratorFile(m))

	imp, call := decoratorOf(testModule+"/gen", m)
	assert.Equal(t, testModule+"/gen/"+dir, imp.Path)
	assert.Equal(t, imp.Alias+".DecorateNew()", call)
}
