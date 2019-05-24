# frozen_string_literal: true

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

require 'json'
require 'vcr'

VCR.configure do |c|
  c.hook_into :webmock
  c.cassette_library_dir = 'inspec-cassettes'
  c.allow_http_connections_when_no_cassette = true

  c.before_record do |i|
    i.request.headers.delete_if { true }
    if auth_call?(i)
      i.request.body = 'AUTH REQUEST'
      i.response.body = "{\n  \"access_token\": \"ya29.c.samsamsamsamsamsamsamsamsa-thisisnintysixcharactersoftexttolooklikeanauthtokenthisisnintysixcharactersoftexttolooklikeanaut\",\n  \"expires_in\": 3600,\n  \"token_type\": \"Bearer\"\n}"
    end
  end
end

def auth_call?(interaction)
  # Auth calls require extra scrubbing, this method is very broad, this is intentional
  interaction.request.uri.include? 'oauth2'
end
