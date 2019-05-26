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
require 'api/product/api_reference'
require 'api/product/version'
require 'google/logger'
require 'compile/core'
require 'json'

module Api
  # Repesents a product to be managed
  class Product < Api::Object::Named
    include Compile::Core

    # Inherited:
    # The name of the product's API capitalised in the appropriate places.
    # This isn't just the API name because it doesn't meaningfully separate
    # words in the api name - "accesscontextmanager" vs "AccessContextManager"
    # Example inputs: "Compute", "AccessContextManager"
    # attr_reader :name

    # The full name of the GCP product; eg "Cloud Bigtable"
    attr_reader :display_name

    attr_reader :objects

    # The list of permission scopes available for the service
    # For example: `https://www.googleapis.com/auth/compute`
    attr_reader :scopes

    # The API versions of this product
    attr_reader :versions

    # The base URL for the service API endpoint
    # For example: `https://www.googleapis.com/compute/v1/`
    attr_reader :base_url

    # The APIs required to be enabled for this product.
    # Usually just the product's API
    attr_reader :apis_required

    attr_reader :async

    def validate
      super
      set_variables @objects, :__product
      check :display_name, type: String
      check :objects, type: Array, item_type: Api::Resource, required: true
      check :scopes, type: Array, item_type: String, required: true
      check :apis_required, type: Array, item_type: Api::Product::ApiReference

      check :async, type: Api::Async

      validate_versions
    end

    # ====================
    # Custom Getters
    # ====================

    # The name of the product's API; "compute", "accesscontextmanager"
    def api_name
      name.downcase
    end

    # The product full name is the "display name" in string form intended for
    # users to read in documentation; "Google Compute Engine", "Cloud Bigtable"
    def product_full_name
      if !display_name.nil?
        display_name
      else
        name.underscore.humanize
      end
    end

    def default_version
      @versions.each do |v|
        return v if v.default
      end

      return @versions.last if @versions.length == 1
    end

    def version_obj(name)
      @versions.each do |v|
        return v if v.name == name
      end

      raise "API version '#{name}' does not exist for product '#{@name}'"
    end

    def version_obj_or_default(name)
      exists_at_version(name) ? version_obj(name) : default_version
    end

    def exists_at_version_or_lower(name)
      # Versions aren't normally going to be empty since products need a
      # base_url. This nil check exists for atypical products, like _bundle.
      return true if @versions.nil?

      name ||= Version::ORDER[0]
      return false unless Version::ORDER.include?(name)

      (0..Version::ORDER.index(name)).each do |i|
        return true if exists_at_version(Version::ORDER[i])
      end
      false
    end

    def exists_at_version(name)
      # Versions aren't normally going to be empty since products need a
      # base_url. This nil check exists for atypical products, like _bundle.
      return true if @versions.nil?

      @versions.any? { |v| v.name == name }
    end

    # Not a conventional setter, so ignore rubocop's warning
    # rubocop:disable Naming/AccessorMethodName
    def set_properties_based_on_version(version)
      @base_url = version.base_url
    end
    # rubocop:enable Naming/AccessorMethodName

    # ====================
    # Debugging Methods
    # ====================

    def to_s
      # relies on the custom to_json definitions
      JSON.pretty_generate(self)
    end

    # Prints a dot notation path to where the field is nested within the parent
    # object when called on a property. eg: parent.meta.label.foo
    # Redefined on Product to terminate the calls up the parent chain.
    def lineage
      name
    end

    def to_json(opts = nil)
      json_out = {}

      instance_variables.each do |v|
        if v == :@objects
          json_out['@resources'] = objects.map { |o| [o.name, o] }.to_h
        elsif instance_variable_get(v) == false || instance_variable_get(v).nil?
          # ignore false or missing because omitting them cleans up result
          # and both are the effective defaults of their types
        else
          json_out[v] = instance_variable_get(v)
        end
      end

      JSON.generate(json_out, opts)
    end

    private

    def validate_versions
      check :versions, type: Array, item_type: Api::Product::Version, required: true

      # Confirm that at most one version is the default
      defaults = 0
      @versions.each do |v|
        defaults += 1 if v.default
      end

      raise "Product '#{@name}' must specify at most one default API version" \
        if defaults > 1

      raise "Product '#{@name}' must specify a default API version" \
        if defaults.zero? && @versions.length > 1
    end
  end
end
