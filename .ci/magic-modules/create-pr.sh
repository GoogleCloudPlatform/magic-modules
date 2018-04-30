#!/bin/bash

# This script configures the git submodule under magic-modules so that it is
# ready to create a new pull request.  It is cloned in a detached-head state,
# but its branch is relevant to the PR creation process, so we want to make
# sure that it's on a branch, and most importantly that that branch tracks
# a branch upstream.

set -e
set -x

shopt -s dotglob
cp -r magic-modules/* magic-modules-with-comment

ORIGINAL_PR_BRANCH="codegen-pr-$(cat ./mm-initial-pr/.git/id)"
pushd magic-modules-with-comment
echo "$ORIGINAL_PR_BRANCH" > ./original_pr_branch_name

# Check out the magic-modules branch with the same name as the current tracked
# branch of the terraform submodule.  All the submodules will be on the the same
# branch name - we pick terraform because it's the first one the magician supported.
BRANCH_NAME="$(git config -f .gitmodules --get submodule.build/terraform.branch)"
IFS="," read -ra PUPPET_PRODUCTS <<< "$PUPPET_MODULES"
git checkout -b "$BRANCH_NAME"

if [ "$BRANCH_NAME" = "$ORIGINAL_PR_BRANCH" ]; then
  DEPENDENCIES=""
  NEWLINE=$'\n'
  # There is no existing PR - this is the first pass through the pipeline and
  # we will need to create a PR using 'hub'.
  if [ -n "$TERRAFORM_REPO_USER" ]; then
    pushd build/terraform

    git log -1 --pretty=%B > ./downstream_body
    echo "" >> ./downstream_body
    echo "<!-- This change is generated by MagicModules. -->" >> ./downstream_body

    git checkout -b "$BRANCH_NAME"
    if TF_PR=$(hub pull-request -b "$TERRAFORM_REPO_USER/terraform-provider-google:master" -F ./downstream_body); then
      DEPENDENCIES="${DEPENDENCIES}depends: $TF_PR ${NEWLINE}"
    else
      echo "Terraform - did not generate a PR."
    fi
    popd
  fi
  
  for PRD in "${PUPPET_PRODUCTS[@]}"; do

    pushd "build/puppet/$PRD"

    git log -1 --pretty=%B > ./downstream_body
    echo "" >> ./downstream_body
    echo "<!-- This change is generated by MagicModules. -->" >> ./downstream_body

    git checkout -b "$BRANCH_NAME"
    if PUP_PR=$(hub pull-request -b "$PUPPET_REPO_USER/puppet-google-$PRD:master" -F ./downstream_body); then
      DEPENDENCIES="${DEPENDENCIES}depends: $PUP_PR ${NEWLINE}"
    else
      echo "Puppet $PRD - did not generate a PR."
    fi
    popd
  done

  if [ -z "$DEPENDENCIES" ]; then
    cat << EOF > ./pr_comment 
I am a robot that works on MagicModules PRs!
I checked the downstream repositories (see README.md for which ones I can write to), and none of them seem to have any changes.

Once this PR is approved, you can feel free to merge it without taking any further steps.
EOF
  else
    cat << EOF > ./pr_comment 
I am a robot that works on MagicModules PRs!

I built this PR into one or more PRs on other repositories, and when those are closed, this PR will also be merged and closed.
$DEPENDENCIES
EOF
  fi

else
  # This is the second-or-more pass through the pipeline - we need to overwrite
  # the codegen-pr-* branch with the new updated code to update the existing
  # PR, rather than create a new one.
  git branch -f "$ORIGINAL_PR_BRANCH"

  if [ -n "$TERRAFORM_REPO_USER" ]; then
    pushd build/terraform
    git branch -f "$ORIGINAL_PR_BRANCH"
    popd
  fi

  for PRD in "${PUPPET_PRODUCTS[@]}"; do
    pushd "build/puppet/$PRD"
    git branch -f "$ORIGINAL_PR_BRANCH"
    popd
  done
  # Note - we're interested in HEAD~1 here, not HEAD, because HEAD is the
  # generated code commit.  :)
  cat << EOF > ./pr_comment
I am (still) a robot that works on MagicModules PRs!

I just wanted to let you know that your changes (as of commit $(git rev-parse --short HEAD~1)) have been included in your existing downstream PRs.
EOF

fi
