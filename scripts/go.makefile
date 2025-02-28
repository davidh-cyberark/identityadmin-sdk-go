# openapi.makefile

# Updated: <2025/02/24 21:31:25>

.PHONY: scripts/go.makefile

GO := $(shell command -v go 2> /dev/null)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: go-installed
go-installed:
ifndef GO
	$(error "GO is not available, please install GO")
endif

.PHONY: vardump
vardump::
	@echo "go.makefile: GO binary: $(GO)"
	@echo "go.makefile: GO version: $(shell $(GO) version)"
