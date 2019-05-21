all: build push

APP_NAME=monstorak
BIN=operator
REPO=quay.io/monstorak/$(APP_NAME)
TAG=$(shell git rev-parse --short HEAD)
FIRST_GOPATH:=$(firstword $(subst :, ,$(shell go env GOPATH)))
# We need jsonnet on Travis; here we default to the user's installed jsonnet binary; if nothing is installed, then install go-jsonnet.
JSONNET_BIN=$(if $(shell which jsonnet 2>/dev/null),$(shell which jsonnet 2>/dev/null),$(FIRST_GOPATH)/bin/jsonnet)
JB_BIN=$(FIRST_GOPATH)/bin/jb
GOJSONTOYAML_BIN=$(FIRST_GOPATH)/bin/gojsontoyaml
GOBINDATA_BIN=$(FIRST_GOPATH)/bin/go-bindata
BINDATA=pkg/manifests/bindata.go

build:  dependencies mixin $(BINDATA)
	operator-sdk build $(REPO):$(TAG)

push:
	docker push $(REPO):$(TAG)

clean:
	docker images -q $(REPO) | xargs --no-run-if-empty docker rmi --force
	rm -rf jsonnet/vendor

dev: build
	operator-sdk up local --namespace=storage-monitoring

$(BINDATA):
	go-bindata  -mode 420 -modtime 1 -pkg manifests -o $@ jsonnet/manifests/...

mixin:
	cd jsonnet && jb update && \
	./build-jsonnet.sh

dependencies: $(JB_BIN) $(JSONNET_BIN) $(GOBINDATA_BIN) $(GOJSONTOYAML_BIN)

$(GOBINDATA_BIN):
	go get -u github.com/jteeuwen/go-bindata/...

$(GOJSONTOYAML_BIN):
	go get -u github.com/brancz/gojsontoyaml

$(JB_BIN):
	go get -u github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb

$(JSONNET_BIN):
	go get -u github.com/google/go-jsonnet/jsonnet
