#!/bin/bash

set -e

REFERENCE=$1

sed -i 's/{{PR_NUMBER}}/'"$REFERENCE"'/g' /teamcityparams.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @/teamcityparams.xml -o build.json

echo $(cat build.json | jq .webUrl)
