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
def self.encode_request(request)
  unless request[:healthCheck].nil?
    # Map one allowed health check into array
    request[:healthChecks] = [request[:healthCheck]]
    request.delete(:healthCheck)
  end
  request
end

def encode_request(resource_request)
  self.class.encode_request(resource_request)
end

# Mask healthChecks into a single element.
# @see self.encode_request for details
def self.decode_request(response, kind)
  response = JSON.parse(response.body)

  return response unless kind == 'compute#targetPool'

  # Map healthChecks[0] => healthCheck
  unless response['healthChecks'].nil?
    response['healthCheck'] = response['healthChecks'][0] \
      unless response['healthChecks'].empty?
    response.delete('healthChecks')
  end

  response
end
