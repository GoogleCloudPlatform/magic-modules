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

module Presubmit
  # An interface class for a Submodule
  class Module
    def add_tests(tests)
      @tests = tests
    end

    # Creates modules with properly initialized tests.
    # This is where per-module and per-product test differences are made.
    class Factory
      def initialize(testers)
        @testers = testers
      end

      def create
        raise 'children must implement'
      end
    end

    # Run all tests and return the results.
    def run
      Hash[
        @tests.map do |test|
          [test, test.run]
        end
      ]
    end
  end
end
