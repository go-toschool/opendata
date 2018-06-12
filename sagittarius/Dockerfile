FROM alpine:3.4

RUN apk add --update ca-certificates

MAINTAINER rodrigo <rf@finciero.com>

COPY bin/aiolos /usr/bin/aiolos

EXPOSE 3000

ENTRYPOINT ["aiolos"]