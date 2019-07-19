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
#
# Google Compute Engine API defines healthChecks as a list but it can only
# take [0, 1] elements. To make it simpler to declare we'll map that to a
# single object and encode/decode as appropriate.
def encode_request(request, module):
    if 'healthCheck' in request:
        request['healthChecks'] = [request['healthCheck']]
        del request['healthCheck']
    return request


# Mask healthChecks into a single element.
# @see encode_request for details
def decode_response(response, module):
    if response['kind'] != 'compute#targetPool':
        return response

    # Map healthChecks[0] => healthCheck
    if 'healthChecks' in response:
        if not response['healthChecks']:
            response['healthCheck'] = response['healthChecks'][0]
            del response['healthChecks']

    return response
