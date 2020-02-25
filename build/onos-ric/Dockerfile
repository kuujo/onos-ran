ARG ONOS_RIC_BASE_VERSION=latest

FROM onosproject/onos-ric-base:$ONOS_RIC_BASE_VERSION as base

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-ric/build/_output/onos-ric /usr/local/bin/onos-ric

ENTRYPOINT ["onos-ric"]
