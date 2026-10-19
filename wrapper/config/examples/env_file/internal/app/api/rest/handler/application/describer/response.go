package describer

import "github.com/xgodev/boost/wrapper/config/examples/env_file/internal/pkg/domain/model/entity"

// Response is the body of GET /application.
type Response struct {
	Name   string `json:"name"`
	MyName string `json:"myName"`
}

func response(app *entity.Application) Response {
	return Response{Name: app.Name, MyName: app.MyName}
}
