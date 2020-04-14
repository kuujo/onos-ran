ARG ONOS_RIC_BASE_VERSION=latest

FROM onosproject/onos-ric-base:$ONOS_RIC_BASE_VERSION as base

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-ric/build/_output/apps/onos-ric-mlb /usr/local/bin/onos-ric-mlb

ENTRYPOINT ["onos-ric-mlb"]
