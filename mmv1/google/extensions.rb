# Copyright 2018 Google Inc.
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
require 'google/string_utils'

# The Ruby String class. It's being extended to add StringUtils methods.
class String
  def underscore
    Google::StringUtils.underscore(self)
  end

  def first_sentence
    Google::StringUtils.first_sentence(self)
  end

  def space_separated
    Google::StringUtils.space_separated(self)
  end

  def plural
    Google::StringUtils.plural(self)
  end

  def title
    Google::StringUtils.title(self)
  end

  def camelize(first_letter = :upper)
    case first_letter
    when :upper
      Google::StringUtils.camelize(self, true)
    when :lower
      Google::StringUtils.camelize(self, false)
    else
      raise ArgumentError, 'Invalid option, use either :upper or :lower.'
    end
  end
end
