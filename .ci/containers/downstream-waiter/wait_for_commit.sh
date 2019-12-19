#!/bin/bash

set -e
if [ $# -lt 3 ]; then
    echo "Usage: $0 (sync-branch) (base-branch) (sha)"
    exit 1
fi

SYNC_BRANCH=$1
BASE_BRANCH=$2
SHA=$3

if git merge-base --is-ancestor $SHA origin/$SYNC_BRANCH; then
    echo "Found $SHA in history of $SYNC_BRANCH - dying to avoid double-generating that commit."
    exit 1
fi

while true; do
    commits="$(git log --pretty=%H origin/$SYNC_BRANCH..origin/$BASE_BRANCH | tail -n 1)"
    if [ "$commits" == "$SHA" ]; then
        break
    else
        echo "git log says waiting on: $commits"
        echo "command says waiting on $SHA"
        git fetch origin $SYNC_BRANCH
    fi
    sleep 5
done
