
.PHONY : default generate test-unit test-integration coverage coverage-report lint build run clean generate-all-rpcs

SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)
ROOT_DIR ?= $(shell git rev-parse --show-toplevel)

# constants
GOLANGCI_VERSION = 1.64.8
DOCKER_REPO = lbsim25
DOCKER_TAG = v1

download: ## Downloads the dependencies
	@go mod download

tidy: ## Cleans up go.mod and go.sum
	@go mod tidy

GO_BUILD = mkdir -pv "$(@)" && go build -ldflags="-w -s" -o "$(@)" ./...

buildcontainerapp: ## Generates a linux binary of the app for the Dockerfile
	GOOS=linux GOARCH=amd64 go build -o app ./cmd/server

docker: ## Builds docker image
	docker buildx build -t $(DOCKER_REPO):$(DOCKER_TAG) .
# Go dependencies versioned through tools.go
GO_DEPENDENCIES = google.golang.org/protobuf/cmd/protoc-gen-go \
				google.golang.org/grpc/cmd/protoc-gen-go-grpc \
				github.com/envoyproxy/protoc-gen-validate \
				github.com/bufbuild/buf/cmd/buf \
                github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking \
                github.com/bufbuild/buf/cmd/protoc-gen-buf-lint \
                github.com/vektra/mockery/v2 \
                github.com/alta/protopatch/cmd/protoc-gen-go-patch
# additional dependencies for grpc-gateway
GO_DEPENDENCIES += github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
				github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

define make-go-dependency
  # target template for go tools, can be referenced e.g. via /bin/<tool>
  bin/$(notdir $(subst /v2,,$1)): $(ROOT_DIR)/go.mod
	GOBIN=$(PWD)/bin go install $1
endef

# this creates a target for each go dependency to be referenced in other targets
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''
