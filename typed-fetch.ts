export interface ClientOptions extends RequestInit {
    baseUrl?: string;

    // Override fetch function (useful for testing)
    fetch?: (input: Request) => Promise<Response>;

    // Global headers -- these will be added to every request
    headers?: Record<string, string>;

    // global body serializer -- allows you to customize how the body is serialized before sending
    // normally not needed unless you are using something like XML instead of JSON
    bodySerializer?: (body: any) => BodyInit | null;

    // global query serializer -- allows you to customize how the query is serialized before sending
    // normally not needed unless you are using some custom array serialization like {foo: [1,2,3,4]} => ?foo=1;2;3;4
    querySerializer?: (query: any) => string;
}

export default function createClient<T>(options?: ClientOptions): T {
    return new ClientImpl(options) as T;
}

///////////////////////////////////////////////////////////////////

// Like RequestInit but with some custom fields
type RequestInitExtended = Omit<RequestInit, "body" | "headers"> & {
    params?: object;
    body?: object | BodyInit | null;
    headers?: Record<string, string>;
    parseAs?: "json" | "text" | "blob" | "arrayBuffer" | "formData";

    // local body serializer -- allows you to customize how the body is serialized before sending
    // normally not needed unless you are using something like XML instead of JSON
    bodySerializer?: (body: any) => BodyInit | null;

    // local query serializer -- allows you to customize how the query is serialized before sending
    // normally not needed unless you are using some custom array serialization like {foo: [1,2,3,4]} => ?foo=1;2;3;4
    querySerializer?: (query: any) => string;
};

type TypedFetchResponse = {
    data: any;
    error: any;
    response: Response;
};

type TypedFetchParams = {
    path?: Record<string, any>;
    query?: Record<string, any>;
    headers?: Record<string, string>;
    cookies?: Record<string, string>;
};

function defaultBodySerializer(contentType: string, body: any): BodyInit | null {
    if (contentType.includes("application/json"))
        return JSON.stringify(body);

    if (contentType.includes("application/x-www-form-urlencoded"))
        return new URLSearchParams(body as Record<string, string>);

    if (contentType.includes("multipart/form-data")) {
        const formData = new FormData();
        for (const [key, value] of Object.entries(body as Record<string, any>))
            formData.append(key, value);
        return formData;
    }

    return body;
}

function defaultQuerySerializer(query: Record<string, any>): string {
    return new URLSearchParams(query).toString();
}

function resolveParams(url: string, init: RequestInitExtended, params: TypedFetchParams, querySerializer: (query: Record<string, any>) => string): string {
    if (params["path"]) {
        for (const [key, value] of Object.entries(params["path"]))
            url = url.replace(`{${key}}`, "" + value);
    }

    if (params["query"]) {
        url += "?" + querySerializer(params["query"]);
    }

    if (params["headers"]) {
        init.headers = { ...init.headers, ...params["headers"] };
    }

    if (params["cookies"]) {
        // Add cookies to the "Cookie" header
        const cookies = Object.entries(params["cookies"]).map(([key, value]) => `${key}=${value}`).join("; ");
        init.headers = { ...init.headers, "Cookie": cookies };
    }

    return url;
}

function resolveBody(init: RequestInitExtended, contentType: string, bodySerializer: (contentType: string, body: any) => BodyInit | null) {
    init.body = bodySerializer(contentType, init.body as any);
}

function resolveHeaders(init: RequestInitExtended, globalHeaders: Record<string, string> = {}) {
    let defaultHeaders: Record<string, string> = {}
    if (init.body) {
        let defaultContentType = "application/json";
        if (init.body instanceof Blob || init.body instanceof File || init.body instanceof ArrayBuffer)
            defaultContentType = "application/octet-stream";

        defaultHeaders["Content-Type"] = defaultContentType;
    }

    // default headers have the lowest priority, followed by globalHeaders, and finally init.headers
    return { ...defaultHeaders, ...globalHeaders, ...init.headers };
}

class ClientImpl {
    #options: ClientOptions;
    #fetchFn: (input: Request) => Promise<Response>;

    constructor(options?: ClientOptions) {
        this.#options = {};
        this.#fetchFn = options?.fetch || globalThis.fetch.bind(globalThis);
        this.#options.baseUrl = options?.baseUrl || ""; // Make sure baseUrl is always a string
        this.#options.headers = options?.headers || {};
    }

    async #fetch(method: string, url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> {
        if (!init)
            init = {};

        init.method = method;

        init.headers = resolveHeaders(init, this.#options.headers);

        if (init?.params) {
            const querySerializer = init?.querySerializer || this.#options.querySerializer || defaultQuerySerializer;
            url = resolveParams(url, init, init.params, querySerializer);
        }

        if (init?.body) {
            const bodySerializer = init?.bodySerializer || this.#options.bodySerializer || defaultBodySerializer;
            const contentType = init.headers["Content-Type"] || "";
            resolveBody(init, init.headers["Content-Type"], bodySerializer);
        }


        const requestUrl = this.#options.baseUrl ? new URL(url, this.#options.baseUrl) : url;
        const request = new Request(requestUrl, init as RequestInit);
        const response = await this.#fetchFn(request);

        init.parseAs = init.parseAs || "json";

        // Return {} for "no content" responses to match openapi-fetch truthy behavior
        if (response.headers.get("Content-Length") === "0") {
            return { data: undefined, error: {}, response };
        }

        // Return {} for "no content" responses to match openapi-fetch truthy behavior
        if (response.status === 204) {
            return { data: {}, error: undefined, response };
        }

        if (response.ok) {
            return { data: await response[init.parseAs](), error: undefined, response };
        } else {
            // Mimic openapi-fetch error handling by falling back to text 
            let error = await response.text();
            try {
                error = JSON.parse(error); // attempt to parse as JSON
            } catch {
                // noop
            }

            return { data: undefined, error, response };
        }
    }

    async GET(url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> { return this.#fetch("GET", url, init); }
    async PUT(url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> { return this.#fetch("PUT", url, init); }
    async POST(url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> { return this.#fetch("POST", url, init); }
    async DELETE(url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> { return this.#fetch("DELETE", url, init); }
    async OPTIONS(url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> { return this.#fetch("OPTIONS", url, init); }
    async HEAD(url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> { return this.#fetch("HEAD", url, init); }
    async PATCH(url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> { return this.#fetch("PATCH", url, init); }
    async TRACE(url: string, init?: RequestInitExtended): Promise<TypedFetchResponse> { return this.#fetch("TRACE", url, init); }
}