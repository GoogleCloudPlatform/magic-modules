#!/bin/bash

set -e
pr_number=$1
mm_commit_sha=$2
echo "PR number: ${pr_number}"
echo "Commit SHA: ${mm_commit_sha}"
github_username=modular-magician

sed -i 's/{{PR_NUMBER}}/'"$pr_number"'/g' /teamcityparams.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @/teamcityparams.xml -o build.json

function update_status {
	local context="beta-provider-vcr-test"
	local post_body=$( jq -n \
		--arg context "${context}" \
		--arg target_url "${1}" \
		--arg state "${2}" \
		'{context: $context, target_url: $target_url, state: $state}')
	echo "Updating status ${context} to ${2} with target_url ${1} for sha ${mm_commit_sha}"
	curl \
	  -X POST \
	  -u "$github_username:$GITHUB_TOKEN" \
	  -H "Accept: application/vnd.github.v3+json" \
	  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/${mm_commit_sha}" \
	  -d "$post_body"
}

build_url=$(cat build.json | jq -r .webUrl)
update_status "${build_url}" "pending"

ID=$(cat build.json | jq .id -r)
curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
STATUS=$(cat poll.json | jq .status -r)
STATE=$(cat poll.json | jq .state -r)
counter=0
while [[ "$STATE" != "finished" ]]; do
	if [ "$counter" -gt "500" ]; then
		echo "Failed to wait for job to finish, exiting"
		# Call this an error because we don't know if the tests failed or not
		update_status "${build_url}" "error"
		# exit 0 because this script didn't have an error; the failure
		# is reported via the Github Status API
		exit 0
	fi
	sleep 30
	curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
	STATUS=$(cat poll.json | jq .status -r)
	STATE=$(cat poll.json | jq .state -r)
	echo "Trying again, State: $STATE Status: $STATUS"
	counter=$((counter + 1))
done

if [ "$STATUS" == "SUCCESS" ]; then
	echo "Tests succeeded."
	update_status "${build_url}" "success"
	exit 0
fi

curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" http://ci-oss.hashicorp.engineering/app/rest/testOccurrences?locator=build:$ID,status:FAILURE --output failed.json -L
set +e
FAILED_TESTS=$(cat failed.json | jq -r '.testOccurrence | map(.name) | join("|")')
ret=$?
if [ $ret -ne 0 ]; then
	echo "Job failed without failing tests"
	update_status "${build_url}" "failure"
	# exit 0 because this script didn't have an error; the failure
	# is reported via the Github Status API
	exit 0
fi
set -e

sed -i 's/{{PR_NUMBER}}/'"$pr_number"'/g' /teamcityparamsrecording.xml
sed -i 's/{{FAILED_TESTS}}/'"$FAILED_TESTS"'/g' /teamcityparamsrecording.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @/teamcityparamsrecording.xml --output record.json
build_url=$(cat record.json | jq -r .webUrl)
comment="I have triggered VCR tests in RECORDING mode for the following tests that failed during VCR: $FAILED_TESTS You can view the result here: $build_url"

curl -H "Authorization: token ${GITHUB_TOKEN}" \
      -d "$(jq -r --arg comment "$comment" -n "{body: \$comment}")" \
      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${pr_number}/comments"


update_status "${build_url}" "pending"

# Reset for checking failed tests
rm poll.json
rm failed.json

ID=$(cat record.json | jq .id -r)
curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
STATUS=$(cat poll.json | jq .status -r)
STATE=$(cat poll.json | jq .state -r)
counter=0
while [[ "$STATE" != "finished" ]]; do
	if [ "$counter" -gt "500" ]; then
		echo "Failed to wait for job to finish, exiting"
		# Call this an error because we don't know if the tests failed or not
		update_status "${build_url}" "error"
		# exit 0 because this script didn't have an error; the failure
		# is reported via the Github Status API
		exit 0
	fi
	sleep 30
	curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
	STATUS=$(cat poll.json | jq .status -r)
	STATE=$(cat poll.json | jq .state -r)
	echo "Trying again, State: $STATE Status: $STATUS"
	counter=$((counter + 1))
done

if [ "$STATUS" == "SUCCESS" ]; then
	echo "Tests succeeded."
	update_status "${build_url}" "success"
	exit 0
fi

curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" http://ci-oss.hashicorp.engineering/app/rest/testOccurrences?locator=build:$ID,status:FAILURE --output failed.json -L
set +e
FAILED_TESTS=$(cat failed.json | jq -r '.testOccurrence | map(.name) | join("|")')
ret=$?
if [ $ret -ne 0 ]; then
	echo "Job failed without failing tests"
	update_status "${build_url}" "failure"
	# exit 0 because this script didn't have an error; the failure
	# is reported via the Github Status API
	exit 0
fi
set -e

comment="Tests failed during RECORDING mode: $FAILED_TESTS Please fix these to complete your PR"

curl -H "Authorization: token ${GITHUB_TOKEN}" \
      -d "$(jq -r --arg comment "$comment" -n "{body: \$comment}")" \
      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${pr_number}/comments"
update_status "${build_url}" "failure"

# exit 0 because this script didn't have an error; the failure
# is reported via the Github Status API
exit 0
