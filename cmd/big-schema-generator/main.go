package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/swaggest/openapi-go/openapi31"
)

func main() {
	openApiSpecPath := flag.String("openapi", "", "Input file path (.json or .yaml)")
	outputPath := flag.String("output", "", "Output file path")
	schemaMultiplier := flag.Int("multiplier", 100, "Number of times to multiply the schema")
	flag.Parse()

	if *openApiSpecPath == "" {
		panic("openapi document path is required")
	}

	reflector, err := openApi31ReflectorFromFile(*openApiSpecPath)
	if err != nil {
		panic(err)
	}

	err = GenerateBigSchema(reflector, *schemaMultiplier)
	if err != nil {
		panic(err)
	}

	schema, err := reflector.Spec.MarshalYAML()
	if err != nil {
		panic(err)
	}

	generatedOutput := string(schema)

	// Generate typed fetch
	if *outputPath != "" {
		err = os.WriteFile(*outputPath, []byte(generatedOutput), 0644)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println(generatedOutput)
	}
}

func openApi31ReflectorFromFile(path string) (*openapi31.Reflector, error) {
	specBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reflector := openapi31.NewReflector()
	if strings.HasSuffix(path, ".json") {
		err = reflector.Spec.UnmarshalJSON(specBytes)
		if err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		err = reflector.Spec.UnmarshalYAML(specBytes)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Unsupported file format: %s", path)
	}

	return reflector, nil
}

func GenerateBigSchema(reflector *openapi31.Reflector, schemaMultiplier int) error {
	// Generate a huge openapi schema given an input schema
	newPathMap := map[string]openapi31.PathItem{}
	for i := 0; i < schemaMultiplier; i++ {
		for path, pathItem := range reflector.Spec.Paths.MapOfPathItemValues {
			tokens := strings.Split(path, "/")
			tokens = append(tokens, fmt.Sprintf("%d", i))
			newPath := strings.Join(tokens, "/")
			newPathMap[newPath] = pathItem
		}
	}

	reflector.Spec.Paths.MapOfPathItemValues = newPathMap
	return nil
}
