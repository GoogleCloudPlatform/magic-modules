#!/bin/bash

set -e

PR_NUMBER=$1

sed -i 's/{{PR_NUMBER}}/'"$PR_NUMBER"'/g' /teamcityparams.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @/teamcityparams.xml -o build.json

# Dont crash here if the curl failed due to authorization
# TODO(slevenick): remove this once this all is stable
set +e
URL=echo $(cat build.json | jq .webUrl)
ret=$?
if [ $ret -ne 0 ]; then
	comment="I have triggered VCR tests based on this PR's diffs. See the results here: $URL"

	curl -H "Authorization: token ${GITHUB_TOKEN}" \
	      -d "$(jq -r --arg comment "$comment" -n "{body: \$comment}")" \
	      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}/comments"
else
  echo "Auth failed"
fi