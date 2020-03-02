# Copyright 2020 Google Inc.
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
require 'api/timeout'

module Api
  # Base class from which other Async classes can inherit.
  class Async < Api::Object
    # Describes an operation
    attr_reader :operation

    # The list of methods where operations are used.
    attr_reader :actions

    def validate
      super

      check :operation, type: Operation, required: true
      check :actions, default: %w[create delete update], type: ::Array, item_type: ::String
    end

    def allow?(method)
      @actions.include?(method.downcase)
    end

    # Base async operation type
    class Operation < Api::Object
      # Contains information about an long-running operation, to make
      # requests for the state of an operation.
      attr_reader :timeouts
      attr_reader :result

      def validate
        check :result, type: Result
        check :timeouts, type: Api::Timeouts
      end
    end

    # Base result class
    class Result < Api::Object
      # Contains information about the result of an Operation

      attr_reader :resource_inside_response

      def validate
        super
        check :resource_inside_response, type: :boolean, default: false
      end
    end
  end

  # Represents an asynchronous operation definition
  class OpAsync < Async
    attr_reader :operation
    attr_reader :result
    attr_reader :status
    attr_reader :error

    # The list of methods where operations are used.
    attr_reader :actions

    def validate
      super

      check :operation, type: Operation, required: true
      check :result, type: Result, required: true
      check :status, type: Status, required: true
      check :error, type: Error, required: true
      check :actions, default: %w[create delete update], type: ::Array, item_type: ::String
    end

    # The main implementation of Operation,
    # corresponding to common GCP Operation resources.
    class Operation < Async::Operation
      attr_reader :kind
      attr_reader :path
      attr_reader :base_url
      attr_reader :wait_ms
      attr_reader :timeouts

      # Use this if the resource includes the full operation url.
      attr_reader :full_url

      def validate
        super

        check :kind, type: String
        check :path, type: String, required: true
        check :base_url, type: String
        check :wait_ms, type: Integer, required: true

        check :full_url, type: String

        conflicts %i[base_url full_url]
      end
    end

    # Represents the results of an Operation request
    class Result < Async::Result
      attr_reader :path

      def validate
        super

        check :path, type: String
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
        check :path, type: String, required: true
        check :allowed, type: Array, item_type: [::String, :boolean], required: true
      end
    end

    # Provides information on how to retrieve errors of the executed operations
    class Error < Api::Object
      attr_reader :path
      attr_reader :message

      def validate
        super
        check :path, type: String, required: true
        check :message, type: String, required: true
      end
    end
  end
end
