`Consus` - service manages a key-value storage.

_In ancient Roman religion, the god Consus was the protector of grains. He was represented by a grain seed_

<a id="table-of-contents" href="#table-of-contents">Table of contents</a>

- [Description](#description)
	- [Architecture](#architecture)
	- [Limitations / things to improve](#limitations--things-to-improve)
	- [Tested on](#tested-on)
	- [Additional notes](#additional-notes)
- [API Docs](#api-docs)
- [Logs](#logs)
- [Monitoring](#monitoring)
- [Development](#development)
	- [Setup](#setup)
	- [Make commands](#make-commands)
	- [Database Migrations](#database-migrations)
- [Test Suites](#test-suites)
	- [Unit tests (testing the code)](#unit-tests-testing-the-code)
	- [Integration tests](#integration-tests)

## Description
### Architecture
```
features / -- integration tests using gherkin
pkg /
     api / -- api package of application
        docs / -- `/docs` endpoint, which basically returns html generated by raml generator
        keys / -- `/keys` GET|PUT|DELETE|HEAD operations with keys
     command / -- commands for running `http` or `migration` in cli
     config / -- configuration
     model / -- models
     routes / -- routing
     storage / -- database layer handlers
resources /
    dev / -- scripts for developers
    docs / -- raml documentation
    migrations / -- migration files
```
### Limitations / things to improve
1. No pagination
2. No sorting
3. No security
4. No cleanup mechanism to flush expired data
5. Medium unit-test coverage - 74.5%

### Tested on
```
MacOS 10.13.6
Go 1.10.*
Docker(Native) 18.06.0-ce, build 0ffa825
```
### Additional notes
If it would be a real application, I would protect the API using [JWT](http://jwt.io/) which will hold information about roles and etc.
`GET|HEAD` endpoints will be available for users granted with READ access and `DELETE|PUT` only for users with WRITE permissions.

To optimize performance I would consider using some sort of caching strategies (maybe it's not really needed if values will be sorted for a short period).

[[table of contents]](#table-of-contents)

## API Docs

Documentation is in [RAML](https://raml.org) format.

Please see the generated result [here](http://localhost:8090/docs/api.html)

RAML sources are located at `resources/docs` directory.

[[table of contents]](#table-of-contents)

## Logs

The application outputs all logs to stdout

[[table of contents]](#table-of-contents)

## Monitoring

The application automatically publishes all available metrics to [/metrics](http://localhost:8090/metrics) endpoint

[[table of contents]](#table-of-contents)

## Development
### Setup
This project requires git, Go and Docker.

As soon as your development environment meets the aforementioned requirements, run the following commands:

1. `make dev-docker-start`
2. `make dev-docker-migration`
3. `build-docs`

### Make commands

```

 Working with native GO
  run                - Runs targets to simply make it work
  deps               - Ensures dependencies using dep and installs several required tools
  deps-dev           - Install dependencies required only for development
  build              - Buld the application
  build-docs         - Build API documentation from RAML files
  install            - Install app using go install
  test               - Run all tests
  test-unit          - Run only unit tests
  test-unit-coverage - Run unit tests with coverage
  test-integration   - Run integration unit tests
  dev-run-http       - Run REST API
  dev-migrate        - Run migrations

 Working with docker containers
  dev-docker-start              - Start docker containers
  dev-docker-stop               - Stop docker containers
  dev-docker-deps               - Install dependencies using docker container
  dev-docker-migration          - Run migration using docker container
  dev-docker-test-unit          - Run all unit tests using docker container
  dev-docker-test-unit-coverage - Run all unit tests with coverage using docker container
  dev-docker-test-integration   - Run all integration tests using docker container
  dev-docker-logs [<CONTAINER>] - Print container logs

```

### Database Migrations

You can run database migrations as follows:

`make dev-docker-migration`

[[table of contents]](#table-of-contents)

## Test Suites

In this project you can test the code by using unit tests or integration tests.

### Unit tests (testing the code)

`make dev-docker-test-unit`

running unit-tests with coverage:

`make dev-docker-test-unit-coverage`

### Integration tests

`make dev-docker-test-integration`

[[table of contents]](#table-of-contents)
