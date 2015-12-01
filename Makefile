include ../includes.mk

BINARY_DEST_DIR=bin
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags '-s' -o $(BINARY_DEST_DIR)/cde cde.go || exit 1