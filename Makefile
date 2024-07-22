# Copyright 2017, 2019, 2020 the Velero contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
ENVIRONMENT ?= dev
PLUGIN_NAME ?= radix-velero-plugin
VERSION  ?=latest

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
PKG := github.com/equinor/$(PLUGIN_NAME)
BIN := radix-velero-plugin

REGISTRY ?= radix$(ENVIRONMENT).azurecr.io
IMAGE ?= $(REGISTRY)/radix-velero-plugin

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# local builds the binary using 'go build' in the local environment.
.PHONY: local
local: build-dirs
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	PKG=$(PKG) \
	BIN=$(BIN) \
	OUTPUT_DIR=$$(pwd)/_output/bin/$(GOOS)/$(GOARCH) \
	./hack/build.sh

# test runs unit tests using 'go test' in the local environment.
.PHONY: test
test:
	CGO_ENABLED=0 go test -v -timeout 60s ./...

.PHONY: lint
lint: bootstrap
	golangci-lint run --max-same-issues 0

# ci is a convenience target for CI builds.
.PHONY: ci
ci: verify-modules local

# container builds a Docker image containing the binary.
.PHONY: container
container:
	docker build -t $(IMAGE):${BRANCH}-$(VERSION) .

# push pushes the Docker image to its registry.
.PHONY: push
push:
	@docker push $(IMAGE):${BRANCH}-$(VERSION)

# modules updates Go module files
.PHONY: modules
modules:
	go mod tidy

# verify-modules ensures Go module files are up to date
.PHONY: verify-modules
verify-modules: modules
	@if !(git diff --quiet HEAD -- go.sum go.mod); then \
		echo "go module files are out of date, please commit the changes to go.mod and go.sum"; exit 1; \
	fi

# build-dirs creates the necessary directories for a build in the local environment.
.PHONY: build-dirs
build-dirs:
	@mkdir -p _output/bin/$(GOOS)/$(GOARCH)

# clean removes build artifacts from the local environment.
.PHONY: clean
clean:
	@echo "cleaning"
	rm -rf _output

HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)

bootstrap:
ifndef HAS_GOLANGCI_LINT
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.58.2
endif
