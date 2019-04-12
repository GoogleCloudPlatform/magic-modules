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
      # GCP kind, e.g. `compute#disk`
      attr_reader :kind
      # URI: relative to `@api.base_url` or absolute
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
      # This is useful in case you need to change the query made for
      # GET requests only. In particular, this is often used
      # to extract an object from a parent object or a collection.
      attr_reader :nested_query

      attr_reader :exclude
      attr_reader :async
      attr_reader :readonly
      # Documentation references
      attr_reader :references
      attr_reader :create_verb
      attr_reader :delete_verb
      attr_reader :update_verb
      attr_reader :input # If true, resource is not updatable as a whole unit
      attr_reader :min_version # Minimum API version this resource is in
      attr_reader :update_mask
      attr_reader :has_self_link

      attr_reader :iam_policy
      attr_reader :exclude_resource
    end

    include Properties

    # Parameters can be overridden via Provider::PropertyOverride
    # A custom getter is used for :parameters instead of `attr_reader`

    # Properties can be overridden via Provider::PropertyOverride
    # A custom getter is used for :properties instead of `attr_reader`

    attr_reader :__product

    # Allows mapping of requests to specific API layout quirks.
    class Wrappers < Api::Object
      attr_reader :create

      def validate
        super
        check :create, type: ::String, required: true
      end
    end

    # Query information for finding resource nested in an returned API object
    # i.e. fine-grained resources
    class NestedQuery < Api::Object
      # A list of keys to traverse in order.
      # i.e. backendBucket --> cdnPolicy.signedUrlKeyNames
      # should be ["cdnPolicy", "signedUrlKeyNames"]
      attr_reader :keys

      # If true, we expect the the nested list to be
      # a list of IDs for the nested resource, rather
      # than a list of nested resource objects
      # i.e. backendBucket.cdnPolicy.signedUrlKeyNames is a list of key names
      # rather than a list of the actual key objects
      attr_reader :is_list_of_ids

      # This is used by Ansible, but may not be necessary.
      attr_reader :kind

      def validate
        super

        check :keys, type: Array, item_type: String, required: true
        check :is_list_of_ids, type: :boolean, default: false

        check :kind, type: String
      end
    end

    # Represents a response from the API that returns a list of objects.
    class ResponseList < Api::Object
      attr_reader :kind
      attr_reader :items

      def validate
        super

        check :items, default: 'items', type: ::String, required: true
        check :kind, type: ::String
      end

      def kind?
        !@kind.nil?
      end
    end

    # Represents a list of documentation links.
    class ReferenceLinks < Api::Object
      # Hash containing
      # name: The title of the link
      # value: The URL to navigate on click
      attr_reader :guides

      # the url of the API guide
      attr_reader :api

      def validate
        super

        check :guides, type: Hash, default: {}, required: true
        check :api, type: String
      end
    end

    # Information about the IAM policy for this resource
    # Several GCP resources have IAM policies that are scoped to
    # and accessed via their parent resource
    # See: https://cloud.google.com/iam/docs/overview
    class IamPolicy < Api::Object
      # boolean of if this binding should be generated
      attr_reader :exclude

      def validate
        super

        check :exclude, type: :boolean, default: false
      end
    end

    def to_s
      JSON.pretty_generate(self)
    end

    def to_json(opts = nil)
      # ignore fields that will contain references to parent resources
      ignored_fields = %i[@__product @__parent @__resource @api_name @collection_url_response]
      json_out = {}

      instance_variables.each do |v|
        json_out[v] = instance_variable_get(v) unless ignored_fields.include? v
      end

      json_out[:@properties] = properties.map { |p| [p.name, p] }.to_h
      json_out[:@parameters] = parameters.map { |p| [p.name, p] }.to_h

      JSON.generate(json_out, opts)
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
      check :async, type: Api::Async
      check :base_url, type: String
      check :create_url, type: String
      check :delete_url, type: String
      check :update_url, type: String
      check :update_mask, type: :boolean
      check :description, type: String, required: true
      check :exclude, type: :boolean
      check :kind, type: String

      check :self_link, type: String
      check :readonly, type: :boolean
      check :references, type: ReferenceLinks

      check :nested_query, type: Api::Resource::NestedQuery
      if @nested_query&.is_list_of_ids && @identity&.length != 1
        raise ':is_list_of_ids = true implies resource`\
              `has exactly one :identity property"'
      end

      check :collection_url_response, default: Api::Resource::ResponseList.new,
                                      type: Api::Resource::ResponseList

      check :create_verb, type: Symbol, default: :POST, allowed: %i[POST PUT]
      check :delete_verb, type: Symbol, default: :DELETE, allowed: %i[POST PUT PATCH DELETE]
      check :update_verb, type: Symbol, default: :PUT, allowed: %i[POST PUT PATCH]

      check :input, type: :boolean
      check :min_version, type: String

      check :has_self_link, type: :boolean, default: false

      set_variables(@parameters, :__resource)
      set_variables(@properties, :__resource)

      check :properties, type: Array, item_type: Api::Type, required: true unless @exclude
      check :parameters, type: Array, item_type: Api::Type unless @exclude

      check :iam_policy, type: Api::Resource::IamPolicy
      check :exclude_resource, type: :boolean, default: false

      check_identity unless @identity.nil?
    end

    def properties
      (@properties || []).reject(&:exclude)
    end

    def parameters
      (@parameters || []).reject(&:exclude)
    end

    def not_in_version?(version)
      version < min_version
    end

    # Recurses through all nested properties and parameters and changes their
    # 'exclude' instance variable if the property is at a version below the
    # one that is passed in.
    def exclude_if_not_in_version!(version)
      @exclude ||= not_in_version? version
      @properties&.each { |p| p.exclude_if_not_in_version!(version) }
      @parameters&.each { |p| p.exclude_if_not_in_version!(version) }

      nil
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
      check :identity, type: Array, item_type: String, required: true

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
          @create_url.split("\n").map(&:strip).compact
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

    # All settable properties in the resource.
    # Fingerprints aren't *really" settable properties, but they behave like one.
    # At Create, they have no value but they can just be read in anyways, and after a Read
    # they will need ot be set in every Update.
    def settable_properties
      all_user_properties.reject { |v| v.output && !v.is_a?(Api::Type::Fingerprint) }
                         .reject(&:url_param_only)
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
        elsif p.nested_properties?
          rrefs.concat(resourcerefs_for_properties(p.nested_properties, original_obj))
        elsif p.is_a? Api::Type::Array
          if p.item_type.is_a? Api::Type::ResourceRef
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
