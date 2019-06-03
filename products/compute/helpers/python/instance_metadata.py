# Copyright 2017 Google Inc.
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
# Mask the fact healthChecks array is actually a single object of type
# HttpHealthCheck.

# TODO(alexstephen): Implement updating metadata on existing resources.

# Expose instance 'metadata' as a simple name/value pair hash. However the API
# defines metadata as a NestedObject with the following layout:
#
# metadata {
#   fingerprint: 'hash-of-last-metadata'
#   items: [
#     {
#       key: 'metadata1-key'
#       value: 'metadata1-value'
#     },
#     ...
#   ]
# }
#
def metadata_encoder(metadata):
    metadata_new = []
    for key in metadata:
        value = metadata[key]
        metadata_new.append({
            "key": key,
            "value": value
        })
    return {
        'items': metadata_new
    }


# Map metadata.items[]{key:,value:} => metadata[key]=value
def metadata_decoder(metadata):
    items = {}
    if 'items' in metadata:
        metadata_items = metadata['items']
        for item in metadata_items:
            items[item['key']] = item['value']
    return items
