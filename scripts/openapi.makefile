# openapi.makefile

# Updated: <2025/02/26 17:06:25>

.PHONY: scripts/openapi.makefile 
include scripts/go.makefile
include scripts/npm.makefile

#####
# Redocly - for linting and generating docs
REDOCLY_CLI := $(NPM) exec -- @redocly/cli

.PHONY: redocly-cli
redocly-cli: | npm-installed  ## install redocly-cli - for linting and generating docs
	$(NPM) list @redocly/cli@latest >/dev/null || $(NPM) install @redocly/cli@latest

redocly-cli-uninstall: | npm-installed  ## uninstall redocly-cli
	$(NPM) uninstall @redocly/cli@latest

#####
# Code generation
OAPI_CODEGEN := $(shell command -v ./bin/oapi-codegen 2> /dev/null)

.PHONY: oapi-codegen-installed
oapi-codegen-installed: ## check if oapi-codegen tool is installed
ifndef OAPI_CODEGEN
	$(error "OAPI_CODEGEN is not installed; try 'make oapi-codegen'")
endif

.PHONY: oapi-codegen
oapi-codegen: bin/oapi-codegen ## install oapi-codegen tool

bin/oapi-codegen: | go-installed
	GOBIN=$(shell pwd)/bin $(GO) install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

.PHONY: vardump distclean
distclean::
	rm -f bin/oapi-codegen
	$(NPM) uninstall @redocly/cli@latest

vardump::
	@echo "openapi.makefile: REDOCLY_CLI: $(REDOCLY_CLI)"
	@echo "openapi.makefile: OAPI_CODEGEN: $(OAPI_CODEGEN)"
