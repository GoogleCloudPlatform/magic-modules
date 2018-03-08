#! /bin/bash
set -e
set -x

pushd "magic-modules"
BRANCH=$(git rev-parse --short HEAD)
git checkout -B "$BRANCH"
# ./branchname is intentionally never committed - it isn't necessary once
# this output is no longer available.
echo "$BRANCH" > ./branchname

cp -r ./ ../magic-modules-branched/
