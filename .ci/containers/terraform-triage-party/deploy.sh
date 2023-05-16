#!/bin/bash
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eu

# Need to set GITHUB_TOKEN environment variable. Use one from the magician.
if [ -z "$GITHUB_TOKEN" ]; then
    echo "Did not provide GITHUB_TOKEN environment variable."
    exit 1
fi

set -x

export PROJECT=terraform-triage
export IMAGE=gcr.io/terraform-triage/triage-party
export SERVICE_NAME=terraform-gcp-triage
export CONFIG_FILE=config.yaml
export TP_VERSION="v1.3.0"

git clone --branch $TP_VERSION --depth 1 https://github.com/google/triage-party
cp $CONFIG_FILE triage-party/config/config.yaml

docker build triage-party/ -t "${IMAGE}"

rm -rf triage-party/

docker push "${IMAGE}" || exit 2

# TODO: add persistence to make it run faster: env vars PERSIST_BACKEND, PERSIST_PATH
gcloud beta run deploy "${SERVICE_NAME}" \
    --project "${PROJECT}" \
    --image "${IMAGE}" \
    --set-env-vars="GITHUB_TOKEN=${GITHUB_TOKEN}" \
    --allow-unauthenticated \
    --region us-central1 \
    --memory 384Mi \
    --platform managed
