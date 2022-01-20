#! /bin/bash

set -e

function clone_repo() {
    SCRATCH_OWNER=modular-magician
    UPSTREAM_BRANCH=master
    if [ "$REPO" == "terraform" ]; then
        if [ "$VERSION" == "ga" ]; then
            UPSTREAM_OWNER=hashicorp
            GH_REPO=terraform-provider-google
            LOCAL_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google
        elif [ "$VERSION" == "beta" ]; then
            UPSTREAM_OWNER=hashicorp
            GH_REPO=terraform-provider-google-beta
            LOCAL_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
        else
            echo "Unrecognized version $VERSION"
            exit 1
        fi
    elif [ "$REPO" == "tf-conversion" ]; then
        # This is here for backwards compatibility and can be removed after Nov 15 2021
        UPSTREAM_OWNER=GoogleCloudPlatform
        GH_REPO=terraform-google-conversion
        LOCAL_PATH=$GOPATH/src/github.com/GoogleCloudPlatform/terraform-google-conversion
    elif [ "$REPO" == "terraform-validator" ]; then
        UPSTREAM_OWNER=GoogleCloudPlatform
        UPSTREAM_BRANCH=main
        GH_REPO=terraform-validator
        LOCAL_PATH=$GOPATH/src/github.com/GoogleCloudPlatform/terraform-validator
    elif [ "$REPO" == "tf-oics" ]; then
        UPSTREAM_OWNER=terraform-google-modules
        GH_REPO=docs-examples
        LOCAL_PATH=$GOPATH/src/github.com/terraform-google-modules/docs-examples
    elif [ "$REPO" == "tf-cloud-docs" ]; then
        UPSTREAM_OWNER=terraform-google-modules
        GH_REPO=terraform-docs-samples
        LOCAL_PATH=$GOPATH/src/github.com/terraform-google-modules/terraform-docs-samples
    else
        echo "Unrecognized repo $REPO"
        exit 1
    fi

    GITHUB_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/$UPSTREAM_OWNER/$GH_REPO
    SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/$SCRATCH_OWNER/$GH_REPO
    mkdir -p "$(dirname $LOCAL_PATH)"
    git clone $GITHUB_PATH $LOCAL_PATH
}

if [ $# -lt 4 ]; then
    echo "Usage: $0 (build|diff|downstream) (terraform|terraform-validator) (ga|beta) (pr number|sha)"
    exit 1
fi
if [ -z "$GITHUB_TOKEN" ]; then
    echo "Did not provide GITHUB_TOKEN environment variable."
    exit 1
fi


COMMAND=$1
REPO=$2
VERSION=$3
REFERENCE=$4

mkdir ../mm-$REPO-$VERSION-$COMMAND
cp -rp ./. ../mm-$REPO-$VERSION-$COMMAND
pushd ../mm-$REPO-$VERSION-$COMMAND

clone_repo

git config --local user.name "Modular Magician"
git config --local user.email "magic-modules@google.com"

# MMv1 now lives inside a sub-folder
pushd mmv1

if [ "$COMMAND" == "head" ]; then
    BRANCH=auto-pr-$REFERENCE
    COMMIT_MESSAGE="New generated code for MM PR $REFERENCE."
elif [ "$COMMAND" == "base" ]; then
    # In this case, there is guaranteed to be a merge commit,
    # and the *left* side of it is the old master branch.
    # the *right* side of it is the code to be merged.
    git checkout HEAD~
    BRANCH=auto-pr-$REFERENCE-old
    COMMIT_MESSAGE="Old generated code for MM PR $REFERENCE."
elif [ "$COMMAND" == "downstream" ]; then
    BRANCH=downstream-pr-$REFERENCE
    COMMIT_MESSAGE="$(git log -1 --pretty=%B "$REFERENCE")"
fi

if [ "$REPO" == "terraform" ]; then
    pushd $LOCAL_PATH
    go mod download
    find . -type f -not -wholename "./.git*" -not -wholename "./.changelog*" -not -name ".travis.yml" -not -name ".golangci.yml" -not -name "CHANGELOG.md" -not -name "GNUmakefile" -not -name "docscheck.sh" -not -name "LICENSE" -not -name "README.md" -not -wholename "./examples*" -not -name ".go-version" -not -name ".hashibot.hcl" -not -name "tools.go"  -exec git rm {} \;
    popd
fi

if [ "$REPO" == "terraform-validator" ] || [ "$REPO" == "tf-conversion" ]; then
    # use terraform generator with validator overrides.
    # Check for tf-conversion is legacy and can be removed after Nov 15 2021
    if [ "$REPO" == "terraform-validator" ] && [ "$COMMAND" == "base" ] && [ ! -d "../.ci/containers/terraform-validator-tester" ]; then
        # Temporary shim to allow building a "base" version; only required until after
        # initial merge. If we're building a base branch on an old mmv1 master (which
        # we can detect by the lack of files added in this PR) the base version will
        # require a `google` folder to exist.
        mkdir -p $LOCAL_PATH/google
    fi
    pushd $LOCAL_PATH
    # clear out the templates as they are copied during
    # generation from mmv1/third_party/validator/tests/data
    rm -rf ./testdata/templates/
    rm -rf ./testdata/generatedconvert/
    git add ./testdata
    find ./test/** -type f -exec git rm {} \;
    popd
    bundle exec compiler -a -e terraform -f validator -o $LOCAL_PATH -v $VERSION
    pushd $LOCAL_PATH
    make build
    export TFV_CREATE_GENERATED_FILES=true
    go test ./test -run "TestAcc.*_generated_offline"
    popd
elif [ "$REPO" == "tf-oics" ]; then
    # use terraform generator with oics override
    bundle exec compiler -a -e terraform -f oics -o $LOCAL_PATH -v $VERSION
elif [ "$REPO" == "tf-cloud-docs" ]; then
    # use terraform generator with cloud docs override
    bundle exec compiler -a -e terraform -f cloud_docs -o $LOCAL_PATH -v $VERSION
else
    if [ "$REPO" == "terraform" ] && [ "$VERSION" == "ga" ]; then
        bundle exec compiler -a -e $REPO -o $LOCAL_PATH -v $VERSION --no-docs
        bundle exec compiler -a -e $REPO -o $LOCAL_PATH -v beta --no-code
        # TODO(slevenick): remove this check when it is safe (~1 month from commit)
        # Previously we had many resources committed to tpgtools that were not
        # ready for generation. Block generation until these are removed
        set +e
        git merge-base --is-ancestor 0be5f0c31a6e69474b14e91b12c0bbc1e550df9c HEAD
        if [ $? == 0 ]; then
            pushd ../
            make tpgtools OUTPUT_PATH=$LOCAL_PATH VERSION=$VERSION
            popd
        fi
        set -e
    else
        bundle exec compiler -a -e $REPO -o $LOCAL_PATH -v $VERSION
        # TODO(slevenick): remove this check when it is safe (~1 month from commit)
        # Previously we had many resources committed to tpgtools that were not
        # ready for generation. Block generation until these are removed
        set +e
        git merge-base --is-ancestor 0be5f0c31a6e69474b14e91b12c0bbc1e550df9c HEAD
        if [ $? == 0 ]; then
            pushd ../
            make tpgtools OUTPUT_PATH=$LOCAL_PATH VERSION=$VERSION
            popd
        fi
        set -e
    fi
fi

popd

pushd $LOCAL_PATH

if [ "$REPO" == "terraform" ]; then
    make generate
fi

git config --local user.name "Modular Magician"
git config --local user.email "magic-modules@google.com"
git add .
git checkout -b $BRANCH

COMMITTED=true
git commit --signoff -m "$COMMIT_MESSAGE" || COMMITTED=false

CHANGELOG=false
if [ "$REPO" == "terraform" ]; then
  CHANGELOG=true
fi

PR_NUMBER=$(curl -L -s -H "Authorization: token ${GITHUB_TOKEN}" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls?state=closed&base=master&sort=updated&direction=desc" | \
    jq -r ".[] | if .merge_commit_sha == \"$REFERENCE\" then .number else empty end")
if [ "$COMMITTED" == "true" ] && [ "$COMMAND" == "downstream" ] && [ "$CHANGELOG" == "true" ]; then
    # Add the changelog entry!
    mkdir -p .changelog/
    curl -L -s -H "Authorization: token ${GITHUB_TOKEN}" \
        "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/$PR_NUMBER" | \
        jq -r .body | \
        sed -e '/```release-note/,/```/!d' \
        > .changelog/$PR_NUMBER.txt
    git add .changelog
    git commit --signoff --amend --no-edit
fi

git push $SCRATCH_PATH $BRANCH -f

if [ "$COMMITTED" == "true" ] && [ "$COMMAND" == "downstream" ]; then
    PR_BODY=$(curl -L -s -H "Authorization: token ${GITHUB_TOKEN}" \
        "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/$PR_NUMBER" | \
        jq -r .body)
    PR_TITLE=$(curl -L -s -H "Authorization: token ${GITHUB_TOKEN}" \
        "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/$PR_NUMBER" | \
        jq -r .title)
    MM_PR_URL=$(curl -L -s -H "Authorization: token ${GITHUB_TOKEN}" \
        "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/$PR_NUMBER" | \
        jq -r .html_url)

    echo "Base: $UPSTREAM_OWNER:$UPSTREAM_BRANCH"
    echo "Head: $SCRATCH_OWNER:$BRANCH"
    NEW_PR_URL=$(hub pull-request -b $UPSTREAM_OWNER:$UPSTREAM_BRANCH -h $SCRATCH_OWNER:$BRANCH -m "$PR_TITLE" -m "$PR_BODY" -m "Derived from $MM_PR_URL")
    if [ $? != 0 ]; then
        exit $?
    fi
    echo "Created PR $NEW_PR_URL"
    NEW_PR_NUMBER=$(echo $NEW_PR_URL | awk -F '/' '{print $NF}')

    # Wait a few seconds, then merge the PR.
    sleep 5
    echo "Merging PR $NEW_PR_URL"
    curl -L -H "Authorization: token ${GITHUB_TOKEN}" \
        -X PUT \
        -d '{"merge_method": "squash"}' \
        "https://api.github.com/repos/$UPSTREAM_OWNER/$GH_REPO/pulls/$NEW_PR_NUMBER/merge"
fi

popd
