#!/bin/sh
export $(cat /tmp/build.env | xargs)
set -ex

curl -sSLf https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl > /usr/bin/kubectl
chmod +x /usr/bin/kubectl

curl -sSLf https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz | tar -zOx docker/docker > /usr/bin/docker
chmod +x /usr/bin/docker

curl -sSLf https://github.com/docker/machine/releases/download/v${DOCKER_MACHINE_VERSION}/docker-machine-Linux-x86_64 > /usr/bin/docker-machine
chmod +x /usr/bin/docker-machine

curl -sSLf https://github.com/theupdateframework/notary/releases/download/v${DOCKER_NOTARY_VERSION}/notary-Linux-amd64 > /usr/bin/notary
chmod +x /usr/bin/notary
