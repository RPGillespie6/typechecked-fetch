package typedfetch

import (
	"fmt"
	"net/http"
	"strings"

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

func getMethodPathMap(reflector *openapi31.Reflector) map[string][]string {
	methodToPaths := map[string][]string{}
	sortedPaths := sortedMapKeys(reflector.Spec.Paths.MapOfPathItemValues)
	for _, path := range sortedPaths {
		item := reflector.Spec.Paths.MapOfPathItemValues[path]
		methods := getPathItemMethods(&item)
		for _, method := range methods {
			if method.Operation != nil {
				methodToPaths[method.Method] = append(methodToPaths[method.Method], path)
			}
		}
	}
	return methodToPaths
}

func generateUrlTypes(reflector *openapi31.Reflector) []string {
	lines := []string{"// URL types"}

	// Generate lines like:
	// type UrlGetPet = '/pet';
	methodToPaths := getMethodPathMap(reflector)
	sortedMethods := sortedMapKeys(methodToPaths)
	for _, method := range sortedMethods {
		paths := methodToPaths[method]
		methodTypes := []string{}
		for _, path := range paths {
			methodType := getUrlTypeName(method, path)
			methodTypes = append(methodTypes, methodType)
			lines = append(lines, fmt.Sprintf("type %s = '%s';", methodType, path))
		}

		// Generate line like:
		// type UrlGet = UrlGetPet | UrlGetPets;
		lines = append(lines, fmt.Sprintf("type %s = %s;", getMethodUrlTypeName(method), strings.Join(methodTypes, " | ")))
		lines = append(lines, "")
	}

	return lines
}

func getMethodUrlTypeName(method string) string {
	return fmt.Sprintf("UrlValid%s", pascalize(method))
}

func getUrlTypeName(method, path string) string {
	return fmt.Sprintf("Url%s", getUniqueEndpointName(method, path))
}

func getUniqueEndpointName(method, path string) string {
	return fmt.Sprintf("%s%s", pascalize(method), pathToVar(path))
}
