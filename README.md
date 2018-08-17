This service manages key-value storage

<a id="table-of-contents" href="#table-of-contents">Table of contents</a>

## API Docs

Documentation is in [RAML](https://raml.org) format.

Please see the generated result [here](http://localhost:8090/docs/api.html)

RAML sources are located at `resources/docs` directory.

[[table of contents]](#table-of-contents)

## Logs

Application exposes all logs to stdout

[[table of contents]](#table-of-contents)

## Monitoring

Application automatically publish all available metics in [/metrics](http://localhost:8090/metrics) endpoint

[[table of contents]](#table-of-contents)

## Development

This project requires git, Go and Docker.

As soon as your development environment meets the aforementioned requirements, run the following commands:

1. `make dev-docker-start`
2. `make dev-docker-migration`
3. `build-docs`

## Database Migrations

You can run database migrations as follows:

`make dev-docker-migration`

[[table of contents]](#table-of-contents)

## Test Suites

In this project you can test the code and the automation scripts.

### Unit tests (testing the code)

You can simply run `make dev-docker-test-unit`.

### Integration tests

You can simply run `make dev-docker-test-integration`.

[[table of contents]](#table-of-contents)
