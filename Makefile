VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

ORG := github.com
OWNER := kairen
PROJECT_NAME := rotate-elasticsearch-index
REPOPATH ?= $(ORG)/$(OWNER)/$(PROJECT_NAME)

GOOS ?= $(shell go env GOOS)

$(shell mkdir -p ./out)

.PHONY: all build out/rotate-index image-build image-push 

all: build 

.PHONY: build
build: out/rotate-index

.PHONY: out/rotate-index
out/rotate-index:
	GOOS=$(GOOS) CGO_ENABLED=0 go build -ldflags="-s -w" -a -o $@ cmd/main.go

image-build: 
	@docker build -t $(OWNER)/$(PROJECT_NAME):$(VERSION) .

image-push:
	@docker push $(OWNER)/$(PROJECT_NAME):$(VERSION)