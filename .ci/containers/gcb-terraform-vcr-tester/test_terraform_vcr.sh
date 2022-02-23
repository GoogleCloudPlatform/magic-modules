#!/bin/bash

set -e
pr_number=$1
mm_commit_sha=$2
echo "PR number: ${pr_number}"
echo "Commit SHA: ${mm_commit_sha}"
github_username=modular-magician
gh_repo=terraform-provider-google-beta

# For testing only.
echo "checking if this is the testing PR for vcr setup"
if [ "$pr_number" == "5703"]; then
  echo "Running tests for new VCR setup"
else
  echo "Skipping new vcr tests: Not testing PR"
  exit 0
fi
