# caddy-s3-proxy

Project Motivation:
- Serving static files through caddy from a private S3 bucket may have different use-cases.
- There is another project [lindenlab/caddy-s3](https://github.com/lindenlab/caddy-s3-proxy) which has sadly not been maintained for a long time.
- The [file system](https://caddyserver.com/docs/caddyfile/directives/fs) backed plugins lack control over utilising custom s3 metadata like caching or content-types.

## Development

### Prerequisites

- Use [mise-en-place](https://mise.jdx.dev/) or install the required tools mentioned in the [mise configuration](./.mise.toml).
- Have [docker compose](https://docs.docker.com/compose/) up and running to spin up the test dependencies.
