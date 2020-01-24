ARG ONOS_BUILD_VERSION=stable

FROM onosproject/golang-build:$ONOS_BUILD_VERSION

RUN go get -u github.com/go-delve/delve/cmd/dlv

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /go/src/github.com/onosproject/onos-ran/
COPY api/ /go/src/github.com/onosproject/onos-ran/api/
COPY cmd/ /go/src/github.com/onosproject/onos-ran/cmd/
COPY pkg/ /go/src/github.com/onosproject/onos-ran/pkg/
COPY vendor/ /go/src/github.com/onosproject/onos-ran/vendor/

RUN cd /go/src/github.com/onosproject/onos-ran && GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}
