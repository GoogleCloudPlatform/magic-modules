#!/usr/bin/env bash

# Example command
# sh scripts/cherry-pick <hash of a post-switchover commit>

set -e
safecommit=$1

backupbranch="$(git rev-parse --symbolic-full-name --abbrev-ref HEAD)"
if [[ $backupbranch != *"-backup"* ]]; then
  echo "\"-backup\" not detected in the branch name \"${backupbranch}\""
  echo "Do you want to still continue with the default branch name \"go-rewrite-convert\"?"
  select yn in "Yes" "No"; do
    case $yn in
        Yes ) backupbranch="go-rewrite-convert"; break;;
        No ) exit;;
    esac
  done
fi
newbranch=${backupbranch%"-backup"}
echo "will use branch \"${newbranch}\""

# prepare temp file commit
git add .
git commit -m "temp file commit" 
currentcommit="$(git rev-parse HEAD)"
echo "committed all changes to ${currentcommit}"

# checkout a new branch from a given post-switchover commit
echo "checking out \"${newbranch}\" at post-switchover commit ${safecommit}"
git checkout -b $newbranch $safecommit

# cherry-pick the previous temp file commit to the new post-switchover branch
echo "cherry-picking ${currentcommit}"
git cherry-pick $currentcommit --no-commit

# overwrite the converted files with the temporary files to produce final diff
echo "cherry-picking ${currentcommit}"
files=`git diff --name-only --diff-filter=A --cached`
for file in $files; do
  echo "moving ${file} to ${file%".temp"}"
  mv $file ${file%".temp"}
done

# stage all changes
git add .
git status