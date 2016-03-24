include ../includes.mk

INSTALL_DIR=$$GOPATH/bin

build: install

install:
	go install || exit 1
	@$(call check-static-binary,$(INSTALL_DIR)/cde)

test: test-unit

test-unit:
	$(GOTEST) . .
