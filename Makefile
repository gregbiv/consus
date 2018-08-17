### --------------------------------------------------------------------------------------------------------------------
### Variables
### (https://www.gnu.org/software/make/manual/html_node/Using-Variables.html#Using-Variables)
### --------------------------------------------------------------------------------------------------------------------

BUILD_DIR ?= build/

NAME=sandbox
REPO=github.com/gregbiv/${NAME}

# Custom local environment file
ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

SRC_DIRS=pkg
GO_PACKAGES := $(shell go list ./pkg...)

BINARY=sandbox
BINARY_SRC=$(REPO)

GO_LINKER_FLAGS=-ldflags="-s -w"

# Path to the project scripts
SCRIPTS_DIR="resources/dev"

# RAML configuration
RAML_BUILD_DIR ?= "resources/docs"

# Docker enviroment vars
DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2> /dev/null)
DOCKER_CONTAINER=http

# Other config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

### --------------------------------------------------------------------------------------------------------------------
### RULES
### (https://www.gnu.org/software/make/manual/html_node/Rule-Introduction.html#Rule-Introduction)
### --------------------------------------------------------------------------------------------------------------------
.PHONY: all clean deps build install

all: usage

usage:
	@echo "---------------------------------------------"
	@echo "Usage:"
	@echo "  usage - Shows this dialog."
	@echo " "
	@echo " Working with native GO"
	@echo "  run              - Runs targets to simply make it work"
	@echo "  deps             - Ensures dependencies using dep and installs several required tools"
	@echo "  deps-dev         - Install dependencies required only for development"
	@echo "  build            - Buld the application"
	@echo "  build-doc        - Build API documentation from RAML files"
	@echo "  install          - Install app using go install"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run only unit tests"
	@echo "  test-integration - Run integration unit tests"
	@echo "  dev-run-http     - Run REST API"
	@echo "  dev-migrate      - Run migrations"
	@echo " "
	@echo " Working with docker containers"
	@echo "  dev-docker-start              - Start docker containers"
	@echo "  dev-docker-stop               - Stop docker containers"
	@echo "  dev-docker-deps               - Install dependencies using docker container"
	@echo "  dev-docker-migration          - Run migration using docker container"
	@echo "  dev-docker-test-unit          - Run all unit tests using docker container"
	@echo "  dev-docker-test-integration   - Run all integration tests using docker container"
	@echo "  dev-docker-logs [<CONTAINER>] - Print container logs"
	@echo " "
	@exit 0

# Make it work
#-----------------------------------------------------------------------------------------------------------------------
run: clean build install

# Installs our project: copies binaries
#-----------------------------------------------------------------------------------------------------------------------
install: build-assets install-bin

install-bin:
	@printf "$(OK_COLOR)==> Installing project$(NO_COLOR)\n"
	go install -v $(BINARY_SRC)

# Building
#-----------------------------------------------------------------------------------------------------------------------
build: build-assets build-bin
build-docs: build-docs-api

build-bin:
	@printf "$(OK_COLOR)==> Building$(NO_COLOR)\n"
	@go build -o ${BUILD_DIR}/${BINARY} ${GO_LINKER_FLAGS} ${BINARY_SRC}

build-assets:
	@printf "$(OK_COLOR)==> Building assets$(NO_COLOR)\n"

	@echo "-> Cleaning up /resources dir"
	@rm -rf ${BUILD_DIR}/resources
	@mkdir -p ${BUILD_DIR}/resources

	@echo " -> Docs"
	@mkdir -p ${BUILD_DIR}/resources/docs
	@cp -rv ./resources/docs ${BUILD_DIR}/resources

	@echo " -> Copying migration files"
	@cp -rv ./resources/migrations/ ${BUILD_DIR}/resources

build-docs-api:
	@printf "$(OK_COLOR)==> API docs$(NO_COLOR)\n"

	@echo "- Generating"
	@mkdir -p "${RAML_BUILD_DIR}"
	@docker run --rm -w "/data/" -v `pwd`:/data mattjtodd/raml2html:7.0.0 raml2html  -i "resources/docs/api.raml" -o "${RAML_BUILD_DIR}/api.html"

# Dependencies
#-----------------------------------------------------------------------------------------------------------------------
deps:
	@git config --global url."https://${GITHUB_TOKEN}@github.com/hellofresh/".insteadOf "https://github.com/hellofresh/"
	@git config --global http.https://gopkg.in.followRedirects true

	@echo "$(OK_COLOR)==> Installing glide dependencies$(NO_COLOR)"
	@glide install

deps-dev:
	@printf "$(OK_COLOR)==> Installing Godog$(NO_COLOR)\n"
	@go get -u github.com/DATA-DOG/godog/cmd/godog

	@printf "$(OK_COLOR)==> Installing CompileDaemon$(NO_COLOR)\n"
	@go get -u github.com/githubnemo/CompileDaemon

	@printf "$(OK_COLOR)==> Installing Linters$(NO_COLOR)\n"
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/golang/lint

# Testing
#-----------------------------------------------------------------------------------------------------------------------
test: test-unit

test-unit:
	@printf "$(OK_COLOR)==> Running tests$(NO_COLOR)\n"
	@go test -race -cover -covermode=atomic $(shell go list ./... | grep -v /vendor/)

test-integration:
	@printf "$(OK_COLOR)==> Running integration tests$(NO_COLOR)\n"
	@godog --tags="~@notImplemented" --strict

test-unit-coverage:
	@printf "$(OK_COLOR)==> Running tests with code coverage$(NO_COLOR)\n"
	@$(SCRIPTS_DIR)/coverage.sh

# Lint
#-----------------------------------------------------------------------------------------------------------------------
lint:
	@echo "$(OK_COLOR)==> Linting... $(NO_COLOR)"
	@golint $(SRC_DIRS)
	@goimports -l -w $(SRC_DIRS)

# Development
#-----------------------------------------------------------------------------------------------------------------------
dev-run-http:
	@CompileDaemon -build="make install" -graceful-kill -command="$(BINARY) http"

# Dev with docker
dev-docker-start:
	@printf "$(OK_COLOR)==> Starting docker containers$(NO_COLOR)\n"
	@docker-compose up -d --build

dev-docker-stop:
	@printf "$(OK_COLOR)==> Stopping docker containers$(NO_COLOR)\n"
	@docker-compose down

dev-docker-deps:
	@printf "$(OK_COLOR)==> Installing dependencies using docker container$(NO_COLOR)\n"
	@docker-compose exec ${DOCKER_CONTAINER} make deps
	@docker-compose exec ${DOCKER_CONTAINER} make deps-dev

dev-docker-migration:
	@printf "$(OK_COLOR)==> Running migration using docker container$(NO_COLOR)\n"
	@docker-compose exec ${DOCKER_CONTAINER} $(BINARY) migrate

dev-docker-test-unit:
	@printf "$(OK_COLOR)==> Running unit test using docker container$(NO_COLOR)\n"
	@printf "Don't forget before run unit test to update deps\n"
	@docker-compose exec ${DOCKER_CONTAINER} make test-unit

dev-docker-test-integration:
	@printf "$(OK_COLOR)==> Running integration test using docker container$(NO_COLOR)\n"
	@printf "Don't forget before run integration test to update deps and run migration if need it\n"
	@docker-compose exec ${DOCKER_CONTAINER} make test-integration

dev-docker-logs:
	@docker-compose logs -f ${CONTAINER}

# House keeping
#-----------------------------------------------------------------------------------------------------------------------
clean:
	@printf "$(OK_COLOR)==> Cleaning project$(NO_COLOR)\n"
	if [ -d ${BUILD_DIR} ] ; then rm -rf ${BUILD_DIR} ; fi
