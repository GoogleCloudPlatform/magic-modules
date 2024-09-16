#!/usr/bin/env bash
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

# run yaml conversion with given .yaml files
bundle exec compiler.rb -e terraform -o $1 -v beta -a --go-yaml-files $yamlstring
go run . --yaml-temp

# convert .erb files with given .erb files
go run . --template-temp $erbstring
go run . --handwritten-temp $erbstring

# add temporary file for all other files that do not need conversion
for i in "${otherlist[@]}"
do
    cp "$i" "${i}.temp" 
done