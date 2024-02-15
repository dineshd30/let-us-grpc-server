# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean

# Build target
.PHONY: build
build:
	$(GOBUILD) -o ./bin/server ./cmd/server

# Run target
.PHONY: run
run: clean build
	./bin/server

# Test target
.PHONY: test
test:
	$(GOTEST) -v ./...

# Clean target
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf ./bin

