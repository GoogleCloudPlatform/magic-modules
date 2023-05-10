#!/bin/bash

set -e

pr_number=${PR_NUMBER:-$1}
github_username=modular-magician
gh_repo=magic-modules
new_branch="auto-pr-$pr_number"
mm_commit_sha=${COMMIT_SHA:-$2}

post_body=$(jq -n \
  --arg owner "$github_username" \
  --arg repo "$gh_repo" \
  --arg branch "$new_branch" \
  --arg sha "$mm_commit_sha" \
  '{
    ref: "main",
    inputs: {
      owner: $owner,
      repo: $repo,
      branch: $branch,
      sha: $sha,
    }
  }')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/actions/workflows/pre-build-validation.yml/dispatches" \
  -d "$post_body"
