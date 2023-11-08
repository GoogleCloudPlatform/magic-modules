#!/bin/bash

set -e

pr_number=$1
mm_commit_sha=$2
build_id=$3
project_id=$4
build_step=$5
gh_repo=$6
github_username=modular-magician


new_branch="auto-pr-$pr_number"
git_remote=https://$github_username:$GITHUB_TOKEN@github.com/$github_username/$gh_repo
local_path=$GOPATH/src/github.com/GoogleCloudPlatform/$gh_repo
mkdir -p "$(dirname $local_path)"
git clone $git_remote $local_path --branch $new_branch --depth 2
pushd $local_path

# Only skip tests if we can tell for sure that no go files were changed
echo "Checking for modified go files"
# get the names of changed files and look for go files
# (ignoring "no matches found" errors from grep)
gofiles=$(git diff --name-only HEAD~1 | { grep "\.go$" || test $? = 1; })
if [[ -z $gofiles ]]; then
    echo "Skipping tests: No go files changed"
    exit 0
else
    echo "Running tests: Go files changed"
fi


TERRAFORM_BINARY=/terraform/$TERRAFORM_VERSION
if test -f "$TERRAFORM_BINARY"; then
    echo "terraform binary $TERRAFORM_BINARY exists on container"
    echo setting terraform to version $TERRAFORM_VERSION
    set x-
    mv /terraform/$TERRAFORM_VERSION /bin/terraform
    set x+
    terraform version
else
    echo "terraform binary $TERRAFORM_BINARY does not exist."
    echo "exiting ..."
	  state="failure"
    post_body=$( jq -n \
      --arg context "${gh_repo}-test-integration-${TERRAFORM_VERSION}" \
      --arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
      --arg state "$state" \
      '{context: $context, target_url: $target_url, state: $state}')
    exit 0
fi

post_body=$( jq -n \
	--arg context "${gh_repo}-test-integration-${TERRAFORM_VERSION}" \
	--arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
	--arg state "pending" \
	'{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"

set +e

go mod edit -replace github.com/hashicorp/terraform-provider-google-beta=github.com/$github_username/terraform-provider-google-beta@$new_branch
go mod tidy

make build

TF_CONFIG_FILE="tf-dev-override.tfrc"

go clean --modcache
go list -json  -m github.com/hashicorp/terraform-provider-google-beta
REPLACE_DIR=`go list -json  -m github.com/hashicorp/terraform-provider-google-beta | jq -r '.Dir // empty'`
VERSION=`go list -json  -m github.com/hashicorp/terraform-provider-google-beta | jq -r .Version`

if [ ! -z "$REPLACE_DIR" ]
then
  pushd $REPLACE_DIR
    go install
  popd
else
  go install github.com/hashicorp/terraform-provider-google-beta@$VERSION
fi

# create terraform configuration file
if ! [ -f $TF_CONFIG_FILE ];then
  cat <<EOF > $TF_CONFIG_FILE
  provider_installation {
    # Developer overrides will stop Terraform from downloading the listed
    # providers their origin provider registries.
    dev_overrides {
        "hashicorp/google-beta" = "$GOPATH/bin"
    }
    # For all other providers, install them directly from their origin provider
    # registries as normal. If you omit this, Terraform will _only_ use
    # the dev_overrides block, and so no other providers will be available.
    # Without this, show "Failed to query available provider packages"
    # at terraform init
    direct{}
  }
EOF
fi

TF_CLI_CONFIG_FILE="${PWD}/${TF_CONFIG_FILE}" go test -run=CLI ./...
exit_code=$?

set -e

if [ $exit_code -ne 0 ]; then
	state="failure"
else
	state="success"
fi

post_body=$( jq -n \
	--arg context "${gh_repo}-test-integration-${TERRAFORM_VERSION}" \
	--arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
	--arg state "${state}" \
	'{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"
