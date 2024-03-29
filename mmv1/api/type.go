// Copyright 2024 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

// Represents a property type
type Type struct {
	NamedObject `yaml:",inline"`

	// TODO: improve the parsing of properties based on type in resource yaml files.
	Type string

	// TODO: set a specific type intead of interface{}
	DefaultValue interface{} `yaml:"default_value"`

	Description string

	Exclude bool

	// Add a deprecation message for a field that's been deprecated in the API
	// use the YAML chomping folding indicator (>-) if this is a multiline
	// string, as providers expect a single-line one w/o a newline.
	DeprecationMessage string `yaml:"deprecation_message"`

	// Add a removed message for fields no longer supported in the API. This should
	// be used for fields supported in one version but have been removed from
	// a different version.
	RemovedMessage string `yaml:"removed_message"`

	// If set value will not be sent to server on sync.
	// For nested fields, this also needs to be set on each descendant (ie. self,
	// child, etc.).
	Output bool

	// If set to true, changes in the field's value require recreating the
	// resource.
	// For nested fields, this only applies at the current level. This means
	// it should be explicitly added to each field that needs the ForceNew
	// behavior.
	Immutable bool

	// url_param_only will not send the field in the resource body and will
	// not attempt to read the field from the API response.
	// NOTE - this doesn't work for nested fields
	UrlParamOnly bool `yaml:"url_param_only"`

	// For nested fields, this only applies within the parent.
	// For example, an optional parent can contain a required child.
	Required bool

	// Additional query Parameters to append to GET calls.
	ReadQueryParams string `yaml:"read_query_params"`

	UpdateVerb string `yaml:"update_verb"`

	UpdateUrl string `yaml:"update_url"`

	// Some updates only allow updating certain fields at once (generally each
	// top-level field can be updated one-at-a-time). If this is set, we group
	// fields to update by (verb, url, fingerprint, id) instead of just
	// (verb, url, fingerprint), to allow multiple fields to reuse the same
	// endpoints.
	UpdateId string `yaml:"update_id"`

	// The fingerprint value required to update this field. Downstreams should
	// GET the resource and parse the fingerprint value while doing each update
	// call. This ensures we can supply the fingerprint to each distinct
	// request.
	FingerprintName string `yaml:"fingerprint_name"`

	// If true, we will include the empty value in requests made including
	// this attribute (both creates and updates).  This rarely needs to be
	// set to true, and corresponds to both the "NullFields" and
	// "ForceSendFields" concepts in the autogenerated API clients.
	SendEmptyValue bool `yaml:"send_empty_value"`

	// [Optional] If true, empty nested objects are sent to / read from the
	// API instead of flattened to null.
	// The difference between this and send_empty_value is that send_empty_value
	// applies when the key of an object is empty; this applies when the values
	// are all nil / default. eg: "expiration: null" vs "expiration: {}"
	// In the case of Terraform, this occurs when a block in config has optional
	// values, and none of them are used. Terraform returns a nil instead of an
	// empty map[string]interface{} like we'd expect.
	AllowEmptyObject bool `yaml:"allow_empty_object"`

	MinVersion string `yaml:"min_version"`

	ExactVersion string `yaml:"exact_version"`

	// A list of properties that conflict with this property. Uses the "lineage"
	// field to identify the property eg: parent.meta.label.foo
	Conflicts []string

	// A list of properties that at least one of must be set.
	AtLeastOneOf []string `yaml:"at_least_one_of"`

	// A list of properties that exactly one of must be set.
	ExactlyOneOf []string `yaml:"exactly_one_of"`

	// A list of properties that are required to be set together.
	RequiredWith []string `yaml:"required_with"`

	// Can only be overridden - we should never set this ourselves.
	NewType string

	// A pattern that maps expected user input to expected API input.
	// TODO: remove?
	Pattern string

	Properties []*Type

	EnumValues []string `yaml:"enum_values"`

	SkipDocsValues bool `yaml:"skip_docs_values"`

	// ====================
	// Array Fields
	// ====================
	ItemType *Type `yaml:"item_type"`
	MinSize  int   `yaml:"min_size"`
	MaxSize  int   `yaml:"max_size"`
	// __name
	ParentName string

	// ====================
	// ResourceRef Fields
	// ====================
	Resource string
	Imports  string

	// ====================
	// Terraform Overrides
	// ====================

	// Adds a DiffSuppressFunc to the schema
	DiffSuppressFunc string `yaml:"diff_suppress_func"`

	StateFunc string `yaml:"state_func"` // Adds a StateFunc to the schema

	Sensitive bool // Adds `Sensitive: true` to the schema

	// Does not set this value to the returned API value.  Useful for fields
	// like secrets where the returned API value is not helpful.
	IgnoreRead bool `yaml:"ignore_read"`

	// Adds a ValidateFunc to the schema
	Validation bool

	// Indicates that this is an Array that should have Set diff semantics.
	UnorderedList bool `yaml:"unordered_list"`

	IsSet bool `yaml:"is_set"` // Uses a Set instead of an Array

	// Optional function to determine the unique ID of an item in the set
	// If not specified, schema.HashString (when elements are string) or
	// schema.HashSchema are used.
	SetHashFunc string `yaml:"set_hash_func"`

	// if true, then we get the default value from the Google API if no value
	// is set in the terraform configuration for this field.
	// It translates to setting the field to Computed & Optional in the schema.
	// For nested fields, this only applies at the current level. This means
	// it should be explicitly added to each field that needs the defaulting
	// behavior.
	DefaultFromApi bool `yaml:"default_from_api"`

	// https://github.com/hashicorp/terraform/pull/20837
	// Apply a ConfigMode of SchemaConfigModeAttr to the field.
	// This should be avoided for new fields, and only used with old ones.
	SchemaConfigModeAttr bool `yaml:"schema_config_mode_attr"`

	// Names of fields that should be included in the updateMask.
	UpdateMaskFields []string `yaml:"update_mask_fields"`

	// For a TypeMap, the expander function to call on the key.
	// Defaults to expandString.
	KeyExpander string `yaml:"key_expander"`

	// For a TypeMap, the DSF to apply to the key.
	KeyDiffSuppressFunc string `yaml:"key_diff_suppress_func"`

	// ====================
	// Map Fields
	// ====================
	// The type definition of the contents of the map.
	ValueType *Type `yaml:"value_type"`

	// While the API doesn't give keys an explicit name, we specify one
	// because in Terraform the key has to be a property of the object.
	//
	// The name of the key. Used in the Terraform schema as a field name.
	KeyName string `yaml:"key_name`

	// A description of the key's format. Used in Terraform to describe
	// the field in documentation.
	KeyDescription string `yaml:"key_description`

	// ====================
	// KeyValuePairs Fields
	// ====================
	IgnoreWrite bool `yaml:"ignore_write`

	// ====================
	// Schema Modifications
	// ====================
	// Schema modifications change the schema of a resource in some
	// fundamental way. They're not very portable, and will be hard to
	// generate so we should limit their use. Generally, if you're not
	// converting existing Terraform resources, these shouldn't be used.
	//
	// With great power comes great responsibility.

	// Flattens a NestedObject by removing that field from the Terraform
	// schema but will preserve it in the JSON sent/retrieved from the API
	//
	// EX: a API schema where fields are nested (eg: `one.two.three`) and we
	// desire the properties of the deepest nested object (eg: `three`) to
	// become top level properties in the Terraform schema. By overriding
	// the properties `one` and `one.two` and setting flatten_object then
	// all the properties in `three` will be at the root of the TF schema.
	//
	// We need this for cases where a field inside a nested object has a
	// default, if we can't spend a breaking change to fix a misshapen
	// field, or if the UX is _much_ better otherwise.
	//
	// WARN: only fully flattened properties are currently supported. In the
	// example above you could not flatten `one.two` without also flattening
	// all of it's parents such as `one`
	FlattenObject bool `yaml:"flatten_object"`

	// ===========
	// Custom code
	// ===========
	// All custom code attributes are string-typed.  The string should
	// be the name of a template file which will be compiled in the
	// specified / described place.

	// A custom expander replaces the default expander for an attribute.
	// It is called as part of Create, and as part of Update if
	// object.input is false.  It can return an object of any type,
	// so the function header *is* part of the custom code template.
	// As with flatten, `property` and `prefix` are available.
	CustomExpand string `yaml:"custom_expand"`

	// A custom flattener replaces the default flattener for an attribute.
	// It is called as part of Read.  It can return an object of any
	// type, and may sometimes need to return an object with non-interface{}
	// type so that the d.Set() call will succeed, so the function
	// header *is* a part of the custom code template.  To help with
	// creating the function header, `property` and `prefix` are available,
	// just as they are in the standard flattener template.
	CustomFlatten string `yaml:"custom_flatten"`

	ResourceMetadata *Resource

	ParentMetadata *Type // is nil for top-level properties
}

const MAX_NAME = 20

func (t *Type) SetDefault(r *Resource) {
	t.ResourceMetadata = r
	if t.UpdateVerb == "" {
		t.UpdateVerb = t.ResourceMetadata.UpdateVerb
	}

	switch {
	case t.IsA("Array"):
		t.ItemType.ParentName = t.Name
		t.ItemType.ParentMetadata = t.ParentMetadata
		t.ItemType.SetDefault(r)
	case t.IsA("Map"):
		t.KeyExpander = "tpgresource.ExpandString"
		t.ValueType.ParentName = t.Name
		t.ValueType.ParentMetadata = t.ParentMetadata
		t.ValueType.SetDefault(r)
	case t.IsA("NestedObject"):
		if t.Name == "" {
			t.Name = t.ParentName
		}

		if t.Description == "" {
			t.Description = "A nested object resource"
		}

		for _, p := range t.Properties {
			p.ParentMetadata = t
			p.SetDefault(r)
		}
	case t.IsA("ResourceRef"):
		if t.Name == "" {
			t.Name = t.Resource
		}

		if t.Description == "" {
			t.Description = fmt.Sprintf("A reference to %s resource", t.Resource)
		}
	default:
	}
}

// super
// check :description, type: ::String, required: true
// check :exclude, type: :boolean, default: false, required: true
// check :deprecation_message, type: ::String
// check :removed_message, type: ::String
// check :min_version, type: ::String
// check :exact_version, type: ::String
// check :output, type: :boolean
// check :required, type: :boolean
// check :send_empty_value, type: :boolean
// check :allow_empty_object, type: :boolean
// check :url_param_only, type: :boolean
// check :read_query_params, type: ::String
// check :immutable, type: :boolean

// raise 'Property cannot be output and required at the same time.' \
//   if @output && @required

// check :update_verb, type: Symbol, allowed: %i[POST PUT PATCH NONE],
//                     default: @__resource&.update_verb

// check :update_url, type: ::String
// check :update_id, type: ::String
// check :fingerprint_name, type: ::String
// check :pattern, type: ::String

// check_default_value_property
// check_conflicts
// check_at_least_one_of
// check_exactly_one_of
// check_required_with

// check :sensitive, type: :boolean, default: false
// check :is_set, type: :boolean, default: false
// check :default_from_api, type: :boolean, default: false
// check :unordered_list, type: :boolean, default: false
// check :schema_config_mode_attr, type: :boolean, default: false

// // technically set as a default everywhere, but only maps will use this.
// check :key_expander, type: ::String, default: 'tpgresource.ExpandString'
// check :key_diff_suppress_func, type: ::String

// check :diff_suppress_func, type: ::String
// check :state_func, type: ::String
// check :validation, type: Provider::Terraform::Validation
// check :set_hash_func, type: ::String

// check :custom_flatten, type: ::String
// check :custom_expand, type: ::String

// raise "'default_value' and 'default_from_api' cannot be both set" \
//   if @default_from_api && !@default_value.nil?
// }

// func (t *Type) to_s() {
// JSON.pretty_generate(self)
// }

// Prints a dot notation path to where the field is nested within the parent
// object. eg: parent.meta.label.foo
// The only intended purpose is to allow better error messages. Some objects
// and at some points in the build this doesn't output a valid output.

// def lineage
func (t Type) Lineage() string {
	if t.ParentMetadata == nil {
		return google.Underscore(t.Name)
	}

	return fmt.Sprintf("%s.%s", t.ParentMetadata.Lineage(), google.Underscore(t.Name))
}

// Prints the access path of the field in the configration eg: metadata.0.labels
// The only intended purpose is to get the value of the labes field by calling d.Get().
// func (t *Type) terraform_lineage() {
func (t Type) TerraformLineage() string {
	if t.ParentMetadata == nil || t.ParentMetadata.FlattenObject {
		return google.Underscore(t.Name)
	}

	return fmt.Sprintf("%s.0.%s", t.ParentMetadata.TerraformLineage(), google.Underscore(t.Name))
}

func (t Type) EnumValuesToString() string {
	var values []string

	for _, val := range t.EnumValues {
		values = append(values, fmt.Sprintf("`%s`", val))
	}

	return strings.Join(values, ", ")
}

// func (t *Type) to_json(opts) {
// ignore fields that will contain references to parent resources and
// those which will be added later
// ignored_fields = %i[@resource @__parent @__resource @api_name @update_verb
//                     @__name @name @properties]
// json_out = {}

// instance_variables.each do |v|
//   if v == :@conflicts && instance_variable_get(v).empty?
//     // ignore empty conflict arrays
//   elsif v == :@at_least_one_of && instance_variable_get(v).empty?
//     // ignore empty at_least_one_of arrays
//   elsif v == :@exactly_one_of && instance_variable_get(v).empty?
//     // ignore empty exactly_one_of arrays
//   elsif v == :@required_with && instance_variable_get(v).empty?
//     // ignore empty required_with arrays
//   elsif instance_variable_get(v) == false || instance_variable_get(v).nil?
//     // ignore false booleans as non-existence indicates falsey
//   elsif !ignored_fields.include? v
//     json_out[v] = instance_variable_get(v)
//   end
// end

// // convert properties to a hash based on name for nested readability
// json_out.merge!(properties&.map { |p| [p.name, p] }.to_h) \
//   if respond_to? 'properties'

// JSON.generate(json_out, opts)
// }

// func (t *Type) check_default_value_property() {
// return if @default_value.nil?

// case self
// when Api::Type::String
//   clazz = ::String
// when Api::Type::Integer
//   clazz = ::Integer
// when Api::Type::Double
//   clazz = ::Float
// when Api::Type::Enum
//   clazz = ::Symbol
// when Api::Type::Boolean
//   clazz = :boolean
// when Api::Type::ResourceRef
//   clazz = [::String, ::Hash]
// else
//   raise "Update 'check_default_value_property' method to support " \
//         "default value for type //{self.class}"
// end

// check :default_value, type: clazz
// }

// Checks that all conflicting properties actually exist.
// This currently just returns if empty, because we don't want to do the check, since
// this list will have a full path for nested attributes.
// func (t *Type) check_conflicts() {
// check :conflicts, type: ::Array, default: [], item_type: ::String

// return if @conflicts.empty?
// }

// Returns list of properties that are in conflict with this property.
// func (t *Type) conflicting() {
func (t Type) Conflicting() []string {
	if t.ResourceMetadata == nil {
		return []string{}
	}

	return t.Conflicts
}

// Checks that all properties that needs at least one of their fields actually exist.
// This currently just returns if empty, because we don't want to do the check, since
// this list will have a full path for nested attributes.
// func (t *Type) check_at_least_one_of() {
// check :at_least_one_of, type: ::Array, default: [], item_type: ::String

// return if @at_least_one_of.empty?
// }

// Returns list of properties that needs at least one of their fields set.
// func (t *Type) at_least_one_of_list() {
func (t Type) AtLeastOneOfList() []string {
	if t.ResourceMetadata == nil {
		return []string{}
	}

	return t.AtLeastOneOf
}

// Checks that all properties that needs exactly one of their fields actually exist.
// This currently just returns if empty, because we don't want to do the check, since
// this list will have a full path for nested attributes.
// func (t *Type) check_exactly_one_of() {
// check :exactly_one_of, type: ::Array, default: [], item_type: ::String

// return if @exactly_one_of.empty?
// }

// Returns list of properties that needs exactly one of their fields set.
// func (t *Type) exactly_one_of_list() {
func (t Type) ExactlyOneOfList() []string {
	if t.ResourceMetadata == nil {
		return []string{}
	}

	return t.ExactlyOneOf
}

// Checks that all properties that needs required with their fields actually exist.
// This currently just returns if empty, because we don't want to do the check, since
// this list will have a full path for nested attributes.
// func (t *Type) check_required_with() {
// check :required_with, type: ::Array, default: [], item_type: ::String

// return if @required_with.empty?
// }

// Returns list of properties that needs required with their fields set.
// func (t *Type) required_with_list() {
func (t Type) RequiredWithList() []string {
	if t.ResourceMetadata == nil {
		return []string{}
	}

	return t.RequiredWith
}

func (t Type) Parent() *Type {
	return t.ParentMetadata
}

// def min_version
func (t Type) MinVersionObj() *product.Version {
	if t.MinVersion != "" {
		return t.ResourceMetadata.ProductMetadata.versionObj(t.MinVersion)
	} else {
		return t.ResourceMetadata.MinVersionObj()
	}
}

// def exact_version
func (t *Type) exactVersionObj() *product.Version {
	if t.ExactVersion == "" {
		return nil
	}

	return t.ResourceMetadata.ProductMetadata.versionObj(t.ExactVersion)
}

// def exclude_if_not_in_version!(version)
func (t *Type) ExcludeIfNotInVersion(version *product.Version) {
	if !t.Exclude {
		if versionObj := t.exactVersionObj(); versionObj != nil {
			t.Exclude = versionObj.CompareTo(version) != 0
		}

		if !t.Exclude {
			t.Exclude = version.CompareTo(t.MinVersionObj()) < 0
		}
	}

	if t.IsA("NestedObject") {
		for _, p := range t.Properties {
			p.ExcludeIfNotInVersion(version)
		}
	} else if t.IsA("Array") && t.ItemType.IsA("NestedObject") {
		t.ItemType.ExcludeIfNotInVersion(version)
	}
}

// Overriding is_a? to enable class overrides.
// Ruby does not let you natively change types, so this is the next best
// thing.

// TODO Q1: check the type of superclasses of property t
// func (t *Type) is_a?(clazz) {
func (t Type) IsA(clazz string) bool {
	if clazz == "" {
		log.Fatalf("class cannot be empty")
	}

	if t.NewType != "" {
		return t.NewType == clazz
	}

	return t.Type == clazz
	// super(clazz)
}

// // Overriding class to enable class overrides.
// // Ruby does not let you natively change types, so this is the next best
// // thing.
// func (t *Type) class() {
//   // return Module.const_get(@new_type) if @new_type

//   // super
// }

// Returns nested properties for this property.
// def nested_properties
func (t Type) NestedProperties() []*Type {
	props := make([]*Type, 0)

	switch {
	case t.IsA("Array"):
		if t.ItemType.IsA("NestedObject") {
			props = google.Reject(t.ItemType.NestedProperties(), func(p *Type) bool {
				return t.Exclude
			})
		}
	case t.IsA("NestedObject"):
		props = t.UserProperties()
	case t.IsA("Map"):
		props = google.Reject(t.ValueType.NestedProperties(), func(p *Type) bool {
			return t.Exclude
		})
	default:
	}
	return props
}

// def removed?
func (t Type) Removed() bool {
	return t.RemovedMessage != ""
}

// def deprecated?
func (t Type) Deprecated() bool {
	return t.DeprecationMessage != ""
}

// // private

// // A constant value to be provided as field
// type Constant struct {
// // < Type
//   value

//   func (t *Type) validate
//     @description = "This is always //{value}."
//     super
//   end
// }

// // Represents a primitive (non-composite) type.
// class Primitive < Type
// end

// // Represents a boolean
// class Boolean < Primitive
// end

// // Represents an integer
// class Integer < Primitive
// end

// // Represents a double
// class Double < Primitive
// end

// // Represents a string
// class String < Primitive
//   func (t *Type) initialize(name = nil)
//     super()

//     @name = name
//   end

//   PROJECT = Api::Type::String.new('project')
//   NAME = Api::Type::String.new('name')
// end

// // Properties that are fetched externally
// class FetchedExternal < Type

//   func (t *Type) validate
//     @conflicts ||= []
//     @at_least_one_of ||= []
//     @exactly_one_of ||= []
//     @required_with ||= []
//   end

//   func (t *Type) api_name
//     name
//   end
// end

// class Path < Primitive
// end

// // Represents a fingerprint.  A fingerprint is an output-only
// // field used for optimistic locking during updates.
// // They are fetched from the GCP response.
// class Fingerprint < FetchedExternal
//   func (t *Type) validate
//     super
//     @output = true if @output.nil?
//   end
// end

// // Represents a timestamp
// class Time < Primitive
// end

// // A base class to tag objects that are composed by other objects (arrays,
// // nested objects, etc)
// class Composite < Type
// end

// // Forwarding declaration to allow defining Array::NESTED_ARRAY_TYPE
// class NestedObject < Composite
// end

// // Forwarding declaration to allow defining Array::RREF_ARRAY_TYPE
// class ResourceRef < Type
// end

// // Represents an array, and stores its items' type
// class Array < Composite
//   item_type
//   min_size
//   max_size

//   func (t *Type) validate
//     super
//     if @item_type.is_a?(NestedObject) || @item_type.is_a?(ResourceRef)
//       @item_type.set_variable(@name, :__name)
//       @item_type.set_variable(@__resource, :__resource)
//       @item_type.set_variable(self, :__parent)
//     end
//     check :item_type, type: [::String, NestedObject, ResourceRef, Enum], required: true

//     unless @item_type.is_a?(NestedObject) || @item_type.is_a?(ResourceRef) \
//         || @item_type.is_a?(Enum) || type?(@item_type)
//       raise "Invalid type //{@item_type}"
//     end

//     check :min_size, type: ::Integer
//     check :max_size, type: ::Integer
//   end

//   func (t *Type) exclude_if_not_in_version!(version)
//     super
//     @item_type.exclude_if_not_in_version!(version) \
//       if @item_type.is_a? NestedObject
//   end

// func (t *Type) nested_properties
// return @item_type.nested_properties.reject(&:exclude) \
// 	if @item_type.is_a?(Api::Type::NestedObject)

// super
// end

// This function is for array field
// def item_type_class
func (t Type) ItemTypeClass() string {
	if !t.IsA("Array") {
		return ""
	}

	return t.ItemType.Type
}

// // Represents an enum, and store is valid values
// class Enum < Primitive
//   values
//   skip_docs_values

//   func (t *Type) validate
//     super
//     check :values, type: ::Array, item_type: [Symbol, ::String, ::Integer], required: true
//     check :skip_docs_values, type: :boolean
//   end

//   func (t *Type) merge(other)
//     result = self.class.new
//     instance_variables.each do |v|
//       result.instance_variable_set(v, instance_variable_get(v))
//     end

//     other.instance_variables.each do |v|
//       if other.instance_variable_get(v).instance_of?(Array)
//         result.instance_variable_set(v, deep_merge(result.instance_variable_get(v),
//                                                     other.instance_variable_get(v)))
//       else
//         result.instance_variable_set(v, other.instance_variable_get(v))
//       end
//     end

//     result
//   end
// end

// // Represents a 'selfLink' property, which returns the URI of the resource.
// class SelfLink < FetchedExternal
//   EXPORT_KEY = 'selfLink'.freeze

//   resource

//   func (t *Type) name
//     EXPORT_KEY
//   end

//   func (t *Type) out_name
//     EXPORT_KEY.underscore
//   end
// end

// // Represents a reference to another resource
// class ResourceRef < Type
//   // The fields which can be overridden in provider.yaml.
//   module Fields
//     resource
//     imports
//   end
//   include Fields

//   func (t *Type) validate
//     super
//     @name = @resource if @name.nil?
//     @description = "A reference to //{@resource} resource" \
//       if @description.nil?

//     return if @__resource.nil? || @__resource.exclude || @exclude

//     check :resource, type: ::String, required: true
//     check :imports, type: ::String, required: TrueClass

//     // TODO: (camthornton) product reference may not exist yet
//     return if @__resource.__product.nil?

//     check_resource_ref_property_exists
//   end

// func (t *Type) resource_ref
func (t Type) ResourceRef() *Resource {
	if !t.IsA("ResourceRef") {
		return nil
	}

	product := t.ResourceMetadata.ProductMetadata
	resources := google.Select(product.Objects, func(obj *Resource) bool {
		return obj.Name == t.Resource
	})

	return resources[0]
}

//   private

//   func (t *Type) check_resource_ref_property_exists
//     return unless defined?(resource_ref.all_user_properties)

//     exported_props = resource_ref.all_user_properties
//     exported_props << Api::Type::String.new('selfLink') \
//       if resource_ref.has_self_link
//     raise "'//{@imports}' does not exist on '//{@resource}'" \
//       if exported_props.none? { |p| p.name == @imports }
//   end
// end

// // An structured object composed of other objects.
// class NestedObject < Composite

//   func (t *Type) validate
//     @description = 'A nested object resource' if @description.nil?
//     @name = @__name if @name.nil?
//     super

//     raise "Properties missing on //{name}" if @properties.nil?

//     @properties.each do |p|
//       p.set_variable(@__resource, :__resource)
//       p.set_variable(self, :__parent)
//     end
//     check :properties, type: ::Array, item_type: Api::Type, required: true
//   end

// Returns all properties including the ones that are excluded
// This is used for PropertyOverride validation
// def all_properties
func (t Type) AllProperties() []*Type {
	return t.Properties
}

// func (t *Type) properties
func (t Type) UserProperties() []*Type {
	if t.IsA("NestedObject") {
		if t.Properties == nil {
			log.Fatalf("Field '{%s}' properties are nil!", t.Lineage())
		}

		return google.Reject(t.Properties, func(p *Type) bool {
			return p.Exclude
		})
	}
	return nil
}

// Returns the list of top-level properties once any nested objects with
// flatten_object set to true have been collapsed
//
//	func (t *Type) root_properties
func (t *Type) RootProperties() []*Type {
	props := make([]*Type, 0)
	for _, p := range t.UserProperties() {
		if p.FlattenObject {
			props = google.Concat(props, p.RootProperties())
		} else {
			props = append(props, p)
		}
	}
	return props
}

//   func (t *Type) exclude_if_not_in_version!(version)
//     super
//     @properties.each { |p| p.exclude_if_not_in_version!(version) }
//   end
// end

// An array of string -> string key -> value pairs, such as labels.
// While this is technically a map, it's split out because it's a much
// simpler property to generate and means we can avoid conditional logic
// in Map.

func NewProperty(name, apiName string, options []func(*Type)) *Type {
	p := &Type{
		NamedObject: NamedObject{
			Name:    name,
			ApiName: apiName,
		},
	}

	for _, option := range options {
		option(p)
	}
	return p
}

func propertyWithType(t string) func(*Type) {
	return func(p *Type) {
		p.Type = t
	}
}

func propertyWithOutput(output bool) func(*Type) {
	return func(p *Type) {
		p.Output = output
	}
}

func propertyWithDescription(description string) func(*Type) {
	return func(p *Type) {
		p.Description = description
	}
}

func propertyWithMinVersion(minVersion string) func(*Type) {
	return func(p *Type) {
		p.MinVersion = minVersion
	}
}

func propertyWithUpdateVerb(updateVerb string) func(*Type) {
	return func(p *Type) {
		p.UpdateVerb = updateVerb
	}
}

func propertyWithUpdateUrl(updateUrl string) func(*Type) {
	return func(p *Type) {
		p.UpdateUrl = updateUrl
	}
}

func propertyWithImmutable(immutable bool) func(*Type) {
	return func(p *Type) {
		p.Immutable = immutable
	}
}

func propertyWithIgnoreWrite(ignoreWrite bool) func(*Type) {
	return func(p *Type) {
		p.IgnoreWrite = ignoreWrite
	}
}

// class KeyValuePairs < Composite
//   // Ignore writing the "effective_labels" and "effective_annotations" fields to API.
//   ignore_write

//   func (t *Type) initialize(name: nil, output: nil, api_name: nil, description: nil, min_version: nil,
//                   ignore_write: nil, update_verb: nil, update_url: nil, immutable: nil)
//     super()

//     @name = name
//     @output = output
//     @api_name = api_name
//     @description = description
//     @min_version = min_version
//     @ignore_write = ignore_write
//     @update_verb = update_verb
//     @update_url = update_url
//     @immutable = immutable
//   end

//   func (t *Type) validate
//     super
//     check :ignore_write, type: :boolean, default: false

//     return if @__resource.__product.nil?

//     product_name = @__resource.__product.name
//     resource_name = @__resource.name

//     if lineage == 'labels' || lineage == 'metadata.labels' ||
//         lineage == 'configuration.labels'
//       if !(is_a? Api::Type::KeyValueLabels) &&
//           // The label value must be empty string, so skip this resource
//           !(product_name == 'CloudIdentity' && resource_name == 'Group') &&

//           // The "labels" field has type Array, so skip this resource
//           !(product_name == 'DeploymentManager' && resource_name == 'Deployment') &&

//           // https://github.com/hashicorp/terraform-provider-google/issues/16219
//           !(product_name == 'Edgenetwork' && resource_name == 'Network') &&

//           // https://github.com/hashicorp/terraform-provider-google/issues/16219
//           !(product_name == 'Edgenetwork' && resource_name == 'Subnet') &&

//           // "userLabels" is the resource labels field
//           !(product_name == 'Monitoring' && resource_name == 'NotificationChannel') &&

//           // The "labels" field has type Array, so skip this resource
//           !(product_name == 'Monitoring' && resource_name == 'MetricDescriptor')
//         raise "Please use type KeyValueLabels for field //{lineage} " \
//               "in resource //{product_name}///{resource_name}"
//       end
//     elsif is_a? Api::Type::KeyValueLabels
//       raise "Please don't use type KeyValueLabels for field //{lineage} " \
//             "in resource //{product_name}///{resource_name}"
//     end

//     if lineage == 'annotations' || lineage == 'metadata.annotations'
//       if !(is_a? Api::Type::KeyValueAnnotations) &&
//           // The "annotations" field has "ouput: true", so skip this eap resource
//           !(product_name == 'Gkeonprem' && resource_name == 'BareMetalAdminClusterEnrollment')
//         raise "Please use type KeyValueAnnotations for field //{lineage} " \
//               "in resource //{product_name}///{resource_name}"
//       end
//     elsif is_a? Api::Type::KeyValueAnnotations
//       raise "Please don't use type KeyValueAnnotations for field //{lineage} " \
//             "in resource //{product_name}///{resource_name}"
//     end
//   end

// def field_min_version
func (t Type) fieldMinVersion() string {
	return t.MinVersion
}

// // An array of string -> string key -> value pairs used specifically for the "labels" field.
// // The field name with this type should be "labels" literally.
// class KeyValueLabels < KeyValuePairs
//   func (t *Type) validate
//     super
//     return unless @name != 'labels'

//     raise "The field //{name} has the type KeyValueLabels, but the field name is not 'labels'!"
//   end
// end

// // An array of string -> string key -> value pairs used for the "terraform_labels" field.
// class KeyValueTerraformLabels < KeyValuePairs
// end

// // An array of string -> string key -> value pairs used for the "effective_labels"
// // and "effective_annotations" fields.
// class KeyValueEffectiveLabels < KeyValuePairs
// end

// // An array of string -> string key -> value pairs used specifically for the "annotations" field.
// // The field name with this type should be "annotations" literally.
// class KeyValueAnnotations < KeyValuePairs
//   func (t *Type) validate
//     super
//     return unless @name != 'annotations'

//     raise "The field //{name} has the type KeyValueAnnotations,\
// but the field name is not 'annotations'!"
//   end
// end

// // Map from string keys -> nested object entries
// class Map < Composite
//   // <provider>.yaml.
//   module Fields
//     // The type definition of the contents of the map.
//     value_type

//     // While the API doesn't give keys an explicit name, we specify one
//     // because in Terraform the key has to be a property of the object.
//     //
//     // The name of the key. Used in the Terraform schema as a field name.
//     key_name

//     // A description of the key's format. Used in Terraform to describe
//     // the field in documentation.
//     key_description
//   end
//   include Fields

//   func (t *Type) validate
//     super
//     check :key_name, type: ::String, required: true
//     check :key_description, type: ::String

//     @value_type.set_variable(@name, :__name)
//     @value_type.set_variable(@__resource, :__resource)
//     @value_type.set_variable(self, :__parent)
//     check :value_type, type: Api::Type::NestedObject, required: true
//     raise "Invalid type //{@value_type}" unless type?(@value_type)
//   end

//   func (t *Type) nested_properties
//     @value_type.nested_properties.reject(&:exclude)
//   end
// end

// // Support for schema ValidateFunc functionality.
// class Validation < Object
//   // Ensures the value matches this regex
//   regex
//   function

//   func (t *Type) validate
//     super

//     check :regex, type: String
//     check :function, type: String
//   end
// end

// func (t *Type) type?(type)
//   type.is_a?(Type) || !get_type(type).nil?
// end

// func (t *Type) get_type(type)
//   Module.const_get(type)
// end

// def property_ns_prefix
func (t Type) PropertyNsPrefix() []string {
	return []string{
		"Google",
		google.Camelize(t.ResourceMetadata.ProductMetadata.Name, "upper"),
		"Property",
	}
}
