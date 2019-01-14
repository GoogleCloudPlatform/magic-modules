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
require 'overrides/runner'

module Provider
  # Settings for the provider
  class Config < Api::Object
    include Compile::Core
    extend Compile::Core

    # Overrides for datasources
    attr_reader :datasources
    attr_reader :files

    # Some tool-specific names may be in use, and they won't all match;
    # For Terraform, some products use the API client name w/o spaces and
    # others use spaces. Eg: "app_engine" vs "appengine".
    attr_reader :legacy_name

    attr_reader :overrides

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

    def self.parse(cfg_file, api = nil, version_name = 'ga')
      raise 'Version passed to the compiler cannot be nil' if version_name.nil?

      # Compile step #1: compile with generic class to instantiate target class
      source = compile(cfg_file)
      config = Google::YamlValidator.parse(source)
      raise "Config #{cfg_file}(#{config.class}) is not a Provider::Config" \
        unless config.class <= Provider::Config

      config.validate
      api = Overrides::Runner.build(api, config.overrides,
                                    config.resource_override,
                                    config.property_override)
      config.spread_api config, api, [], '' unless api.nil?
      config.validate
      api.validate
      [api, config]
    end

    def provider
      raise "#{self.class}#provider not implemented"
    end

    def self.next_version(version)
      [Gem::Version.new(version).bump, 0].join('.')
    end

    def validate
      super

      default_value_property :overrides, Overrides::ResourceOverrides.new

      check_optional_property :files, Provider::Config::Files
      check_property :overrides, Overrides::ResourceOverrides
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

    def resource_override
      Overrides::ResourceOverride
    end

    def property_override
      Overrides::PropertyOverride
    end
  end
end
