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
require 'google/string_utils'

module Api
  # Represents a property type
  class Type < Api::Object::Named
    # The list of properties (attr_reader) that can be overridden in
    # <provider>.yaml.
    module Fields
      include Api::Object::Named::Properties

      attr_reader :description
      attr_reader :exclude

      attr_reader :output # If set value will not be sent to server on sync
      attr_reader :input # If set to true value is used only on creation
      attr_reader :field
      attr_reader :required
      attr_reader :update_verb
      attr_reader :update_url
    end

    include Fields

    attr_reader :__resource
    attr_reader :__parent # is nil for top-level properties

    MAX_NAME = 20

    def validate
      super
      @exclude ||= false

      check_property :description, ::String
      check_property :exclude, :boolean

      check_optional_property :output, :boolean
      check_optional_property :field, ::String
      check_optional_property :required, :boolean

      raise 'Property cannot be output and required at the same time.' \
        if @output && @required

      check_optional_property_oneof_default \
        :update_verb, %i[POST PUT PATCH NONE], @__resource&.update_verb, Symbol
      check_optional_property :update_url, ::String
    end

    def type
      self.class.name.split('::').last
    end

    def property_type
      property_ns_prefix.concat([type]).join('::')
    end

    def requires
      File.join(
        'google',
        @__resource.__product.prefix[1..-1],
        'property',
        type
      ).downcase
    end

    def field_name
      @field || @name
    end

    def parent
      @__parent
    end

    private

    # Shrinks a long composite type name into something that can barely be
    # read by humans.
    #
    # E.g.: Google::Compute::Property::AutoscalerCustomMetricUtilizationsArray
    #   --> Google::Compute::Property::Autos.....Custo.Metri.Utili.......Arr..
    #   --> Google::Compute::Property::AutosCustoMetriUtiliArr
    def shrink_type_name(type)
      name_parts = shrink_type_name_parts(type)

      # Isolate the Google common prefix
      name_parts = name_parts.drop(property_ns_prefix.size)
      num_parts = name_parts.flatten.size
      shrunk_names = recurse_shrink_name(name_parts,
                                         (1.0 * MAX_NAME / num_parts).round)
      type_name = Google::StringUtils.camelize(shrunk_names.flatten.join('_'),
                                               :upper)
      property_ns_prefix.concat([type_name])
    end

    def recurse_shrink_name(name, max_size)
      return name[0, max_size] unless name.is_a?(::Array)
      name.map { |part| recurse_shrink_name(part, max_size) }
    end

    def shrink_type_name_parts(type)
      type.map do |t|
        if t.is_a?(::Array)
          t.map { |u| Google::StringUtils.underscore(u).split('_') }
        else
          Google::StringUtils.underscore(t).split('_')
        end
      end
    end

    # A constant value to be provided as field
    class Constant < Type
      attr_reader :value

      def validate
        @description = "This is always #{value}."
        super
      end
    end

    # Represents a primitive (non-composite) type.
    class Primitive < Type
    end

    # Represents a boolean
    class Boolean < Primitive
    end

    # Represents an integer
    class Integer < Primitive
    end

    # Represents a double
    class Double < Primitive
    end

    # Represents a string
    class String < Primitive
      def initialize(name)
        @name = name
      end

      PROJECT = Api::Type::String.new('project')
      NAME = Api::Type::String.new('name')
    end

    # Represents a timestamp
    class Time < Primitive
    end

    # A base class to tag objects that are composed by other objects (arrays,
    # nested objects, etc)
    class Composite < Type
    end

    # Forwarding declaration to allow defining Array::NESTED_ARRAY_TYPE
    class NestedObject < Composite
    end

    # Forwarding declaration to allow defining Array::RREF_ARRAY_TYPE
    class ResourceRef < Type
    end

    # Represents an array, and stores its items' type
    class Array < Composite
      attr_reader :item_type
      attr_reader :max_size

      STRING_ARRAY_TYPE = [Api::Type::Array, Api::Type::String].freeze
      NESTED_ARRAY_TYPE = [Api::Type::Array, Api::Type::NestedObject].freeze
      RREF_ARRAY_TYPE = [Api::Type::Array, Api::Type::ResourceRef].freeze

      def validate
        super
        if @item_type.is_a?(NestedObject) || @item_type.is_a?(ResourceRef)
          @item_type.set_variable(@name, :__name)
          @item_type.set_variable(@__resource, :__resource)
          @item_type.set_variable(self, :__parent)
        end
        check_property :item_type, [::String, NestedObject, ResourceRef]
        unless @item_type.is_a?(NestedObject) || @item_type.is_a?(ResourceRef) \
          || type?(@item_type)
          raise "Invalid type #{@item_type}"
        end

        check_optional_property :max_size, ::Integer
      end

      def item_type_class
        return Api::Type::NestedObject if @item_type.is_a? NestedObject
        return Api::Type::ResourceRef if @item_type.is_a? ResourceRef
        get_type("Api::Type::#{@item_type}")
      end

      def property_class
        if @item_type.is_a?(NestedObject) || @item_type.is_a?(ResourceRef)
          type = @item_type.property_class
        else
          type = property_ns_prefix
          type << get_type(@item_type).new(@name).type
        end
        type = shrink_type_name(type)
        class_name = type.pop
        type << "#{class_name}Array"
      end

      def property_type
        property_class.join('::')
      end

      def property_file
        File.join(
          'google', @__resource.__product.prefix[1..-1], 'property',
          [get_type(@item_type).new(@name).type, 'array'].join('_')
        ).downcase
      end

      # Returns the file that implements this property
      def requires
        if @item_type.is_a?(NestedObject) || @item_type.is_a?(ResourceRef)
          return @item_type.requires
        end
        [property_file]
      end
    end

    # Represents an enum, and store is valid values
    class Enum < Primitive
      attr_reader :values

      def validate
        super
        check_property :values, ::Array
      end
    end

    # Properties that are fetched externally
    class FetchedExternal < Type
      attr_writer :resource
    end

    # Represents a 'selfLink' property, which returns the URI of the resource.
    class SelfLink < FetchedExternal
      EXPORT_KEY = 'selfLink'.freeze

      attr_reader :resource

      def name
        EXPORT_KEY
      end

      def out_name
        Google::StringUtils.underscore(EXPORT_KEY)
      end

      def field_name
        name
      end
    end

    # Represents a reference to another resource
    class ResourceRef < Type
      ALLOWED_WITHOUT_PROPERTY = [SelfLink::EXPORT_KEY].freeze

      attr_reader :resource
      attr_reader :imports

      def out_type
        resource_ref.out_name
      end

      def validate
        super
        @name = @resource if @name.nil?
        @description = "A reference to #{@resource} resource"

        return if @__resource.nil? || @__resource.exclude

        check_property :resource, ::String
        check_property :imports, ::String
        check_resource_ref_exists
        check_resource_ref_property_exists
      end

      def property
        props = resource_ref.exported_properties
                            .select { |prop| prop.name == @imports }
        return props.first unless props.empty?
        raise "#{@imports} does not exist on #{@resource}" if props.empty?
      end

      def resource_ref
        product = @__resource.__product
        resources = product.objects.select { |obj| obj.name == @resource }
        raise "Unknown item type '#{@resource}'" if resources.empty?
        resources[0]
      end

      def property_class
        type = property_ns_prefix
        type << [@resource, @imports, 'Ref']
        shrink_type_name(type)
      end

      def property_type
        property_class.join('::')
      end

      def property_file
        File.join('google', @__resource.__product.prefix[1..-1], 'property',
                  "#{resource}_#{@imports}").downcase
      end

      def requires
        [property_file]
      end

      private

      def check_resource_ref_exists
        product = @__resource.__product
        resources = product.objects.select { |obj| obj.name == @resource }
        raise "Missing '#{@resource}'" \
          if resources.empty? || resources[0].exclude
      end

      def check_resource_ref_property_exists
        exported_props = resource_ref.exported_properties
        raise "'#{@imports}' does not exist on '#{@resource}'" \
          if exported_props.none? { |p| p.name == @imports }
      end
    end

    # An structured object composed of other objects.
    class NestedObject < Composite
      # A custom getter is used for :properties instead of `attr_reader`

      def validate
        @description = 'A nested object resource' if @description.nil?
        @name = @__name if @name.nil?
        super
        @properties.each do |p|
          p.set_variable(@__resource, :__resource)
          p.set_variable(self, :__parent)
        end
        check_property_list :properties, Api::Type
      end

      def property_class
        type = property_ns_prefix
        type << [@__resource.name, @name]
        shrink_type_name(type)
      end

      def property_type
        property_class.join('::')
      end

      def property_file
        File.join(
          'google', @__resource.__product.prefix[1..-1], 'property',
          [@__resource.name, Google::StringUtils.underscore(@name)].join('_')
        ).downcase
      end

      def requires
        [property_file].concat(properties.map(&:requires))
      end

      # Returns all properties including the ones that are excluded
      # This is used for PropertyOverride validation
      def all_properties
        @properties
      end

      def properties
        @properties.reject(&:exclude)
      end
    end

    # Represents an array of name=value pairs, and stores its items' type
    class NameValues < Composite
      attr_reader :key_type
      attr_reader :value_type

      def validate
        check_property :key_type, ::String
        check_property :value_type, ::String
        raise "Invalid type #{@key_type}" unless type?(@key_type)
        raise "Invalid type #{@value_type}" unless type?(@value_type)
      end
    end

    def type?(type)
      !get_type(type).nil?
    end

    def get_type(type)
      Module.const_get(type)
    end

    def property_ns_prefix
      [
        'Google',
        Google::StringUtils.camelize(@__resource.__product.prefix[1..-1],
                                     :upper),
        'Property'
      ]
    end
  end
end
