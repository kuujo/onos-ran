export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_RIC_VERSION := latest
ONOS_RIC_HO_VERSION := latest
ONOS_RIC_MLB_VERSION := latest
ONOS_BUILD_VERSION := v0.6.4
ONOS_PROTOC_VERSION := v0.6.4

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/onos-ric ./cmd/onos-ric
	CGO_ENABLED=1 go build -o build/_output/apps/onos-ric-ho ./cmd/apps/onos-ric-ho
	CGO_ENABLED=1 go build -o build/_output/apps/onos-ric-mlb ./cmd/apps/onos-ric-mlb

generate: # @HELP generate store interfaces and implementations
generate:
	go run github.com/onosproject/onos-ric/cmd/onos-ric-generate ./build/generate/stores.yaml
	mockgen github.com/onosproject/onos-ric/pkg/southbound E2 > test/mocks/southbound/e2_mock.go

test: # @HELP run the unit tests and source code validation
test: build deps linters
	CGO_ENABLED=1 go test -race github.com/onosproject/onos-ric/pkg/...
	CGO_ENABLED=1 go test -race github.com/onosproject/onos-ric/cmd/...
	CGO_ENABLED=1 go test -race github.com/onosproject/onos-ric/api/...

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./build/bin/coveralls-coverage

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run --timeout 30m

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

protos: # @HELP compile the protobuf files (using protoc-go Docker)
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-ric \
		-w /go/src/github.com/onosproject/onos-ric \
		--entrypoint build/bin/compile-protos.sh \
		onosproject/protoc-go:${ONOS_PROTOC_VERSION}

onos-ric-base-docker: # @HELP build onos-ric base Docker image
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		--build-arg ONOS_MAKE_TARGET=build \
		-t onosproject/onos-ric-base:${ONOS_RIC_VERSION}
	@rm -rf vendor

onos-ric-docker: # @HELP build onos-ric Docker image
onos-ric-docker: onos-ric-base-docker
	docker build . -f build/onos-ric/Dockerfile \
		--build-arg ONOS_RIC_BASE_VERSION=${ONOS_RIC_VERSION} \
		-t onosproject/onos-ric:${ONOS_RIC_VERSION}

onos-ric-ho-docker: # @HELP build onos-ric-ho Docker image
onos-ric-ho-docker: onos-ric-base-docker
	docker build . -f build/apps/onos-ric-ho/Dockerfile \
		--build-arg ONOS_RIC_BASE_VERSION=${ONOS_RIC_HO_VERSION} \
		-t onosproject/onos-ric-ho:${ONOS_RIC_HO_VERSION}

onos-ric-mlb-docker: # @HELP build onos-ric-mlb Docker image
onos-ric-mlb-docker: onos-ric-base-docker
	docker build . -f build/apps/onos-ric-mlb/Dockerfile \
		--build-arg ONOS_RIC_BASE_VERSION=${ONOS_RIC_MLB_VERSION} \
		-t onosproject/onos-ric-mlb:${ONOS_RIC_MLB_VERSION}

images: # @HELP build all Docker images
images: build onos-ric-docker onos-ric-ho-docker onos-ric-mlb-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-ric:${ONOS_RIC_VERSION}
	kind load docker-image onosproject/onos-ric-ho:${ONOS_RIC_VERSION}
	kind load docker-image onosproject/onos-ric-mlb:${ONOS_RIC_VERSION}

all: build images

publish: # @HELP publish version on github and dockerhub
	./../build-tools/publish-version ${VERSION} onosproject/onos-ric onosproject/onos-ric-ho onosproject/onos-ric-mlb

bumponosdeps: # @HELP update "onosproject" go dependencies and push patch to git. Add a version to dependency to make it different to $VERSION
	./../build-tools/bump-onos-deps ${VERSION}

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-ric/onos-ric ./cmd/onos/onos
	go clean -testcache github.com/onosproject/onos-ric/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
