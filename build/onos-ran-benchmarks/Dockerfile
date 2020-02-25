FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/onos-ran-benchmarks/_output/bin/onos-ran-benchmarks /usr/local/bin/onos-ran-benchmarks

ENTRYPOINT ["onos-ran-benchmarks"]
