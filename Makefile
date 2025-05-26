# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=motion-snapshot-server
BINARY_PATH=bin/$(BINARY_NAME)

all: test build

build:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -o $(BINARY_PATH) -v ./src/...

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

run:
	$(GOBUILD) -o $(BINARY_PATH) -v ./src/...
	./$(BINARY_PATH)

deps:
	$(GOGET) -v ./...

.PHONY: all build test clean run deps