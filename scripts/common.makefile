# common.makefile

# Updated: <2025/02/26 17:07:34>
.PHONY: scripts/common.makefile

VERSION := $(shell if [ -f VERSION ]; then cat VERSION; else printf "v0.0.1"; fi)
NEXTVERSION := $(shell echo "$(VERSION)" | awk -F. '{print $$1"."$$2"."$$3+1}')

help: ## show help
	@echo "The following build targets have help summaries:"
	@gawk 'BEGIN{FS=":.*[#][#]"} /[#][#]/ && !/^#/ {h[$$1":"]=$$2}END{n=asorti(h,d);for (i=1;i<=n;i++){printf "%-26s%s\n", d[i], h[d[i]]}}' $(MAKEFILE_LIST)
	@echo

versionbump:  ## increment BUILD number in VERSION file
	echo "$(VERSION)" | awk -F. '{print $$1"."$$2"."$$3+1}' > VERSION

minorversionbump:  ## increment MINOR number in VERSION file
	echo "$(VERSION)" | awk -F. '{print $$1"."$$2+1"."0}' > VERSION

majorversionbump:  ## increment MAJOR number in VERSION file
	echo "$(VERSION)" | awk -F. '{print $$1+1"."0"."0}' > VERSION

.PHONY: help versionbump

vardump::  ## echo make variables
	@echo "common.makefile: VERSION: $(VERSION)"
	@echo "common.makefile: NEXTVERSION: $(NEXTVERSION)"

clean:: ## clean ephemeral build resources

realclean:: clean  ## clean all resources that can be re-made (implies clean)

distclean:: realclean  ## clean all resources that can be installed or re-made (implies realclean)

.PHONY: vardump clean realclean
