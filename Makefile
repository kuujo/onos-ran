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
	CGO_ENABLED=1 go build -o build/_output/onos-ran ./cmd/onos-ran
	CGO_ENABLED=1 go build -o build/_output/apps/onos-ran-ho ./cmd/apps/onos-ran-ho
	CGO_ENABLED=1 go build -o build/_output/apps/onos-ran-mlb ./cmd/apps/onos-ran-mlb

build-debug: # @HELP build the Go binaries and run all validations (default)
build-debug:
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/onos-ran-debug ./cmd/onos-ran
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/apps/onos-ran-ho-debug ./cmd/apps/onos-ran-ho
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/apps/onos-ran-mlb-debug ./cmd/apps/onos-ran-mlb

build-plugins: # @HELP build plugin binaries
build-plugins: $(MODELPLUGINS)

test: # @HELP run the unit tests and source code validation
test: build deps linters license_check
	CGO_ENABLED=1 go test -race github.com/onosproject/onos-ran/pkg/...
	CGO_ENABLED=1 go test -race github.com/onosproject/onos-ran/cmd/...
	CGO_ENABLED=1 go test -race github.com/onosproject/onos-ran/api/...

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
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-ran \
		-w /go/src/github.com/onosproject/onos-ran \
		--entrypoint build/bin/compile-protos.sh \
		onosproject/protoc-go:stable

onos-ran-base-docker: # @HELP build onos-ran base Docker image
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		--build-arg ONOS_MAKE_TARGET=build \
		-t onosproject/onos-ran-base:${ONOS_RAN_VERSION}
	@rm -rf vendor

onos-ran-docker: onos-ran-base-docker # @HELP build onos-ran Docker image
	docker build . -f build/onos-ran/Dockerfile \
		--build-arg ONOS_RAN_BASE_VERSION=${ONOS_RAN_VERSION} \
		-t onosproject/onos-ran:${ONOS_RAN_VERSION}

onos-ran-ho-docker: onos-ran-base-docker # @HELP build onos-ran-ho Docker image
	docker build . -f build/apps/onos-ran-ho/Dockerfile \
		--build-arg ONOS_RAN_BASE_VERSION=${ONOS_RAN_HO_VERSION} \
		-t onosproject/onos-ran-ho:${ONOS_RAN_HO_VERSION}

onos-ran-mlb-docker: onos-ran-base-docker # @HELP build onos-ran-mlb Docker image
	docker build . -f build/apps/onos-ran-mlb/Dockerfile \
		--build-arg ONOS_RAN_BASE_VERSION=${ONOS_RAN_MLB_VERSION} \
		-t onosproject/onos-ran-mlb:${ONOS_RAN_MLB_VERSION}

onos-ran-tests-docker: # @HELP build onos-ran tests Docker image
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/onos-ran-tests/_output/bin/onos-ran-tests ./cmd/onos-ran-tests
	docker build . -f build/onos-ran-tests/Dockerfile -t onosproject/onos-ran-tests:${ONOS_RAN_VERSION}

onos-ran-benchmarks-docker: # @HELP build onos-ran benchmarks Docker image
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/onos-ran-benchmarks/_output/bin/onos-ran-benchmarks ./cmd/onos-ran-benchmarks
	docker build . -f build/onos-ran-benchmarks/Dockerfile -t onosproject/onos-ran-benchmarks:${ONOS_RAN_VERSION}

images: # @HELP build all Docker images
images: build onos-ran-docker onos-ran-ho-docker onos-ran-mlb-docker onos-ran-tests-docker onos-ran-benchmarks-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-ran:${ONOS_RAN_VERSION}
	kind load docker-image onosproject/onos-ran-ho:${ONOS_RAN_VERSION}
	kind load docker-image onosproject/onos-ran-mlb:${ONOS_RAN_VERSION}
	kind load docker-image onosproject/onos-ran-tests:${ONOS_RAN_VERSION}
	kind load docker-image onosproject/onos-ran-benchmarks:${ONOS_RAN_VERSION}


all: build images

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-ran/onos-ran ./cmd/onos/onos
	go clean -testcache github.com/onosproject/onos-ran/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
