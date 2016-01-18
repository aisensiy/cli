include ../includes.mk

BINARY_DEST_DIR=.
build:
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o $(BINARY_DEST_DIR)/cde cde.go || exit 1
	@$(call check-static-binary,$(BINARY_DEST_DIR)/cde)

install: build
	cp $(BINARY_DEST_DIR)/cde $$GOPATH/bin

test: test-unit

test-unit:
	$(GOTEST) . .