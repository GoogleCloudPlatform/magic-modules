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

def encode_request(resource_request, module):
    """Structures the request as accountId + rest of request"""
    account_id = resource_request['name'].split('@')[0]
    del resource_request['name']
    return {
        'accountId': account_id,
        'serviceAccount': resource_request
    }


def decode_response(response, module):
    """Unstructures the request from accountId + rest of request"""
    if 'name' not in response:
        return response
    response['name'] = response['name'].split('/')[-1]
    return response

