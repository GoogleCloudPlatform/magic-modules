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

# Ensure we can load 'google/' libraries
auth_libraries = File.expand_path('../../../resources', __dir__)
$LOAD_PATH.unshift(auth_libraries) unless $LOAD_PATH.include?(auth_libraries)

require 'google/authorization'

module Google
  # Module that holds gauth_credential_serviceaccount_for_function
  module Functions
    def gauth_credential_serviceaccount_for_function(path, scopes)
      ::Google::Authorization.new.for!(scopes).from_service_account_json!(path)
    end
  end
end
