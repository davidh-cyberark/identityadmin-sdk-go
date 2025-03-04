# npm.makefile

# Updated: <2025/03/03 20:13:35>

.PHONY: scripts/npm.makefile

NPM := $(shell command -v npm 2> /dev/null)

NPMBINSUBDIR := .bin
NPMROOTDIR := $(shell npm root)
NPMBINDIR := $(NPMROOTDIR)/$(NPMBINSUBDIR)
NPMBINRELDIR := ./$(shell basename "$(NPMROOTDIR)")/$(NPMBINSUBDIR)


npm-installed:
ifndef NPM
	$(error "npm is not available, please install npm")
endif

.PHONY: vardump
vardump::
	@echo "npm.makefile: NPM: $(NPM)"
	@echo "npm.makefile: NPMROOTDIR: $(NPMROOTDIR)"
	@echo "npm.makefile: NPMBINDIR: $(NPMBINDIR)"
	@echo "npm.makefile: NPMBINRELDIR: $(NPMBINRELDIR)"