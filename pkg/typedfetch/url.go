package typedfetch

import (
	"fmt"
	"net/http"

	"github.com/swaggest/openapi-go/openapi31"
)

type OperationMethodTuple struct {
	Operation *openapi31.Operation
	Method    string
}

func getPathItemMethods(item *openapi31.PathItem) []OperationMethodTuple {
	return []OperationMethodTuple{
		{Operation: item.Get, Method: http.MethodGet},
		{Operation: item.Put, Method: http.MethodPut},
		{Operation: item.Post, Method: http.MethodPost},
		{Operation: item.Delete, Method: http.MethodDelete},
		{Operation: item.Options, Method: http.MethodOptions},
		{Operation: item.Head, Method: http.MethodHead},
		{Operation: item.Patch, Method: http.MethodPatch},
		{Operation: item.Trace, Method: http.MethodTrace},
	}
}

func getUniqueEndpointName(method, path string) string {
	return fmt.Sprintf("%s%s", pascalize(method), pathToVar(path))
}
