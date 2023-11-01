# Makefile for building and pushing a Go application

# Variables
APP_NAME := truffle
VERSION := 0.0.1

# Targets
.PHONY: build
build: clean
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) cmd/main.go

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f bin/$(APP_NAME)

.PHONY: all
all: build

.PHONY: push
push: build
	@echo "Pushing $(APP_NAME) to a cabral server"
	@scp bin/$(APP_NAME) cabral:
	@echo "Pushing $(APP_NAME) to a rondo server"
	@scp bin/$(APP_NAME) rondo:



# Usage: make release VERSION=x.y.z
.PHONY: release
release: all
	@echo "Tagging version $(VERSION)..."
	@git tag $(VERSION)
	@git push origin $(VERSION)