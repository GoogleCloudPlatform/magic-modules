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
if $(echo $USER | fgrep -wq -e megan07 -e rambleraptor -e SirGitsalot -e slevenick -e c2thorn -e rileykarson -e melinath -e ScottSuarez -e shuyama1); then
  echo "User is on the list, not assigning."
  exit 0
fi

# This is where you add people to the random-assignee rotation.  This list
# might not equal the list above.
ASSIGNEE=$(shuf -n 1 <(printf "rileykarson\nc2thorn\nscottsuarez\nshuyama1\nmegan07"))

comment=$(cat << EOF
Hello!  I am a robot who works on Magic Modules PRs.

I've detected that you're a community contributor. @$ASSIGNEE, a repository maintainer, has been assigned to assist you and help review your changes. 

<details>
  <summary>:question: First time contributing? Click here for more details</summary>

---

Your assigned reviewer will help review your code by: 
* Ensuring it's backwards compatible, covers common error cases, etc.
* Summarizing the change into a user-facing changelog note.
* Passes tests, either our "VCR" suite, a set of presubmit tests, or with manual test runs.

You can help make sure that review is quick by running local tests and ensuring they're passing in between each push you make to your PR's branch. Also, try to leave a comment with each push you make, as pushes generally don't generate emails.

If your reviewer doesn't get back to you within a week after your most recent change, please feel free to leave a comment on the issue asking them to take a look! In the absence of a dedicated review dashboard most maintainers manage their pending reviews through email, and those will sometimes get lost in their inbox.

---

</details>

EOF
)

curl -H "Authorization: token ${GITHUB_TOKEN}" \
      -d "$(jq -r --arg comment "$comment" -n "{body: \$comment}")" \
      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}/comments"
curl -H "Authorization: token ${GITHUB_TOKEN}" \
      -d "$(jq -r --arg assignee "$ASSIGNEE" -n "{reviewers: [\$assignee], team_reviewers: []}")" \
      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/${PR_NUMBER}/requested_reviewers"
