SHELL = /bin/bash

.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	yarn --cwd front build
	statik -f -src front/build
	go build -o cra-example
