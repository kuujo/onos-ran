ARG ONOS_RAN_BASE_VERSION=latest

FROM onosproject/onos-ran-base:$ONOS_RAN_BASE_VERSION as base

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-ran/build/_output/onos-ran /usr/local/bin/onos-ran
COPY --from=base /go/src/github.com/onosproject/onos-ran/build/_output/onos-ran-ho /usr/local/bin/onos-ran-ho
COPY --from=base /go/src/github.com/onosproject/onos-ran/build/_output/onos-ran-mlb /usr/local/bin/onos-ran-mlb

ENTRYPOINT ["onos-ran"]
