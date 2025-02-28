# npm.makefile

# Updated: <2025/02/24 21:12:02>

.PHONY: scripts/npm.makefile

NPM := $(shell command -v npm 2> /dev/null)
NPMDIR := $(shell npm root)

npm-installed:
ifndef NPM
	$(error "npm is not available, please install npm")
endif

.PHONY: vardump
vardump::
	@echo "npm.makefile: NPM: $(NPM)"
	@echo "npm.makefile: NPMDIR: $(NPMDIR)"
