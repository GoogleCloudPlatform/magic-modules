#!/bin/bash
set -e
git restore . || true
git fetch
git checkout -B rewrite-branch origin/main
for commit in $(git rev-list --reverse origin/main..origin/container_flatten_bk); do
    if ! git cherry-pick -n $commit; then
        echo "Conflict, resolving favoring origin/main for target files and our changes else"
        git checkout --theirs .
    fi
    git checkout origin/main -- mmv1/provider/terraform_tgc_next.go 2>/dev/null || true
    git checkout origin/main -- mmv1/third_party/tgc_next/Makefile 2>/dev/null || true
    git add .
    GIT_AUTHOR_NAME="$(git log -1 --format="%an" $commit)" \
    GIT_AUTHOR_EMAIL="$(git log -1 --format="%ae" $commit)" \
    GIT_AUTHOR_DATE="$(git log -1 --format="%ad" $commit)" \
    GIT_COMMITTER_NAME="$(git log -1 --format="%cn" $commit)" \
    GIT_COMMITTER_EMAIL="$(git log -1 --format="%ce" $commit)" \
    GIT_COMMITTER_DATE="$(git log -1 --format="%cd" $commit)" \
    git commit -m "$(git log -1 --format=%B $commit)" || echo "Empty commit skipped"
done
git checkout container_flatten_bk
git reset --hard rewrite-branch
git push -f origin container_flatten_bk
