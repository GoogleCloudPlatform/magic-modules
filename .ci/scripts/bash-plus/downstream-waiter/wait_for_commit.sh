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
    if [ "$BASE_BRANCH" != "main" ]; then
        SYNC_HEAD="$(git rev-parse --short origin/$SYNC_BRANCH)"
        BASE_PARENT="$(git rev-parse --short $SHA~)"
        if [ "$SYNC_HEAD" == "$BASE_PARENT" ]; then
            break;
        else
            echo "sync branch is at: $SYNC_HEAD"
            echo "current commit is $SHA"
            git fetch origin $SYNC_BRANCH
        fi
    else 
        commits="$(git log --pretty=%H origin/$SYNC_BRANCH..origin/$BASE_BRANCH | tail -n 1)"
        if [ "$commits" == "$SHA" ]; then
            break
        else
            echo "git log says waiting on: $commits"
            echo "command says waiting on $SHA"
            git fetch origin $SYNC_BRANCH
        fi
    fi
    sleep 5
done