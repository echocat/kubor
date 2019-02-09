FROM alpine:latest

ARG  image
ARG  version
COPY var/docker/resources   /
COPY dist/kubor-linux-amd64 /usr/bin/kubor

RUN chmod +x /usr/bin/kubor \
    && apk add --no-cache curl make tar zip gzip \
    && rm -rf \
        /usr/share/man \
        /tmp/* \
        /var/cache/apk \
    && mkdir -p /root/.kubor/binaries \
    && ln -s /usr/bin/kubor /root/.kubor/binaries/kubor-linux-amd64-${version} \
    && mkdir -p /usr/lib/kubor \
    && echo -n "${image}" > /usr/lib/kubor/docker-image \
    && echo -n "${version}" > /usr/lib/kubor/docker-version
