# Copyright 2019 Google Inc.
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
  class Product < Api::Object::Named
    # A version of the API for a given product / API group
    # In GCP, different product versions are generally ordered where alpha is
    # a superset of beta, and beta a superset of GA. Each version will have a
    # different version url.
    class Version < Api::Object
      include Comparable

      attr_reader :base_url
      attr_reader :name

      ORDER = %w[ga beta alpha private].freeze

      def validate
        super
        check :base_url, type: String, required: true
        check :name, type: String, allowed: ORDER, required: true
      end

      def to_s
        str = "#{name}: #{base_url}"
        str
      end

      def <=>(other)
        ORDER.index(name) <=> ORDER.index(other.name) if other.is_a?(Version)
      end
    end
  end
end
