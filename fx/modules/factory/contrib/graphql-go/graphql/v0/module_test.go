package graphql

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModuleProvidesHandlerFromSchemaInGraph(t *testing.T) {
	boost.Start()

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"hello": &graphql.Field{
					Type:    graphql.String,
					Resolve: func(graphql.ResolveParams) (interface{}, error) { return "world", nil },
				},
			},
		}),
	})
	require.NoError(t, err)

	var h *handler.Handler
	app := fxtest.New(t, Module(), fx.Supply(&schema), fx.Populate(&h))
	app.RequireStart()
	defer app.RequireStop()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"query":"{hello}"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"data":{"hello":"world"}}`, rec.Body.String())
}
