FROM alpine:latest

COPY .docker/install-generic.sh /tmp/install-generic.sh
COPY .docker/build.env /tmp/build.env
COPY dist/kubor-linux-amd64 /usr/bin/kubor

RUN chmod +x /usr/bin/kubor \
    && apk add --no-cache curl make tar zip gzip \
    && sh /tmp/install-generic.sh \
    && rm -rf \
        /usr/share/man \
        /tmp/* \
        /var/cache/apk
