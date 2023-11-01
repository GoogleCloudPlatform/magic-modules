#!/bin/bash

set -e

pr_number=${PR_NUMBER:-$1}
mm_commit_sha=${COMMIT_SHA:-$2}
gh_repo=terraform-google-conversion
github_username=modular-magician
new_branch="auto-pr-$pr_number"

post_body=$(jq -n \
  --arg owner "$github_username" \
  --arg branch "$new_branch" \
  --arg repo "$gh_repo" \
  --arg sha "$mm_commit_sha" \
  '{
    ref: "tgc-units",
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
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/actions/workflows/test-tgc.yml/dispatches" \
  -d "$post_body"
