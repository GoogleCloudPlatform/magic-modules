#!/bin/bash

# The vast majority of this file is a direct copy of generate-terraform.sh.  We could factor out all that
# code into a shared library, but I don't think we need to do that.  This is an inherently temporary file,
# until TPG 3.0.0 is released, which is in the relatively near future.  The cost of the copy is that
# we need to maintain both files - but the last change to that file was several months ago and I expect
# we're looking at 1 - 2 changes that need to be made in both places.  The cost of not copying it is
# an extra few hours of work now, and some minor readability issues.

set -x
set -e
source "$(dirname "$0")/helpers.sh"

# Create $GOPATH structure - in order to successfully run Terraform codegen, we need to run
# it with a correctly-set-up $GOPATH.  It calls out to `goimports`, which means that
# we need to have all the dependencies correctly downloaded.
export GOPATH="${PWD}/go"
mkdir -p "${GOPATH}/src/github.com/$GITHUB_ORG"

for mm_dir in magic-modules-branched magic-modules-previous; do

    pushd $mm_dir
    # delete the symlink if it exists
    rm "${GOPATH}/src/github.com/$GITHUB_ORG/$PROVIDER_NAME" || true
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
        find . -type f -not -wholename "./.git*" -not -wholename "./vendor*" -not -name ".travis.yml" -not -name ".golangci.yml" -not -name "CHANGELOG.md" -not -name GNUmakefile -not -name LICENSE -not -name README.md -not -wholename "./examples*" -not -name "go.mod" -not -name "go.sum" -not -name "staticcheck.conf" -not -name  ".hashibot.hcl"  -exec git rm {} \;
    fi

    popd

    pushd $mm_dir

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

    if [ "$mm_dir" == "magic-modules-branched" ] ; then
        TERRAFORM_COMMIT_MSG="$(cat .git/title)"
    else
        TERRAFORM_COMMIT_MSG="Old generated base as of $(git rev-parse HEAD)."
    fi

    BRANCH_NAME="$(cat branchname)"

    pushd "build/$SHORT_NAME"

    # These config entries will set the "committer".
    git config --global user.email "magic-modules@google.com"
    git config --global user.name "Modular Magician"

    git add -A

    git commit -m "$TERRAFORM_COMMIT_MSG" --author="$COMMIT_AUTHOR" || true  # don't crash if no changes
    git checkout -B "$BRANCH_NAME"

    popd
    popd

done

mkdir "./terraform-diff/$VERSION"

git clone "magic-modules-branched/build/$SHORT_NAME" "./terraform-diff/$VERSION/new"
git clone "magic-modules-previous/build/$SHORT_NAME" "./terraform-diff/$VERSION/old"
