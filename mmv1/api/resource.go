// Copyright 2024 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package api

import (
	"bytes"
	"fmt"
	"log"
	"maps"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/golang/glog"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/utils"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

const RELATIVE_MAGICIAN_LOCATION = "mmv1/"
const GITHUB_BASE_URL = "https://github.com/GoogleCloudPlatform/magic-modules/tree/main/" + RELATIVE_MAGICIAN_LOCATION

type Resource struct {
	Name string

	// original value of :name before the provider override happens
	// same as :name if not overridden in provider
	ApiName string `yaml:"api_name,omitempty"`

	// [Required] A description of the resource that's surfaced in provider
	// documentation.
	Description string

	// [Required] Reference links provided in
	// downstream documentation. Expected to follow the format as follows:
	//
	//	references:
	//  	guides:
	//			'Guide name': 'official_documentation_url'
	//		api: 'rest_api_reference_url/version'
	//
	References resource.ReferenceLinks `yaml:"references,omitempty"`

	// [Required] The GCP "relative URI" of a resource, relative to the product
	// base URL. It can often be inferred from the `create` path.
	BaseUrl string `yaml:"base_url,omitempty"`

	// ====================
	// Common Configuration
	// ====================
	//
	// [Optional] The minimum API version this resource is in. Defaults to ga.
	MinVersion string `yaml:"min_version,omitempty"`

	// [Optional] If set to true, don't generate the resource.
	Exclude bool `yaml:"exclude,omitempty"`

	// [Optional] If set to true, the resource is not able to be updated.
	Immutable bool `yaml:"immutable,omitempty"`

	// [Optional] If set to true, this resource uses an update mask to perform
	// updates. This is typical of newer GCP APIs.
	UpdateMask bool `yaml:"update_mask,omitempty"`

	// [Optional] If set to true, the object has a `self_link` field. This is
	// typical of older GCP APIs.
	HasSelfLink bool `yaml:"has_self_link,omitempty"`

	// [Optional] The validator "relative URI" of a resource, relative to the product
	// base URL. Specific to defining the resource as a CAI asset.
	CaiBaseUrl string `yaml:"cai_base_url,omitempty"`

	// ====================
	// URL / HTTP Configuration
	// ====================
	//
	// [Optional] The "identity" URL of the resource. Defaults to:
	// * base_url when the create_verb is POST
	// * self_link when the create_verb is PUT  or PATCH
	SelfLink string `yaml:"self_link,omitempty"`

	// [Optional] The URL used to creating the resource. Defaults to:
	// * collection url when the create_verb is POST
	// * self_link when the create_verb is PUT or PATCH
	CreateUrl string `yaml:"create_url,omitempty"`

	// [Optional] The URL used to delete the resource. Defaults to the self
	// link.
	DeleteUrl string `yaml:"delete_url,omitempty"`

	// [Optional] The URL used to update the resource. Defaults to the self
	// link.
	UpdateUrl string `yaml:"update_url,omitempty"`
	// [Optional] The HTTP verb used during create. Defaults to POST.
	CreateVerb string `yaml:"create_verb,omitempty"`

	// [Optional] The HTTP verb used during read. Defaults to GET.
	ReadVerb string `yaml:"read_verb,omitempty"`

	// [Optional] The HTTP verb used during update. Defaults to PUT.
	UpdateVerb string `yaml:"update_verb,omitempty"`

	// [Optional] The HTTP verb used during delete. Defaults to DELETE.
	DeleteVerb string `yaml:"delete_verb,omitempty"`

	// [Optional] Additional Query Parameters to append to GET. Defaults to ""
	ReadQueryParams string `yaml:"read_query_params,omitempty"`

	// ====================
	// Collection / Identity URL Configuration
	// ====================
	//
	// [Optional] This is the name of the list of items
	// within the collection (list) json. Will default to the
	// camelcase plural name of the resource.
	CollectionUrlKey string `yaml:"collection_url_key,omitempty"`

	// [Optional] An ordered list of names of parameters that uniquely identify
	// the resource.
	// Generally, it's safe to leave empty, in which case it defaults to `name`.
	// Other values are normally useful in cases where an object has a parent
	// and is identified by some non-name value, such as an ip+port pair.
	// If you're writing a fine-grained resource (eg with nested_query) a value
	// must be set.
	Identity []string `yaml:"identity,omitempty"`

	// [Optional] (Api::Resource::NestedQuery) This is useful in case you need
	// to change the query made for GET requests only. In particular, this is
	// often used to extract an object from a parent object or a collection.
	// Note that if both nested_query and custom_code.decoder are provided,
	// the decoder will be included within the code handling the nested query.
	NestedQuery *resource.NestedQuery `yaml:"nested_query,omitempty"`

	// ====================
	// IAM Configuration
	// ====================
	//
	// [Optional] (Api::Resource::IamPolicy) Configuration of a resource's
	// resource-specific IAM Policy.
	IamPolicy *resource.IamPolicy `yaml:"iam_policy,omitempty"`

	// [Optional] If set to true, don't generate the resource itself; only
	// generate the IAM policy.
	// TODO rewrite: rename?
	ExcludeResource bool `yaml:"exclude_resource,omitempty"`

	// [Optional] GCP kind, e.g. `compute//disk`
	Kind string `yaml:"kind,omitempty"`

	// [Optional] If set to true, indicates that a resource is not configurable
	// such as GCP regions.
	Readonly bool `yaml:"readonly,omitempty"`

	// ====================
	// Terraform Overrides
	// ====================
	// [Optional] If non-empty, overrides the full filename prefix
	// i.e. google/resource_product_{{resource_filename_override}}.go
	// i.e. google/resource_product_{{resource_filename_override}}_test.go
	FilenameOverride string `yaml:"filename_override,omitempty"`

	// If non-empty, overrides the full given resource name.
	// i.e. 'google_project' for resourcemanager.Project
	// Use Provider::Terraform::Config.legacy_name to override just
	// product name.
	// Note: This should not be used for vanity names for new products.
	// This was added to handle preexisting handwritten resources that
	// don't match the natural generated name exactly, and to support
	// services with a mix of handwritten and generated resources.
	LegacyName string `yaml:"legacy_name,omitempty"`

	// The Terraform resource id format used when calling //setId(...).
	// For instance, `{{name}}` means the id will be the resource name.
	IdFormat string `yaml:"id_format,omitempty"`

	// Override attribute used to handwrite the formats for generating regex strings
	// that match templated values to a self_link when importing, only necessary when
	// a resource is not adequately covered by the standard provider generated options.
	// Leading a token with `%`
	// i.e. {{%parent}}/resource/{{resource}}
	// will allow that token to hold multiple /'s.
	//
	// Expected to be formatted as follows:
	//
	//	import_format:
	//		- example_import_one
	//		- example_import_two
	//
	ImportFormat []string `yaml:"import_format,omitempty"`

	CustomCode resource.CustomCode `yaml:"custom_code,omitempty"`

	Docs resource.Docs `yaml:"docs,omitempty"`

	// This block inserts entries into the customdiff.All() block in the
	// resource schema -- the code for these custom diff functions must
	// be included in the resource constants or come from tpgresource
	CustomDiff []string `yaml:"custom_diff,omitempty"`

	// Lock name for a mutex to prevent concurrent API calls for a given
	// resource.
	Mutex string `yaml:"mutex,omitempty"`

	// Examples in documentation. Backed by generated tests, and have
	// corresponding OiCS walkthroughs.
	Examples []resource.Examples

	// If true, generates product operation handling logic.
	AutogenAsync bool `yaml:"autogen_async,omitempty"`

	// If true, resource is not importable
	ExcludeImport bool `yaml:"exclude_import,omitempty"`

	// If true, resource should be autogenerated as a data source
	Datasource *resource.Datasource `yaml:"datasource,omitempty"`

	// If true, skip sweeper generation for this resource
	ExcludeSweeper bool `yaml:"exclude_sweeper,omitempty"`

	// Override sweeper settings
	Sweeper resource.Sweeper `yaml:"sweeper,omitempty"`

	Timeouts *Timeouts `yaml:"timeouts,omitempty"`

	// An array of function names that determine whether an error is retryable.
	ErrorRetryPredicates []string `yaml:"error_retry_predicates,omitempty"`

	// An array of function names that determine whether an error is not retryable.
	ErrorAbortPredicates []string `yaml:"error_abort_predicates,omitempty"`

	// Optional attributes for declaring a resource's current version and generating
	// state_upgrader code to the output .go file from files stored at
	// mmv1/templates/terraform/state_migrations/
	// used for maintaining state stability with resources first provisioned on older api versions.
	SchemaVersion int `yaml:"schema_version,omitempty"`

	// From this schema version on, state_upgrader code is generated for the resource.
	// When unset, state_upgrade_base_schema_version defauts to 0.
	// Normally, it is not needed to be set.
	StateUpgradeBaseSchemaVersion int `yaml:"state_upgrade_base_schema_version,omitempty"`

	StateUpgraders bool `yaml:"state_upgraders,omitempty"`

	// Do not apply the default attribution label
	ExcludeAttributionLabel bool `yaml:"exclude_attribution_label,omitempty"`

	// This block inserts the named function and its attribute into the
	// resource schema -- the code for the migrate_state function must
	// be included in the resource constants or come from tpgresource
	// included for backwards compatibility as an older state migration method
	// and should not be used for new resources.
	MigrateState string `yaml:"migrate_state,omitempty"`

	// Set to true for resources that are unable to be deleted, such as KMS keyrings or project
	// level resources such as firebase project
	ExcludeDelete bool `yaml:"exclude_delete,omitempty"`

	// Set to true for resources that are unable to be read from the API, such as
	// public ca external account keys
	ExcludeRead bool `yaml:"exclude_read,omitempty"`

	// Set to true for resources that wish to disable automatic generation of default provider
	// value customdiff functions
	// TODO rewrite: 1 instance used
	ExcludeDefaultCdiff bool `yaml:"exclude_default_cdiff,omitempty"`

	// This enables resources that get their project via a reference to a different resource
	// instead of a project field to use User Project Overrides
	SupportsIndirectUserProjectOverride bool `yaml:"supports_indirect_user_project_override,omitempty"`

	// If true, the resource's project field can be specified as either the short form project
	// id or the long form projects/project-id. The extra projects/ string will be removed from
	// urls and ids. This should only be used for resources that previously supported long form
	// project ids for backwards compatibility.
	LegacyLongFormProject bool `yaml:"legacy_long_form_project,omitempty"`

	// Function to transform a read error so that handleNotFound recognises
	// it as a 404. This should be added as a handwritten fn that takes in
	// an error and returns one.
	ReadErrorTransform string `yaml:"read_error_transform,omitempty"`

	// If true, resources that failed creation will be marked as tainted. As a consequence
	// these resources will be deleted and recreated on the next apply call. This pattern
	// is preferred over deleting the resource directly in post_create_failure hooks.
	TaintResourceOnFailedCreate bool `yaml:"taint_resource_on_failed_create,omitempty"`

	// Add a deprecation message for a resource that's been deprecated in the API.
	DeprecationMessage string `yaml:"deprecation_message,omitempty"`

	Async *Async

	// Tag autogen resources so that we can track them. In the future this will
	// control if a resource is continuously generated from public OpenAPI docs
	AutogenStatus string `yaml:"autogen_status"`

	// If true, this resource generates with the new plugin framework resource template
	FrameworkResource bool `yaml:"plugin_framework,omitempty"`

	// The three groups of []*Type fields are expected to be strictly ordered within a yaml file
	// in the sequence of Virtual Fields -> Parameters -> Properties

	// Virtual fields are Terraform-only fields that control Terraform's
	// behaviour. They don't map to underlying API fields (although they
	// may map to parameters), and will require custom code to be added to
	// control them.
	//
	// Virtual fields are similar to url_param_only fields in that they create
	// a schema entry which is not read from or submitted to the API. However
	// virtual fields are meant to provide toggles for Terraform-specific behavior in a resource
	// (eg: delete_contents_on_destroy) whereas url_param_only fields _should_
	// be used for url construction.
	//
	// Both are resource level fields and do not make sense, and are also not
	// supported, for nested fields. Nested fields that shouldn't be included
	// in API payloads are better handled with custom expand/encoder logic.
	VirtualFields []*Type `yaml:"virtual_fields,omitempty"`

	Parameters []*Type

	Properties []*Type

	ProductMetadata *Product `yaml:"-"`

	// The version name provided by the user through CI
	TargetVersionName string `yaml:"-"`

	// The compiler to generate the downstream files, for example "terraformgoogleconversion-codegen".
	Compiler string `yaml:"-"`

	// The API "resource type kind" used for this resource e.g., "Function".
	// If this is not set, then :name is used instead, which is strongly
	// preferred wherever possible. Its main purpose is for supporting
	// fine-grained resources and legacy resources.
	ApiResourceTypeKind string `yaml:"api_resource_type_kind,omitempty"`

	// The API URL patterns used by this resource that represent variants e.g.,
	// "folders/{folder}/feeds/{feed}". Each pattern must match the value
	// defined in the API exactly. The use of `api_variant_patterns` is only
	// meaningful when the resource type has multiple parent types available.
	// This is commonly used for resources that have a project, folder, and
	// organization variant, however most resources do not need it.
	ApiVariantPatterns []string `yaml:"api_variant_patterns,omitempty"`

	ImportPath     string `yaml:"-"`
	SourceYamlFile string `yaml:"-"`

	// ====================
	// TGC
	// ====================
	TGCResource `yaml:",inline"`
}

type TGCResource struct {
	// If true, exclude resource from Terraform Validator
	// (i.e. terraform-provider-conversion)
	ExcludeTgc bool `yaml:"exclude_tgc,omitempty"`

	// If true, include resource in the new package of TGC (terraform-provider-conversion)
	IncludeInTGCNext bool `yaml:"include_in_tgc_next_DO_NOT_USE,omitempty"`

	// Name of the hcl resource block used in TGC
	TgcHclBlockName string `yaml:"tgc_hcl_block_name,omitempty"`

	// The resource kind in CAI.
	// If this is not set, then :name is used instead.
	// For example: compute.googleapis.com/Address has Address for CaiResourceKind,
	// and compute.googleapis.com/GlobalAddress has GlobalAddress for CaiResourceKind.
	// But they have the same api resource type: address
	CaiResourceKind string `yaml:"cai_resource_kind,omitempty"`

	// If true, the Terraform custom encoder is not applied during tfplan2cai
	TGCIgnoreTerraformEncoder bool `yaml:"tgc_ignore_terraform_encoder,omitempty"`

	// [Optional] The parameter that uniquely identifies the resource.
	// Generally, it shouldn't be set when the identity can be decided.
	// Otherswise, it should be set.
	CaiIdentity string `yaml:"cai_identity,omitempty"`
}

func (r *Resource) UnmarshalYAML(unmarshal func(any) error) error {
	type resourceAlias Resource
	aliasObj := (*resourceAlias)(r)

	err := unmarshal(aliasObj)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resource) SetDefault(product *Product) {
	if r.CreateVerb == "" {
		r.CreateVerb = "POST"
	}
	if r.ReadVerb == "" {
		r.ReadVerb = "GET"
	}
	if r.DeleteVerb == "" {
		r.DeleteVerb = "DELETE"
	}
	if r.UpdateVerb == "" {
		r.UpdateVerb = "PUT"
	}

	if r.ApiName == "" {
		r.ApiName = r.Name
	}
	if r.CollectionUrlKey == "" {
		key := r.Name
		if r.ApiResourceTypeKind != "" {
			key = r.ApiResourceTypeKind
		}
		r.CollectionUrlKey = google.Camelize(google.Plural(key), "lower")
	}
	if r.IdFormat == "" {
		r.IdFormat = r.SelfLinkUri()
	}

	if len(r.VirtualFields) > 0 {
		for _, f := range r.VirtualFields {
			f.ClientSide = true
		}
	}

	r.ProductMetadata = product
	for _, property := range r.AllProperties() {
		property.SetDefault(r)
	}
	for _, vf := range r.VirtualFields {
		vf.SetDefault(r)
	}
	if r.IamPolicy != nil && r.IamPolicy.MinVersion == "" {
		r.IamPolicy.MinVersion = r.MinVersion
	}
	if r.Timeouts == nil {
		r.Timeouts = NewTimeouts()
	}

}

func (r *Resource) Validate() {
	if r.Name == "" {
		log.Fatalf("Missing `name` for resource")
	}

	if r.NestedQuery != nil && r.NestedQuery.IsListOfIds && len(r.Identity) != 1 {
		log.Fatalf("`is_list_of_ids: true` implies resource has exactly one `identity` property")
	}

	// Ensures we have all properties defined
	for _, i := range r.Identity {
		hasIdentify := slices.ContainsFunc(r.AllUserProperties(), func(p *Type) bool {
			return p.Name == i
		})
		if !hasIdentify {
			log.Fatalf("Missing property/parameter for identity %s", i)
		}
	}

	if r.Description == "" {
		log.Fatalf("Missing `description` for resource %s", r.Name)
	}

	if !r.Exclude {
		if len(r.Properties) == 0 {
			log.Fatalf("Missing `properties` for resource %s", r.Name)
		}
	}

	allowed := []string{"POST", "PUT", "PATCH"}
	if !slices.Contains(allowed, r.CreateVerb) {
		log.Fatalf("Value on `create_verb` should be one of %#v", allowed)
	}

	allowed = []string{"GET", "POST"}
	if !slices.Contains(allowed, r.ReadVerb) {
		log.Fatalf("Value on `read_verb` should be one of %#v", allowed)
	}

	allowed = []string{"POST", "PUT", "PATCH", "DELETE"}
	if !slices.Contains(allowed, r.DeleteVerb) {
		log.Fatalf("Value on `delete_verb` should be one of %#v", allowed)
	}

	allowed = []string{"POST", "PUT", "PATCH"}
	if !slices.Contains(allowed, r.UpdateVerb) {
		log.Fatalf("Value on `update_verb` should be one of %#v", allowed)
	}

	for _, property := range r.AllProperties() {
		property.Validate(r.Name)
	}

	if r.IamPolicy != nil {
		r.IamPolicy.Validate(r.Name)
	}

	if r.NestedQuery != nil {
		r.NestedQuery.Validate(r.Name)
	}

	for _, example := range r.Examples {
		example.Validate(r.Name)
	}

	if r.Async != nil {
		r.Async.Validate()
	}
}

// ====================
// Custom Getters and Setters
// ====================

// Returns all properties and parameters including the ones that are
// excluded. This is used for PropertyOverride validation
func (r Resource) AllProperties() []*Type {
	return google.Concat(r.Properties, r.Parameters)
}

func (r Resource) AllPropertiesInVersion() []*Type {
	return google.Reject(google.Concat(r.Properties, r.Parameters), func(p *Type) bool {
		return p.Exclude
	})
}

func (r Resource) PropertiesWithExcluded() []*Type {
	return r.Properties
}

func (r Resource) UserProperites() []*Type {
	return google.Reject(r.Properties, func(p *Type) bool {
		return p.Exclude
	})
}

func (r Resource) UserParameters() []*Type {
	return google.Reject(r.Parameters, func(p *Type) bool {
		return p.Exclude
	})
}

func (r Resource) UserVirtualFields() []*Type {
	return google.Reject(r.VirtualFields, func(p *Type) bool {
		return p.Exclude
	})
}

func (r Resource) ServiceVersion() string {
	if r.CaiBaseUrl != "" {
		return extractVersionFromBaseUrl(r.CaiBaseUrl)
	}
	return extractVersionFromBaseUrl(r.BaseUrl)
}

func extractVersionFromBaseUrl(baseUrl string) string {
	parts := strings.Split(baseUrl, "/")
	// starts with v...
	if parts[0] != "" && parts[0][0] == 'v' {
		return parts[0]
	}
	// starts with /v...
	if parts[0] == "" && parts[1][0] == 'v' {
		return parts[1]
	}
	return ""
}

// Return the user-facing properties in client tools; this ends up meaning
// both properties and parameters but without any that are excluded due to
// version mismatches or manual exclusion
func (r Resource) AllUserProperties() []*Type {
	return google.Concat(r.UserProperites(), r.UserParameters())
}

func (r Resource) RequiredProperties() []*Type {
	return google.Select(r.AllUserProperties(), func(p *Type) bool {
		return p.Required
	})
}

func (r Resource) AllNestedProperties(props []*Type) []*Type {
	nested := props
	for _, prop := range props {
		if nestedProperties := prop.NestedProperties(); !prop.FlattenObject && nestedProperties != nil {
			nested = google.Concat(nested, r.AllNestedProperties(nestedProperties))
		}
	}

	return nested
}

func (r Resource) SensitiveProps() []*Type {
	props := r.AllNestedProperties(r.RootProperties())
	return google.Select(props, func(p *Type) bool {
		return p.Sensitive
	})
}

func (r Resource) WriteOnlyProps() []*Type {
	props := r.AllNestedProperties(r.RootProperties())
	return google.Select(props, func(p *Type) bool {
		return p.WriteOnlyLegacy
	})
}

func (r Resource) SensitivePropsToString() string {
	var props []string

	for _, prop := range r.SensitiveProps() {
		props = append(props, fmt.Sprintf("`%s`", prop.Lineage()))
	}

	return strings.Join(props, ", ")
}

func (r Resource) WriteOnlyPropsToString() string {
	var props []string

	for _, prop := range r.WriteOnlyProps() {
		props = append(props, fmt.Sprintf("`%s`", prop.Lineage()))
	}

	return strings.Join(props, ", ")
}

// All settable properties in the resource.
// Fingerprints aren't *really" settable properties, but they behave like one.
// At Create, they have no value but they can just be read in anyways, and after a Read
// they will need to be set in every Update.
func (r Resource) SettableProperties() []*Type {
	props := make([]*Type, 0)

	props = google.Reject(r.AllUserProperties(), func(v *Type) bool {
		return v.Output && !v.IsA("Fingerprint") && !v.IsA("KeyValueEffectiveLabels")
	})

	props = google.Reject(props, func(v *Type) bool {
		return v.UrlParamOnly
	})

	props = google.Reject(props, func(v *Type) bool {
		return v.IsA("KeyValueLabels") || v.IsA("KeyValueAnnotations")
	})

	return props
}

func (r Resource) IsSettableProperty(t *Type) bool {
	return slices.Contains(r.SettableProperties(), t)
}

func (r Resource) UnorderedListProperties() []*Type {
	return google.Select(r.SettableProperties(), func(t *Type) bool {
		return t.UnorderedList
	})
}

// Properties that will be returned in the API body
func (r Resource) GettableProperties() []*Type {
	return google.Reject(r.AllUserProperties(), func(v *Type) bool {
		return v.UrlParamOnly
	})
}

// Returns the list of top-level properties once any nested objects with flatten_object
// set to true have been collapsed
func (r Resource) RootProperties() []*Type {
	props := make([]*Type, 0)

	for _, p := range r.AllUserProperties() {
		if p.FlattenObject {
			props = google.Concat(props, p.RootProperties())
		} else {
			props = append(props, p)
		}
	}
	return props
}

// Returns a sorted list of all "leaf" properties, meaning properties that have
// no children.
func (r Resource) LeafProperties() []*Type {
	types := r.AllNestedProperties(google.Concat(r.RootProperties(), r.UserVirtualFields()))

	// Remove types that have children, because we only want "leaf" fields
	types = slices.DeleteFunc(types, func(t *Type) bool {
		nestedProperties := t.NestedProperties()
		return len(nestedProperties) > 0
	})

	// Sort types by lineage
	slices.SortFunc(types, func(a, b *Type) int {
		if a.MetadataLineage() < b.MetadataLineage() {
			return -1
		}
		return 1
	})

	return types
}

// Return the product-level async object, or the resource-specific one
// if one exists.
func (r Resource) GetAsync() *Async {
	if r.Async != nil {
		return r.Async
	}

	return r.ProductMetadata.Async
}

// Return the resource-specific identity properties, or a best guess of the
// `name` value for the resource.
func (r Resource) GetIdentity() []*Type {
	props := r.AllUserProperties()

	if r.Identity != nil {
		identities := google.Select(props, func(p *Type) bool {
			return slices.Contains(r.Identity, p.Name)
		})

		slices.SortFunc(identities, func(a, b *Type) int {
			return slices.Index(r.Identity, a.Name) - slices.Index(r.Identity, b.Name)
		})

		return identities
	}

	return google.Select(props, func(p *Type) bool {
		return p.Name == "name"
	})
}

func (r *Resource) AddLabelsRelatedFields(props []*Type, parent *Type) []*Type {
	for _, p := range props {
		if p.IsA("KeyValueLabels") {
			props = r.addLabelsFields(props, parent, p)
		} else if p.IsA("KeyValueAnnotations") {
			props = r.addAnnotationsFields(props, parent, p)
		} else if p.IsA("NestedObject") && len(p.AllProperties()) > 0 {
			p.Properties = r.AddLabelsRelatedFields(p.AllProperties(), p)
		}
	}
	return props
}

func (r *Resource) addLabelsFields(props []*Type, parent *Type, labels *Type) []*Type {
	if parent == nil || parent.FlattenObject {
		if r.ExcludeAttributionLabel {
			r.CustomDiff = append(r.CustomDiff, "tpgresource.SetLabelsDiffWithoutAttributionLabel")
		} else {
			r.CustomDiff = append(r.CustomDiff, "tpgresource.SetLabelsDiff")
		}
	} else if parent.Name == "metadata" {
		r.CustomDiff = append(r.CustomDiff, "tpgresource.SetMetadataLabelsDiff")
	}

	terraformLabelsField := buildTerraformLabelsField("labels", parent, labels)
	effectiveLabelsField := buildEffectiveLabelsField("labels", labels)
	props = append(props, terraformLabelsField, effectiveLabelsField)

	// The effective_labels field is used to write to API, instead of the labels field.
	labels.IgnoreWrite = true
	labels.Description = fmt.Sprintf("%s\n\n%s", labels.Description, getLabelsFieldNote(labels.Name))

	if parent == nil {
		labels.Immutable = false
	}

	return props
}

func (r *Resource) HasLabelsField() bool {
	for _, p := range r.Properties {
		if p.Name == "labels" {
			return true
		}
	}
	return false
}

func (r *Resource) addAnnotationsFields(props []*Type, parent *Type, annotations *Type) []*Type {

	// The effective_annotations field is used to write to API,
	// instead of the annotations field.
	annotations.IgnoreWrite = true
	annotations.Description = fmt.Sprintf("%s\n\n%s", annotations.Description, getLabelsFieldNote(annotations.Name))

	if parent == nil {
		r.CustomDiff = append(r.CustomDiff, "tpgresource.SetAnnotationsDiff")
	} else if parent.Name == "metadata" {
		r.CustomDiff = append(r.CustomDiff, "tpgresource.SetMetadataAnnotationsDiff")
	}

	effectiveAnnotationsField := buildEffectiveLabelsField("annotations", annotations)
	props = append(props, effectiveAnnotationsField)
	return props
}

func buildEffectiveLabelsField(name string, labels *Type) *Type {
	description := fmt.Sprintf("All of %s (key/value pairs) present on the resource in GCP, "+
		"including the %s configured through Terraform, other clients and services.", name, name)

	t := "KeyValueEffectiveLabels"

	n := fmt.Sprintf("effective%s", strings.Title(name))

	options := []func(*Type){
		propertyWithType(t),
		propertyWithOutput(true),
		propertyWithDescription(description),
		propertyWithMinVersion(labels.fieldMinVersion()),
		propertyWithUpdateVerb(labels.UpdateVerb),
		propertyWithUpdateUrl(labels.UpdateUrl),
		propertyWithImmutable(labels.Immutable),
	}
	return NewProperty(n, name, options)
}

func buildTerraformLabelsField(name string, parent *Type, labels *Type) *Type {
	description := fmt.Sprintf("The combination of %s configured directly on the resource\n"+
		" and default %s configured on the provider.", name, name)

	immutable := false
	if parent != nil {
		immutable = labels.Immutable
	}

	n := fmt.Sprintf("terraform%s", strings.Title(name))

	options := []func(*Type){
		propertyWithType("KeyValueTerraformLabels"),
		propertyWithOutput(true),
		propertyWithDescription(description),
		propertyWithMinVersion(labels.fieldMinVersion()),
		propertyWithIgnoreWrite(true),
		propertyWithUpdateUrl(labels.UpdateUrl),
		propertyWithImmutable(immutable),
	}
	return NewProperty(n, name, options)
}

// Check if the resource has root "labels" field
func (r Resource) RootLabels() bool {
	for _, p := range r.RootProperties() {
		if p.IsA("KeyValueLabels") {
			return true
		}
	}
	return false
}

// Return labels fields that should be added to ImportStateVerifyIgnore
func (r Resource) IgnoreReadLabelsFields(props []*Type) []string {
	fields := make([]string, 0)
	for _, p := range props {
		if p.IsA("KeyValueLabels") ||
			p.IsA("KeyValueTerraformLabels") ||
			p.IsA("KeyValueAnnotations") {
			fields = append(fields, p.TerraformLineage())
		} else if p.IsA("NestedObject") && len(p.AllProperties()) > 0 {
			fields = google.Concat(fields, r.IgnoreReadLabelsFields(p.AllProperties()))
		}
	}
	return fields
}

func getLabelsFieldNote(title string) string {
	return fmt.Sprintf(
		"**Note**: This field is non-authoritative, and will only manage the %s present "+
			"in your configuration.\n"+
			"Please refer to the field `effective_%s` for all of the %s present on the resource.",
		title, title, title)
}

func (r Resource) StateMigrationFile() string {
	return fmt.Sprintf("templates/terraform/state_migrations/%s_%s.go.tmpl", google.Underscore(r.ProductMetadata.Name), google.Underscore(r.Name))
}

// ====================
// Version-related methods
// ====================
func (r Resource) MinVersionObj() *product.Version {
	if r.MinVersion != "" {
		return r.ProductMetadata.versionObj(r.MinVersion)
	} else {
		return r.ProductMetadata.lowestVersion()
	}
}

func (r Resource) NotInVersion(version *product.Version) bool {
	return version.CompareTo(r.MinVersionObj()) < 0
}

// Recurses through all nested properties and parameters and changes their
// 'exclude' instance variable if the property is at a version below the
// one that is passed in.
func (r *Resource) ExcludeIfNotInVersion(version *product.Version) {
	if !r.Exclude {
		r.Exclude = r.NotInVersion(version)
	}

	if r.Properties != nil {
		for _, p := range r.Properties {
			p.ExcludeIfNotInVersion(version)
		}
	}

	if r.Parameters != nil {
		for _, p := range r.Parameters {
			p.ExcludeIfNotInVersion(version)
		}
	}
}

// ====================
// URL-related methods
// ====================

// Returns the "self_link_url" which is generally really the resource's GET
// URL. In older resources generally, this was the self_link value & was the
// product.base_url + resource.base_url + '/name'
// In newer resources there is much less standardisation in terms of value.
// Generally for them though, it's the product.base_url + resource.name
func (r Resource) SelfLinkUrl() string {
	s := []string{r.ProductMetadata.BaseUrl, r.SelfLinkUri()}
	return strings.Join(s, "")
}

// Returns the partial uri / relative path of a resource. In newer resources,
// this is the name. This fn is named self_link_uri for consistency, but
// could otherwise be considered to be "path"
func (r Resource) SelfLinkUri() string {
	// If the terms in this are not snake-cased, this will require
	// an override in Terraform.
	if r.SelfLink != "" {
		return r.SelfLink
	}

	return strings.Join([]string{r.BaseUrl, "{{name}}"}, "/")
}

func (r Resource) CollectionUrl() string {
	s := []string{r.ProductMetadata.BaseUrl, r.collectionUri()}
	return strings.Join(s, "")
}

func (r Resource) collectionUri() string {
	return r.BaseUrl
}

func (r Resource) CreateUri() string {
	if r.CreateUrl != "" {
		return r.CreateUrl
	}

	if r.CreateVerb == "" || r.CreateVerb == "POST" {
		return r.collectionUri()
	}

	return r.SelfLinkUri()
}

func (r Resource) UpdateUri() string {
	if r.UpdateUrl != "" {
		return r.UpdateUrl
	}

	return r.SelfLinkUri()
}

func (r Resource) DeleteUri() string {
	if r.DeleteUrl != "" {
		return r.DeleteUrl
	}

	return r.SelfLinkUri()
}

func (r Resource) ResourceName() string {
	return fmt.Sprintf("%s%s", r.ProductMetadata.Name, r.Name)
}

// Filter the properties to keep only the ones don't have custom update
// method and group them by update url & verb.
func propertiesWithoutCustomUpdate(properties []*Type) []*Type {
	return google.Select(properties, func(p *Type) bool {
		return p.UpdateUrl == "" || p.UpdateVerb == "" || p.UpdateVerb == "NOOP"
	})
}

func (r Resource) UpdateBodyProperties() []*Type {
	updateProp := propertiesWithoutCustomUpdate(r.SettableProperties())
	if r.UpdateVerb == "PATCH" {
		updateProp = google.Reject(updateProp, func(p *Type) bool {
			return p.Immutable
		})
	}
	return updateProp
}

// Handwritten TF Operation objects will be shaped like accessContextManager
// while the Google Go Client will have a name like accesscontextmanager
func (r Resource) ClientNamePascal() string {
	clientName := r.ProductMetadata.ClientName
	if clientName == "" {
		clientName = r.ProductMetadata.Name
	}

	return google.Camelize(clientName, "upper")
}

func (r Resource) PackageName() string {
	return strings.ToLower(r.ProductMetadata.Name)
}

// In order of preference, use TF override,
// general defined timeouts, or default Timeouts
func (r Resource) GetTimeouts() *Timeouts {
	timeoutsFiltered := r.Timeouts
	if timeoutsFiltered == nil {
		if async := r.GetAsync(); async != nil && async.Operation != nil {
			timeoutsFiltered = async.Operation.Timeouts
		}

		if timeoutsFiltered == nil {
			timeoutsFiltered = NewTimeouts()
		}
	}

	return timeoutsFiltered
}

func (r Resource) HasProject() bool {
	return strings.Contains(r.BaseUrl, "{{project}}") || strings.Contains(r.CreateUrl, "{{project}}")
}

func (r Resource) IncludeProjectForOperation() bool {
	return strings.Contains(r.BaseUrl, "{{project}}") || (r.GetAsync().IsA("OpAsync") && r.GetAsync().IncludeProject)
}

func (r Resource) HasRegion() bool {
	found := false
	for _, p := range r.Parameters {
		if p.Name == "region" && p.IgnoreRead {
			found = true
			break
		}
	}
	return found && strings.Contains(r.BaseUrl, "{{region}}")
}

func (r Resource) HasZone() bool {
	found := false
	for _, p := range r.Parameters {
		if p.Name == "zone" && p.IgnoreRead {
			found = true
			break
		}
	}
	return found && strings.Contains(r.BaseUrl, "{{zone}}")
}

// resource functions needed for template that previously existed in terraform.go
// but due to how files are being inherited here it was easier to put in here
// taken wholesale from tpgtools
func (r Resource) Updatable() bool {
	if !r.Immutable {
		return true
	}
	for _, p := range r.AllPropertiesInVersion() {
		if p.UpdateUrl != "" {
			return true
		}
	}
	return false
}

// ====================
// Debugging Methods
// ====================

// Prints a dot notation path to where the field is nested within the parent
// object when called on a property. eg: parent.meta.label.foo
// Redefined on Resource to terminate the calls up the parent chain.
func (r Resource) Lineage() string {
	return r.Name
}

func (r Resource) TerraformName() string {
	if r.LegacyName != "" {
		return r.LegacyName
	}
	return fmt.Sprintf("google_%s_%s", r.ProductMetadata.TerraformName(), google.Underscore(r.Name))
}

func (r Resource) ImportIdFormatsFromResource() []string {
	return ImportIdFormats(r.ImportFormat, r.Identity, r.BaseUrl)
}

// Returns a list of import id formats for a given resource. If an id
// contains provider-default values, this fn will return formats both
// including and omitting the value.
//
// If a resource has an explicit import_format value set, that will be the
// base import url used. Next, the values of `identity` will be used to
// construct a URL. Finally, `{{name}}` will be used by default.
//
// For instance, if the resource base url is:
//
//	projects/{{project}}/global/networks
//
// It returns 3 formats:
// a) self_link: projects/{{project}}/global/networks/{{name}}
// b) short id: {{project}}/{{name}}
// c) short id w/o defaults: {{name}}
func ImportIdFormats(importFormat, identity []string, baseUrl string) []string {
	var idFormats []string
	if len(importFormat) == 0 {
		underscoredBaseUrl := baseUrl

		if len(identity) == 0 {
			idFormats = []string{fmt.Sprintf("%s/{{name}}", underscoredBaseUrl)}
		} else {
			var transformedIdentity []string
			for _, id := range identity {
				transformedIdentity = append(transformedIdentity, fmt.Sprintf("{{%s}}", id))
			}
			identityPath := strings.Join(transformedIdentity, "/")
			idFormats = []string{fmt.Sprintf("%s/%s", underscoredBaseUrl, google.Underscore(identityPath))}
		}
	} else {
		idFormats = importFormat
	}

	// short id: {{project}}/{{zone}}/{{name}}
	fieldMarkers := regexp.MustCompile(`{{[[:word:]]+}}`).FindAllString(idFormats[0], -1)
	shortIdFormat := strings.Join(fieldMarkers, "/")

	// short ids without fields with provider-level defaults:

	// without project
	fieldMarkers = slices.DeleteFunc(fieldMarkers, func(s string) bool { return s == "{{project}}" })
	shortIdDefaultProjectFormat := strings.Join(fieldMarkers, "/")

	// without project or location
	fieldMarkers = slices.DeleteFunc(fieldMarkers, func(s string) bool { return s == "{{region}}" })
	fieldMarkers = slices.DeleteFunc(fieldMarkers, func(s string) bool { return s == "{{zone}}" })
	shortIdDefaultFormat := strings.Join(fieldMarkers, "/")

	// If the id format can include `/` characters we cannot allow short forms such as:
	// `{{project}}/{{%name}}` as there is no way to differentiate between
	// project-name/resource-name and resource-name/with-slash
	if !strings.Contains(idFormats[0], "%") {
		idFormats = append(idFormats, shortIdFormat, shortIdDefaultProjectFormat, shortIdDefaultFormat)
	}

	slices.SortFunc(idFormats, func(a, b string) int {
		i := strings.Count(a, "/")
		j := strings.Count(b, "/")
		if i == j {
			return strings.Count(a, "{{") - strings.Count(b, "{{")
		}
		return i - j
	})
	slices.Reverse(idFormats)

	// Remove duplicates from idFormats
	uniq := make([]string, len(idFormats))
	uniq[0] = idFormats[0]
	i := 1
	j := 1
	for j < len(idFormats) {
		format := idFormats[j]
		if format != uniq[i-1] {
			uniq[i] = format
			i++
		}
		j++
	}

	uniq = google.Reject(slices.Compact(uniq), func(i string) bool {
		return i == ""
	})
	return uniq
}

// IgnoreReadProperties returns a sorted slice of property names (snake_case) that should be ignored when reading.
// This is useful for downstream code that needs to iterate over these properties.
func (r Resource) IgnoreReadProperties(e resource.Examples) []string {
	var props []string
	for _, tp := range r.AllUserProperties() {
		if tp.UrlParamOnly || tp.IsA("ResourceRef") {
			props = append(props, google.Underscore(tp.Name))
		}
	}
	props = append(props, e.IgnoreReadExtra...)
	props = append(props, r.IgnoreReadLabelsFields(r.PropertiesWithExcluded())...)
	props = append(props, ignoreReadFields(r.AllUserProperties())...)

	slices.Sort(props)
	return props
}

// IgnoreReadPropertiesToString returns the ignore read properties as a Go-syntax string slice.
// This is a wrapper around IgnoreReadProperties for backwards compatibility.
func (r Resource) IgnoreReadPropertiesToString(e resource.Examples) string {
	props := r.IgnoreReadProperties(e)
	if len(props) > 0 {
		return fmt.Sprintf("[]string{%s}", strings.Join(quoteStrings(props), ", "))
	}
	return ""
}

// quoteStrings returns a new slice with each string quoted.
func quoteStrings(strs []string) []string {
	quoted := make([]string, len(strs))
	for i, s := range strs {
		quoted[i] = fmt.Sprintf("\"%s\"", s)
	}
	return quoted
}

func ignoreReadFields(props []*Type) []string {
	var fields []string
	for _, tp := range props {
		if tp.IgnoreRead && !tp.UrlParamOnly && !tp.IsA("ResourceRef") {
			fields = append(fields, tp.TerraformLineage())
		} else if tp.IsA("NestedObject") && tp.AllProperties() != nil {
			fields = append(fields, ignoreReadFields(tp.AllProperties())...)
		}
	}
	return fields
}

func (r *Resource) SetCompiler(t string) {
	r.Compiler = fmt.Sprintf("%s-codegen", strings.ToLower(t))
}

// Returns the id format of an object, or self_link_uri if none is explicitly defined
// We prefer the long name of a resource as the id so that users can reference
// resources in a standard way, and most APIs accept short name, long name or self_link
func (r Resource) GetIdFormat() string {
	idFormat := r.IdFormat
	if idFormat == "" {
		idFormat = r.SelfLinkUri()
	}
	return idFormat
}

// Returns true if the Type is in the ID format and false otherwise.
func (r Resource) InPostCreateComputed(prop Type) bool {
	fields := map[string]struct{}{}
	for _, f := range r.ExtractIdentifiers(r.GetIdFormat()) {
		fields[f] = struct{}{}
	}
	for _, f := range r.ExtractIdentifiers(r.SelfLinkUri()) {
		fields[f] = struct{}{}
	}
	_, ok := fields[google.Underscore(prop.Name)]
	return ok
}

// Returns true if at least one of the fields in the ID format is computed
func (r Resource) HasPostCreateComputedFields() bool {
	fields := map[string]struct{}{}
	for _, f := range r.ExtractIdentifiers(r.GetIdFormat()) {
		fields[f] = struct{}{}
	}
	for _, f := range r.ExtractIdentifiers(r.SelfLinkUri()) {
		fields[f] = struct{}{}
	}
	for _, p := range r.GettableProperties() {
		// Skip fields not in the id format
		if _, ok := fields[google.Underscore(p.Name)]; !ok {
			continue
		}
		if (p.Output || p.DefaultFromApi) && !p.IgnoreRead {
			return true
		}
	}
	return false
}

// ====================
// Template Methods
// ====================
// Functions used to create slices of resource properties that could not otherwise be called from within generating templates.
func (r Resource) ReadProperties() []*Type {
	return google.Reject(r.GettableProperties(), func(p *Type) bool {
		return p.IgnoreRead
	})
}

func (r Resource) FlattenedProperties() []*Type {
	return google.Select(r.ReadProperties(), func(p *Type) bool {
		return p.FlattenObject
	})
}

func (r Resource) IsInIdentity(t Type) bool {
	for _, i := range r.GetIdentity() {
		if i.Name == t.Name {
			return true
		}
	}
	return false
}

// ====================
// Iam Methods
// ====================
func (r Resource) IamParentResourceName() string {
	var parentResourceName string

	if r.IamPolicy != nil {
		parentResourceName = r.IamPolicy.ParentResourceAttribute
	}

	if parentResourceName == "" {
		parentResourceName = google.Underscore(r.Name)
	}

	return parentResourceName
}

// For example: "projects/{{project}}/schemas/{{name}}"
func (r Resource) IamResourceUri() string {
	var resourceUri string
	if r.IamPolicy != nil {
		resourceUri = r.IamPolicy.BaseUrl
	}
	if resourceUri == "" {
		resourceUri = r.SelfLinkUri()
	}
	return resourceUri
}

// For example: "projects/%s/schemas/%s"
func (r Resource) IamResourceUriFormat() string {
	return regexp.MustCompile(`\{\{%?(\w+)\}\}`).ReplaceAllString(r.IamResourceUri(), "%s")
}

// For example: the uri "projects/{{project}}/schemas/{{name}}"
// The paramerters are "project", "schema".
func (r Resource) IamResourceParams() []string {
	resourceUri := strings.ReplaceAll(r.IamResourceUri(), "{{name}}", fmt.Sprintf("{{%s}}", r.IamParentResourceName()))

	return r.ExtractIdentifiers(resourceUri)
}

func (r Resource) IsInIamResourceParams(param string) bool {
	return slices.Contains(r.IamResourceParams(), param)
}

// For example: for the uri "projects/{{project}}/schemas/{{name}}",
// the string qualifiers are "u.project, u.schema"
func (r Resource) IamResourceUriStringQualifiers() string {
	var transformed []string
	for _, param := range r.IamResourceParams() {
		transformed = append(transformed, fmt.Sprintf("u.%s", google.Camelize(param, "lower")))
	}
	return strings.Join(transformed[:], ", ")
}

// For example, for the url "projects/{{project}}/schemas/{{schema}}",
// the identifiers are "project", "schema".
func (r Resource) ExtractIdentifiers(url string) []string {
	matches := regexp.MustCompile(`\{\{%?(\w+)\}\}`).FindAllStringSubmatch(url, -1)
	var result []string
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func (r Resource) IamImportFormats() []string {
	var importFormat []string

	if r.IamPolicy != nil {
		importFormat = r.IamPolicy.ImportFormat
	}
	if len(importFormat) == 0 {
		importFormat = r.ImportFormat
	}
	return importFormat
}

// For example, "projects/{{project}}/schemas/{{name}}", "{{project}}/{{name}}", "{{name}}"
func (r Resource) RawImportIdFormatsFromIam() []string {
	return ImportIdFormats(r.IamImportFormats(), r.Identity, r.BaseUrl)
}

// For example, projects/(?P<project>[^/]+)/schemas/(?P<schema>[^/]+)", "(?P<project>[^/]+)/(?P<schema>[^/]+)", "(?P<schema>[^/]+)
func (r Resource) ImportIdRegexesFromIam() string {
	var transformed []string

	importIdFormats := r.RawImportIdFormatsFromIam()
	for _, s := range importIdFormats {
		s = google.Format2Regex(s)
		s = strings.ReplaceAll(s, "<name>", fmt.Sprintf("<%s>", r.IamParentResourceName()))
		transformed = append(transformed, s)
	}

	return strings.Join(slices.Compact(transformed[:]), "\", \"")
}

// For example, "projects/{{project}}/schemas/{{name}}", "{{project}}/{{name}}", "{{name}}"
func (r Resource) ImportIdFormatsFromIam() []string {
	importIdFormats := r.RawImportIdFormatsFromIam()
	var transformed []string
	for _, s := range importIdFormats {
		transformed = append(transformed, strings.ReplaceAll(s, "%", ""))
	}
	return transformed
}

// For example, projects/{{project}}/schemas/{{schema}}
func (r Resource) FirstIamImportIdFormat() string {
	importIdFormats := r.ImportIdFormatsFromIam()
	if len(importIdFormats) == 0 {
		return ""
	}
	first := importIdFormats[0]
	first = strings.ReplaceAll(first, "{{name}}", fmt.Sprintf("{{%s}}", google.Underscore(r.Name)))
	return first
}

func (r Resource) IamTerraformName() string {
	return fmt.Sprintf("%s_iam", r.TerraformName())
}

func (r Resource) IamSelfLinkIdentifiers() []string {
	var selfLink string
	if r.IamPolicy != nil {
		selfLink = r.IamPolicy.SelfLink
	}
	if selfLink == "" {
		selfLink = r.SelfLinkUrl()
	}

	return r.ExtractIdentifiers(selfLink)
}

// Returns the resource properties that are idenfifires in Iam resource when generating the docs.
// The "project" and "organization" properties are excluded, as they are handled seperated in the docs.
func (r Resource) IamResourceProperties() []*Type {
	urlProperties := make([]*Type, 0)
	for _, param := range r.IamResourceParams() {
		if param == "project" || param == "organization" {
			continue
		}

		found := false
		for _, p := range r.AllUserProperties() {
			if param == google.Underscore(p.Name) {
				urlProperties = append(urlProperties, p)
				found = true
				break
			}
		}
		if !found {
			urlProperties = append(urlProperties, &Type{Name: param})
		}
	}

	return urlProperties
}

// Returns the attributes from the selflink url
func (r Resource) IamAttributes() []string {
	var attributes []string
	ids := r.IamSelfLinkIdentifiers()
	for i, p := range ids {
		var attribute string
		if i == len(ids)-1 {
			attribute = r.IamPolicy.ParentResourceAttribute
			if attribute == "" {
				attribute = p
			}
		} else {
			attribute = p
		}
		attributes = append(attributes, attribute)
	}
	return attributes
}

// Since most resources define a "basic" config as their first example,
// we can reuse that config to create a resource to test IAM resources with.
func (r Resource) FirstTestExample() resource.Examples {
	examples := google.Reject(r.Examples, func(e resource.Examples) bool {
		return e.ExcludeTest
	})
	examples = google.Reject(examples, func(e resource.Examples) bool {
		return (r.ProductMetadata.VersionObjOrClosest(r.TargetVersionName).CompareTo(r.ProductMetadata.VersionObjOrClosest(e.MinVersion)) < 0)
	})

	return examples[0]
}

func (r Resource) ExamplePrimaryResourceId() string {
	examples := google.Reject(r.Examples, func(e resource.Examples) bool {
		return e.ExcludeTest
	})
	examples = google.Reject(examples, func(e resource.Examples) bool {
		return (r.ProductMetadata.VersionObjOrClosest(r.TargetVersionName).CompareTo(r.ProductMetadata.VersionObjOrClosest(e.MinVersion)) < 0)
	})

	if len(examples) == 0 {
		examples = google.Reject(r.Examples, func(e resource.Examples) bool {
			return (r.ProductMetadata.VersionObjOrClosest(r.TargetVersionName).CompareTo(r.ProductMetadata.VersionObjOrClosest(e.MinVersion)) < 0)
		})
	}
	return examples[0].PrimaryResourceId
}

func (r Resource) IamParentSourceType() string {
	t := r.IamPolicy.ParentResourceType
	if t == "" {
		t = r.TerraformName()
	}
	return t
}

func (r Resource) IamImportFormat() string {
	var importFormat string
	if len(r.IamPolicy.ImportFormat) > 0 {
		importFormat = r.IamPolicy.ImportFormat[0]
	} else {
		importFormat = r.IamPolicy.SelfLink
		if importFormat == "" {
			importFormat = r.SelfLinkUrl()
		}
	}

	importFormat = regexp.MustCompile(`\{\{%?(\w+)\}\}`).ReplaceAllString(importFormat, "%s")
	return strings.ReplaceAll(importFormat, r.ProductMetadata.BaseUrl, "")
}

func (r Resource) IamImportQualifiersForTest() string {
	var importFormat string
	if len(r.IamPolicy.ImportFormat) > 0 {
		importFormat = r.IamPolicy.ImportFormat[0]
	} else {
		importFormat = r.IamPolicy.SelfLink
		if importFormat == "" {
			importFormat = r.SelfLinkUrl()
		}
	}

	params := r.ExtractIdentifiers(importFormat)
	var importQualifiers []string
	for i, param := range params {
		if param == "project" {
			if i != len(params)-1 {
				// If the last parameter is project then we want to create a new project to use for the test, so don't default from the environment
				if r.IamPolicy.TestProjectName == "" {
					importQualifiers = append(importQualifiers, "envvar.GetTestProjectFromEnv()")
				} else {
					importQualifiers = append(importQualifiers, `context["project_id"]`)
				}
			}
		} else if param == "zone" && r.IamPolicy.SubstituteZoneValue {
			importQualifiers = append(importQualifiers, "envvar.GetTestZoneFromEnv()")
		} else if param == "region" || param == "location" {
			example := r.FirstTestExample()
			if example.RegionOverride == "" {
				importQualifiers = append(importQualifiers, "envvar.GetTestRegionFromEnv()")
			} else {
				importQualifiers = append(importQualifiers, fmt.Sprintf("\"%s\"", example.RegionOverride))
			}
		} else if param == "universe_domain" {
			importQualifiers = append(importQualifiers, "envvar.GetTestUniverseDomainFromEnv()")
		} else {
			break
		}
	}

	if len(importQualifiers) == 0 {
		return ""
	}

	return strings.Join(importQualifiers, ", ")
}

func (r Resource) OrderProperties(props []*Type) []*Type {
	req := google.Select(props, func(p *Type) bool {
		return p.Required
	})
	slices.SortFunc(req, CompareByName)
	rest := google.Reject(props, func(p *Type) bool {
		return p.Output || p.Required
	})
	slices.SortFunc(rest, CompareByName)
	output := google.Select(props, func(p *Type) bool {
		return p.Output
	})
	slices.SortFunc(output, CompareByName)
	returnProps := google.Concat(req, rest)
	return google.Concat(returnProps, output)
}

func CompareByName(a, b *Type) int {
	return strings.Compare(a.Name, b.Name)
}

func (r Resource) GetPropertyUpdateMasksGroupKeys(properties []*Type) []string {
	keys := []string{}
	for _, prop := range properties {
		if prop.FlattenObject {
			k := r.GetPropertyUpdateMasksGroupKeys(prop.Properties)
			keys = append(keys, k...)
		} else {
			keys = append(keys, google.Underscore(prop.Name))
		}
	}
	return keys
}

func (r Resource) GetPropertyUpdateMasksGroups(properties []*Type, maskPrefix string) map[string][]string {
	maskGroups := map[string][]string{}
	for _, prop := range properties {
		if prop.FlattenObject {
			maps.Copy(maskGroups, r.GetPropertyUpdateMasksGroups(prop.Properties, prop.ApiName+"."))
		} else if len(prop.UpdateMaskFields) > 0 {
			maskGroups[google.Underscore(prop.Name)] = prop.UpdateMaskFields
		} else {
			maskGroups[google.Underscore(prop.Name)] = []string{maskPrefix + prop.ApiName}
		}
	}
	return maskGroups
}

// Formats whitespace in the style of the old Ruby generator's descriptions in documentation
func (r Resource) FormatDocDescription(desc string, indent bool) string {
	if desc == "" {
		return ""
	}
	returnString := desc
	if indent {
		returnString = strings.ReplaceAll(returnString, "\n\n", "\n")
		returnString = strings.ReplaceAll(returnString, "\n", "\n  ")

		// fix removing for ruby -> go transition diffs
		returnString = strings.ReplaceAll(returnString, "\n  \n  **Note**: This field is non-authoritative,", "\n\n  **Note**: This field is non-authoritative,")

		return fmt.Sprintf("\n  %s", strings.TrimSuffix(returnString, "\n  "))
	}
	return strings.TrimSuffix(returnString, "\n")
}

func (r Resource) CustomTemplate(templatePath string, appendNewline bool) string {
	output := ExecuteTemplate(&r, templatePath, appendNewline)
	if !appendNewline {
		output = strings.TrimSuffix(output, "\n")
	}
	return output
}

func ExecuteTemplate(e any, templatePath string, appendNewline bool) string {
	templates := []string{
		templatePath,
		"templates/terraform/expand_resource_ref.tmpl",
		"templates/terraform/custom_flatten/bigquery_table_ref.go.tmpl",
		"templates/terraform/flatten_property_method.go.tmpl",
		"templates/terraform/expand_property_method.go.tmpl",
		"templates/terraform/update_mask.go.tmpl",
		"templates/terraform/nested_query.go.tmpl",
		"templates/terraform/unordered_list_customize_diff.go.tmpl",
	}
	templateFileName := filepath.Base(templatePath)

	tmpl, err := template.New(templateFileName).Funcs(google.TemplateFunctions).ParseFiles(templates...)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, templateFileName, e); err != nil {
		glog.Exit(err)
	}

	rs := contents.String()

	if !strings.HasSuffix(rs, "\n") && appendNewline {
		rs = fmt.Sprintf("%s\n", rs)
	}

	return rs
}

// Returns the key of the list of resources in the List API response
// Used to get the list of resources to sweep
func (r Resource) ResourceListKey() string {
	var k string
	if r.NestedQuery != nil && len(r.NestedQuery.Keys) > 0 {
		k = r.NestedQuery.Keys[0]
	}

	if k == "" {
		k = r.CollectionUrlKey
	}

	return k
}

func (r Resource) ListUrlTemplate() string {
	return strings.Replace(r.CollectionUrl(), "zones/{{zone}}", "aggregated", 1)
}

func (r Resource) DeleteUrlTemplate() string {
	return fmt.Sprintf("%s%s", r.ProductMetadata.BaseUrl, r.DeleteUri())
}

func (r Resource) LastNestedQueryKey() string {
	if r.NestedQuery == nil {
		return ""
	}
	len := len(r.NestedQuery.Keys)
	return r.NestedQuery.Keys[len-1]
}

func (r Resource) FirstIdentityProp() *Type {
	idProps := r.GetIdentity()
	if len(idProps) == 0 {
		return nil
	}

	return idProps[0]
}

type UpdateGroup struct {
	UpdateUrl       string
	UpdateVerb      string
	UpdateId        string
	FingerprintName string
}

func (r Resource) propertiesWithCustomUpdate(properties []*Type) []*Type {
	return google.Reject(properties, func(p *Type) bool {
		return p.UpdateUrl == "" || p.UpdateVerb == "" || p.UpdateVerb == "NOOP" ||
			p.IsA("KeyValueTerraformLabels") || p.IsA("KeyValueLabels")
	})
}

func (r Resource) PropertiesByCustomUpdate(properties []*Type) map[UpdateGroup][]*Type {
	customUpdateProps := r.propertiesWithCustomUpdate(properties)
	groupedCustomUpdateProps := map[UpdateGroup][]*Type{}
	for _, prop := range customUpdateProps {
		groupedProperty := UpdateGroup{UpdateUrl: prop.UpdateUrl,
			UpdateVerb:      prop.UpdateVerb,
			UpdateId:        prop.UpdateId,
			FingerprintName: prop.FingerprintName}
		groupedCustomUpdateProps[groupedProperty] = append(groupedCustomUpdateProps[groupedProperty], prop)
	}
	return groupedCustomUpdateProps
}

func (r Resource) PropertiesByCustomUpdateGroups() []UpdateGroup {
	customUpdateProps := r.propertiesWithCustomUpdate(r.RootProperties())
	var updateGroups []UpdateGroup
	for _, prop := range customUpdateProps {
		groupedProperty := UpdateGroup{UpdateUrl: prop.UpdateUrl,
			UpdateVerb:      prop.UpdateVerb,
			UpdateId:        prop.UpdateId,
			FingerprintName: prop.FingerprintName}

		if slices.Contains(updateGroups, groupedProperty) {
			continue
		}
		updateGroups = append(updateGroups, groupedProperty)
	}
	sort.Slice(updateGroups, func(i, j int) bool {
		a := updateGroups[i]
		b := updateGroups[j]
		if a.UpdateVerb != b.UpdateVerb {
			return a.UpdateVerb > b.UpdateVerb
		}
		return a.UpdateId < b.UpdateId
	})
	return updateGroups
}

func (r Resource) FieldSpecificUpdateMethods() bool {
	return (len(r.PropertiesByCustomUpdate(r.RootProperties())) > 0)
}

func (r Resource) CustomUpdatePropertiesByKey(properties []*Type, updateUrl string, updateId string, fingerprintName string, updateVerb string) []*Type {
	groupedProperties := r.PropertiesByCustomUpdate(properties)
	groupedProperty := UpdateGroup{UpdateUrl: updateUrl,
		UpdateVerb:      updateVerb,
		UpdateId:        updateId,
		FingerprintName: fingerprintName}
	return google.Reject(groupedProperties[groupedProperty], func(p *Type) bool {
		return p.UrlParamOnly
	})
}

func (r Resource) PropertyNamesToStrings(properties []*Type) []string {
	var propertyNames []string
	for _, prop := range properties {
		propertyNames = append(propertyNames, google.Underscore(prop.Name))
	}
	return propertyNames
}

func (r Resource) IsExcluded() bool {
	return r.Exclude || r.ExcludeResource
}

func (r Resource) TestExamples() []resource.Examples {
	return google.Reject(google.Reject(r.Examples, func(e resource.Examples) bool {
		return e.ExcludeTest
	}), func(e resource.Examples) bool {
		return e.MinVersion != "" && slices.Index(product.ORDER, r.TargetVersionName) < slices.Index(product.ORDER, e.MinVersion)
	})
}

func (r Resource) VersionedProvider(exampleVersion string) bool {
	var vp string
	if exampleVersion != "" {
		vp = exampleVersion
	} else if r.MinVersion == "" {
		vp = r.ProductMetadata.lowestVersion().Name
	} else {
		vp = r.MinVersion
	}
	return vp != "" && vp != "ga"
}

func (r Resource) StateUpgradersCount() []int {
	var nums []int
	for i := r.StateUpgradeBaseSchemaVersion; i < r.SchemaVersion; i++ {
		nums = append(nums, i)
	}
	return nums
}

func (r Resource) CaiProductBaseUrl() string {
	version := r.ProductMetadata.VersionObjOrClosest(r.TargetVersionName)
	baseUrl := version.CaiBaseUrl
	if baseUrl == "" {
		baseUrl = version.BaseUrl
	}
	return baseUrl
}

// Gets the CAI product legacy base url.
// For example, https://www.googleapis.com/compute/v1/ for compute
func (r Resource) CaiProductLegacyBaseUrl() string {
	version := r.ProductMetadata.VersionObjOrClosest(r.TargetVersionName)
	baseUrl := version.CaiLegacyBaseUrl
	if baseUrl == "" {
		baseUrl = version.CaiBaseUrl
	}
	if baseUrl == "" {
		baseUrl = version.BaseUrl
	}
	return baseUrl
}

// Returns the Cai product backend name from the version base url
// base_url: https://accessapproval.googleapis.com/v1/ -> accessapproval
func (r Resource) CaiProductBackendName(caiProductBaseUrl string) string {
	backendUrl := strings.Split(strings.Split(caiProductBaseUrl, "://")[1], ".googleapis.com")[0]
	return strings.ToLower(backendUrl)
}

// Returns the asset type for this resource.
func (r Resource) CaiAssetType() string {
	baseURL := r.CaiProductBaseUrl()
	productBackendName := r.CaiProductBackendName(baseURL)
	return fmt.Sprintf("%s.googleapis.com/%s", productBackendName, r.CaiResourceName())
}

// DefineAssetTypeForResourceInProduct marks the AssetType constant for this resource as defined.
// It returns true if this is the first time it's been called for this resource,
// and false otherwise, preventing duplicate definitions.
func (r Resource) DefineAssetTypeForResourceInProduct() bool {
	if r.ProductMetadata.ResourcesWithCaiAssetType == nil {
		r.ProductMetadata.ResourcesWithCaiAssetType = make(map[string]struct{}, 1)
	}
	if _, alreadyDefined := r.ProductMetadata.ResourcesWithCaiAssetType[r.CaiResourceType()]; alreadyDefined {
		return false
	}
	r.ProductMetadata.ResourcesWithCaiAssetType[r.CaiResourceType()] = struct{}{}
	return true
}

// Gets the Cai asset name template, which could include version
// For example: //monitoring.googleapis.com/v3/projects/{{project}}/services/{{service_id}}
func (r Resource) rawCaiAssetNameTemplate(productBackendName string) string {
	caiBaseUrl := ""
	caiId := ""
	if r.CaiIdentity != "" {
		caiId = r.CaiIdentity
	} else {
		caiId = r.getCaiId()
	}
	caiIdTemplate := fmt.Sprintf("{{%s}}", caiId)
	if r.CaiBaseUrl != "" {
		if caiId == "" || strings.Contains(r.CaiBaseUrl, caiIdTemplate) {
			caiBaseUrl = r.CaiBaseUrl
		} else {
			caiBaseUrl = fmt.Sprintf("%s/%s", r.CaiBaseUrl, caiIdTemplate)
		}
	}
	if caiBaseUrl == "" {
		caiBaseUrl = r.SelfLink
	}
	if caiBaseUrl == "" {
		if caiId == "" || strings.Contains(r.BaseUrl, caiIdTemplate) {
			caiBaseUrl = r.BaseUrl
		} else {
			caiBaseUrl = fmt.Sprintf("%s/%s", r.BaseUrl, caiIdTemplate)
		}
	}
	return fmt.Sprintf("//%s.googleapis.com/%s", productBackendName, caiBaseUrl)
}

// Guesses the identifier of the resource, as "name" is not always the identifier
// For example, the cai identifier is feed_id in google_cloud_asset_folder_feed
func (r Resource) getCaiId() string {
	for _, p := range r.AllUserProperties() {
		if p.Name == "name" && !p.Output {
			return "name"
		}
	}

	// Get the last identifier extracted from selfLink
	id := r.getCandidateCaiId(r.SelfLink)
	if id != "" {
		return id
	}

	// Get the last identifier extracted from createUrl
	id = r.getCandidateCaiId(r.CreateUrl)
	if id != "" {
		return id
	}

	return ""
}

// Extracts the last identifier from the url, if it is not computed,
// then it is the candidate identifier
func (r Resource) getCandidateCaiId(url string) string {
	identifiers := r.ExtractIdentifiers(url)
	if len(identifiers) > 0 {
		id := identifiers[len(identifiers)-1]
		for _, p := range r.AllUserProperties() {
			if google.Underscore(p.Name) == id && !p.Output {
				return id
			}
		}
	}
	return ""
}

// Gets the Cai asset name template, which doesn't include version
// For example: //monitoring.googleapis.com/projects/{{project}}/services/{{service_id}}
func (r Resource) CaiAssetNameTemplate(productBackendName string) string {
	template := r.rawCaiAssetNameTemplate(productBackendName)
	versionRegex, err := regexp.Compile(`\/(v\d[^\/]*)\/`)
	if err != nil {
		log.Fatalf("Cannot compile the regular expression: %v", err)
	}

	return versionRegex.ReplaceAllString(template, "/")
}

// Gets the Cai API version
func (r Resource) CaiApiVersion(productBackendName, caiProductBaseUrl string) string {
	template := r.rawCaiAssetNameTemplate(productBackendName)

	versionRegex, err := regexp.Compile(`\/(v\d[^\/]*)\/`)
	if err != nil {
		log.Fatalf("Cannot compile the regular expression: %v", err)
	}

	apiVersion := strings.ReplaceAll(versionRegex.FindString(template), "/", "")
	if apiVersion != "" {
		return apiVersion
	}

	splits := strings.Split(caiProductBaseUrl, "/")
	for i := 0; i < len(splits); i++ {
		if splits[len(splits)-1-i] != "" {
			return splits[len(splits)-1-i]
		}
	}
	return ""
}

// For example: the uri "projects/{{project}}/schemas/{{name}}"
// The paramerter is "schema" as "project" is not returned.
func (r Resource) CaiIamResourceParams() []string {
	resourceUri := strings.ReplaceAll(r.IamResourceUri(), "{{name}}", fmt.Sprintf("{{%s}}", r.IamParentResourceName()))

	return google.Reject(r.ExtractIdentifiers(resourceUri), func(param string) bool {
		return param == "project"
	})
}

// Gets the Cai IAM asset name template
// For example: //monitoring.googleapis.com/v3/projects/{{project}}/services/{{service_id}}
func (r Resource) CaiIamAssetNameTemplate(productBackendName string) string {
	iamImportFormat := r.IamImportFormats()
	if len(iamImportFormat) > 0 {
		name := strings.ReplaceAll(iamImportFormat[0], "{{name}}", fmt.Sprintf("{{%s}}", r.IamParentResourceName()))
		name = strings.ReplaceAll(name, "%", "")
		return fmt.Sprintf("//%s.googleapis.com/%s", productBackendName, name)
	}

	caiBaseUrl := r.CaiBaseUrl

	if caiBaseUrl == "" {
		caiBaseUrl = r.SelfLink
	}
	if caiBaseUrl == "" {
		caiBaseUrl = r.BaseUrl
	}
	return fmt.Sprintf("//%s.googleapis.com/%s/{{%s}}", productBackendName, caiBaseUrl, r.IamParentResourceName())
}

func urlContainsOnlyAllowedKeys(templateURL string, allowedKeys []string) bool {
	// Create regex to match anything between {{ and }}
	re := regexp.MustCompile(`{{\s*([^}]+)\s*}}`)

	// Find all matches in the template URL
	matches := re.FindAllStringSubmatch(templateURL, -1)

	// Create a map of allowed keys for O(1) lookup
	allowedKeysMap := make(map[string]bool)
	for _, key := range allowedKeys {
		allowedKeysMap[key] = true
	}

	// Check each found key against the allowed keys
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		// Trim spaces from the key
		key := strings.TrimSpace(match[1])

		// If the key isn't in our allowed list, return false
		if !allowedKeysMap[key] {
			return false
		}
	}

	return true
}

func (r *Resource) ShouldGenerateSingularDataSource() bool {

	if r.Datasource == nil {
		return false
	}

	return r.Datasource.Generate
}

func (r *Resource) ShouldGenerateSingularDataSourceTests() bool {
	if r.Datasource == nil {
		return false
	}
	return !r.Datasource.ExcludeTest
}

func (r Resource) ShouldDatasourceSetLabels() bool {
	for _, p := range r.Properties {
		if p.Name == "labels" && p.Type == "KeyValueLabels" {
			return true
		}
	}
	return false
}

func (r Resource) ShouldDatasourceSetAnnotations() bool {
	for _, p := range r.Properties {
		if p.Name == "annotations" && p.Type == "KeyValueAnnotations" {
			return true
		}
	}
	return false
}

// DatasourceOptionalFields returns a list of fields from the resource's URI
// that should be marked as "Required".
func (r Resource) DatasourceRequiredFields() []string {
	requiredFields := []string{}
	uriParts := strings.Split(r.IdFormat, "/")

	for _, part := range uriParts {
		if strings.HasPrefix(part, "{{") && strings.HasSuffix(part, "}}") {
			field := strings.TrimSuffix(strings.TrimPrefix(part, "{{"), "}}")
			if field != "region" && field != "project" && field != "zone" {
				requiredFields = append(requiredFields, field)
			}
		}
	}
	return requiredFields
}

// DatasourceOptionalFields returns a list of fields from the resource's URI
// that should be marked as "Optional".
func (r Resource) DatasourceOptionalFields() []string {
	optionalFields := []string{}
	uriParts := strings.Split(r.IdFormat, "/")

	for _, part := range uriParts {
		if strings.HasPrefix(part, "{{") && strings.HasSuffix(part, "}}") {
			field := strings.TrimSuffix(strings.TrimPrefix(part, "{{"), "}}")
			if field == "region" || field == "project" || field == "zone" {
				optionalFields = append(optionalFields, field)
			}
		}
	}
	return optionalFields
}

func (r Resource) ShouldGenerateSweepers() bool {
	if !r.ExcludeSweeper && !utils.IsEmpty(r.Sweeper) {
		return true
	}

	allowedKeys := []string{"project", "region", "location", "zone", "billing_account"}
	if !urlContainsOnlyAllowedKeys(r.ListUrlTemplate(), allowedKeys) {
		return false
	}
	if r.ExcludeSweeper || r.CustomCode.CustomDelete != "" || r.CustomCode.PreDelete != "" || r.CustomCode.PostDelete != "" || r.ExcludeDelete {
		return false
	}
	return true
}

func (r Resource) GithubURL() string {
	return GITHUB_BASE_URL + r.SourceYamlFile
}

func (r Resource) CodeHeader(templatePath string) string {
	templateUrl := GITHUB_BASE_URL + templatePath

	return fmt.Sprintf(`// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: %s
//     Template:      %s
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------`, r.GithubURL(), templateUrl)
}

func (r Resource) MarkdownHeader(templatePath string) string {
	return strings.Replace(r.CodeHeader(templatePath), "//", "#", -1)
}

// TGC Methods
// ====================
// Lists fields that test.BidirectionalConversion should ignore
func (r Resource) TGCTestIgnorePropertiesToStrings(e resource.Examples) []string {
	props := []string{
		"depends_on",
		"count",
		"for_each",
		"provider",
		"lifecycle",
	}
	for _, tp := range r.VirtualFields {
		props = append(props, google.Underscore(tp.Name))
	}
	for _, tp := range r.AllNestedProperties(r.RootProperties()) {
		if tp.UrlParamOnly {
			props = append(props, google.Underscore(tp.Name))
		} else if tp.IsMissingInCai {
			props = append(props, tp.MetadataLineage())
		}
	}
	props = append(props, e.TGCTestIgnoreExtra...)

	slices.Sort(props)
	return props
}

// Filters out computed properties during cai2hcl
func (r Resource) ReadPropertiesForTgc() []*Type {
	return google.Reject(r.AllUserProperties(), func(v *Type) bool {
		return v.Output || v.UrlParamOnly || v.TGCIgnoreRead
	})
}

// OutputFieldSetStr returns a Go-syntax string representation of a set
// containing all the output properties for a resource.
// The property names are converted to snake_case.
// This is useful for generating code that requires a map literal of field names.
func (r Resource) OutputFieldSetStr() string {
	fieldNames := make(map[string]struct{})
	for _, tp := range r.AllUserProperties() {
		if tp.Output {
			fieldNames[google.Underscore(tp.Name)] = struct{}{}
		}
	}
	return fmt.Sprintf("%#v", fieldNames)
}

// For example, the CAI resource type with product of "google_compute_autoscaler" is "ComputeAutoscalerAssetType".
// The CAI resource type with product of "google_compute_region_autoscaler" is also "ComputeAutoscalerAssetType".
func (r Resource) CaiResourceType() string {
	return fmt.Sprintf("%s%s", r.ProductMetadata.Name, r.CaiResourceName())
}

// The API resource type of the resource. Normally, it is the resource name.
// Rarely, it is the API "resource type kind" or CAI "resource kind"
// For example, the CAI resource type of "google_compute_autoscaler" is "Autoscaler".
// The CAI resource type of "google_compute_region_autoscaler" is also "Autoscaler".
func (r Resource) CaiResourceName() string {
	if r.CaiResourceKind != "" {
		return r.CaiResourceKind
	}
	if r.ApiResourceTypeKind != "" {
		return r.ApiResourceTypeKind
	}
	return r.Name
}

func (r Resource) IsTgcCompiler() bool {
	return r.Compiler == "terraformgoogleconversionnext-codegen"
}
