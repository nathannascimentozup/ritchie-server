#!/bin/bash

if expr "$CIRCLE_BRANCH" : 'qa' >/dev/null; then
  export DOCKER_REGISTRY="$DOCKER_REGISTRY_QA"

elif expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then
  export DOCKER_REGISTRY="$DOCKER_REGISTRY_PROD"
else echo ""
fi
