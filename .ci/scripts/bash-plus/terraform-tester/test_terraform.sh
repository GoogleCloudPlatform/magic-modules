#!/bin/bash

set -e

version=${VERSION:-$1} # fallback for old variable based declaration
pr_number=${PR_NUMBER:-$2}
mm_commit_sha=${COMMIT_SHA:-$3}
github_username=modular-magician

if [ "$version" == "ga" ]; then
    gh_repo=terraform-provider-google
elif [ "$version" == "beta" ]; then
    gh_repo=terraform-provider-google-beta
else
    echo "no repo, dying."
    exit 1
fi

new_branch="auto-pr-$pr_number"

post_body=$(jq -n \
  --arg owner "$github_username" \
  --arg branch "$new_branch" \
  --arg repo "$gh_repo" \
  --arg sha "$mm_commit_sha" \
  '{
    ref: "main",
    inputs: {
      owner: $owner,
      repo: $repo,
      branch: $branch,
      sha: $sha
    }
  }')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/actions/workflows/test-tpg.yml/dispatches" \
  -d "$post_body"
