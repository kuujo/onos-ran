FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/onos-ric-tests/_output/bin/onos-ric-tests /usr/local/bin/onos-ric-tests

ENTRYPOINT ["onos-ric-tests"]
