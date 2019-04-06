SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET = $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL := $(TARGET)

VERSION = 0.1.0
BUILD = $(shell git rev-parse HEAD)

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

GOPATH = $(shell go env GOPATH)
GOLINT = $(GOPATH)/bin/golint
LDFLAGS = -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

VERBOSE = 0
QUIET   = $(if $(filter 1,$VERBOSE),,@)
BULLET  = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all bench build check clean fmt lint install vet

all: check install

$(TARGET): $(SRC)
	$(info $(BULLET) building executable $(TARGET)…)
	$(QUIET) go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

bench:
	@go test -bench=.

clean:
	$(info $(BULLET) cleaning…)
	$(QUIET) rm -f $(TARGET)

install:
	$(info $(BULLET) installing executable $(TARGET)…)
	$(QUIET) go install $(LDFLAGS)

check: fmt vet lint

fmt:
	$(info $(BULLET) running gofmt…)
	$(QUIET) @ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		gofmt -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

$(GOLINT):
	$(info $(BULLET) installing golint…)
	$(QUIET) go get -u golang.org/x/lint/golint

lint: $(GOLINT)
	$(info $(BULLET) running golint…)
	$(QUIET) ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret

vet:
	$(info $(BULLET) running vet…)
	$(QUIET) go vet ${SRC}
