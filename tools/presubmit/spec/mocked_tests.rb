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

# Mocked out tests for spec purposes.
module Presubmit
  # Mocked out class with a successful result.
  class SuccessfulTest < Test
    def run
      Presubmit::Results.new(self, 0, 'success')
    end
  end

  # Mocked out class with a failed result.
  class FailedTest < Test
    def run
      Presubmit::Results.new(self, 1, 'failed')
    end
  end
end
