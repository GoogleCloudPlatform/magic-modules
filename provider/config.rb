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
require 'compile/core'
require 'provider/resource_override'
require 'provider/resource_overrides'

module Provider
  # Settings for the provider
  class Config < Api::Object
    include Compile::Core
    extend Compile::Core

    attr_reader :overrides
    # Overrides for datasources
    attr_reader :datasources
    attr_reader :properties # TODO(nelsonjr): Remove this once bug 193 is fixed.
    attr_reader :tests
    attr_reader :files
    attr_reader :changelog
    # Product names are complicated in MagicModules.  They are given by
    # product.prefix, which is in the format 'g<nameofproduct>', e.g.
    # gcompute or gresourcemanager.  This is munged in many places.
    # Some examples:
    #   - prefix[1:-1] ('compute' / 'resourcemanager') for the
    #     directory to fetch chef / puppet examples.
    #   - camelCase(prefix[1:-1]) for resource namespaces.
    #   - TitleCase(prefix[1:-1]) for resource names in terraform.
    #   - prefix[1:-1] again, for working with libraries directly.
    # This override does not change any of those inner workings, but
    # instead is passed directly to the template as `product_ns` if
    # set.  Otherwise, the normal logic applies.
    attr_reader :name

    # List of files to copy or compile into target module
    class Files < Api::Object
      attr_reader :compile
      attr_reader :copy

      def validate
        super
        check_optional_property :compile, Hash
        check_optional_property :copy, Hash
      end
    end

    # Identifies all changes releted to a release of the compiled artifact.
    class Changelog < Api::Object
      attr_reader :version
      attr_reader :date
      attr_reader :general
      attr_reader :features
      attr_reader :fixes

      def validate
        super
        check_property :version, String
        check_property :date, Time
        check_optional_property :general, String
        check_property_list :features, String
        check_property_list :fixes, String

        raise "Required general/features/fixes for change #{@version}." \
          if @general.nil? && @features.nil? && @fixes.nil?
      end
    end

    def self.parse(cfg_file, api = nil, _version_name = nil)
      # Compile step #1: compile with generic class to instantiate target class
      source = compile(cfg_file)
      config = Google::YamlValidator.parse(source)
      raise "Config #{cfg_file}(#{config.class}) is not a Provider::Config" \
        unless config.class <= Provider::Config
      # Config must be validated so items are properly setup for next compile
      config.validate
      # Compile step #2: Now that we have the target class, compile with that
      # class features
      source = config.compile(cfg_file)
      config = Google::YamlValidator.parse(source)
      config.default_overrides
      config.spread_api config, api, [], '' unless api.nil?
      config.validate
      config.cfg_file = cfg_file
      config
    end

    def provider
      raise "#{self.class}#provider not implemented"
    end

    def self.next_version(version)
      [Gem::Version.new(version).bump, 0].join('.')
    end

    def validate
      super

      default_overrides

      check_optional_property :files, Provider::Config::Files
      check_property :overrides, Provider::ResourceOverrides
      check_property_list :changelog, Provider::Config::Changelog \
        unless @changelog.nil?

      @datasources.__is_data_source = true unless @datasources.nil?
    end

    # Provides the API object to any type that requires, e.g. for validation
    # purposes, such as Api::Resource::HashArray which enforces that the keys
    # are necessarily objects defined in the API.
    def spread_api(object, api, visited, indent)
      object.instance_variables.each do |var|
        var_value = object.instance_variable_get(var)
        next if visited.include?(var_value)
        visited << var_value
        var_value.consume_api api if var_value.respond_to?(:consume_api)
        var_value.consume_config api, self \
          if var_value.respond_to?(:consume_config)
        spread_api(var_value, api, visited, indent)
      end
    end

    # TODO(nelsonjr): Investigate why we need to call default_overrides twice.
    def default_overrides
      @overrides ||= Provider::ResourceOverrides.new
    end
  end
end
