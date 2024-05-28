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
  class Resource < Api::NamedObject
    # The list of properties (attr_reader) that can be overridden in
    # <provider>.yaml.
    module Properties
      include Api::NamedObject::Properties

      # [Required] A description of the resource that's surfaced in provider
      # documentation.
      attr_accessor :description
      # [Required] (Api::Resource::ReferenceLinks) Reference links provided in
      # downstream documentation.
      attr_reader :references
      # [Required] The GCP "relative URI" of a resource, relative to the product
      # base URL. It can often be inferred from the `create` path.
      attr_accessor :base_url

      # ====================
      # Common Configuration
      # ====================
      #
      # [Optional] The minimum API version this resource is in. Defaults to ga.
      attr_reader :min_version
      # [Optional] If set to true, don't generate the resource.
      attr_reader :exclude
      # [Optional] If set to true, the resource is not able to be updated.
      attr_accessor :immutable
      # [Optional] If set to true, this resource uses an update mask to perform
      # updates. This is typical of newer GCP APIs.
      attr_accessor :update_mask
      # [Optional] If set to true, the object has a `self_link` field. This is
      # typical of older GCP APIs.
      attr_reader :has_self_link

      # [Optional] The validator "relative URI" of a resource, relative to the product
      # base URL. Specific to defining the resource as a CAI asset.
      attr_reader :cai_base_url

      # ====================
      # URL / HTTP Configuration
      # ====================
      #
      # [Optional] The "identity" URL of the resource. Defaults to:
      # * base_url when the create_verb is :POST
      # * self_link when the create_verb is :PUT  or :PATCH
      attr_accessor :self_link
      # [Optional] The URL used to creating the resource. Defaults to:
      # * collection url when the create_verb is :POST
      # * self_link when the create_verb is :PUT or :PATCH
      attr_accessor :create_url
      # [Optional] The URL used to delete the resource. Defaults to the self
      # link.
      attr_accessor :delete_url
      # [Optional] The URL used to update the resource. Defaults to the self
      # link.
      attr_accessor :update_url
      # [Optional] The HTTP verb used during create. Defaults to :POST.
      attr_reader :create_verb
      # [Optional] The HTTP verb used during read. Defaults to :GET.
      attr_reader :read_verb
      # [Optional] The HTTP verb used during update. Defaults to :PUT.
      attr_accessor :update_verb
      # [Optional] The HTTP verb used during delete. Defaults to :DELETE.
      attr_reader :delete_verb
      # [Optional] Additional Query Parameters to append to GET. Defaults to ""
      attr_reader :read_query_params
      # ====================
      # Collection / Identity URL Configuration
      # ====================
      #
      # [Optional] This is the name of the list of items
      # within the collection (list) json. Will default to the
      # camelcase plural name of the resource.
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

      # ====================
      # Terraform Overrides
      # ====================

      # [Optional] If non-empty, overrides the full filename prefix
      # i.e. google/resource_product_{{resource_filename_override}}.go
      # i.e. google/resource_product_{{resource_filename_override}}_test.go
      attr_reader :filename_override

      # If non-empty, overrides the full given resource name.
      # i.e. 'google_project' for resourcemanager.Project
      # Use Provider::Terraform::Config.legacy_name to override just
      # product name.
      # Note: This should not be used for vanity names for new products.
      # This was added to handle preexisting handwritten resources that
      # don't match the natural generated name exactly, and to support
      # services with a mix of handwritten and generated resources.
      attr_reader :legacy_name

      # The Terraform resource id format used when calling #setId(...).
      # For instance, `{{name}}` means the id will be the resource name.
      attr_accessor :id_format
      # Override attribute used to handwrite the formats for generating regex strings
      # that match templated values to a self_link when importing, only necessary when
      # a resource is not adequately covered by the standard provider generated options.
      # Leading a token with `%`
      # i.e. {{%parent}}/resource/{{resource}}
      # will allow that token to hold multiple /'s.
      attr_accessor :import_format
      attr_reader :custom_code
      attr_reader :docs

      # This block inserts entries into the customdiff.All() block in the
      # resource schema -- the code for these custom diff functions must
      # be included in the resource constants or come from tpgresource
      attr_reader :custom_diff

      # Lock name for a mutex to prevent concurrent API calls for a given
      # resource.
      attr_reader :mutex

      # Examples in documentation. Backed by generated tests, and have
      # corresponding OiCS walkthroughs.
      attr_reader :examples

      # Virtual fields on the Terraform resource. Usage and differences from url_param_only
      # are documented in provider/terraform/virtual_fields.rb
      attr_reader :virtual_fields

      # TODO(alexstephen): Deprecate once all resources using autogen async.
      # If true, generates product operation handling logic.
      attr_accessor :autogen_async

      # If true, resource is not importable
      attr_reader :exclude_import

      # If true, exclude resource from Terraform Validator
      # (i.e. terraform-provider-conversion)
      attr_reader :exclude_tgc

      # If true, skip sweeper generation for this resource
      attr_reader :skip_sweeper

      # Override sweeper settings
      attr_reader :sweeper

      attr_reader :timeouts

      # An array of function names that determine whether an error is retryable.
      attr_reader :error_retry_predicates

      # An array of function names that determine whether an error is not retryable.
      attr_reader :error_abort_predicates

      # Optional attributes for declaring a resource's current version and generating
      # state_upgrader code to the output .go file from files stored at
      # mmv1/templates/terraform/state_migrations/
      # used for maintaining state stability with resources first provisioned on older api versions.
      attr_reader :schema_version
      # From this schema version on, state_upgrader code is generated for the resource.
      # When unset, state_upgrade_base_schema_version defauts to 0.
      # Normally, it is not needed to be set.
      attr_reader :state_upgrade_base_schema_version
      attr_reader :state_upgraders
      # This block inserts the named function and its attribute into the
      # resource schema -- the code for the migrate_state function must
      # be included in the resource constants or come from tpgresource
      # included for backwards compatibility as an older state migration method
      # and should not be used for new resources.
      attr_reader :migrate_state

      # Set to true for resources that are unable to be deleted, such as KMS keyrings or project
      # level resources such as firebase project
      attr_reader :skip_delete

      # Set to true for resources that are unable to be read from the API, such as
      # public ca external account keys
      attr_reader :skip_read

      # Set to true for resources that wish to disable automatic generation of default provider
      # value customdiff functions
      attr_reader :skip_default_cdiff

      # This enables resources that get their project via a reference to a different resource
      # instead of a project field to use User Project Overrides
      attr_reader :supports_indirect_user_project_override

      # If true, the resource's project field can be specified as either the short form project
      # id or the long form projects/project-id. The extra projects/ string will be removed from
      # urls and ids. This should only be used for resources that previously supported long form
      # project ids for backwards compatibility.
      attr_reader :legacy_long_form_project

      # Function to transform a read error so that handleNotFound recognises
      # it as a 404. This should be added as a handwritten fn that takes in
      # an error and returns one.
      attr_reader :read_error_transform

      # If true, resources that failed creation will be marked as tainted. As a consequence
      # these resources will be deleted and recreated on the next apply call. This pattern
      # is preferred over deleting the resource directly in post_create_failure hooks.
      attr_reader :taint_resource_on_failed_create

      # Add a deprecation message for a resource that's been deprecated in the API.
      attr_reader :deprecation_message
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
      check :cai_base_url, type: String, required: false
      check :create_url, type: String
      check :delete_url, type: String
      check :update_url, type: String
      check :read_query_params, type: String
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

      check :collection_url_key, default: @name.plural.camelize(:lower)

      check :create_verb, type: Symbol, default: :POST, allowed: %i[POST PUT PATCH]
      check :read_verb, type: Symbol, default: :GET, allowed: %i[GET POST]
      check :delete_verb, type: Symbol, default: :DELETE, allowed: %i[POST PUT PATCH DELETE]
      check :update_verb, type: Symbol, default: :PUT, allowed: %i[POST PUT PATCH]

      check :immutable, type: :boolean
      check :min_version, type: String

      check :has_self_link, type: :boolean, default: false

      set_variables(@parameters, :__resource)
      set_variables(@properties, :__resource)

      check :properties, type: Array, item_type: Api::Type, required: true unless @exclude
      check :parameters, type: Array, item_type: Api::Type unless @exclude

      check :iam_policy, type: Api::Resource::IamPolicy
      check :exclude_resource, type: :boolean, default: false

      @examples ||= []

      check :filename_override, type: String
      check :legacy_name, type: String
      check :id_format, type: String
      check :examples, item_type: Provider::Terraform::Examples, type: Array, default: []
      check :virtual_fields,
            item_type: Api::Type,
            type: Array,
            default: []

      check :custom_code, type: Provider::Terraform::CustomCode,
                          default: Provider::Terraform::CustomCode.new
      check :sweeper, type: Provider::Terraform::Sweeper, default: Provider::Terraform::Sweeper.new
      check :docs, type: Provider::Terraform::Docs, default: Provider::Terraform::Docs.new
      check :import_format, type: Array, item_type: String, default: []
      check :autogen_async, type: :boolean, default: false
      check :exclude_import, type: :boolean, default: false
      check :custom_diff, type: Array, item_type: String, default: []
      check :timeouts, type: Api::Timeouts
      check :error_retry_predicates, type: Array, item_type: String
      check :error_abort_predicates, type: Array, item_type: String
      check :schema_version, type: Integer
      check :state_upgrade_base_schema_version, type: Integer, default: 0
      check :state_upgraders, type: :boolean, default: false
      check :migrate_state, type: String
      check :skip_delete, type: :boolean, default: false
      check :skip_read, type: :boolean, default: false
      check :skip_default_cdiff, type: :boolean, default: false
      check :supports_indirect_user_project_override, type: :boolean, default: false
      check :legacy_long_form_project, type: :boolean, default: false
      check :read_error_transform, type: String
      check :taint_resource_on_failed_create, type: :boolean, default: false
      check :skip_sweeper, type: :boolean, default: false
      check :deprecation_message, type: ::String

      validate_identity unless @identity.nil?
    end

    # ====================
    # Custom Getters and Setters
    # ====================

    # Returns all properties and parameters including the ones that are
    # excluded. This is used for PropertyOverride validation
    def all_properties
      ((@properties || []) + (@parameters || []))
    end

    def properties_with_excluded
      @properties || []
    end

    def properties
      (@properties || []).reject(&:exclude)
    end

    attr_writer :properties

    def parameters
      (@parameters || []).reject(&:exclude)
    end

    attr_writer :parameters

    # Return the user-facing properties in client tools; this ends up meaning
    # both properties and parameters but without any that are excluded due to
    # version mismatches or manual exclusion
    def all_user_properties
      properties + parameters
    end

    def required_properties
      all_user_properties.select(&:required)
    end

    def all_nested_properties(props)
      nested = props
      props.each do |prop|
        if !prop.flatten_object && !prop.nested_properties.nil?
          nested += all_nested_properties(prop.nested_properties)
        end
      end
      nested
    end

    # All settable properties in the resource.
    # Fingerprints aren't *really" settable properties, but they behave like one.
    # At Create, they have no value but they can just be read in anyways, and after a Read
    # they will need to be set in every Update.
    def settable_properties
      props = all_user_properties.reject do |v|
        v.output && !v.is_a?(Api::Type::Fingerprint) && !v.is_a?(Api::Type::KeyValueEffectiveLabels)
      end
      props = props.reject(&:url_param_only)
      props.reject do |v|
        v.is_a?(Api::Type::KeyValueLabels) || v.is_a?(Api::Type::KeyValueAnnotations)
      end
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

    def sensitive_props
      all_nested_properties(root_properties).select(&:sensitive)
    end

    # Return the product-level async object, or the resource-specific one
    # if one exists.
    def async
      return @__product.async unless @async

      @async
    end

    attr_writer :async

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

    def add_labels_related_fields(props, parent)
      props.each do |p|
        if p.is_a? Api::Type::KeyValueLabels
          add_labels_fields(props, parent, p)
        elsif p.is_a? Api::Type::KeyValueAnnotations
          add_annotations_fields(props, parent, p)
        elsif (p.is_a? Api::Type::NestedObject) && !p.all_properties.nil?
          p.properties = add_labels_related_fields(p.all_properties, p)
        end
      end
      props
    end

    def add_labels_fields(props, parent, labels)
      @custom_diff ||= []
      if parent.nil? || parent.flatten_object
        @custom_diff.append('tpgresource.SetLabelsDiff')
      elsif parent.name == 'metadata'
        @custom_diff.append('tpgresource.SetMetadataLabelsDiff')
      end

      props << build_terraform_labels_field('labels', parent, labels)
      props << build_effective_labels_field('labels', labels)

      # The effective_labels field is used to write to API, instead of the labels field.
      labels.ignore_write = true
      labels.description = "#{labels.description}\n\n#{get_labels_field_note(labels.name)}"
      return unless parent.nil?

      labels.immutable = false
    end

    def add_annotations_fields(props, parent, annotations)
      # The effective_annotations field is used to write to API,
      # instead of the annotations field.
      annotations.ignore_write = true
      note = get_labels_field_note(annotations.name)
      annotations.description = "#{annotations.description}\n\n#{note}"

      @custom_diff ||= []
      if parent.nil?
        @custom_diff.append('tpgresource.SetAnnotationsDiff')
      elsif parent.name == 'metadata'
        @custom_diff.append('tpgresource.SetMetadataAnnotationsDiff')
      end

      props << build_effective_labels_field('annotations', annotations)
    end

    def build_effective_labels_field(name, labels)
      description = "All of #{name} (key/value pairs)\
 present on the resource in GCP, including the #{name} configured through Terraform,\
 other clients and services."

      Api::Type::KeyValueEffectiveLabels.new(
        name: "effective#{name.capitalize}",
        output: true,
        api_name: name,
        description:,
        min_version: labels.field_min_version,
        update_verb: labels.update_verb,
        update_url: labels.update_url,
        immutable: labels.immutable
      )
    end

    def build_terraform_labels_field(name, parent, labels)
      description = "The combination of #{name} configured directly on the resource
 and default #{name} configured on the provider."

      immutable = if parent.nil?
                    false
                  else
                    labels.immutable
                  end

      Api::Type::KeyValueTerraformLabels.new(
        name: "terraform#{name.capitalize}",
        output: true,
        api_name: name,
        description:,
        min_version: labels.field_min_version,
        ignore_write: true,
        update_url: labels.update_url,
        immutable:
      )
    end

    # Check if the resource has root "labels" field
    def root_labels?
      root_properties.each do |p|
        return true if p.is_a? Api::Type::KeyValueLabels
      end
      false
    end

    # Return labels fields that should be added to ImportStateVerifyIgnore
    def ignore_read_labels_fields(props)
      fields = []
      props.each do |p|
        if (p.is_a? Api::Type::KeyValueLabels) ||
           (p.is_a? Api::Type::KeyValueTerraformLabels) ||
           (p.is_a? Api::Type::KeyValueAnnotations)
          fields << p.terraform_lineage
        elsif (p.is_a? Api::Type::NestedObject) && !p.all_properties.nil?
          fields.concat(ignore_read_labels_fields(p.all_properties))
        end
      end
      fields
    end

    # Return ignore_read fields that should be added to ImportStateVerifyIgnore
    def ignore_read_fields(props)
      fields = []
      props.each do |p|
        if p.ignore_read && !p.url_param_only && !p.is_a?(Api::Type::ResourceRef)
          fields << p.terraform_lineage
        elsif (p.is_a? Api::Type::NestedObject) && !p.all_properties.nil?
          fields.concat(ignore_read_fields(p.all_properties))
        end
      end
      fields
    end

    def get_labels_field_note(title)
      "**Note**: This field is non-authoritative, and will only manage the #{title} present " \
"in your configuration.
Please refer to the field `effective_#{title}` for all of the #{title} present on the resource."
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
        # If the terms in this are not snake-cased, this will require
        # an override in Terraform.
        @self_link
      end
    end

    def collection_url
      [@__product.base_url, collection_uri].flatten.join
    end

    def collection_uri
      @base_url
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

    def delete_uri
      if @delete_url.nil?
        self_link_uri
      else
        @delete_url
      end
    end

    def resource_name
      __product.name + name
    end

    # Filter the properties to keep only the ones don't have custom update
    # method and group them by update url & verb.
    def properties_without_custom_update(properties)
      properties.select do |p|
        p.update_url.nil? || p.update_verb.nil? || p.update_verb == :NOOP
      end
    end

    def update_body_properties
      update_prop = properties_without_custom_update(settable_properties)
      update_prop = update_prop.reject(&:immutable) if update_verb == :PATCH
      update_prop
    end

    # Handwritten TF Operation objects will be shaped like accessContextManager
    # while the Google Go Client will have a name like accesscontextmanager
    def client_name_pascal
      client_name = __product.client_name || __product.name
      client_name.camelize(:upper)
    end

    # In order of preference, use TF override,
    # general defined timeouts, or default Timeouts
    def timeouts
      timeouts_filtered = @timeouts
      timeouts_filtered ||= async&.operation&.timeouts
      timeouts_filtered ||= Api::Timeouts.new
      timeouts_filtered
    end

    def project?
      base_url.include?('{{project}}') || create_url&.include?('{{project}}')
    end

    def region?
      base_url.include?('{{region}}') && parameters.any? { |p| p.name == 'region' && p.ignore_read }
    end

    def zone?
      base_url.include?('{{zone}}') && parameters.any? { |p| p.name == 'zone' && p.ignore_read }
    end

    def merge(other)
      result = self.class.new
      instance_variables.each do |v|
        result.instance_variable_set(v, instance_variable_get(v))
      end

      other.instance_variables.each do |v|
        if other.instance_variable_get(v).instance_of?(Array)
          result.instance_variable_set(v, deep_merge(result.instance_variable_get(v),
                                                     other.instance_variable_get(v)))
        else
          result.instance_variable_set(v, other.instance_variable_get(v))
        end
      end

      result
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

      json_out.merge!(properties.to_h { |p| [p.name, p] })
      json_out.merge!(parameters.to_h { |p| [p.name, p] })

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
  end
end
