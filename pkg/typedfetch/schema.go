package typedfetch

import (
	"fmt"
	"strings"
)

func jsonTypeToTypescriptType(schema map[string]any) (string, error) {
	ref, ok := schema["$ref"].(string)
	if ok {
		if !strings.HasPrefix(ref, "#/components/schemas/") {
			return "", fmt.Errorf("unsupported ref, expected #/components/schemas/: %v", ref)
		}

		componentName := getComponentSchemaTypeName(strings.TrimPrefix(ref, "#/components/schemas/"))
		return componentName, nil
	}

	componentType, ok := schema["type"].(string)
	if !ok || !isValidJsonType(componentType) {
		return "", fmt.Errorf("invalid type: %v", componentType)
	}

	switch componentType {
	case "object":
		return jsonObjectToTypescriptType(schema)
	case "array":
		return jsonArrayToTypescriptType(schema)
	case "string":
		return jsonStringToTypescriptType(schema)
	case "number", "integer":
		return "number", nil
	case "boolean":
		return "boolean", nil
	}

	return "", fmt.Errorf("unsupported type: %v", componentType)
}

func jsonObjectToTypescriptType(schema map[string]any) (string, error) {
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		properties = map[string]any{}
	}

	additionalProperties, ok := schema["additionalProperties"]
	hasAdditionalProps := ok

	if len(properties) == 0 && !hasAdditionalProps {
		return "", fmt.Errorf("missing properties or additionalProperties: %v", schema["properties"])
	}

	requiredProps, err := getRequiredProps(schema)
	if err != nil {
		return "", err
	}

	lines := []string{"{"}
	sortedProperties := sortedMapKeys(properties)
	for _, property := range sortedProperties {
		propItem := properties[property]
		optional := "?"
		if itemInSlice(requiredProps, property) {
			optional = ""
		}

		propSchema, ok := propItem.(map[string]any)
		if !ok {
			return "", fmt.Errorf("invalid property schema: %v", propItem)
		}

		propType, err := jsonTypeToTypescriptType(propSchema)
		if err != nil {
			return "", err
		}

		docString := getDocString(propSchema)
		if docString != "" {
			lines = append(lines, fmt.Sprintf("    %s", docString))
		}
		lines = append(lines, fmt.Sprintf("    %s%s: %s;", property, optional, propType))
	}

	// https://swagger.io/docs/specification/data-models/dictionaries/
	if hasAdditionalProps {
		anyAdditionalProps, ok := additionalProperties.(bool)     // additionalProperties: true
		anyAdditionalProps2, ok2 := additionalProperties.(string) // additionalProperties: {} or additionalProperties: ""
		if (ok && anyAdditionalProps) || (ok2 && anyAdditionalProps2 == "") {
			lines = append(lines, "    [key: string]: any;")
		} else {
			additionalPropertiesSchema, ok := additionalProperties.(map[string]any)
			if !ok {
				return "", fmt.Errorf("invalid additionalProperties: %v", schema["additionalProperties"])
			}

			propType, err := jsonTypeToTypescriptType(additionalPropertiesSchema)
			if err != nil {
				return "", err
			}

			lines = append(lines, fmt.Sprintf("    [key: string]: %s;", propType))
		}
	}

	lines = append(lines, "}")
	return strings.Join(lines, "\n"), nil
}

func jsonArrayToTypescriptType(schema map[string]any) (string, error) {
	items, ok := schema["items"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("missing items: %v", schema["items"])
	}

	itemType, err := jsonTypeToTypescriptType(items)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s[]", itemType), nil
}

func jsonStringToTypescriptType(schema map[string]any) (string, error) {
	enum, ok := schema["enum"].([]any)
	if !ok {
		return "string", nil
	}

	enumValues := []string{}
	for _, value := range enum {
		valueString, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("expected enum value to be a string: %v", value)
		}

		enumValues = append(enumValues, fmt.Sprintf("'%s'", valueString))
	}

	return strings.Join(enumValues, " | "), nil
}

func getRequiredProps(schema map[string]any) ([]string, error) {
	requiredProps := []string{}
	if _, ok := schema["required"].([]any); ok {
		for _, prop := range schema["required"].([]any) {
			propString, ok := prop.(string)
			if !ok {
				return nil, fmt.Errorf("expected required property to be a string: %v", prop)
			}

			requiredProps = append(requiredProps, propString)
		}
	}

	return requiredProps, nil
}

func getStringProp(schema map[string]any, propName string) string {
	str, ok := schema[propName].(string)
	if ok {
		return str
	}
	return ""
}

func getDescription(schema map[string]any) string {
	return getStringProp(schema, "description")
}

func getExample(schema map[string]any) string {
	return getStringProp(schema, "example")
}

func buildDocString(description, example string) string {
	if description == "" && example == "" {
		return ""
	}

	if description != "" && example != "" {
		return fmt.Sprintf("/** %s; Example: %s */", description, example)
	} else if example != "" {
		return fmt.Sprintf("/** Example: %s */", example)
	}

	return fmt.Sprintf("/** %s */", description)
}

func getDocString(schema map[string]any) string {
	description := getDescription(schema)
	example := getExample(schema)
	return buildDocString(description, example)
}
