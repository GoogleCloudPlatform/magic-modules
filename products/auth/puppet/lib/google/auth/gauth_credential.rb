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

require 'google/authorization'

module Google
  module Auth
    # A helper class to allocate credential for Google Cloud Platform access.
    class GAuthCredential
      def self.serviceaccount_for_function(credential_file, scopes)
        fn_name = :gauth_credential_serviceaccount_for_function
        Puppet::Parser::Compiler.new(Puppet::Node.new(:function))
                                .context_overrides[:global_scope]
                                .call_function(fn_name,
                                               [credential_file, scopes])
      end
    end
  end
end
