package typedfetch

import (
	"fmt"
	"strings"

	"github.com/swaggest/openapi-go/openapi31"
)

func generateClient(reflector *openapi31.Reflector) ([]string, error) {
	clientInterfaceLookups := map[string][]string{}

	sortedPaths := sortedMapKeys(reflector.Spec.Paths.MapOfPathItemValues)
	for _, path := range sortedPaths {
		item := reflector.Spec.Paths.MapOfPathItemValues[path]
		methods := getPathItemMethods(&item)
		for _, method := range methods {
			if method.Operation == nil {
				continue
			}

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

			// Generate the client interface
			initRequired := paramInfo.Required || bodyInfo.Required
			requestTypeName := getRequestTypeName(method.Method, path)
			responseDataTypeName := getResponseDataTypeName(method.Method, path)
			responseErrTypeName := getResponseErrTypeName(method.Method, path)

			if _, ok := clientInterfaceLookups[method.Method]; !ok {
				clientInterfaceLookups[method.Method] = []string{}
			}

			initOptionalQ := "?"
			if initRequired {
				initOptionalQ = ""
			}

			lookupLine := fmt.Sprintf("\"%s\": {init%s: %s, response: FetchResponse<%s, %s>}",
				path,
				initOptionalQ,
				requestTypeName,
				responseDataTypeName,
				responseErrTypeName,
			)

			clientInterfaceLookups[method.Method] = append(clientInterfaceLookups[method.Method], lookupLine)
		}
	}

	// Generate the client interface
	typeLookupLines := []string{}
	clientInterfaceLines := []string{}
	for method, lines := range clientInterfaceLookups {
		typeLookupTypeName := getLookupTypeName(method)
		typeLookupLines = append(typeLookupLines, fmt.Sprintf("type %s = {\n    %s\n};", typeLookupTypeName, strings.Join(lines, ",\n    ")))
		typeLookupLines = append(typeLookupLines, "")

		clientInterfaceLines = append(clientInterfaceLines,
			fmt.Sprintf(`%s: ClientMethod<%s>;`,
				method,
				typeLookupTypeName,
			))
	}

	lines := []string{
		strings.TrimSpace(fmt.Sprintf(`
// Response Generics

type DataResponse<D> = { data: D; error: undefined; response: Response; };
type ErrorResponse<E> = { data: undefined; error: E; response: Response; };
type FetchResponse<D, E> = DataResponse<D> | ErrorResponse<E>;

// Generics Type Lookups
// These are lookup tables for each method type (GET, POST, etc) to match the url to its payload

%s

/* 
    We could just generate a bunch of overloaded functions in the interface like:

    GET(url: "/pet", init?: FetchRequestGetPet): Promise<FetchResponse<GetPetResponse, GetPetError>>;
    GET(url: "/store/order/{orderId}", init: FetchRequestGetPet2): Promise<FetchResponse<GetPet2Response, GetPet2Error>>;
    GET(url: "/user/login", init?: FetchRequestGetPet3): Promise<FetchResponse<GetPet3Response, GetPet3Error>>;
    etc
    
    This approach is valid and correctly type checks request and response shape.
    However, for some reason, VSCode intellisense doesn't work well with overloads; 
    it gets confused as to which overload you want and won't list out request body properties, etc.

    So instead, we just use 1 generic function per method type using a lookup table to match the url to its payload.
    This seems to work better with VSCode intellisense...

    ...However, it requires some evil TypeScript magic to make a single generic function behave like the group of overloads above

    With generics, it's much trickier to allow both init and init? in the same function signature
    These two helpers below are the only solution I was able to find to achieve that, at the cost of readability...

    TODO: Revisit this in the future to see if overloads work better with intellisense, because it's a much simpler solution
    and doesn't require the evil TypeScript magic
*/

// https://stackoverflow.com/questions/52984808/is-there-a-way-to-get-all-required-properties-of-a-typescript-object
// Example: OptionalKeys<{a: string, b?: number}> = "b"
type OptionalKeys<T extends object> = keyof { [K in keyof T as {} extends Pick<T, K> ? K : never]: any }

// https://stackoverflow.com/questions/77714794/how-to-use-void-to-make-function-parameters-optional-when-using-generics
// Basically, this is a way to achieve i.e. GET(url, init?) if init is optional, and GET(url, init) if init is required
type ClientMethod<Lookup extends Record<string, any>> = <Url extends keyof Lookup>(
    url: Url,
    ...[init]: "init" extends OptionalKeys<Lookup[Url]> ? [init?: Lookup[Url]["init"]] : [init: Lookup[Url]["init"]]
) => Lookup[Url]["response"];

// end evil TypeScript magic

export interface Client {
    %s
}
`, strings.Join(typeLookupLines, "\n"), strings.Join(clientInterfaceLines, "\n    "))),
	}

	return lines, nil
}

func getLookupTypeName(method string) string {
	return fmt.Sprintf("%sTypesLookup", pascalize(method))
}
