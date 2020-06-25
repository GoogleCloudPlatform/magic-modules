#!/bin/bash
if [[ -z "$GITHUB_TOKEN" ]]; then
    echo "Did not provide GITHUB_TOKEN environment variable."
    exit 1
fi
if [[ $# -lt 1 ]]; then
    echo "Usage: $0 pr-number"
    exit 1
fi
PR_NUMBER=$1

set -x

ASSIGNEE=$(curl -H "Authorization: token ${GITHUB_TOKEN}" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/${PR_NUMBER}/requested_reviewers" | jq .users[0].login)
  
if [[ "$ASSIGNEE" == "null" || -z "$ASSIGNEE" ]] ; then 
  ASSIGNEE=$(curl -H "Authorization: token ${GITHUB_TOKEN}" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/${PR_NUMBER}/reviews" | jq .[0].user.login)
fi

if [[ "$ASSIGNEE" == "null"  || -z "$ASSIGNEE" ]] ; then 
  echo "Issue is not assigned."
else
  echo "Issue is assigned, not assigning."
  exit 0
fi

USER=$(curl -H "Authorization: token ${GITHUB_TOKEN}" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}" | jq .user.login)

# This is where you add users who do not need to have an assignee chosen for
# them.
if $(echo $USER | fgrep -wq -e ndmckinley -e danawillow -e emilymye -e megan07 -e paddycarver -e rambleraptor -e SirGitsalot -e slevenick -e c2thorn -e rileykarson); then
  echo "User is on the list, not assigning."
  exit 0
fi

# This is where you add people to the random-assignee rotation.  This list
# might not equal the list above.
ASSIGNEE=$(shuf -n 1 <(printf "danawillow\nemilymye\nrileykarson\nslevenick\nc2thorn\nndmckinley"))

comment=$(cat << EOF
Hello!  I am a robot who works on Magic Modules PRs.

I have detected that you are a community contributor, so your PR will be assigned to someone with a commit-bit on this repo for initial review.

Thanks for your contribution!  A human will be with you soon.

@$ASSIGNEE, please review this PR or find an appropriate assignee.
EOF
)

curl -H "Authorization: token ${GITHUB_TOKEN}" \
      -d "$(jq -r --arg comment "$comment" -n "{body: \$comment}")" \
      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}/comments"
curl -H "Authorization: token ${GITHUB_TOKEN}" \
      -d "$(jq -r --arg assignee "$ASSIGNEE" -n "{reviewers: [\$assignee], team_reviewers: []}")" \
      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/${PR_NUMBER}/requested_reviewers"
