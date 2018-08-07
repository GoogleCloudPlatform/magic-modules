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

require 'puppet'
require 'puppet/pops'
require 'puppet/functions'

require 'google/authorization'

Puppet::Functions.create_function(
  :gauth_credential_serviceaccount_for_function
) do
  dispatch :gauth_credential_serviceaccount_for_function do
    param 'String', :path
    param 'Array', :scopes
  end

  def gauth_credential_serviceaccount_for_function(path, scopes)
    Google::Authorization.new
                         .for!(scopes)
                         .from_service_account_json!(path)
  end
end
