#!/bin/sh
for submodule in resourcemanager auth compute sql storage spanner dns pubsub container logging; do
  echo "Checking submodule $submodule"
  version=$(jq -r .version "puppet-$submodule-forge/metadata.json")
  name=$(jq -r .name "puppet-$submodule-forge/metadata.json" | sed 's!-!/!')
  echo "Found $name @ $version."
  jq ".dependencies = [.dependencies[] | if .name == \"$name\" then .version_requirement = \">= $version\" else . end]" magic-modules/build/puppet/_bundle/metadata.json > /tmp/metadata.json
  mv /tmp/metadata.json magic-modules/build/puppet/_bundle/metadata.json
done

cp -r -v ./magic-modules/build/puppet/_bundle/. ./release-bundle/
