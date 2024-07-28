# typed-fetch

typed-fetch is intended to be a drop-in replacement for [openapi-typescript/openapi-fetch](https://github.com/openapi-ts/openapi-typescript) but with a much simpler TypeScript implementation that I believe results in fewer issues and enhanced readability.

## Quickstart

```bash
# Generate TypeScript types from OpenAPI document
typed-fetch --openapi examples/petstore-openapi.yaml --output petstore-openapi.ts
```

```ts
// Use the generated library in your .ts files
import { createClient } from "./petstore-openapi";

const client = createClient({ baseUrl: "https://petstore.swagger.io/v2" });

const { data, error } = await client.POST("/store/order", {
    body: {
        id: 10,
        petId: 198772,
        quantity: 7,
        shipDate: "2021-07-07T00:00:00.000Z",
        status: "approved",
        complete: true
    }
});
```

## Installation

You can download pre-built binaries from [Releases](https://github.com/RPGillespie6/typed-fetch/releases).

If you have Go installed, you can download and install the latest version directly with:

```bash
go install github.com/RPGillespie6/typed-fetch@latest
```

Currently this utility is not being published to npm, but it's a future possibility.

# Overview

Type-checked [fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) calls using OpenAPI + TypeScript.

Mostly API-compatible with [openapi-fetch](https://github.com/openapi-ts/openapi-typescript).

Why create this library, if it's nearly identical to [openapi-fetch](https://github.com/openapi-ts/openapi-typescript)? It's because openapi-fetch uses a complex generics/constraints-based TypeScript implementation for type checking, which in my opinion makes it *very* difficult to understand, test, and maintain. The goal of typed-fetch, on the other hand, is to generate simple, straightforward types given an OpenAPI document so that any generalist programmer could inspect the generated TypeScript and easiy understand it -- [see for yourself, no PhD in TypeScript required](examples/petstore-openapi.ts). 

TypeScript generated with this utility is composed almost entirely of type definitions which are stripped out at compile time, resulting in an extremely lightweight fetch wrapper. As a result, typed-fetch should have the same size and performance characteristics as openapi-fetch.

Features:
- Generated TypeScript definitions are at least an order of magnitude simpler and more straightforward than openapi-fetch, which means you can easily jump to and inspect type definitions in your favorite IDE without fear
- Arbitrary combinations of required and optional parameters in request bodies are correctly type-checked (broken in openapi-fetch as of July 2024)
- View your OpenAPI documentation in VSCode when you hover on functions and property names
- Like esbuild, typed-fetch is written in golang, so it's lightning fast

Limitations:
- Only OpenAPI 3.1+ specifications officially supported (see [migration guide](https://www.openapis.org/blog/2021/02/16/migrating-from-openapi-3-0-to-3-1-0)), though most OpenAPI 3.0 documents will also work.
- Some OpenAPI 3 features not currently implemented, especially if it's unclear how it would map to a TypeScript API.

Missing functionality?

Please open an issue with the following 2 things:
- Snippet of valid OpenAPI 3 yaml
- Expected TypeScript to be generated

Also, I'm open to the idea of migrating this utility/approach *into* openapi-fetch if the maintainer(s) there are on-board; I'd rather there be one great tool everyone uses than further fragment the space.