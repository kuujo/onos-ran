FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/onos-ran-tests/_output/bin/onos-ran-tests /usr/local/bin/onos-ran-tests

ENTRYPOINT ["onos-ran-tests"]
