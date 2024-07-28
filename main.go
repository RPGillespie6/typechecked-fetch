package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/RPGillespie6/typed-fetch/pkg/typedfetch"
	"github.com/swaggest/openapi-go/openapi31"
)

func main() {
	openApiSpecPath := flag.String("openapi", "", "Input file path (.json or .yaml)")
	outputPath := flag.String("output", "", "Output file path")
	flag.Parse()

	if *openApiSpecPath == "" {
		panic("openapi document path is required")
	}

	reflector, err := openApi31ReflectorFromFile(*openApiSpecPath)
	if err != nil {
		panic(err)
	}

	generatedOutput, err := typedfetch.GenerateTypedFetch(reflector)
	if err != nil {
		panic(err)
	}

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
