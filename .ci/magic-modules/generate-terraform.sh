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
mkdir -p "${GOPATH}/src/github.com/$GITHUB_ORG"

pushd magic-modules-branched
ln -s "${PWD}/build/$SHORT_NAME/" "${GOPATH}/src/github.com/$GITHUB_ORG/$PROVIDER_NAME"
popd

pushd "${GOPATH}/src/github.com/$GITHUB_ORG/$PROVIDER_NAME"

# Other orgs are not fully-generated.  This may be transitional - if this causes pain,
# try vendoring into third-party, as with TPG and TPGB.
if [ "$GITHUB_ORG" = "terraform-providers" ]; then
    # This line removes every file which is not specified here.
    # If you add files to Terraform which are not generated, you have to add them here.
    # It uses the somewhat obtuse 'find' command.  To explain:
    # "find .": all files and directories recursively under the current directory, subject to matchers.
    # "-type f": all regular real files, i.e. not directories.
    # "-not": do the opposite of the next thing, always used with another matcher.
    # "-wholename": entire relative path - including directory names - matches following wildcard.
    # "-name": filename alone matches following string.  e.g. -name README.md matches ./README.md *and* ./foo/bar/README.md
    # "-exec": for each file found, execute the command following until the literal ';'
    find . -type f -not -wholename "./.git*" -not -wholename "./vendor*" -not -name ".travis.yml" -not -name ".golangci.yml" -not -name "CHANGELOG.md" -not -name "GNUmakefile" -not -name "docscheck.sh" -not -name "LICENSE" -not -name "README.md" -not -wholename "./examples*" -not -name "go.mod" -not -name "go.sum" -not -name "staticcheck.conf" -not -name ".go-version" -not -name ".hashibot.hcl" -not -name "tools.go"  -exec git rm {} \;
fi

popd

pushd magic-modules-branched

# Choose the author of the most recent commit as the downstream author
# Note that we don't use the last submitted commit, we use the primary GH email
# of the GH PR submitted. If they've enabled a private email, we'll actually
# use their GH noreply email which isn't compatible with CLAs.
COMMIT_AUTHOR="$(git log --pretty="%an <%ae>" -n1 HEAD)"

if [ -n "$OVERRIDE_PROVIDER" ] && [ "$OVERRIDE_PROVIDER" != "null" ]; then
  bundle exec compiler -a -e terraform -f "$OVERRIDE_PROVIDER" -o "${GOPATH}/src/github.com/$GITHUB_ORG/$PROVIDER_NAME/"
else
  bundle exec compiler -a -e terraform -o "${GOPATH}/src/github.com/$GITHUB_ORG/$PROVIDER_NAME/" -v "$VERSION"
fi

TERRAFORM_COMMIT_MSG="$(cat .git/title)"

pushd "build/$SHORT_NAME"

# These config entries will set the "committer".
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

git add -A

git commit -m "$TERRAFORM_COMMIT_MSG" --author="$COMMIT_AUTHOR" || true  # don't crash if no changes
git checkout -B "$(cat ../../branchname)"

apply_patches "$PATCH_DIR/$GITHUB_ORG/$PROVIDER_NAME" "$TERRAFORM_COMMIT_MSG" "$COMMIT_AUTHOR" "master"

popd
popd

git clone "magic-modules-branched/build/$SHORT_NAME" "./terraform-generated/$VERSION"
