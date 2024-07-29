package typedfetch

import (
	"fmt"
	"strings"

	"github.com/swaggest/openapi-go/openapi31"
)

func generateResponseTypes(reflector *openapi31.Reflector, op *openapi31.Operation, method, path string) ([]string, error) {
	lines := []string{}

	// Data = first non-error response or default
	dataInfo, err := getResponseObj(reflector, op, method, path, []string{"2"})
	if err != nil {
		return nil, err
	}

	dataResponseTypeName := getResponseDataTypeName(method, path)
	dataResponseLines, err := generateResponseType(method, path, dataResponseTypeName, dataInfo)
	if err != nil {
		return nil, err
	}
	lines = append(lines, dataResponseLines...)

	// Error = first error response or default
	errInfo, err := getResponseObj(reflector, op, method, path, []string{"4", "5"})
	if err != nil {
		return nil, err
	}

	errResponseTypeName := getResponseErrTypeName(method, path)
	errResponseLines, err := generateResponseType(method, path, errResponseTypeName, errInfo)
	if err != nil {
		return nil, err
	}
	lines = append(lines, errResponseLines...)

	return lines, nil
}

func generateResponseType(method, path, typeName string, response *openapi31.Response) ([]string, error) {
	lines := []string{}

	preferredContentTypes := []string{"application/json", "multipart/form-data", "application/x-www-form-urlencoded", "application/octet-stream"}
	contentType := ""
	for _, preferredContentType := range preferredContentTypes {
		if _, ok := response.Content[preferredContentType]; ok {
			contentType = preferredContentType
			break
		}
	}

	if contentType == "" {
		// default to the first content type
		for contentType = range response.Content {
			break
		}
	}

	var responseType string

	if contentType == "" {
		// Empty response
		responseType = "{}"
	} else {
		// Generate the body type
		var err error
		responseType, err = jsonTypeToTypescriptType(response.Content[contentType].Schema)
		if err != nil {
			return nil, fmt.Errorf("%s %s, %s (%s): %v", method, path, typeName, contentType, err)
		}
	}

	responseDecl := fmt.Sprintf("type %s = %s;", typeName, responseType)
	lines = append(lines, responseDecl)

	return lines, nil
}

func getResponseObj(reflector *openapi31.Reflector, op *openapi31.Operation, method, path string, codePrefixes []string) (*openapi31.Response, error) {
	if op.Responses == nil {
		return nil, fmt.Errorf("operation %s %s has no responses", method, path)
	}

	// Find the first matching response
	for _, codePrefix := range codePrefixes {
		for code, responseOrRef := range op.Responses.MapOfResponseOrReferenceValues {
			if strings.HasPrefix(code, codePrefix) {
				resolvedResponse, err := resolveResponseOrReference(reflector, &responseOrRef)
				if err != nil {
					return nil, err
				}

				return resolvedResponse, nil
			}
		}
	}

	// If no matching response was found, use the default response
	if op.Responses.Default != nil {
		resolvedResponse, err := resolveResponseOrReference(reflector, op.Responses.Default)
		if err != nil {
			return nil, err
		}

		return resolvedResponse, nil
	}

	// No matching response or default response found, return an empty response
	response := &openapi31.Response{
		Content: map[string]openapi31.MediaType{},
	}

	return response, nil
}

func resolveResponseOrReference(reflector *openapi31.Reflector, responseOrRef *openapi31.ResponseOrReference) (*openapi31.Response, error) {
	if responseOrRef.Reference != nil {
		return resolveRefResponse(responseOrRef.Reference.Ref, reflector)
	}

	return responseOrRef.Response, nil
}

func resolveRefResponse(ref string, reflector *openapi31.Reflector) (*openapi31.Response, error) {
	if !strings.HasPrefix(ref, "#/components/responses/") {
		return nil, fmt.Errorf("reference %s is not a response", ref)
	}

	responseName := strings.TrimPrefix(ref, "#/components/responses/")
	responseOrReference, ok := reflector.Spec.Components.Responses[responseName]
	if !ok {
		return nil, fmt.Errorf("response %s not found", responseName)
	}

	if responseOrReference.Reference != nil {
		return resolveRefResponse(responseOrReference.Reference.Ref, reflector)
	}

	return responseOrReference.Response, nil
}

func getResponseDataTypeName(method, path string) string {
	return fmt.Sprintf("ResponseData%s", getUniqueEndpointName(method, path))
}

func getResponseErrTypeName(method, path string) string {
	return fmt.Sprintf("ResponseError%s", getUniqueEndpointName(method, path))
}
