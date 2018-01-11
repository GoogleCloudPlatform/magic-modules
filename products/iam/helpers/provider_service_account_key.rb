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

# Format the response to match Puppet's expectations
def self.decode_response(response)
  response = JSON.parse(response.body)
  if response.key? 'privateKeyData'
    require 'base64'
    require 'byebug'
    @resource[:file]
    Base64.decode(response['privateKeyData'])
  end
  response
end

def decode_response(response)
  self.class.decode_response(response)
end
