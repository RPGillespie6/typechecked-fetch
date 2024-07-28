package typedfetch

import (
	"fmt"

	"github.com/swaggest/openapi-go/openapi31"
)

type RequestBodyInfo struct {
	TsType       string
	Required     bool
	Included     bool
	ResolvedBody *openapi31.RequestBody
}

func getRequestBodyInfo(reflector *openapi31.Reflector, op *openapi31.Operation, method, path string) (*RequestBodyInfo, error) {
	bodyType := getRequestBodyTypeName(method, path)
	bodyRequired := false
	bodyIncluded := op.RequestBody != nil
	var resolvedBody *openapi31.RequestBody

	if bodyIncluded {
		if op.RequestBody.Reference != nil {
			body, err := resolveRefRequestBody(op.RequestBody.Reference.Ref, reflector)
			if err != nil {
				return nil, err
			}

			resolvedBody = body
		} else if op.RequestBody.RequestBody != nil {
			resolvedBody = op.RequestBody.RequestBody
		} else {
			return nil, fmt.Errorf("request body is nil")
		}

		if resolvedBody.Required != nil && *resolvedBody.Required {
			bodyRequired = true
		}
	}

	return &RequestBodyInfo{
		TsType:       bodyType,
		Required:     bodyRequired,
		Included:     bodyIncluded,
		ResolvedBody: resolvedBody,
	}, nil
}

func generateBodyType(method, path string, bodyInfo *RequestBodyInfo) ([]string, error) {
	lines := []string{}

	if !bodyInfo.Included {
		return lines, nil
	}

	preferredContentTypes := []string{"application/json", "multipart/form-data", "application/x-www-form-urlencoded", "application/octet-stream"}
	contentType := ""
	for _, preferredContentType := range preferredContentTypes {
		if _, ok := bodyInfo.ResolvedBody.Content[preferredContentType]; ok {
			contentType = preferredContentType
			break
		}
	}

	if contentType == "" {
		// default to the first content type
		for contentType = range bodyInfo.ResolvedBody.Content {
			break
		}
	}

	if contentType == "" {
		return nil, fmt.Errorf("no content type found for request body")
	}

	// TODO: register the content type in map?

	// Generate the body type
	bodyType, err := jsonTypeToTypescriptType(bodyInfo.ResolvedBody.Content[contentType].Schema)
	if err != nil {
		return nil, err
	}

	bodyDecl := fmt.Sprintf("type %s = %s;", getRequestBodyTypeName(method, path), bodyType)

	lines = append(lines, bodyDecl)

	return lines, nil
}
