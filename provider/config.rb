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
require 'provider/compiler'
require 'compile/core'

module Provider
  # Settings for the provider
  class Config < Api::Object
    include Compile::Core
    extend Compile::Core

    attr_reader :objects
    attr_reader :examples
    attr_reader :properties # TODO(nelsonjr): Remove this once bug 193 is fixed.
    attr_reader :tests
    attr_reader :test_data
    attr_reader :files
    attr_reader :style
    attr_reader :changelog
    attr_reader :functions

    # A custom client side function provided by the module.
    class Function < Api::Object::Named
      attr_reader :description
      attr_reader :arguments
      attr_reader :requires
      attr_reader :code
      attr_reader :helpers
      attr_reader :examples
      attr_reader :notes

      def validate
        super
        check_property_list :requires, @requires, String
        check_property :code, String
        check_property_list :arguments, @arguments,
                            Provider::Config::Function::Argument
        check_property :helpers, String unless @helpers.nil?
      end

      # An argument required by the function being provided by the module.
      class Argument < Api::Object::Named
        attr_reader :description
        attr_reader :type

        def validate
          super
          check_property :description, String
          check_property :type, String
        end
      end
    end

    # Operating system supported by the module
    class OperatingSystem < Api::Object::Named
      attr_reader :versions

      def validate
        super
        check_property :versions
      end

      def all_versions
        [@name, @versions.join(', ')].join(' ')
      end
    end

    # Reference to a module required by the module
    class Requirements < Api::Object::Named
      attr_reader :versions

      def self.create(name, versions)
        Requirements.new(name, versions)
      end

      def validate
        super
        check_property :versions
      end

      private

      def initialize(name, versions)
        @name = name
        @versions = versions
      end
    end

    # Reference to a module required by the module
    class TestData < Api::Object
      attr_reader :network

      # A dummy class that identifies the property as deliberately unused.
      class NONE < TestData
        include Api::Object::MissingObject

        def validate
          @network = Api::Resource::HashArray.new
          super
        end
      end

      def validate
        super
        check_property :network, Api::Resource::HashArray
      end
    end

    # List of files to copy or compile into target module
    class Files < Api::Object
      attr_reader :compile
      attr_reader :copy
      attr_reader :permissions

      # A dummy class that identifies the property as deliberately unused.
      class NONE < Files
        include Api::Object::MissingObject
      end

      def validate
        super
        check_property :compile, Hash unless @compile.nil?
        check_property :copy, Hash unless @copy.nil?
        check_property_list \
          :permissions, @permissions, Provider::Config::Permission
      end
    end

    # Represents a permission to be set at the generated module
    class Permission < Api::Object
      attr_reader :path
      attr_reader :acl

      def validate
        super
        check_property :path, String
        check_property :acl, String
      end
    end

    # Identifies a location where a code style exception happened. This is used
    # to guide the compiler to produce linter correct code, i.e. adding the
    # necessary guards to avoid violations.
    class StyleException < Api::Object::Named
      attr_reader :pinpoints

      def validate
        super
        check_property :pinpoints, Array
        check_property_list :pinpoints, @pinpoints, Hash
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
        check_property :general, String unless @general.nil?
        check_property_list :features, @features, String
        check_property_list :fixes, @fixes, String

        raise "Required general/features/fixes for change #{@version}." \
          if @general.nil? && @features.nil? && @fixes.nil?
      end
    end

    def self.parse(cfg_file, api = nil)
      # Compile step #1: compile with generic class to instantiate target class
      source = compile(cfg_file)
      config = Google::YamlValidator.parse(source)
      raise "Config #{cfg_file}(#{config.class}) is not a Provider::Config" \
        unless config.class <= Provider::Config
      # Compile step #2: Now that we have the target class, compile with that
      # class features
      source = config.compile(cfg_file, 0)
      config = Google::YamlValidator.parse(source)
      config.spread_api config, api, [], '' unless api.nil?
      config.validate
      config
    end

    def provider
      raise "#{self.class}#provider not implemented"
    end

    def validate
      super
      check_property :examples, Api::Resource::HashArray
      check_property :files, Provider::Config::Files
      check_property :objects, Api::Resource::HashArray
      check_property :test_data, Provider::Config::TestData
      check_property :tests, Api::Resource::HashArray unless @tests.nil?
      check_property_list :style, @style, Provider::Config::StyleException
      check_property_list :changelog, @changelog, Provider::Config::Changelog
      check_property_list :functions, @functions, Provider::Config::Function
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
        spread_api(var_value, api, visited, indent)
      end
    end
  end
end
