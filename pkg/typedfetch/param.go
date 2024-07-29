package typedfetch

import (
	"fmt"

	"github.com/swaggest/openapi-go/openapi31"
)

type ParamInfo struct {
	TsType         string
	Required       bool
	Included       bool
	ResolvedParams []*openapi31.Parameter
}

func getParamInfo(reflector *openapi31.Reflector, op *openapi31.Operation, method, path string) (*ParamInfo, error) {
	paramType := getRequestParamTypeName(method, path)
	paramRequired := false
	paramIncluded := len(op.Parameters) > 0
	resolvedParams := []*openapi31.Parameter{}

	for _, param := range op.Parameters {
		if param.Reference != nil {
			refParam, err := resolveRefParameter(param.Reference.Ref, reflector)
			if err != nil {
				return nil, err
			}

			resolvedParams = append(resolvedParams, refParam)
		} else if param.Parameter != nil {
			resolvedParams = append(resolvedParams, param.Parameter)
		} else {
			return nil, fmt.Errorf("parameter is nil")
		}
	}

	// Check if the param is required
	for _, param := range resolvedParams {
		if param.Required != nil && *param.Required {
			paramRequired = true
			break
		}
	}

	return &ParamInfo{
		TsType:         paramType,
		Required:       paramRequired,
		Included:       paramIncluded,
		ResolvedParams: resolvedParams,
	}, nil
}

func generateParamType(method, path string, paramInfo *ParamInfo) ([]string, error) {
	lines := []string{}

	if !paramInfo.Included {
		return lines, nil
	}

	// Generate the param type
	lines = append(lines, fmt.Sprintf("type %s = {", getRequestParamTypeName(method, path)))

	paramInMap := map[openapi31.ParameterIn][]*openapi31.Parameter{}
	for _, param := range paramInfo.ResolvedParams {
		paramInMap[param.In] = append(paramInMap[param.In], param)
	}

	for in, params := range paramInMap {
		inRequired := false
		inLines := []string{}
		for _, param := range params {
			paramRequired := param.Required != nil && *param.Required
			paramType, err := jsonTypeToTypescriptType(param.Schema)
			if err != nil {
				return nil, err
			}

			paramRequiredQ := ""
			if !paramRequired {
				paramRequiredQ = "?"
			}

			if param.Description != nil {
				docString := buildDocString(*param.Description, "")
				if docString != "" {
					inLines = append(inLines, "        "+docString)
				}
			}

			inLines = append(inLines, fmt.Sprintf("        %s%s: %s;", param.Name, paramRequiredQ, paramType))

			if paramRequired {
				inRequired = true
			}
		}

		inRequiredQ := ""
		if !inRequired {
			inRequiredQ = "?"
		}

		lines = append(lines, fmt.Sprintf("    %s%s: {", in, inRequiredQ))
		lines = append(lines, inLines...)
		lines = append(lines, "    };")
	}

	lines = append(lines, "}")

	return lines, nil
}
