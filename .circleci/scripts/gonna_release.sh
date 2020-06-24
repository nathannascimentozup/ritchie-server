#!/bin/sh

DEPLOYED_VERSION=$(curl -s https://commons-repo.ritchiecli.io/stable-server.txt)

DIFF_RESULT=$(git --no-pager log --oneline beta...${DEPLOYED_VERSION} 2>/dev/null)

if [ -z "$DIFF_RESULT" ]
then
      echo "ABORT"
else
      echo "RELEASE"
fi