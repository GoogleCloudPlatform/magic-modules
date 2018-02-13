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
  # An interface class for a Presubmit test.
  class Test
    def initialize(mod)
      @mod = mod
    end

    # This will run a test and return its results.
    def run
      raise 'This must be implemented by the test'
    end
  end

  # A class that stores the results of a test.
  class Results
    attr_reader :tester
    attr_reader :status
    attr_reader :output

    def initialize(tester, status, output)
      @tester = tester
      @status = status
      @output = output
    end

    def success?
      @status.zero?
    end

    def failed?
      !success?
    end

    def warning?
      false
    end
  end
end
