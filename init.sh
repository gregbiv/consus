#!/bin/bash

make deps
make dev-docker-start
make build-docs
make dev-docker-migration
