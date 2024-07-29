package typedfetch

import (
	"fmt"

	"github.com/swaggest/openapi-go/openapi31"
)

func generateOperationTypes(reflector *openapi31.Reflector) ([]string, error) {
	lines := []string{
		"// Request/Response types",
		"",
	}

	sortedPaths := sortedMapKeys(reflector.Spec.Paths.MapOfPathItemValues)
	for _, path := range sortedPaths {
		item := reflector.Spec.Paths.MapOfPathItemValues[path]
		methods := getPathItemMethods(&item)
		for _, method := range methods {
			if method.Operation == nil {
				continue
			}

			lines = append(lines, fmt.Sprintf("// %s %s", method.Method, path))

			// Generate the param type
			paramInfo, err := getParamInfo(reflector, method.Operation, method.Method, path)
			if err != nil {
				return nil, err
			}

			// Generate the body type
			bodyInfo, err := getRequestBodyInfo(reflector, method.Operation, method.Method, path)
			if err != nil {
				return nil, err
			}

			requestLines, err := generateRequestTypes(method.Method, path, paramInfo, bodyInfo)
			if err != nil {
				return nil, err
			}
			lines = append(lines, requestLines...)

			// Generate the response types
			responseLines, err := generateResponseTypes(reflector, method.Operation, method.Method, path)
			if err != nil {
				return nil, err
			}
			lines = append(lines, responseLines...)
			lines = append(lines, "")
		}
	}

	return lines, nil
}
