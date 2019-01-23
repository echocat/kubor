FROM alpine:latest

COPY var/docker/resources   /
COPY dist/kubor-linux-amd64 /usr/bin/kubor

RUN chmod +x /usr/bin/kubor \
    && apk add --no-cache curl make tar zip gzip \
    && rm -rf \
        /usr/share/man \
        /tmp/* \
        /var/cache/apk
