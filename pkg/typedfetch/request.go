package typedfetch

import (
	"fmt"
)

func generateRequestTypes(method, path string, paramInfo *ParamInfo, bodyInfo *RequestBodyInfo) ([]string, error) {
	lines := []string{}

	paramLines, err := generateParamType(method, path, paramInfo)
	if err != nil {
		return nil, err
	}
	lines = append(lines, paramLines...)

	bodyLines, err := generateBodyType(method, path, bodyInfo)
	if err != nil {
		return nil, err
	}
	lines = append(lines, bodyLines...)

	// Generate the request type
	// Example:
	// type FetchRequestGetFoo = RequestInit & { params?: RequestParamGetFoo; };
	// type FetchRequestGetFoo2 = RequestInit;
	// type FetchRequestPostBar = Omit<RequestInit, "body"> & { params: RequestParamPostBar; body: ComponentSchemaAddress; };
	requestTypeLines, err := generateRequestType(method, path, paramInfo.TsType, paramInfo.Required, paramInfo.Included, bodyInfo.TsType, bodyInfo.Required, bodyInfo.Included)
	if err != nil {
		return nil, err
	}
	lines = append(lines, requestTypeLines...)

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

	lines = append(lines, fmt.Sprintf("type %s = %s%s & RequestInitExtended;", getRequestTypeName(method, path), baseType, intersectionType))

	return lines, nil
}

func getRequestTypeName(method, path string) string {
	return fmt.Sprintf("Request%s", getUniqueEndpointName(method, path))
}

func getRequestParamTypeName(method, path string) string {
	return fmt.Sprintf("Param%s", getUniqueEndpointName(method, path))
}

func getRequestBodyTypeName(method, path string) string {
	return fmt.Sprintf("Body%s", getUniqueEndpointName(method, path))
}
