.PHONY: help run build run-lint run-tests run-example

###############################################################################
#
#	Utils
#
###############################################################################

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


TARGET_MAX_CHAR_NUM=25
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

###############################################################################
#
#	Tests
#
###############################################################################

TEST_PATH ?= ./...

## Run the echo server.
run:
	go run main.go

## Build the echo server in bin/.
build:
	go build -o bin/echo main.go

	.PHONY: help run-lint run-tests run-example

## Run the linter (static code analysis with golangci-lint).
run-lint:
	$(warning This command requires `golangci-lint`)
	golangci-lint run

## Run the unit tests (not done).
run-tests:
	$(warning Not implemented yet)