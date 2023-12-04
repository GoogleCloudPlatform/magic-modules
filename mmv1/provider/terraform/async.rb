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

module Provider
  class Terraform
    # Async implementation for polling in Terraform
    class PollAsync < Api::Async
      # Details how to poll for an eventually-consistent resource state.

      # Function to call for checking the Poll response for
      # creating and updating a resource
      attr_reader :check_response_func_existence

      # Function to call for checking the Poll response for
      # deleting a resource
      attr_reader :check_response_func_absence

      # Custom code to get a poll response, if needed.
      # Will default to same logic as Read() to get current resource
      attr_reader :custom_poll_read

      # If true, will suppress errors from polling and default to the
      # result of the final Read()
      attr_reader :suppress_error

      # Number of times the desired state has to occur continuously
      # during polling before returning a success
      attr_reader :target_occurrences

      def validate
        super

        check :check_response_func_existence, type: String, required: true
        check :check_response_func_absence, type: String,
                                            default: 'transport_tpg.PollCheckForAbsence'
        check :custom_poll_read, type: String
        check :suppress_error, type: :boolean, default: false
        check :target_occurrences, type: Integer, default: 1
      end
    end
  end
end
