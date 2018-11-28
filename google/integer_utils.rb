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

module Google
  # A helper to convert integer to 1_000_222_333 Ruby underscore notation.
  class IntegerUtils
    def self.underscore(value)
      return '0' if value.zero?

      result = []
      while value.positive?
        value, part = value.divmod(1000)
        result << format('%03d', part) if value.positive?
        result << part if value.zero?
      end
      result.reverse.join('_')
    end
  end
end
