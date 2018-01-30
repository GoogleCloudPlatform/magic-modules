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
  # An object available in the product
  # rubocop:disable Metrics/ClassLength
  class Resource < Api::Object::Named
    attr_reader :description
    attr_reader :kind
    attr_reader :base_url
    attr_reader :self_link
    attr_reader :self_link_query
    # identity: an array with items that uniquely identify the resource.
    # default=name
    attr_reader :identity
    attr_reader :parameters
    attr_reader :properties
    attr_reader :exclude
    attr_reader :virtual
    attr_reader :async
    attr_reader :readonly
    attr_reader :exports
    attr_reader :label_override
    attr_reader :transport
    attr_reader :references
    attr_reader :create_verb

    attr_reader :__product

    # Allows overriding snowflake transport requests
    class Transport < Api::Object
      attr_reader :encoder
      attr_reader :decoder

      def validate
        super
        check_property :encoder, ::String unless @encoder.nil?
        check_property :decoder, ::String unless @decoder.nil?
      end
    end

    # Allows mapping of requests to specific API layout quirks.
    class Wrappers < Api::Object
      attr_reader :create

      def validate
        super
        check_property :create, ::String
      end
    end

    # Represents a response from the API that returns a list of objects.
    class ResponseList < Api::Object
      attr_reader :kind
      attr_reader :items

      def validate
        super
        check_property :kind, String
        check_property :items, String
      end

      def kind?
        !@kind.nil?
      end
    end

    # Represents a list of documentation links.
    class ReferenceLinks < Api::Object
      attr_reader :guides
      attr_reader :api

      def validate
        super
        check_property :guides, Hash unless @guide.nil?
        check_property :api, String unless @guide.nil?
      end
    end

    # Represents a hierarchy that has an object as its key. For example, when
    # creating test data, we'll do it per type, so it would look like this in
    # the provider.yaml file:
    #
    # test_data: !ruby/object:Api::Resource::HashArray
    #   Object1:
    #     - data1
    #     - data2
    #   Object2:
    #     - data3
    #     - data4
    class HashArray < Api::Object
      def consume_api(api)
        @__api = api
      end

      def validate
        return unless @__objects.nil? # allows idempotency of calling validate
        convert_findings_to_hash
        ensure_keys_are_objects unless @__api.nil?
        super
      end

      def [](index)
        @__objects[index]
      end

      def each
        return enum_for(:each) unless block_given?
        @__objects.each { |o| yield o }
        self
      end

      def select
        return enum_for(:select) unless block_given?
        @__objects.select { |o| yield o }
        self
      end

      def fetch(key, *args)
        # *args only holds default value. Needs to mimic ::Hash
        if args.empty?
          # KeyErorr will be thrown if key not found
          @__objects&.fetch(key)
        else
          # args[0] will be returned if key not found
          @__objects&.fetch(key, args[0])
        end
      end

      def key?(key)
        @__objects&.key?(key)
      end

      private

      # Converts every variable into @__objects
      def convert_findings_to_hash
        @__objects = {}
        instance_variables.each do |var|
          next if var.id2name.start_with?('@__')
          @__objects[var.id2name[1..-1]] = instance_variable_get(var)
          remove_instance_variable(var)
        end
      end

      def ensure_keys_are_objects
        @__objects.each_key do |type|
          next unless @__api.objects.select { |o| o.name == type }.empty?
          raise [
            "Object #{type} is not a valid type.",
            "Allowed types are: #{@__api.objects.map(&:name)}"
          ].join(' ')
        end
      end
    end

    def out_name
      [@__product.prefix, Google::StringUtils.underscore(@name)].join('_')
    end

    # rubocop:disable Lint/DuplicateMethods
    def identity
      props = all_user_properties
      if @identity.nil?
        props.select { |p| p.name == Api::Type::String::NAME.name }
      else
        props.select { |p| @identity.include?(p.name) }
      end
    end
    # rubocop:enable Lint/DuplicateMethods

    # 'identity' is already taken by Ruby.
    def __identity
      @identity
    end

    # Main data validation. As the validation code is simple, but long due to
    # the number of properties, it is okay to ignore Rubocop warnings about
    # method size and complexity.
    #
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/MethodLength
    # rubocop:disable Metrics/PerceivedComplexity
    def validate
      super
      check_property :async, Api::Async unless @async.nil?
      check_property :base_url, String unless @exclude
      check_property :description, String
      check_property :exclude, :boolean unless @exclude.nil?
      check_property :kind, String unless @kind.nil?
      check_property :parameters, Array unless @parameters.nil?
      check_property :properties, Array unless @exclude
      check_property :exports, Array unless @exports.nil?
      check_property :self_link, String unless @self_link.nil?
      check_property :self_link_query, Api::Resource::ResponseList \
        unless @self_link_query.nil?
      check_property :virtual, :boolean unless @virtual.nil?
      check_property :readonly, :boolean unless @readonly.nil?
      check_property :label_override, String unless @label_override.nil?
      check_property :transport, Transport unless @transport.nil?
      check_property :references, ReferenceLinks unless @references.nil?
      set_variables(@parameters, :__resource)
      set_variables(@properties, :__resource)
      check_property_list :parameters, @parameters, Api::Type
      check_property_list :properties, @properties, Api::Type

      check_property :create_verb, Symbol if @create_verb
      raise "Invalid create verb #{@create_verb}" \
        unless @create_verb.nil? || %i[POST PUT].include?(@create_verb)

      check_identity unless @identity.nil?
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/PerceivedComplexity

    def all_user_properties
      (properties || []) + (parameters || [])
    end

    def required_properties
      all_user_properties.select(&:required)
    end

    def exported_properties
      return [] if @exports.nil?
      from_api = @exports.select { |e| e.is_a?(Api::Type::FetchedExternal) }
                         .each { |e| e.resource = self }
      prop_names = @exports - from_api
      all_user_properties.select { |p| prop_names.include?(p.name) }
                         .concat(from_api)
                         .sort_by(&:name)
    end

    # TODO(alexstephen): Update test_constants to use this function.
    # Returns all of the properties that are a part of the self_link or
    # collection URLs
    def uri_properties
      [@base_url, @__product.base_url].map do |url|
        parts = url.scan(/\{\{(.*?)\}\}/).flatten
        parts << 'name'
        parts.delete('project')
        parts.map { |pt| all_user_properties.select { |p| p.name == pt }[0] }
      end.flatten
    end

    def check_identity
      check_property :identity, Array

      # Ensures we have all properties defined
      @identity.each do |i|
        raise "Missing property/parameter for identity #{i}" \
          if all_user_properties.select { |p| p.name == i }.empty?
      end
    end

    # Returns all resourcerefs at any depth
    def all_resourcerefs
      resourcerefs_for_properties(all_user_properties, self)
    end

    def kind?
      !@kind.nil?
    end

    def encoder?
      return false if @transport.nil?
      return false if @transport.encoder.nil?
      true
    end

    private

    # Given an array of properties, return all ResourceRefs contained within
    # Requires:
    #   props- a list of props
    #   original_object - the original object containing props. This is to
    #                     avoid self-referencing objects.
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/MethodLength
    # rubocop:disable Metrics/PerceivedComplexity
    def resourcerefs_for_properties(props, original_obj)
      rrefs = []
      props.each do |p|
        # We need to recurse on ResourceRefs to get all levels
        # We do not want to recurse on resourcerefs of type self to avoid
        # infinite loop.
        if p.is_a? Api::Type::ResourceRef
          # We want to avoid a circular reference
          # This reference may be the next reference or have some number of refs
          # in between it.
          next if p.resource_ref == original_obj
          next if p.resource_ref == p.__resource
          rrefs << p
          rrefs.concat(resourcerefs_for_properties(p.resource_ref
                                                    .required_properties,
                                                   original_obj))
        elsif p.is_a? Api::Type::NestedObject
          rrefs.concat(resourcerefs_for_properties(p.properties, original_obj))
        elsif p.is_a? Api::Type::Array
          if p.item_type.is_a? Api::Type::NestedObject
            rrefs.concat(resourcerefs_for_properties(p.item_type.properties,
                                                     original_obj))
          elsif p.item_type.is_a? Api::Type::ResourceRef
            rrefs << p.item_type
            rrefs.concat(resourcerefs_for_properties(p.item_type.resource_ref
                                                      .required_properties,
                                                     original_obj))
          end
        end
      end
      rrefs.uniq
    end

    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/PerceivedComplexity
  end
  # rubocop:enable Metrics/ClassLength
end
