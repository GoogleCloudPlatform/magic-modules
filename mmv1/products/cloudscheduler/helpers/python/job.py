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
# Verify that firewall network names match what the API returns

def encode_request(request, module):
    request['name'] = "projects/%s/locations/%s/jobs/%s" %(module.params['project'], module.params['region'], module.params['name'])
    return request

def decode_request(response, module):
    if 'name' in response:
        response['name'] = response['name'].split('/')[-1]

    return response
