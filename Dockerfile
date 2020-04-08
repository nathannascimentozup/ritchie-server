ARG BUILD_IMAGE=docker:dind
FROM $BUILD_IMAGE

RUN apk update && apk add --virtual build-dependencies build-base gcc go git python3-dev libffi-dev openssl-dev curl

RUN pip3 install --no-cache-dir --upgrade pip
RUN pip3 install --no-cache-dir docker-compose

RUN docker-compose version
