FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/onos-ric-benchmarks/_output/bin/onos-ric-benchmarks /usr/local/bin/onos-ric-benchmarks

ENTRYPOINT ["onos-ric-benchmarks"]
