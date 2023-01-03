#TARGETS := $(shell ls hack | grep -v \\.sh | grep -v \\.txt)
TARGET=plugin-server
GO=go
GO_MODULE=GO111MODULE=on
BIN_DIR=bin/
ALPINE_COMPILE_FLAGS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
#$(TARGETS):
#	./hack/$@
GOOS=linux GOARCH=amd64
PACKAGE_PREFIX=github.com/Interstellarss/lighthouse-plugin/cmd/


plugin-server:
	$(GO_MODULE) $(ALPINE_COMPILE_FLAGS) $(GO) build -o $(BIN_DIR)$@ $(PACKAGE_PREFIX)$@
