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

def self.encode_request(request)
  metadata_encoder(request[:metadata]) unless request[:metadata].nil?
  request
end

def encode_request(resource_request)
  self.class.encode_request(resource_request)
end

def self.decode_response(response, kind)
  response = JSON.parse(response.body)
  return response unless kind == 'compute#instance'

  metadata_decoder(response['metadata']) unless response['metadata'].nil?
  response
end
