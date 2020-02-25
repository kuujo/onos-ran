export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_RAN_VERSION := latest
ONOS_RAN_HO_VERSION := latest
ONOS_RAN_MLB_VERSION := latest
ONOS_RAN_DEBUG_VERSION := debug
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/onos-ric ./cmd/onos-ric
	CGO_ENABLED=1 go build -o build/_output/apps/onos-ric-ho ./cmd/apps/onos-ric-ho
	CGO_ENABLED=1 go build -o build/_output/apps/onos-ric-mlb ./cmd/apps/onos-ric-mlb

build-debug: # @HELP build the Go binaries and run all validations (default)
build-debug:
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/onos-ric-debug ./cmd/onos-ric
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/apps/onos-ric-ho-debug ./cmd/apps/onos-ric-ho
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/apps/onos-ric-mlb-debug ./cmd/apps/onos-ric-mlb

build-plugins: # @HELP build plugin binaries
build-plugins: $(MODELPLUGINS)

test: # @HELP run the unit tests and source code validation
test: build deps linters license_check
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
		onosproject/protoc-go:stable

onos-ric-base-docker: # @HELP build onos-ric base Docker image
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		--build-arg ONOS_MAKE_TARGET=build \
		-t onosproject/onos-ric-base:${ONOS_RAN_VERSION}
	@rm -rf vendor

onos-ric-docker: onos-ric-base-docker # @HELP build onos-ric Docker image
	docker build . -f build/onos-ric/Dockerfile \
		--build-arg ONOS_RAN_BASE_VERSION=${ONOS_RAN_VERSION} \
		-t onosproject/onos-ric:${ONOS_RAN_VERSION}

onos-ric-ho-docker: onos-ric-base-docker # @HELP build onos-ric-ho Docker image
	docker build . -f build/apps/onos-ric-ho/Dockerfile \
		--build-arg ONOS_RAN_BASE_VERSION=${ONOS_RAN_HO_VERSION} \
		-t onosproject/onos-ric-ho:${ONOS_RAN_HO_VERSION}

onos-ric-mlb-docker: onos-ric-base-docker # @HELP build onos-ric-mlb Docker image
	docker build . -f build/apps/onos-ric-mlb/Dockerfile \
		--build-arg ONOS_RAN_BASE_VERSION=${ONOS_RAN_MLB_VERSION} \
		-t onosproject/onos-ric-mlb:${ONOS_RAN_MLB_VERSION}

onos-ric-tests-docker: # @HELP build onos-ric tests Docker image
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/onos-ric-tests/_output/bin/onos-ric-tests ./cmd/onos-ric-tests
	docker build . -f build/onos-ric-tests/Dockerfile -t onosproject/onos-ric-tests:${ONOS_RAN_VERSION}

onos-ric-benchmarks-docker: # @HELP build onos-ric benchmarks Docker image
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/onos-ric-benchmarks/_output/bin/onos-ric-benchmarks ./cmd/onos-ric-benchmarks
	docker build . -f build/onos-ric-benchmarks/Dockerfile -t onosproject/onos-ric-benchmarks:${ONOS_RAN_VERSION}

images: # @HELP build all Docker images
images: build onos-ric-docker onos-ric-ho-docker onos-ric-mlb-docker onos-ric-tests-docker onos-ric-benchmarks-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-ric:${ONOS_RAN_VERSION}
	kind load docker-image onosproject/onos-ric-ho:${ONOS_RAN_VERSION}
	kind load docker-image onosproject/onos-ric-mlb:${ONOS_RAN_VERSION}
	kind load docker-image onosproject/onos-ric-tests:${ONOS_RAN_VERSION}
	kind load docker-image onosproject/onos-ric-benchmarks:${ONOS_RAN_VERSION}

all: build images

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
