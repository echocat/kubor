FROM alpine:latest

RUN apk add --no-cache curl make tar zip gzip \
    && rm -rf \
        /usr/share/man \
        /tmp/* \
        /var/cache/apk

COPY dist/kubor-linux-amd64 /usr/bin/kubor
RUN chmod +x /usr/bin/kubor
