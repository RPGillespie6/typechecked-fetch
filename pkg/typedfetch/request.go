package typedfetch

import (
	"fmt"

	"github.com/swaggest/openapi-go/openapi31"
)

func generateRequestTypes(reflector *openapi31.Reflector) ([]string, error) {
	lines := []string{
		"// Request types",
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

			paramLines, err := generateParamType(method.Method, path, paramInfo)
			if err != nil {
				return nil, err
			}
			lines = append(lines, paramLines...)

			// Generate the body type
			bodyInfo, err := getRequestBodyInfo(reflector, method.Operation, method.Method, path)
			if err != nil {
				return nil, err
			}

			bodyLines, err := generateBodyType(method.Method, path, bodyInfo)
			if err != nil {
				return nil, err
			}
			lines = append(lines, bodyLines...)

			// Generate the request type
			// Example:
			// type FetchRequestGetFoo = RequestInit & { params?: RequestParamGetFoo; };
			// type FetchRequestGetFoo2 = RequestInit;
			// type FetchRequestPostBar = Omit<RequestInit, "body"> & { params: RequestParamPostBar; body: ComponentSchemaAddress; };
			requestTypeLines, err := generateRequestType(method.Method, path, paramInfo.TsType, paramInfo.Required, paramInfo.Included, bodyInfo.TsType, bodyInfo.Required, bodyInfo.Included)
			if err != nil {
				return nil, err
			}
			lines = append(lines, requestTypeLines...)
		}
	}

	return lines, nil
}

func generateRequestType(method, path, paramType string, paramRequired, paramIncluded bool, bodyType string, bodyRequired, bodyIncluded bool) ([]string, error) {
	lines := []string{}

	paramRequiredQ := ""
	if !paramRequired {
		paramRequiredQ = "?"
	}

	bodyRequiredQ := ""
	if !bodyRequired {
		bodyRequiredQ = "?"
	}

	paramDecl := ""
	if paramIncluded {
		paramDecl = fmt.Sprintf("params%s: %s;", paramRequiredQ, paramType)
	}

	bodyDecl := ""
	if bodyIncluded {
		bodyDecl = fmt.Sprintf("body%s: %s;", bodyRequiredQ, bodyType)
	}

	baseType := `Omit<RequestInit, 'body'>`
	if !bodyIncluded {
		baseType = "RequestInit"
	}

	intersectionType := ""
	if paramDecl != "" || bodyDecl != "" {
		intersectionType = fmt.Sprintf(" & { %s %s }", paramDecl, bodyDecl)
	}

	lines = append(lines, fmt.Sprintf("type Request%s = %s%s;", getUniqueEndpointName(method, path), baseType, intersectionType))
	lines = append(lines, "")

	return lines, nil
}

func getRequestParamTypeName(method, path string) string {
	return fmt.Sprintf("Param%s", getUniqueEndpointName(method, path))
}

func getRequestBodyTypeName(method, path string) string {
	return fmt.Sprintf("Body%s", getUniqueEndpointName(method, path))
}
