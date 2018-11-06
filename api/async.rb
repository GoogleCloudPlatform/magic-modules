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

require 'api/object'

module Api
  # Represents an asynchronous operation definition
  class Async < Api::Object
    attr_reader :operation
    attr_reader :result
    attr_reader :status
    attr_reader :error

    def validate
      super

      check_property :operation, Operation
      check_property :result, Result
      check_property :status, Status
      check_property :error, Error
    end

    # Represents the operations (requests) issues to watch for completion
    class Operation < Api::Object
      attr_reader :kind
      attr_reader :path
      attr_reader :base_url
      attr_reader :wait_ms
      attr_reader :timeouts

      def validate
        super

        @timeouts ||= Timeouts.new

        check_property :kind, String
        check_property :path, String
        check_property :base_url, String
        check_property :wait_ms, Integer
        check_property :timeouts, Timeouts
      end

      # Provides timeout information for the different operation types
      class Timeouts < Api::Object
        # Default timeout for all operation types is 4 minutes. This can be
        # overridden for each resource.
        DEFAULT_INSERT_TIMEOUT_SEC = 4 * 60
        DEFAULT_UPDATE_TIMEOUT_SEC = 4 * 60
        DEFAULT_DELETE_TIMEOUT_SEC = 4 * 60

        attr_reader :insert_sec
        attr_reader :update_sec
        attr_reader :delete_sec

        def validate
          super

          @insert_sec ||= DEFAULT_INSERT_TIMEOUT_SEC
          @update_sec ||= DEFAULT_UPDATE_TIMEOUT_SEC
          @delete_sec ||= DEFAULT_DELETE_TIMEOUT_SEC

          check_property :insert_sec, Integer
          check_property :update_sec, Integer
          check_property :delete_sec, Integer
        end
      end
    end

    # Represents the results of an Operation request
    class Result < Api::Object
      attr_reader :path
      attr_reader :resource_inside_response

      def validate
        super
        default_value_property :resource_inside_response, false

        check_optional_property :path, String
        check_optional_property :resource_inside_response, :boolean
      end
    end

    # Provides information to parse the result response to check operation
    # status
    class Status < Api::Object
      attr_reader :path
      attr_reader :complete
      attr_reader :allowed

      def validate
        super
        check_property :path, String
        check_property :allowed, Array
      end
    end

    # Provides information on how to retrieve errors of the executed operations
    class Error < Api::Object
      attr_reader :path
      attr_reader :message

      def validate
        super
        check_property :path, String
        check_property :message, String
      end
    end
  end
end
