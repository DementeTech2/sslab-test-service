include .env

PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables.
GOBASE=$(shell pwd)
GOPATH="$(GOBASE)/vendor:$(GOBASE)"
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

run:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

build:
	rm -rf ./bin
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)


