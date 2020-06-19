#!/bin/bash

if expr "$CIRCLE_BRANCH" : 'qa' || expr "$CIRCLE_BRANCH" : '^beta-.*' >/dev/null; then
  export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID_QA"
  export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY_QA"
  export DOCKER_AWS_REGION="sa-east-1"

elif expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then
  export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID_PROD"
  export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY_PROD"
  export DOCKER_AWS_REGION="sa-east-1"
else echo ""
fi