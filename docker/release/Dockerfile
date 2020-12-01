FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD docker-registry-mirror /

CMD ["/docker-registry-mirror"]