# It's necessary to set this because some environments don't link sh -> bash.

include ./make/verbose.mk
.DEFAULT_GOAL := build

ifeq ("$(GOPATH)","")
$(error "ERROR: GOPATH env variable not set")
endif

FIRST_GOPATH:=$(firstword $(subst :, ,$(shell go env GOPATH)))
export PATH := $(FIRST_GOPATH)/bin:$(PATH)
$(info PATH: $(PATH))
GIT_COMMIT_ID := $(shell git rev-parse HEAD)
BUILD_TIME = `date -u '+%Y-%m-%dT%H:%M:%SZ'`

# setting bash shell with the new path
SHELL := env PATH=$(PATH) /bin/bash

include ./make/find-tools.mk
include ./make/go.mk
include ./make/lint.mk
include ./make/test.mk

$(info Q: $(Q))
$(info GPP: $(GO_PACKAGE_PATH))
$(info GCID: ${GIT_COMMIT_ID})

APP_NAME=monstorak-operator
BIN=operator
REPO=quay.io/monstorak/$(APP_NAME)
TAG=$(shell git rev-parse --short HEAD)
# We need jsonnet on Travis; here we default to the user's installed jsonnet binary; if nothing is installed, then install go-jsonnet.
JSONNET_BIN=$(if $(shell which jsonnet 2>/dev/null),$(shell which jsonnet 2>/dev/null),$(FIRST_GOPATH)/bin/jsonnet)
JB_BIN=$(FIRST_GOPATH)/bin/jb
GOJSONTOYAML_BIN=$(FIRST_GOPATH)/bin/gojsontoyaml
GOBINDATA_BIN=$(FIRST_GOPATH)/bin/go-bindata
BINDATA=pkg/manifests/bindata.go

all: build push

build:  dependencies mixin $(BINDATA) ./vendor
	operator-sdk build $(REPO):$(TAG)

./vendor: Gopkg.toml Gopkg.lock
	$(Q)dep ensure ${V_FLAG} -vendor-only

./out/operator: ./vendor $(shell find . -path ./vendor -prune -o -name '*.go' -print)
	$(Q)CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
                go build ${V_FLAG} \
                -ldflags "-X ${GO_PACKAGE_PATH}/cmd/manager.Commit=${GIT_COMMIT_ID} -X ${GO_PACKAGE_PATH}/cmd/manager.BuildTime=${BUILD_TIME}" \
                -o ./out/operator \
                cmd/manager/main.go

.PHONY: copy-crds
## Copy CRD files to latest OLM manifests directory
copy-crds:
	mkdir -p ./manifests/monstorak-operator
	$(Q)cp ./deploy/crds/*.yaml ./manifests/monstorak-operator/

push:
	docker push $(REPO):$(TAG)

clean:
	docker images -q $(REPO) | xargs --no-run-if-empty docker rmi --force
	rm -rf jsonnet/vendor

dev: build
	# operator-sdk up local --namespace=storage-monitoring

$(BINDATA):
	go-bindata  -mode 420 -modtime 1 -pkg manifests -o $@ jsonnet/manifests/...

mixin:
	cd jsonnet && ${JB_BIN} update && \
	./build-jsonnet.sh

dependencies: $(JB_BIN) $(JSONNET_BIN) $(GOBINDATA_BIN) $(GOJSONTOYAML_BIN)

$(GOBINDATA_BIN):
	go get -u github.com/jteeuwen/go-bindata/...

$(GOJSONTOYAML_BIN):
	go get -u github.com/brancz/gojsontoyaml

$(JB_BIN):
	go get -u github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb

$(JSONNET_BIN):
	go get github.com/google/go-jsonnet/cmd/jsonnet
