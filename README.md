# caddy-s3-proxy

## Scope

Provide a read-only S3 backend for [Caddy](https://caddyserver.com/) to serve static files from a private S3 bucket.

## Motivation:

- Serving static files through caddy from a private S3 bucket may have different use-cases.
- There is another project [lindenlab/caddy-s3](https://github.com/lindenlab/caddy-s3-proxy) which has sadly not been maintained for a long time and causes [memory-leaks in production](https://github.com/lindenlab/caddy-s3-proxy/issues/64).
- The [file system](https://caddyserver.com/docs/caddyfile/directives/fs) backed plugins lack control over utilising custom s3 metadata like caching or content-types.

## Development

### Prerequisites

- Use [mise-en-place](https://mise.jdx.dev/) or install the required tools mentioned in the [mise configuration](./.mise.toml).
- Have [docker compose](https://docs.docker.com/compose/) up and running to spin up the test dependencies.
