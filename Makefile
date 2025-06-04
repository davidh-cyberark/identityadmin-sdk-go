# Makefile  -*-Makefile-*-

BINDIR := ./bin

.PHONY: Makefile
include scripts/common.makefile
include scripts/go.makefile
include scripts/openapi.makefile
export

OPENAPI_SPECS_FILES := Authentication-and-Authorization.yaml RoleManagement.yaml
OPENAPI_SPECS := $(addprefix api/,$(OPENAPI_SPECS_FILES))

api/identity-combined.yaml: $(OPENAPI_SPECS) | redocly-cli
	$(REDOCLY_CLI) join $(OPENAPI_SPECS) -o $@

identity/identity-types.gen.go: api/identity-combined.yaml | oapi-codegen-installed
	$(OAPI_CODEGEN) -generate types -package identity $< > $@

identity/identity-client.gen.go: api/identity-combined.yaml | oapi-codegen-installed
	$(OAPI_CODEGEN) -generate client -package identity $< > $@

$(BINDIR)/identity-client: VERSION identity/identity-client.gen.go identity/identity-types.gen.go identity/*.go examples/identity-client/main.go
	$(GO) build -o $@ $(LDFLAGS) examples/identity-client/main.go

.PHONY: identity-client
identity-client: $(BINDIR)/identity-client ## Build the identity-client binary

.PHONY: deps clean vardump
deps: oapi-codegen redocly-cli ## Install dependencies

clean::
	rm -f api/identity-combined.yaml
	rm -f identity/*.gen.go
	rm -f $(BINDIR)/identity-client

vardump::
	@echo "Makefile: OPENAPI_SPECS: $(OPENAPI_SPECS)"
