#!/bin/bash

set -e

gh_repo=magic-modules
new_branch="auto-pr-$PR_NUMBER"

if [ $PR_NUMBER == 7874 ]; then
  post_body=$(jq -n \
    --arg owner "GoogleCloudPlatform" \
    --arg repo "$gh_repo" \
    --arg branch "$new_branch" \
    --arg sha "$COMMIT_SHA" \
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
    -u "modular-magician:$GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/actions/workflows/pre-build-validation.yml/dispatches" \
    -d "$post_body"
fi
