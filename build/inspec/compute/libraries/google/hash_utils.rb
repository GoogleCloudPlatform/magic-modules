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

require 'google/string_utils'

module Google
  # Helper class to process and mutate hashes.
  class HashUtils
    # Converts all keys to symbols
    def self.camelize_keys(source)
      result = source.clone
      # rubocop:disable Performance/HashEachMethods
      result.keys.each do |k|
        result[Google::StringUtils.camelize(k.to_s)] = result.delete(k)
      end
      # rubocop:enable Performance/HashEachMethods
      result
    end

    def self.symbolize_keys(source)
      result = source.clone
      # rubocop:disable Performance/HashEachMethods
      result.keys.each do |k|
        result[Google::StringUtils.symbolize(k)] = result.delete(k)
      end
      # rubocop:enable Performance/HashEachMethods
      result
    end

    # Allows fetching objects within a tree path.
    def self.navigate(source, path, default = nil)
      key = path.take(1)[0]
      path = path.drop(1)
      return default unless source.key?(key)
      result = source.fetch(key)
      return HashUtils.navigate(result, path, default) unless path.empty?
      return result if path.empty?
    end

    # Converts a path in the form a/b/c/d into %w(a b c d)
    def self.path2navigate(path)
      "%w[#{path.split('/').join(' ')}]"
    end
  end
end
