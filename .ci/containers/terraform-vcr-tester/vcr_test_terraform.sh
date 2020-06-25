#!/bin/bash

set -e

PR_NUMBER=$1

USER=$(curl -H "Authorization: token ${GITHUB_TOKEN}" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}" | jq .user.login)

# Only run tests for safe users
if $(echo $USER | fgrep -wq -e ndmckinley -e danawillow -e emilymye -e megan07 -e paddycarver -e rambleraptor -e SirGitsalot -e slevenick -e c2thorn -e rileykarson); then
  echo "User is on the list, not skipping."
else
	echo "User is not on the list, skipping."
  exit 0
fi

sed -i 's/{{PR_NUMBER}}/'"$PR_NUMBER"'/g' /teamcityparams.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @/teamcityparams.xml -o build.json

URL=$(cat build.json | jq .webUrl)
ret=$?
if [ $ret -ne 0 ]; then
  echo "Auth failed"
  exit 1
else
	comment="I have triggered VCR tests based on this PR's diffs. See the results here: $URL"

	curl -H "Authorization: token ${GITHUB_TOKEN}" \
	      -d "$(jq -r --arg comment "$comment" -n "{body: \$comment}")" \
	      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}/comments"
fi

ID=$(cat build.json | jq .id -r)
curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
STATUS=$(cat poll.json | jq .status -r)
STATE=$(cat poll.json | jq .state -r)
counter=0
while [[ "$STATE" != "finished" ]]; do
	if [ "$counter" -gt "50" ]; then
		echo "Failed to wait for job to finish, exiting"
		exit 1
	fi
	sleep 30
	curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
	STATUS=$(cat poll.json | jq .status -r)
	STATE=$(cat poll.json | jq .state -r)
	echo "Trying again, State: $STATE Status: $STATUS"
	counter=$((counter + 1))
done

if [ "$STATUS" == "SUCCESS" ]
	echo "Tests succeeded."
	exit 0
fi

curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" http://ci-oss.hashicorp.engineering/app/rest/testOccurrences?locator=build:$ID,status:FAILURE --output failed.json -L
set +e
FAILED_TESTS=$(cat failed.json | jq -r '.testOccurrence | map(.name) | join("|")')
ret=$?
if [ $ret -ne 0 ]; then
  echo "Job failed without failing tests"
  exit 1
else

sed -i 's/{{PR_NUMBER}}/'"$PR_NUMBER"'/g' /teamcityparamsrecording.xml
sed -i 's/{{FAILED_TESTS}}/'"$FAILED_TESTS"'/g' /teamcityparamsrecording.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @/teamcityparamsrecording.xml