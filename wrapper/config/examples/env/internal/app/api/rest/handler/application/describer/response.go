package describer

import "github.com/xgodev/boost/wrapper/config/examples/env/internal/pkg/domain/model/entity"

// Response is the body of GET /application.
type Response struct {
	Name string `json:"name"`
}

func response(app *entity.Application) Response {
	return Response{Name: app.Name}
}
