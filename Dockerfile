FROM alpine:latest

COPY dist/kubor-linux-amd64 /usr/bin/kubor

RUN apk add --no-cache curl make tar zip gzip \
    && curl -sSL https://storage.googleapis.com/kubernetes-release/release/v1.13.0/bin/linux/amd64/kubectl > /usr/bin/kubectl \
    && chmod +x /usr/bin/kubectl
    && chmod +x /usr/bin/kubor \
    && rm -rf \
        /usr/share/man \
        /tmp/* \
        /var/cache/apk
