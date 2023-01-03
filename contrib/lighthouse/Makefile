TARGET=lighthouse
GO=go
GO_MODULE=GO111MODULE=on
BIN_DIR=bin/
ALPINE_COMPILE_FLAGS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
#$(TARGETS):
#	./hack/$@
GOOS=linux GOARCH=amd64
PACKAGE_PREFIX=github.com/Interstellarss/lighthouse/cmd/


lighthouse:
	$(GO_MODULE) $(ALPINE_COMPILE_FLAGS) $(GO) build -o $(BIN_DIR)$@ $(PACKAGE_PREFIX)$@