# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -v
GOCLEAN=$(GOCMD) clean
GOFMT=gofmt -d -s
GOGET=$(GOCMD) get
BINARY_NAME=bakinbacon

LINUX_BINARY=$(BINARY_NAME)-linux-amd64
MAC_BINARY=$(BINARY_NAME)-darwin-amd64
WINDOWS_BASE=$(BINARY_NAME)-windows-amd64
WINDOWS_BINARY=$(WINDOWS_BASE).exe

GIT_COMMIT := $(shell git rev-list -1 HEAD | cut -c 1-6)

all: build

build: $(LINUX_BINARY)
$(LINUX_BINARY):
	$(GOBUILD) -o $(LINUX_BINARY) -ldflags "-X main.commitHash=$(GIT_COMMIT)"

dist: $(LINUX_BINARY)
	tar -cvzf $(LINUX_BINARY).tar.gz $(LINUX_BINARY)

windows: $(WINDOWS_BINARY)
$(WINDOWS_BINARY):
	PWD=$(shell pwd)
	docker run --rm -it -v golang-windows-cache:/go/pkg -v $(PWD):/go/src/bakinbacon -e GOCACHE=/go/pkg/.cache x1unix/go-mingw /bin/sh -c "cd /go/src/bakinbacon && go build -v -o $(WINDOWS_BINARY) -ldflags '-X main.commitHash=$(GIT_COMMIT)'"

windows-dist: $(WINDOWS_BINARY)
	tar -cvzf $(WINDOWS_BASE).tar.gz $(WINDOWS_BINARY)

fmt: 
	$(GOFMT) baconclient/ nonce/ notifications/ storage/ util/ webserver/ *.go

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

ui:
	npm --prefix webserver/ install
	npm --prefix webserver/ run build
