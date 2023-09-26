#!/bin/bash

set -e
if [ $# -lt 3 ]; then
    echo "Usage: $0 (sync-branch) (base-branch) (sha)"
    exit 1
fi

SYNC_BRANCH_PREFIX=$1
BASE_BRANCH=$2
SHA=$3

if [ "$BASE_BRANCH" == "main" ]; then
    SYNC_BRANCH=$SYNC_BRANCH_PREFIX
else
    SYNC_BRANCH=$SYNC_BRANCH_PREFIX-$BASE_BRANCH
fi

echo "SYNC_BRANCH: $SYNC_BRANCH"

if git merge-base --is-ancestor $SHA origin/$SYNC_BRANCH; then
    echo "Found $SHA in history of $SYNC_BRANCH - dying to avoid double-generating that commit."
    exit 1
fi

while true; do
    SYNC_HEAD="$(git rev-parse --short origin/$SYNC_BRANCH)"
    BASE_PARENT="$(git rev-parse --short origin/$BASE_BRANCH~)"
    if [ "$SYNC_HEAD" == "$BASE_PARENT" ]; then
        break;
    else
        echo "sync branch at: $SYNC_HEAD"
        echo "base branch at: $BASE_BRANCH"
        git fetch origin $SYNC_BRANCH
    fi
    sleep 5
done