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
require 'api/resource/iam_policy'
require 'api/resource/nested_query'
require 'api/resource/reference_links'
require 'google/string_utils'

module Api
  # An object available in the product
  class Resource < Api::Object::Named
    # The list of properties (attr_reader) that can be overridden in
    # <provider>.yaml.
    module Properties
      include Api::Object::Named::Properties

      # [Required] A description of the resource that's surfaced in provider
      # documentation.
      attr_reader :description
      # [Required] (Api::Resource::ReferenceLinks) Reference links provided in
      # downstream documentation.
      attr_reader :references
      # [Required] The GCP "relative URI" of a resource, relative to the product
      # base URL. It can often be inferred from the `create` path.
      attr_reader :base_url

      # ====================
      # Common Configuration
      # ====================
      #
      # [Optional] The minimum API version this resource is in. Defaults to ga.
      attr_reader :min_version
      # [Optional] If set to true, don't generate the resource.
      attr_reader :exclude
      # [Optional] If set to true, the resource is not able to be updated.
      attr_reader :input
      # [Optional] If set to true, this resource uses an update mask to perform
      # updates. This is typical of newer GCP APIs.
      attr_reader :update_mask
      # [Optional] If set to true, the object has a `self_link` field. This is
      # typical of older GCP APIs.
      attr_reader :has_self_link

      # ====================
      # URL / HTTP Configuration
      # ====================
      #
      # [Optional] The "identity" URL of the resource. Defaults to:
      # * base_url when the create_verb is :POST
      # * self_link when the create_verb is :PUT  or :PATCH
      attr_reader :self_link
      # [Optional] The URL used to creating the resource. Defaults to:
      # * collection url when the create_verb is :POST
      # * self_link when the create_verb is :PUT or :PATCH
      attr_reader :create_url
      # [Optional] The URL used to delete the resource. Defaults to the self
      # link.
      attr_reader :delete_url
      # [Optional] The URL used to update the resource. Defaults to the self
      # link.
      attr_reader :update_url
      # [Optional] The HTTP verb used during create. Defaults to :POST.
      attr_reader :create_verb
      # [Optional] The HTTP verb used during read. Defaults to :GET.
      attr_reader :read_verb
      # [Optional] The HTTP verb used during update. Defaults to :PUT.
      attr_reader :update_verb
      # [Optional] The HTTP verb used during delete. Defaults to :DELETE.
      attr_reader :delete_verb
      # ====================
      # Collection / Identity URL Configuration
      # ====================
      #
      # [Optional] This is the name of the list of items
      # within the collection (list) json. Will default to the
      # camelcase pluralize name of the resource.
      attr_reader :collection_url_key
      # [Optional] An ordered list of names of parameters that uniquely identify
      # the resource.
      # Generally, it's safe to leave empty, in which case it defaults to `name`.
      # Other values are normally useful in cases where an object has a parent
      # and is identified by some non-name value, such as an ip+port pair.
      # If you're writing a fine-grained resource (eg with nested_query) a value
      # must be set.
      attr_reader :identity

      # [Optional] (Api::Resource::NestedQuery) This is useful in case you need
      # to change the query made for GET requests only. In particular, this is
      # often used to extract an object from a parent object or a collection.
      # Note that if both nested_query and custom_code.decoder are provided,
      # the decoder will be included within the code handling the nested query.
      attr_reader :nested_query

      # ====================
      # IAM Configuration
      # ====================
      #
      # [Optional] (Api::Resource::IamPolicy) Configuration of a resource's
      # resource-specific IAM Policy.
      attr_reader :iam_policy
      # [Optional] If set to true, don't generate the resource itself; only
      # generate the IAM policy.
      attr_reader :exclude_resource

      # [Optional] GCP kind, e.g. `compute#disk`
      attr_reader :kind
      # [Optional] If set to true, indicates that a resource is not configurable
      # such as GCP regions.
      attr_reader :readonly
    end

    include Properties

    # Parameters can be overridden via Provider::PropertyOverride
    # A custom getter is used for :parameters instead of `attr_reader`

    # Properties can be overridden via Provider::PropertyOverride
    # A custom getter is used for :properties instead of `attr_reader`

    attr_reader :__product

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

      check :collection_url_key, default: @name.pluralize.camelize(:lower)

      check :create_verb, type: Symbol, default: :POST, allowed: %i[POST PUT PATCH]
      check :read_verb, type: Symbol, default: :GET, allowed: %i[GET POST]
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

      validate_identity unless @identity.nil?
    end

    # ====================
    # Custom Getters
    # ====================

    # Returns all properties and parameters including the ones that are
    # excluded. This is used for PropertyOverride validation
    def all_properties
      ((@properties || []) + (@parameters || []))
    end

    def properties
      (@properties || []).reject(&:exclude)
    end

    def parameters
      (@parameters || []).reject(&:exclude)
    end

    # Return the user-facing properties in client tools; this ends up meaning
    # both properties and parameters but without any that are excluded due to
    # version mismatches or manual exclusion
    def all_user_properties
      properties + parameters
    end

    def required_properties
      all_user_properties.select(&:required)
    end

    # Returns all resourcerefs at any depth
    def all_resourcerefs
      resourcerefs_for_properties(all_user_properties, self)
    end

    # All settable properties in the resource.
    # Fingerprints aren't *really" settable properties, but they behave like one.
    # At Create, they have no value but they can just be read in anyways, and after a Read
    # they will need ot be set in every Update.
    def settable_properties
      all_user_properties.reject { |v| v.output && !v.is_a?(Api::Type::Fingerprint) }
                         .reject(&:url_param_only)
    end

    # Properties that will be returned in the API body
    def gettable_properties
      all_user_properties.reject(&:url_param_only)
    end

    # Returns the list of top-level properties once any nested objects with flatten_object
    # set to true have been collapsed
    def root_properties
      all_user_properties.flat_map do |p|
        if p.flatten_object
          p.root_properties
        else
          p
        end
      end
    end

    # Return the product-level async object, or the resource-specific one
    # if one exists.
    def async
      return @__product.async unless @async

      @async
    end

    # Return the resource-specific identity properties, or a best guess of the
    # `name` value for the resource.
    def identity
      props = all_user_properties
      if @identity.nil?
        props.select { |p| p.name == Api::Type::String::NAME.name }
      else
        props.select { |p| @identity.include?(p.name) }.sort_by { |p| @identity.index p.name }
      end
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

    # ====================
    # Version-related methods
    # ====================

    def min_version
      if @min_version.nil?
        @__product.lowest_version
      else
        @__product.version_obj(@min_version)
      end
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

    # ====================
    # URL-related methods
    # ====================

    # Returns the "self_link_url" which is generally really the resource's GET
    # URL. In older resources generally, this was the self_link value & was the
    # product.base_url + resource.base_url + '/name'
    # In newer resources there is much less standardisation in terms of value.
    # Generally for them though, it's the product.base_url + resource.name
    def self_link_url
      [@__product.base_url, self_link_uri].flatten.join
    end

    # Returns the partial uri / relative path of a resource. In newer resources,
    # this is the name. This fn is named self_link_uri for consistency, but
    # could otherwise be considered to be "path"
    def self_link_uri
      if @self_link.nil?
        [@base_url, '{{name}}'].join('/')
      else
        @self_link
      end
    end

    def collection_url
      [@__product.base_url, collection_uri].flatten.join
    end

    def collection_uri
      @base_url
    end

    def async_operation_url
      [@__product.base_url, async_operation_uri].flatten.join
    end

    def async_operation_uri
      raise 'Not an async resource' if async.nil?

      async.operation.base_url
    end

    def full_create_url
      [@__product.base_url, create_uri].flatten.join
    end

    def create_uri
      if @create_url.nil?
        if @create_verb.nil? || @create_verb == :POST
          collection_uri
        else
          self_link_uri
        end
      else
        @create_url
      end
    end

    def full_delete_url
      [@__product.base_url, delete_uri].flatten.join
    end

    def delete_uri
      if @delete_url.nil?
        self_link_uri
      else
        @delete_url
      end
    end

    # ====================
    # Debugging Methods
    # ====================

    # Prints a dot notation path to where the field is nested within the parent
    # object when called on a property. eg: parent.meta.label.foo
    # Redefined on Resource to terminate the calls up the parent chain.
    def lineage
      name
    end

    def to_s
      JSON.pretty_generate(self)
    end

    def to_json(opts = nil)
      # ignore fields that will contain references to parent resources
      ignored_fields = %i[@__product @__parent @__resource @api_name
                          @properties @parameters]
      json_out = {}

      instance_variables.each do |v|
        json_out[v] = instance_variable_get(v) unless ignored_fields.include? v
      end

      json_out.merge!(properties.map { |p| [p.name, p] }.to_h)
      json_out.merge!(parameters.map { |p| [p.name, p] }.to_h)

      JSON.generate(json_out, opts)
    end

    private

    def validate_identity
      check :identity, type: Array, item_type: String, required: true

      # Ensures we have all properties defined
      @identity.each do |i|
        raise "Missing property/parameter for identity #{i}" \
          if all_user_properties.select { |p| p.name == i }.empty?
      end
    end

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
