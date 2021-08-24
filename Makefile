# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -v
GOCLEAN=$(GOCMD) clean
GOFMT=gofmt -d -s
GOGET=$(GOCMD) get
BINARY_NAME=bakinbacon

GIT_COMMIT := $(shell git rev-list -1 HEAD | cut -c 1-6)

all: build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -ldflags "-X main.commitHash=$(GIT_COMMIT)"

windows:
	PWD=$(shell pwd)
	docker run --rm -it -v golang-windows-cache:/go/pkg -v $(PWD):/go/src/bakinbacon x1unix/go-mingw /bin/sh -c "cd /go/src/bakinbacon && go build -v -o $(BINARY_NAME)-win-x64.exe -ldflags '-X main.commitHash=$(GIT_COMMIT)'"

fmt: 
	$(GOFMT) baconclient/ nonce/ notifications/ storage/ util/ webserver/ *.go

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

ui:
	npm --prefix webserver/ install
	npm --prefix webserver/ run build
