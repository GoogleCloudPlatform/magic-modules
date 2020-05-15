#!/bin/bash

set -e

REFERENCE=$1
TEAMCITY_TOKEN=$2

sed -i 's/{{PR_NUMBER}}/'"$REFERENCE"'/g' teamcityparams.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @teamcityparams.xml -o build.json

echo $(cat build.json | jq .webUrl)

#BUILD_ID=$(cat build.json | jq .id)
#STATE=$(cat build.json | jq .state)
#while [ "$STATE" != "finished" ]
#do
#	curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue/id:$BUILD_ID -o job_status.json
#	STATE=$(cat job_status.json | jq .state)
#done

