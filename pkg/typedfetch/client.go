package typedfetch

import (
	"strings"

	"github.com/swaggest/openapi-go/openapi31"
)

func generateClient(reflector *openapi31.Reflector) []string {
	lines := []string{}
	lines = append(lines, strings.TrimSpace(`

type DataResponse<D> = { data: D; response: Response; };
type ErrorResponse<E> = { error: E; response: Response; };
type FetchResponse<D, E> = DataResponse<D> | ErrorResponse<E>;

interface ClientOptions extends RequestInit {
    baseUrl?: string;
    
	// Override fetch function (useful for testing)
	fetch?: (input: Request) => Promise<Response>;
	
	// global body serializer -- allows you to customize how the body is serialized before sending
	// normally not needed unless you are using something like XML instead of JSON
    bodySerializer?: (body: any) => BodyInit | null; 
    
	// global query serializer -- allows you to customize how the query is serialized before sending
	// normally not needed unless you are using some custom array serialization like {foo: [1,2,3,4]} => ?foo=1;2;3;4
	querySerializer?: (query: any) => string;
}

interface Client {
    // GET(url: string, init?: RequestInit): Promise<FetchResponse<Foo, Err>>;
	// POST(url: string, init?: BarRequestInit): Promise<FetchResponse<Bar, Err>>;
}

// Client Implementation

// If not specified, default to application/json
// Used to deduce body serializer and to set content-type header
const contentTypeMap: Record<string, string> = {
	// "/some/binary/url": "application/octet-stream",
}

class ClientImpl {
    options: ClientOptions;

    constructor(options: ClientOptions) {
        options.baseUrl = options.baseUrl || ""; // Make sure baseUrl is always a string
        this.options = options;
    }

    #fetch(input: Request): Promise<Response> {
        return (this.options.fetch || fetch)(input);
    }
}

export function createClient(options: ClientOptions): Client {
    return new ClientImpl(options);
}
`))

	return lines
}
