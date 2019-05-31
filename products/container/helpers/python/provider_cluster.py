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

# Google Container Engine API has its own layout for the create method,
# defined like this:
#
# {
#   'cluster': {
#     ... cluster data
#   }
# }
#
# Format the request to match the expected input by the API
def encode_request(resource_request, module):
    return {
        'cluster': resource_request
    }

# Deletes the default node pool on default creation.
def delete_default_node_pool(module):
    auth = GcpSession(module, 'container')
    link = "https://container.googleapis.com/v1/projects/%s/locations/%s/clusters/%s/nodePools/default-pool" % \
        (module.params['project'], module.params['location'], module.params['name'])
    return wait_for_operation(module, auth.delete(link))
