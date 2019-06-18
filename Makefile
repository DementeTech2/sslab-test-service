
GOBASE=$(shell pwd)
_GOPATH="${GOPATH}:$(GOBASE)"
GOBIN=$(GOBASE)/bin
APPNAME=ssllabtestservice

run:
	GOPATH=$(_GOPATH) GOBIN=$(GOBIN) go run $(APPNAME)

build:
	rm -rf ./bin
	GOPATH=$(_GOPATH) GOBIN=$(GOBIN) go install $(APPNAME)

deps:
	GOPATH=$(_GOPATH) GOBIN=$(GOBIN) go get ./src/...

build_docs:
	apidoc -i src/ -o docs/
