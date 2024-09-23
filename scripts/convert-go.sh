#!/usr/bin/env bash

# Example command (in a pre-switchover commit)
# sh scripts/convert-go.sh <provider output directory> <comma separated list of files that have changed in the PR> 

set -e
outputPath=$1
files=$2

yamllist=()
erblist=()
otherlist=()
yamlstring=""
erbstring=""

# Read all input files
IFS=',' read -ra file <<< "$files"
for i in "${file[@]}"; do
  filename=$(basename -- "$i")
  extension="${filename##*.}"
  if [[ $extension == "yaml" ]]; then
    yamllist+=($i)
    if [[ $yamlstring == "" ]]; then
        yamlstring=$i
    else
        yamlstring="${yamlstring},${i}"
    fi
  elif [[ $extension == "erb" ]]; then
    erblist+=($i)
    if [[ $erbstring == "" ]]; then
        erbstring=$i
    else
        erbstring="${erbstring},${i}"
    fi
  else
    otherlist+=($i)
  fi
done

pushd mmv1

if [[ $yamlstring != "" ]]; then
  # run yaml conversion with given .yaml files
  bundle exec compiler.rb -e terraform -o $1 -v beta -a --go-yaml-files $yamlstring
  go run . --yaml-temp
  for i in `find . -name "*.temp" -type f`; do
    echo "removing go/ paths in ${i}"
    perl -pi -e 's/go\///g' $i
  done

fi


if [[ $erbstring != "" ]]; then
  # convert .erb files with given .erb files
  go run . --template-temp $erbstring
  go run . --handwritten-temp $erbstring
fi
popd

# add temporary file for all other files that do not need conversion
for i in "${otherlist[@]}"
do
    cp "$i" "${i}.temp" 
done