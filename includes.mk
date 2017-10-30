SHELL=/bin/bash
GO=go
GOFMT=gofmt -l
GOLINT=golint
GOTEST=$(GO) test --cover --race -v
GOVET=$(GO) vet
GIT_SHA=$(shell git rev-parse --short=8 HEAD)
VERSION=$(shell if [ "$$CDE_VERSION" == "" ]; then echo "git-$$(git rev-parse --short=8 HEAD)"; else echo "$$CDE_VERSION"; fi) 
ifndef BUILD_TAG
  BUILD_TAG = git-$(GIT_SHA)
endif

define check-static-binary
  if file $(1) | egrep -q "(statically linked|Mach-O)"; then \
	echo -n ""; \
  else \
	echo -e '\033[0;31m'"The binary file $(1) is not statically linked. Build canceled" '\033[0m' ; \
	exit 1; \
  fi
endef

setup-gotools:
	@go get github.com/tools/godep
