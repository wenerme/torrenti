### BEGIN COMMON
COLOR 	:= "\e[1;36m%s\e[0m\n"
RED 	:=   "\e[1;31m%s\e[0m\n"

ifdef GOROOT
PATH 	:= $(GOROOT)/bin:$(PATH)
endif

GOBIN 	:= $(if $(shell go env GOBIN),$(shell go env GOBIN),$(GOPATH)/bin)
PATH 	:= $(GOBIN):$(PATH)

GOOS 	?= $(shell go env GOOS)
GOARCH 	?= $(shell go env GOARCH)
GOPATH 	?= $(shell go env GOPATH)

MODDIR	:= $(shell dirname $(shell go env GOMOD))
MODNAME := $(shell basename $(MODDIR))

DOCKER_REPO := $(or $(DOCKER_REPO),cr.incos.cloud/sparkle)
DOCKER_TAG 	:= $(or $(DOCKER_TAG),$(shell git rev-parse --abbrev-ref HEAD))
BUILD_TAG	:= $(or $(BUILD_TAG),$(shell git rev-parse --abbrev-ref HEAD))

GOFLAGS	:= -v -trimpath -ldflags "-s -w"

-include local.Makefile
-include ignored.Makefile

info:
	@echo "GOOS=$(GOOS)"
	@echo "GOARCH=$(GOARCH)"
	@echo "GOPROXY=`go env GOPROXY`"
	@echo "DOCKER_REPO=$(DOCKER_REPO)"
	@echo "DOCKER_TAG=$(DOCKER_TAG)"

##### Bazel #####

git-add-bazel:
	git add --ignore-error '**/BUILD.bazel' WORKSPACE deps.bzl BUILD.bazel

##### Container #####

image: GOOS:=linux
image: build # build image
	@ls cmd | xargs -n1 -I {} sh -c 'set -x;echo Building {}; docker build -t {}:$(BUILD_TAG) build/{}'

image-ls: image
	@[ ! -e BUILD.bazel ] || ls bazel-bin/build/*-layer.tar | xargs -n1 -I {} sh -c 'echo {};tar -tvf {}'

push: ## push image to registry
	@ls cmd | xargs -n1 -I {} sh -c 'set -x;echo Pushing {}; docker tag {}:$(BUILD_TAG) $(DOCKER_REPO)/{}:$(DOCKER_TAG); docker push $(DOCKER_REPO)/{}:$(DOCKER_TAG)'

##### Bazel Go #####

gazelle: ## bazel gazelle
	@[ ! -e BUILD.bazel ] || bazel run --noshow_progress --noshow_loading_progress //:gazelle
	git add --ignore-error '**/BUILD.bazel'

gazelle-update-repos: tidy  ## bazel gazelle update-repos
	@[ ! -e BUILD.bazel ] || bazel run --noshow_progress --noshow_loading_progress //:gazelle-update-repos

update: gazelle-update-repos gazelle

.PHONY: build
build: ## build binary
	@ls cmd | xargs -n1 -I {} sh -c 'set -x;echo Building {}; GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -o build/{}/bin/{} ./cmd/{}'

.PHONY: build
test: ## bazel test
	@[ ! -e BUILD.bazel ] || bazel test $(shell bazel query 'kind(go_test, //...)')

##### Golang #####

.PHONY: lint
lint: ## lint
	golangci-lint run

.PHONY: fmt
fmt: tidy ## tidy,format and imports
	gofumpt -w `find . -type f -name '*.go' -not -path "./vendor/*"`
	goimports -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: tidy
tidy: ## go mod tidy
	go mod tidy

.PHONY: gen
gen: ## generate
	[ -e buf.gen.yaml ] && buf generate || true
	$(MAKE) fmt

.PHONY: go-test-cover
go-test-cover: ## run test & generate coverage
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

update-dependencies: ## update go dependencies
	@printf $(COLOR) "Update dependencies..."
	@go get -u -t $(PINNED_DEPENDENCIES) ./...
	@go mod tidy
	@$(MAKE) gazelle-update-repos

ensure-no-changes: ## ensure git doesn't have any changes
	@printf $(COLOR) "Check for local changes..."
	@printf $(COLOR) "========================================================================"
	@git diff --name-status --exit-code || (printf $(COLOR) "========================================================================"; printf $(RED) "Above files are not regenerated properly. Regenerate them and try again."; exit 1)

clean: ## cleanup build
	rm -rf build/*/bin/*
	rm -rf bazel-bin/*
	[ ! -e BUILD.bazel ] || bazel clean

.PHONY: help
.DEFAULT_GOAL := help
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
### END COMMON
