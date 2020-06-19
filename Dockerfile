ARG BUILD_IMAGE=docker:dind
FROM $BUILD_IMAGE

RUN apk update && apk add --update --no-cache \
    build-base \
    python3-dev \
    python3 \
    gcc \
    go \
    git \
    curl \
    libffi-dev \
    openssl-dev \
    curl \
    && curl -O https://bootstrap.pypa.io/get-pip.py \
    && python3 get-pip.py

RUN pip install --no-cache-dir --upgrade pip
RUN pip install --no-cache-dir docker-compose

RUN docker-compose version
