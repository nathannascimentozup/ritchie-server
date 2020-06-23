#!/bin/bash

if   expr "$CIRCLE_BRANCH" : 'qa' || expr "$CIRCLE_BRANCH" : 'beta' >/dev/null; then export RELEASE_VERSION="qa-${CIRCLE_BUILD_NUM}"
elif expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then export RELEASE_VERSION=$(echo "$CIRCLE_BRANCH"| cut -d '-' -f 2-)
else echo ""
fi