#!/bin/sh

cat <<EOH
# Copyright 2018 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

---
EOH

# Fetch list of all Google APIs
API_LIST=$(curl -s https://www.googleapis.com/discovery/v1/apis)

# Load list of excluded APIs (focus on Google Cloud)
EXCLUDED_APIS=$(cat excluded-apis.txt)

# Build list of API product discovery URLs
URL_LIST=$(echo "$API_LIST" | grep discoveryRestUrl | awk '{print $2}')

for url in $URL_LIST
do
  product=$(echo "$url" | sed 's|"https://||' | awk -F\. '{print $1}')
  version=$(echo "$url" | awk -F= '{print $2}')
  if [ "$product" = "www" ]; then
    product=$(echo "$url" | sed 's|"https://www.googleapis.com/discovery/v1/apis/||' | awk -F/ '{print $1}')
    version=$(echo "$url" | sed 's|"https://www.googleapis.com/discovery/v1/apis/||' | awk -F/ '{print $2}')
  fi
  version=$(echo "$version" | sed 's|"||g' | sed 's|,||g')

  # Skip over anything in the excluded list
  exclude=0
  for skipped in $EXCLUDED_APIS
  do
    if [ "$product" = "$skipped" ]; then
      exclude=1
    fi
  done

  if [ "$exclude" = 0 ]; then
    url=$(echo "$url" | sed 's|"||g' | sed 's|,||g')
    echo "- url: $url"
    echo "  product: $product"
    echo "  version: $version"
    echo "  filename: products/$product/api.yaml"
  fi
done
