# typed-fetch

typed-fetch is intended to be a drop-in replacement for [openapi-fetch](https://github.com/openapi-ts/openapi-typescript) but using a simplified TypeScript implementation that I believe results in fewer issues and superior usability/readability/debuggability.

## Quickstart

```bash
# Generate TypeScript types from OpenAPI document
typed-fetch --openapi examples/petstore-openapi.yaml --output petstore-openapi.d.ts
```

```ts
// Use the generated library in your .ts files
import type { Client as PetstoreClient } from "./petstore-openapi"; // petstore-openapi.d.ts
import { createClient } from "./typed-fetch"; // typed-fetch.ts file from root of this repo

const client = createClient<PetstoreClient>({ baseUrl: "https://petstore.swagger.io/v2" });

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

Mostly API-compatible with [openapi-fetch](https://github.com/openapi-ts/openapi-typescript). Doesn't include extra bells and whistles like middleware support since native fetch doesn't support middlewares either.

Why create this library, if it's nearly identical to [openapi-fetch](https://github.com/openapi-ts/openapi-typescript)? It's because openapi-fetch uses a complex generics/constraints-based TypeScript implementation for type checking, which in my opinion makes it *very* difficult to understand, test, and debug because it [requires knowledge of esoteric TypeScript behavior](https://github.com/openapi-ts/openapi-typescript/issues/1778#issuecomment-2276217668). 

The goal of typed-fetch, on the other hand, is to generate simple, straightforward "dumb" types given an OpenAPI document so that any generalist programmer could inspect the generated TypeScript and easily understand it -- [see for yourself, no TypeScript black belt required](examples/petstore-openapi.d.ts). Note there are 5 lines of "evil" TypeScript required to make this library play more nicely with VSCode intellisense, but they could be removed if VSCode had better intellisense for overloaded functions.

Like [openapi-fetch](https://github.com/openapi-ts/openapi-typescript), TypeScript generated with this utility is composed entirely of type definitions which are stripped out at compile time, resulting in an extremely lightweight fetch wrapper. Currently weighs in at **1.4 KiB** minified (**700 bytes** if minified and compressed) which means typed-fetch should meet or exceed the size and performance characteristics of openapi-fetch.

Features:
- Generated TypeScript definitions are *at least* an order of magnitude simpler and more straightforward than openapi-fetch, which means you don't have to be a TypeScript ninja to contribute to or debug issues with the type checking.
- Arbitrary combinations of required and optional parameters in request bodies are correctly type-checked (broken in openapi-fetch as of August 2024 - check if [this issue](https://github.com/openapi-ts/openapi-typescript/issues/1769) is still open)
- Like esbuild, typed-fetch is written in golang, so it's lightning fast

Limitations:
- Only OpenAPI 3.1+ specifications "officially" supported (see [migration guide](https://www.openapis.org/blog/2021/02/16/migrating-from-openapi-3-0-to-3-1-0)), though all OpenAPI 3.0 documents I've tested so far also work.
- Some of the more obscure OpenAPI 3 features are not currently implemented (polymorphism, links, callbacks, etc), and I don't plan to implement them unless there's both a strong use case and a clean way to map them to *both* fetch *and* TypeScript.

# Missing functionality?

Please open an issue with the following 3 things:
- Snippet of valid OpenAPI 3 yaml
- Expected behavior
- Actual behavior

# Note to openapi-fetch maintainers

Feel free to copy any approach and/or techniques you see in this repo. I personally think the approach I implemented here in go could be fairly easily replicated in openapi-typescript, and would result in a more-maintainable and more-contributor-friendly openapi-fetch implementation. I'd be happy to help if you're interested which would obsolete this repo.