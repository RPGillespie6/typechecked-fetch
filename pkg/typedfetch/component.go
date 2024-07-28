package typedfetch

import (
	"fmt"
	"strings"

	"github.com/swaggest/openapi-go/openapi31"
)

func generateComponentSchemaTypes(reflector *openapi31.Reflector) ([]string, error) {
	lines := []string{
		"// Component types",
		"",
	}

	// For each explicit component, generate the type
	sortedComponents := sortedMapKeys(reflector.Spec.Components.Schemas)
	for _, component := range sortedComponents {
		item := reflector.Spec.Components.Schemas[component]
		componentName := getComponentSchemaTypeName(component)
		typeDecl, err := jsonTypeToTypescriptType(item)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", componentName, err)
		}

		docString := getDocString(item)
		if docString != "" {
			lines = append(lines, docString)
		}
		lines = append(lines, fmt.Sprintf("type %s = %s", componentName, typeDecl))
		lines = append(lines, "")
	}

	return lines, nil
}

func resolveRefParameter(ref string, reflector *openapi31.Reflector) (*openapi31.Parameter, error) {
	if !strings.HasPrefix(ref, "#/components/parameters/") {
		return nil, fmt.Errorf("reference %s is not a parameter", ref)
	}

	parameterName := strings.TrimPrefix(ref, "#/components/parameters/")
	parameterOrReference, ok := reflector.Spec.Components.Parameters[parameterName]
	if !ok {
		return nil, fmt.Errorf("parameter %s not found", parameterName)
	}

	if parameterOrReference.Reference != nil {
		return resolveRefParameter(parameterOrReference.Reference.Ref, reflector)
	}

	return parameterOrReference.Parameter, nil
}

func resolveRefRequestBody(ref string, reflector *openapi31.Reflector) (*openapi31.RequestBody, error) {
	if !strings.HasPrefix(ref, "#/components/requestBodies/") {
		return nil, fmt.Errorf("reference %s is not a requestBody", ref)
	}

	requestBodyName := strings.TrimPrefix(ref, "#/components/requestBodies/")
	requestBodyOrReference, ok := reflector.Spec.Components.RequestBodies[requestBodyName]
	if !ok {
		return nil, fmt.Errorf("requestBody %s not found", requestBodyName)
	}

	if requestBodyOrReference.Reference != nil {
		return resolveRefRequestBody(requestBodyOrReference.Reference.Ref, reflector)
	}

	return requestBodyOrReference.RequestBody, nil
}

func getComponentSchemaTypeName(name string) string {
	return fmt.Sprintf("ComponentSchema%s", capitalize(name))
}
