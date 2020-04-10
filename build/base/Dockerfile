ARG ONOS_BUILD_VERSION=v0.5.0

FROM onosproject/golang-build:$ONOS_BUILD_VERSION

RUN go get -u github.com/go-delve/delve/cmd/dlv

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /go/src/github.com/onosproject/onos-ric/
COPY api/ /go/src/github.com/onosproject/onos-ric/api/
COPY cmd/ /go/src/github.com/onosproject/onos-ric/cmd/
COPY pkg/ /go/src/github.com/onosproject/onos-ric/pkg/
COPY vendor/ /go/src/github.com/onosproject/onos-ric/vendor/

RUN cd /go/src/github.com/onosproject/onos-ric && GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}
