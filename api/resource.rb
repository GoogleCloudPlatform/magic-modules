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
  class Resource < Api::Object::Named
    # The list of properties (attr_reader) that can be overridden in
    # <provider>.yaml.
    module Properties
      include Api::Object::Named::Properties

      attr_reader :description
      attr_reader :kind
      attr_reader :base_url
      # URL to use for creating the resource. If not specified, the
      # collection url (when create_verb is default or :POST) or
      # self_link (when create_verb is :PUT) is used instead.
      attr_reader :create_url
      # URL to use to delete the resource. If not specified, the
      # self link is used.
      attr_reader :delete_url
      # URL to use for updating the resource. If not specified, the self link
      # will be used. This currently can only be used with Terraform resources.
      # TODO(#302): Add support for the other providers.
      attr_reader :update_url
      attr_reader :self_link
      # This is useful in case you need to change the query made for
      # GET/DELETE requests only.  In particular, this is often used
      # to add query parameters.
      attr_reader :self_link_query
      # This is the type of response from the collection URL. It contains
      # the name of the list of items within the json, as well as the
      # type that this list should be. This is of type Api::Resource::ResponseList
      attr_reader :collection_url_response
      # This is an array with items that uniquely identify the resource.
      # This is useful in case an API returns a list result and we need
      # to fetch the particular resource we're interested in from that
      # list.  Otherwise, it's safe to leave empty.
      # If empty, we assume that `name` is the identifier.
      attr_reader :identity
      attr_reader :exclude
      attr_reader :async
      attr_reader :readonly
      attr_reader :exports
      attr_reader :transport
      attr_reader :references
      attr_reader :create_verb
      attr_reader :delete_verb
      attr_reader :update_verb
      attr_reader :input # If true, resource is not updatable as a whole unit
      attr_reader :min_version # Minimum API version this resource is in
    end

    include Properties

    # Parameters can be overridden via Provider::PropertyOverride
    # A custom getter is used for :parameters instead of `attr_reader`

    # Properties can be overridden via Provider::PropertyOverride
    # A custom getter is used for :properties instead of `attr_reader`

    attr_reader :__product

    # Allows overriding snowflake transport requests
    class Transport < Api::Object
      attr_reader :encoder
      attr_reader :decoder

      def validate
        super
        check_optional_property :encoder, ::String
        check_optional_property :decoder, ::String
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
        default_value_property :items, 'items'

        check_optional_property :kind, String
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

        @guides ||= {}

        check_property :guides, Hash

        check_optional_property :api, String
      end
    end

    # Represents a hierarchy that has an object as its key.
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
        @__objects.select { |k, v| yield k, v }
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

      def keys
        @__objects.keys
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
      [@__product.prefix, @name.underscore].join('_')
    end

    def identity
      props = all_user_properties
      if @identity.nil?
        props.select { |p| p.name == Api::Type::String::NAME.name }
      else
        props.select { |p| @identity.include?(p.name) }
      end
    end

    # 'identity' is already taken by Ruby.
    def __identity
      @identity
    end

    # Main data validation. As the validation code is simple, but long due to
    # the number of properties, it is okay to ignore Rubocop warnings about
    # method size and complexity.
    #
    def validate
      super
      check_optional_property :async, Api::Async
      check_optional_property :base_url, String
      check_optional_property :create_url, String
      check_optional_property :delete_url, String
      check_optional_property :update_url, String
      check_property :description, String
      check_optional_property :exclude, :boolean
      check_optional_property :kind, String
      check_optional_property :parameters, Array
      check_optional_property :exports, Array
      check_optional_property :self_link, String
      check_optional_property :self_link_query, Api::Resource::ResponseList
      check_optional_property :readonly, :boolean
      check_optional_property :transport, Transport
      check_optional_property :references, ReferenceLinks

      default_value_property :collection_url_response,
                             Api::Resource::ResponseList.new
      check_property :collection_url_response, Api::Resource::ResponseList
      default_value_property :collection_url_response,
                             Api::Resource::ResponseList.new

      check_property_oneof_default :create_verb, %i[POST PUT], :POST, Symbol
      check_property_oneof_default \
        :delete_verb, %i[POST PUT PATCH DELETE], :DELETE, Symbol
      check_property_oneof_default \
        :update_verb, %i[POST PUT PATCH], :PUT, Symbol
      check_optional_property :input, :boolean
      check_optional_property :min_version, String

      set_variables(@parameters, :__resource)
      set_variables(@properties, :__resource)

      check_property_list :parameters, Api::Type

      check_property :properties, Array unless @exclude
      check_property_list :properties, Api::Type

      check_identity unless @identity.nil?
    end

    def properties
      (@properties || []).reject(&:exclude)
    end

    def parameters
      (@parameters || []).reject(&:exclude)
    end

    def exclude_if_not_in_version(version)
      @exclude ||= version < min_version
      @properties&.each { |p| p.exclude_if_not_in_version(version) }
      @parameters&.each { |p| p.exclude_if_not_in_version(version) }

      @exclude
    end

    # Returns all properties and parameters including the ones that are
    # excluded. This is used for PropertyOverride validation
    def all_properties
      ((@properties || []) + (@parameters || []))
    end

    def all_user_properties
      properties + parameters
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
      check_property_list :identity, String

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
      !@transport&.encoder.nil?
    end

    def decoder?
      !@transport&.decoder.nil?
    end

    # Returns true if this resource needs access to the saved API response
    # This response is stored in the @fetched variable
    # Requires:
    #   config: The config for an object
    #   object: An Api::Resource object
    def save_api_results?
      exported_properties.any? { |p| p.is_a? Api::Type::FetchedExternal } \
        || access_api_results
    end

    def min_version
      if @min_version.nil?
        @__product.default_version
      else
        @__product.version_obj(@min_version)
      end
    end

    # Returns self link in two parts - base_url + product_url
    def self_link_url
      base_url = @__product.base_url.split("\n").map(&:strip).compact
      if @self_link.nil?
        [base_url, [@base_url, '{{name}}'].join('/')]
      else
        self_link = @self_link.split("\n").map(&:strip).compact
        [base_url, self_link]
      end
    end

    def collection_url
      [
        @__product.base_url.split("\n").map(&:strip).compact,
        @base_url.split("\n").map(&:strip).compact
      ]
    end

    def async_operation_url
      raise 'Not an async resource' if @async.nil?
      [@__product.base_url, @async.operation.base_url]
    end

    def default_create_url
      if @create_verb.nil? || @create_verb == :POST
        collection_url
      elsif @create_verb == :PUT
        self_link_url
      else
        raise "unsupported create verb #{@create_verb}"
      end
    end

    def full_create_url
      if @create_url.nil?
        default_create_url
      else
        [
          @__product.base_url.split("\n").map(&:strip).compact,
          @create_url
        ]
      end
    end

    def full_delete_url
      if @delete_url.nil?
        self_link_url
      else
        [
          @__product.base_url.split("\n").map(&:strip).compact,
          @delete_url
        ]
      end
    end

    # A regex to check if a full URL was returned or just a shortname.
    def regex_url
      self_link_url.join.gsub('{{project}}', '.*')
                   .gsub('{{name}}', '[a-z1-9\-]*')
                   .gsub('{{zone}}', '[a-z1-9\-]*')
    end

    private

    # Given an array of properties, return all ResourceRefs contained within
    # Requires:
    #   props- a list of props
    #   original_object - the original object containing props. This is to
    #                     avoid self-referencing objects.
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
  end
end
