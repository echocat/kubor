#!/bin/sh

imageName=${1}
if test -z ${imageName}; then
    echo "Usage: $0 <imageName>" 1>&2
    exit 1
fi
export $(cat .docker/build.env | xargs)

set -ex

docker run --rm ${imageName} kubor -v 2>&1               | grep "$TRAVIS_BRANCH"
docker run --rm ${imageName} kubectl version 2>&1        | grep "GitVersion:\"v${KUBECTL_VERSION}\""
docker run --rm ${imageName} docker version 2>&1         | grep "Version:           ${DOCKER_VERSION}"
docker run --rm ${imageName} docker-machine version 2>&1 | grep "version ${DOCKER_MACHINE_VERSION},"
docker run --rm ${imageName} notary version 2>&1         | grep "Version:    ${DOCKER_NOTARY_VERSION}"
