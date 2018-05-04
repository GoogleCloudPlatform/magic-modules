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
require 'google/logger'
require 'compile/core'

module Api
  # Repesents a product to be managed
  class Product < Api::Object::Named
    attr_reader :objects
    attr_reader :prefix
    attr_reader :scopes
    attr_reader :versions

    include Compile::Core

    def validate
      super
      set_variables @objects, :__product
      check_property :objects, Array
      check_property_list :objects, Api::Resource
      check_property :prefix, String
      check_property :scopes, ::Array
      check_property_list :scopes, String

      check_versions
    end

    # Represents a version of the API for this product
    class Version < Api::Object
      attr_reader :base_url
      attr_reader :default
      attr_reader :name

      def validate
        super
        check_property :base_url, String
        check_property :name, String
        check_optional_property :default, [TrueClass, FalseClass]
      end
    end

    def default_version
      @versions.each do |v|
        return v if v.default
      end
    end

    private

    def check_versions
      check_property :versions, Array
      check_property_list :versions, Api::Product::Version

      # Confirm that at exactly one version is the default
      defaults = 0
      @versions.each do |v|
        defaults += 1 if v.default
      end

      raise "Product '#{@name}' must specify exactly one default API version" \
        if defaults != 1
    end
  end
end
