#!/usr/bin/env bash

# Example command
# sh scripts/cherry-pick <hash of a post-switchover commit>

set -e
safecommit=$1

# prepare temp file commit
git add .
git commit -m "temp file commit" 
currentcommit="$(git rev-parse HEAD)"

# checkout a new branch from a given post-switchover commit
git checkout -b go-rewrite-convert $safecommit

echo $currentcommit

# cherry-pick the previous temp file commit to the new post-switchover branch
git cherry-pick $currentcommit --no-commit

# overwrite the converted files with the temporary files to produce final diff
files=`git diff --name-only --diff-filter=A --cached`
for file in $files; do
  mv $file ${file%".temp"}
done