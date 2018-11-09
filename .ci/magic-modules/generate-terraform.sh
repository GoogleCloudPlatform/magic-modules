#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It outputs "terraform-generated", a non-submodule git repo containing the generated terraform code.

set -x
set -e
source "$(dirname "$0")/helpers.sh"
PATCH_DIR="$(pwd)/patches"

# Create $GOPATH structure - in order to successfully run Terraform codegen, we need to run
# it with a correctly-set-up $GOPATH.  It calls out to `goimports`, which means that
# we need to have all the dependencies correctly downloaded.
export GOPATH="${PWD}/go"
mkdir -p "${GOPATH}/src/github.com/terraform-providers"

pushd magic-modules-branched
ln -s "${PWD}/build/$SHORT_NAME/" "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME"
popd

pushd "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME"

go get
# This line removes every file which is not specified here.
# If you add files to Terraform which are not generated, you have to add them here.
# It uses the somewhat obtuse 'find' command.  To explain:
# "find .": all files and directories recursively under the current directory, subject to matchers.
# "-type f": all regular real files, i.e. not directories.
# "-not": do the opposite of the next thing, always used with another matcher.
# "-wholename": entire relative path - including directory names - matches following wildcard.
# "-name": filename alone matches following string.  e.g. -name README.md matches ./README.md *and* ./foo/bar/README.md
# "-exec": for each file found, execute the command following until the literal ';'
find . -type f -not -wholename "./.git*" -not -wholename "./vendor*" -not -name ".travis.yml" -not -name "CHANGELOG.md" -not -name GNUmakefile -not -name LICENSE -not -name README.md -not -wholename "./examples*" -not -name "main.go" -not -wholename "./scripts*" -not -wholename "./version*" -exec git rm {} \;
popd

pushd magic-modules-branched
LAST_COMMIT_AUTHOR="$(git log --pretty="%an <%ae>" -n1 HEAD)"
bundle install

# Build all terraform products
bundle exec compiler -a -e terraform -o "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME/" -v "$VERSION"

# This command can crash - if that happens, the script should not fail.
set +e
TERRAFORM_COMMIT_MSG="$(python .ci/magic-modules/extract_from_pr_description.py --tag "$SHORT_NAME" < .git/body)"
set -e
if [ -z "$TERRAFORM_COMMIT_MSG" ]; then
  TERRAFORM_COMMIT_MSG="Magic Modules changes."
fi

pushd "build/$SHORT_NAME"
# These config entries will set the "committer".
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

git add -A
# Set the "author" to the commit's real author.
git commit -m "$TERRAFORM_COMMIT_MSG" --author="$LAST_COMMIT_AUTHOR" || true  # don't crash if no changes
git checkout -B "$(cat ../../branchname)"

apply_patches "$PATCH_DIR/terraform-providers/$PROVIDER_NAME" "$TERRAFORM_COMMIT_MSG" "$LAST_COMMIT_AUTHOR" "2.0.0"

popd
popd

git clone "magic-modules-branched/build/$SHORT_NAME" "./terraform-generated/$VERSION"
