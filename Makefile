# Copyright 2014 The Cockroach Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied. See the License for the specific language governing
# permissions and limitations under the License. See the AUTHORS file
# for names of contributors.
#
# Author: Andrew Bonventre (andybons@gmail.com)
# Author: Shawn Morel (shawnmorel@gmail.com)
# Author: Spencer Kimball (spencer.kimball@gmail.com)

# Cockroach build rules.
GO ?= go
# Allow setting of go build flags from the command line.
GOFLAGS := 
# Set to 1 to use static linking for all builds (including tests).
STATIC := $(STATIC)
# The cockroach image to be used for starting Docker containers
# during acceptance tests. Usually cockroachdb/cockroach{,-dev}
# depending on the context.
COCKROACH_IMAGE :=

RUN := run

# TODO(pmattis): Figure out where to clear the CGO_* variables when
# building "release" binaries.
export CGO_CFLAGS :=-g
export CGO_CXXFLAGS :=-g
export CGO_LDFLAGS :=-g

PKG        := "./..."
TESTS      := ".*"
TESTFLAGS  := -logtostderr -timeout 10s
RACEFLAGS  := -logtostderr -timeout 1m
BENCHFLAGS := -logtostderr -timeout 5m

ifeq ($(STATIC),1)
GOFLAGS  += -a -tags netgo -ldflags '-extldflags "-lm -lstdc++ -static"'
endif

.PHONY: all
all: build test

.PHONY: build
build: LDFLAGS += -X github.com/cockroachdb/cockroach/util.buildTag "$(shell git describe --dirty)"
build: LDFLAGS += -X github.com/cockroachdb/cockroach/util.buildTime "$(shell date -u '+%Y/%m/%d %H:%M:%S')"
build: LDFLAGS += -X github.com/cockroachdb/cockroach/util.buildDeps "$(shell GOPATH=${GOPATH} build/depvers.sh)"
build:
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -v -i -o cockroach

# Similar to "testrace", we want to cache the build before running the
# tests.
.PHONY: test
test:
	$(GO) test $(GOFLAGS) -i $(PKG)
	$(GO) test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

# "go test -i" builds dependencies and installs them into GOPATH/pkg, but does not run the
# tests. Run it as a part of "testrace" since race-enabled builds are not covered by
# "make build", and so they would be built from scratch every time (including the
# slow-to-compile cgo packages).
.PHONY: testrace
testrace:
	$(GO) test $(GOFLAGS) -race -i $(PKG)
	$(GO) test $(GOFLAGS) -race -run $(TESTS) $(PKG) $(RACEFLAGS)

.PHONY: bench
bench:
	$(GO) test $(GOFLAGS) -run $(TESTS) -bench $(TESTS) $(PKG) $(BENCHFLAGS)

# Build, but do not run the tests. This is used to verify the deployable
# Docker image which comes without the build environment. See ./build/deploy
# for details.
.PHONY: testbuild
testbuild:
	for p in $(shell $(GO) list $(PKG)); do \
	  $(GO) test $(GOFLAGS) -c -i $$p || exit $?; \
	done


.PHONY: coverage
coverage: build
	$(GO) test $(GOFLAGS) -cover -run $(TESTS) $(PKG) $(TESTFLAGS)

.PHONY: acceptance
acceptance:
# The first `stop` stops and cleans up any containers from previous runs.
	(cd $(RUN) && export COCKROACH_IMAGE="$(COCKROACH_IMAGE)" && \
	  ../build/build-docker-dev.sh && \
	  ./local-cluster.sh stop && \
	  ./local-cluster.sh start && \
	  ./local-cluster.sh stop)

.PHONY: errcheck
errcheck:
	errcheck -ignore='os:Close,net:Close,code.google.com/p/biogo.store/interval:.*,io:Write,bytes:Write.*' $(PKG)

.PHONY: clean
clean:
	$(GO) clean -i github.com/cockroachdb/...
	find . -name '*.test' -type f -exec rm -f {} \;
	rm -rf build/deploy/build
# List all of the dependencies which are not part of the standard
# library or cockroachdb/cockroach.
.PHONY: listdeps
listdeps:
	@go list -f '{{range .Deps}}{{printf "%s\n" .}}{{end}}' ./... | \
	  sort | uniq | egrep '[^/]+\.[^/]+/' | \
	  egrep -v 'github.com/cockroachdb/cockroach'


GITHOOKS := $(subst githooks/,.git/hooks/,$(wildcard githooks/*))
.git/hooks/%: githooks/%
	@echo installing $<
	@rm -f $@
	@mkdir -p $(dir $@)
	@ln -s ../../$(basename $<) $(dir $@)

# Update the git hooks and run the bootstrap script whenever either
# one changes.
.bootstrap: $(GITHOOKS) build/devbase/godeps.sh
	@build/devbase/godeps.sh
	@touch $@

-include .bootstrap
