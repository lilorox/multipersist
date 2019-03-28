NAME   = multipersist
GOPATH = $(shell go env GOPATH)
GOLINT = $(GOPATH)/bin/golint

VERBOSE = 0
QUIET   = $(if $(filter 1,$VERBOSE),,@)
BULLET  = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all fmt lint
all: fmt lint ; $(info $(BULLET) building executable…)
	$(QUIET) go build -o $(NAME) *.go

fmt: ; $(info $(BULLET) running gofmt…)
	$(QUIET) @ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		gofmt -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

$(GOLINT): ; $(info $(BULLET) building golint…)
	$(QUIET) go get -u golang.org/x/lint/golint

lint: $(GOLINT); $(info $(BULLET) runnining golint…)
	$(QUIET) ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret
