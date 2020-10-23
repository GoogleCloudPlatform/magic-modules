#!/bin/bash

set -e
PR_NUMBER=$1

USER=$(curl -H "Authorization: token ${GITHUB_TOKEN}" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}" | jq -r .user.login)

# Only run tests for safe users
if $(echo $USER | fgrep -wq -e ndmckinley -e danawillow -e emilymye -e megan07 -e paddycarver -e rambleraptor -e SirGitsalot -e slevenick -e c2thorn -e rileykarson -e melinath -e scottsuarez); then
	echo "User is on the list, not skipping."
else
	echo "Checking GCP org membership"
	GCP_MEMBER=$(curl -sw '%{http_code}' -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/orgs/GoogleCloudPlatform/members/$USER -o /dev/null)
	if [ "$GCP_MEMBER" != "404" ]; then
		echo "User is a GCP org member, continuing"
	else
		echo "User is not a GCP org member"
		exit 0
	fi
fi

# Pass PR number to runner, which expects it
sh /run_vcr_tests.sh $PR_NUMBER