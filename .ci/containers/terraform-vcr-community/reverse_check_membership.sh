#!/bin/bash

set -e
pr_number=$1
mm_commit_sha=$2

USER=$(curl -H "Authorization: token ${GITHUB_TOKEN}" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${pr_number}" | jq -r .user.login)

# This image runs tests for community PRs. This script reverses check_membership.sh to exit without running tests
# for users for who tests are automatically run.
if $(echo $USER | fgrep -wq -e megan07 -e slevenick -e c2thorn -e rileykarson -e melinath -e ScottSuarez -e shuyama1 -e trodge -e roaks3); then
	echo "User is on the list, skipping."
	exit 0
else
	echo "Checking GCP org membership"
	GCP_MEMBER=$(curl -Lsw '%{http_code}' -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/orgs/GoogleCloudPlatform/members/$USER -o /dev/null)
	if [ "$GCP_MEMBER" != "404" ]; then
		echo "User is a GCP org member, exiting"
		exit 0
	else
		echo "Checking googlers org membership"
		GOOGLERS_MEMBER=$(curl -Lsw '%{http_code}' -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/orgs/googlers/members/$USER -o /dev/null)
		if [ "$GOOGLERS_MEMBER" != "404" ]; then
			echo "User is a googlers org member, exiting"
			exit 0
		else
			echo "User is not a GCP org member or a googlers org member"
		fi
	fi
fi

# Pass args through to runner
bash /run_vcr_tests.sh $pr_number $mm_commit_sha
