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

# Format the request to match the expected input by the API
def self.encode_request(resource_request)
  account_id = resource_request[:name].split('@').first
  resource_request.delete(:name)
  {
    'accountId' => account_id,
    'serviceAccount' => resource_request
  }
end

def encode_request(resource_request)
  self.class.encode_request(resource_request)
end

# Format the response to match Puppet's expectations
def self.decode_response(response)
  response = JSON.parse(response.body)
  return response unless response.key? 'name'
  response['name'] = response['name'].split('/').last
  response
end

def decode_response(response)
  self.class.decode_response(response)
end
