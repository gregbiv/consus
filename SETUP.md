# Guten Tag, Native Instruments!

Consus - the key-value service.

## Setup instructions

```
# Single script to create docker containers, create schema and generate documentation
./init.sh
```

## Or you can setup everything manually

```
# Install vendor dependecies
make deps
```
```
# Build docker container
make dev-docker-start
```
```
# Create schema (ensure postrgresql container is up and running)
make dev-docker-migration
```
```
# Generate documentation
make build-docs
```
```
# Run tests
make dev-docker-test-unit (ensure consus container is up and running)
```
```
# Run integration tests
make dev-docker-test-integration
```
