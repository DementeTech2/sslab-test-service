
GOBASE=$(shell pwd)
_GOPATH="${GOPATH}:$(GOBASE)"
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard src/*.go)

run:
	GOPATH=$(_GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

build:
	rm -rf ./bin
	GOPATH=$(_GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

deps:
	GOPATH=$(_GOPATH) GOBIN=$(GOBIN) go get ./src/...


