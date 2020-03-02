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

require 'api/async'
require 'provider/abstract_core'

module Provider
  class Terraform < Provider::AbstractCore
    # Async implementation for polling in Terraform
    class PollAsync < Api::Async
      # Details how to poll for an eventually-consistent resource state.

      # Function to call for checking the Poll response
      attr_reader :check_response_func

      # Custom code to get a poll response, if needed.
      # Will default to same logic as Read() to get current resource
      attr_reader :custom_poll_read

      # If true, will suppress errors from polling and default to the
      # result of the final Read()
      attr_reader :suppress_error

      def validate
        super

        check :check_response_func, type: String, required: true
        check :custom_poll_read, type: String
        check :suppress_error, type: :boolean, default: false
      end
    end
  end
end
