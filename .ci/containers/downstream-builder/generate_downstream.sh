#! /bin/bash

set -e

function clone_repo() {
    if [ "$REPO" == "terraform" ]; then
        if [ "$VERSION" == "ga" ]; then
            GITHUB_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/terraform-providers/terraform-provider-google
            SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-provider-google
            LOCAL_PATH=$GOPATH/src/github.com/terraform-providers/terraform-provider-google
        elif [ "$VERSION" == "beta" ]; then
            GITHUB_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/terraform-providers/terraform-provider-google-beta
            SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-provider-google-beta
            LOCAL_PATH=$GOPATH/src/github.com/terraform-providers/terraform-provider-google-beta
        else
            echo "Unrecognized version $VERSION"
            exit 1
        fi
    elif [ "$REPO" == "tf-conversion" ]; then
        GITHUB_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/GoogleCloudPlatform/terraform-google-conversion
        SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-google-conversion
        LOCAL_PATH=$GOPATH/src/github.com/GoogleCloudPlatform/terraform-google-conversion
    elif [ "$REPO" == "ansible" ]; then
        GITHUB_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/ansible-collections/ansible_collections_google
        SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/ansible_collections_google
        LOCAL_PATH=$PWD/../ansible
    elif [ "$REPO" == "inspec" ]; then
        GITHUB_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/inspec-gcp
        SCRATCH_PATH=$GITHUB_PATH
        LOCAL_PATH=$PWD/../inspec
    else
        echo "Unrecognized repo $REPO"
        exit 1
    fi
    mkdir -p "$(dirname $LOCAL_PATH)"
    git clone $GITHUB_PATH $LOCAL_PATH
}

if [ $# -lt 4 ]; then
    echo "Usage: $0 (build|diff) (terraform|tf-conversion|ansible|inspec) (ga|beta) (pr number)"
    exit 1
fi
if [ -z "$GITHUB_TOKEN" ]; then
    echo "Did not provide GITHUB_TOKEN environment variable."
    exit 1
fi

COMMAND=$1
REPO=$2
VERSION=$3
PR_NUMBER=$4

clone_repo

git config --global user.name "Modular Magician"
git config --global user.email "magic-modules@google.com"

if [ "$COMMAND" == "head" ]; then
    BRANCH=auto-pr-$PR_NUMBER
elif [ "$COMMAND" == "base" ]; then
    git checkout HEAD~
    BRANCH=auto-pr-$PR_NUMBER-old
fi

if [ "$REPO" == "terraform" ]; then
    pushd $LOCAL_PATH
    go get -v

    if [ "$REPO" == "terraform" ]; then
        find . -type f -not -wholename "./.git*" -not -wholename "./vendor*" -not -name ".travis.yml" -not -name ".golangci.yml" -not -name "CHANGELOG.md" -not -name "GNUmakefile" -not -name "docscheck.sh" -not -name "LICENSE" -not -name "README.md" -not -wholename "./examples*" -not -name "go.mod" -not -name "go.sum" -not -name "staticcheck.conf" -not -name ".go-version" -not -name ".hashibot.hcl" -not -name "tools.go"  -exec git rm {} \;
    fi
    popd
elif [ "$REPO" == "tf-conversion" ]; then
    pushd $LOCAL_PATH
    go get -v ./google
    popd
fi

if [ "$REPO" == "tf-conversion" ]; then
    bundle exec compiler -a -e terraform -f validator -o $LOCAL_PATH -v $VERSION
else
    bundle exec compiler -a -e $REPO -o $LOCAL_PATH -v $VERSION
fi


pushd $LOCAL_PATH
git add .
git checkout -b $BRANCH
git commit -m "New generated code for PR $PR_NUMBER." || true
git push $SCRATCH_PATH $BRANCH -f
popd
