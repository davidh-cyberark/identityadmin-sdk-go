# openapi.makefile

# Updated: <2025/06/04 17:38:20>

.PHONY: scripts/openapi.makefile 
include scripts/go.makefile
include scripts/npm.makefile

#####
# Redocly - for linting and generating docs
# REDOCLY_CLI := $(NPM) exec -- @redocly/cli
REDOCLY_CLI := $(NPMBINRELDIR)/redocly

.PHONY: redocly-cli
redocly-cli: $(REDOCLY_CLI)  ## install redocly-cli - for linting and generating docs

$(REDOCLY_CLI): | npm-installed  ## check if redocly-cli is installed
	$(NPM) list @redocly/cli >/dev/null || $(NPM) install @redocly/cli

redocly-cli-uninstall: | npm-installed  ## uninstall redocly-cli
	$(NPM) uninstall @redocly/cli

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
	GOBIN=$(shell pwd)/bin $(GO) install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

.PHONY: vardump distclean
distclean::
	rm -f bin/oapi-codegen
	$(NPM) uninstall @redocly/cli@latest

vardump::
	@echo "openapi.makefile: REDOCLY_CLI: $(REDOCLY_CLI)"
	@echo "openapi.makefile: REDOCLY_CLI version: $(shell $(REDOCLY_CLI) --version)"
	@echo "openapi.makefile: OAPI_CODEGEN: $(OAPI_CODEGEN)"
