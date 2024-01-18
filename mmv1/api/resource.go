// Copyright 2017 Google Inc.
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

package main

type Resource struct {

	// [Required] A description of the resource that's surfaced in provider
	// documentation.
	description string
	// [Required] (Api::Resource::ReferenceLinks) Reference links provided in
	// downstream documentation.
	references ReferenceLinks
	// [Required] The GCP "relative URI" of a resource, relative to the product
	// base URL. It can often be inferred from the `create` path.
	base_url string

	// ====================
	// Common Configuration
	// ====================
	//
	// [Optional] The minimum API version this resource is in. Defaults to ga.
	min_version string
	// [Optional] If set to true, don't generate the resource.
	exclude bool
	// [Optional] If set to true, the resource is not able to be updated.
	immutable bool
	// [Optional] If set to true, this resource uses an update mask to perform
	// updates. This is typical of newer GCP APIs.
	update_mask bool
	// [Optional] If set to true, the object has a `self_link` field. This is
	// typical of older GCP APIs.
	has_self_link bool

	// [Optional] The validator "relative URI" of a resource, relative to the product
	// base URL. Specific to defining the resource as a CAI asset.
	cai_base_url string

	// ====================
	// URL / HTTP Configuration
	// ====================
	//
	// [Optional] The "identity" URL of the resource. Defaults to:
	// * base_url when the create_verb is :POST
	// * self_link when the create_verb is :PUT  or :PATCH
	self_link string
	// [Optional] The URL used to creating the resource. Defaults to:
	// * collection url when the create_verb is :POST
	// * self_link when the create_verb is :PUT or :PATCH
	create_url string
	// [Optional] The URL used to delete the resource. Defaults to the self
	// link.
	delete_url string
	// [Optional] The URL used to update the resource. Defaults to the self
	// link.
	update_url string
	// [Optional] The HTTP verb used during create. Defaults to :POST.
	create_verb string
	// [Optional] The HTTP verb used during read. Defaults to :GET.
	read_verb string
	// [Optional] The HTTP verb used during update. Defaults to :PUT.
	update_verb string
	// [Optional] The HTTP verb used during delete. Defaults to :DELETE.
	delete_verb string
	// [Optional] Additional Query Parameters to append to GET. Defaults to ""
	read_query_params string
	// ====================
	// Collection / Identity URL Configuration
	// ====================
	//
	// [Optional] This is the name of the list of items
	// within the collection (list) json. Will default to the
	// camelcase plural name of the resource.
	collection_url_key string
	// [Optional] An ordered list of names of parameters that uniquely identify
	// the resource.
	// Generally, it's safe to leave empty, in which case it defaults to `name`.
	// Other values are normally useful in cases where an object has a parent
	// and is identified by some non-name value, such as an ip+port pair.
	// If you're writing a fine-grained resource (eg with nested_query) a value
	// must be set.
	identity []string

	// [Optional] (Api::Resource::NestedQuery) This is useful in case you need
	// to change the query made for GET requests only. In particular, this is
	// often used to extract an object from a parent object or a collection.
	// Note that if both nested_query and custom_code.decoder are provided,
	// the decoder will be included within the code handling the nested query.
	nested_query

	// ====================
	// IAM Configuration
	// ====================
	//
	// [Optional] (Api::Resource::IamPolicy) Configuration of a resource's
	// resource-specific IAM Policy.
	iam_policy IamPolicy
	// [Optional] If set to true, don't generate the resource itself; only
	// generate the IAM policy.
	exclude_resource bool

	// [Optional] GCP kind, e.g. `compute//disk`
	kind string
	// [Optional] If set to true, indicates that a resource is not configurable
	// such as GCP regions.
	readonly bool

	// ====================
	// Terraform Overrides
	// ====================

	// [Optional] If non-empty, overrides the full filename prefix
	// i.e. google/resource_product_{{resource_filename_override}}.go
	// i.e. google/resource_product_{{resource_filename_override}}_test.go
	filename_override string

	// If non-empty, overrides the full given resource name.
	// i.e. 'google_project' for resourcemanager.Project
	// Use Provider::Terraform::Config.legacy_name to override just
	// product name.
	// Note: This should not be used for vanity names for new products.
	// This was added to handle preexisting handwritten resources that
	// don't match the natural generated name exactly, and to support
	// services with a mix of handwritten and generated resources.
	legacy_name string

	// The Terraform resource id format used when calling //setId(...).
	// For instance, `{{name}}` means the id will be the resource name.
	id_format string
	// Override attribute used to handwrite the formats for generating regex strings
	// that match templated values to a self_link when importing, only necessary when
	// a resource is not adequately covered by the standard provider generated options.
	// Leading a token with `%`
	// i.e. {{%parent}}/resource/{{resource}}
	// will allow that token to hold multiple /'s.
	import_format []string
	custom_code
	docs

	// This block inserts entries into the customdiff.All() block in the
	// resource schema -- the code for these custom diff functions must
	// be included in the resource constants or come from tpgresource
	custom_diff []string

	// Lock name for a mutex to prevent concurrent API calls for a given
	// resource.
	mutex string

	// Examples in documentation. Backed by generated tests, and have
	// corresponding OiCS walkthroughs.
	examples Examples

	// Virtual fields on the Terraform resource. Usage and differences from url_param_only
	// are documented in provider/terraform/virtual_fields.rb
	virtual_fields

	// TODO(alexstephen): Deprecate once all resources using autogen async.
	// If true, generates product operation handling logic.
	autogen_async bool

	// If true, resource is not importable
	exclude_import bool

	// If true, exclude resource from Terraform Validator
	// (i.e. terraform-provider-conversion)
	exclude_tgc bool

	// If true, skip sweeper generation for this resource
	skip_sweeper bool

	timeouts Timeouts

	// An array of function names that determine whether an error is retryable.
	error_retry_predicates []string

	// An array of function names that determine whether an error is not retryable.
	error_abort_predicates []string

	// Optional attributes for declaring a resource's current version and generating
	// state_upgrader code to the output .go file from files stored at
	// mmv1/templates/terraform/state_migrations/
	// used for maintaining state stability with resources first provisioned on older api versions.
	schema_version int
	// From this schema version on, state_upgrader code is generated for the resource.
	// When unset, state_upgrade_base_schema_version defauts to 0.
	// Normally, it is not needed to be set.
	state_upgrade_base_schema_version int
	state_upgraders                   bool
	// This block inserts the named function and its attribute into the
	// resource schema -- the code for the migrate_state function must
	// be included in the resource constants or come from tpgresource
	// included for backwards compatibility as an older state migration method
	// and should not be used for new resources.
	migrate_state string

	// Set to true for resources that are unable to be deleted, such as KMS keyrings or project
	// level resources such as firebase project
	skip_delete bool

	// Set to true for resources that are unable to be read from the API, such as
	// public ca external account keys
	skip_read bool

	// Set to true for resources that wish to disable automatic generation of default provider
	// value customdiff functions
	skip_default_cdiff bool

	// This enables resources that get their project via a reference to a different resource
	// instead of a project field to use User Project Overrides
	supports_indirect_user_project_override bool

	// If true, the resource's project field can be specified as either the short form project
	// id or the long form projects/project-id. The extra projects/ string will be removed from
	// urls and ids. This should only be used for resources that previously supported long form
	// project ids for backwards compatibility.
	legacy_long_form_project bool

	// Function to transform a read error so that handleNotFound recognises
	// it as a 404. This should be added as a handwritten fn that takes in
	// an error and returns one.
	read_error_transform string

	// If true, resources that failed creation will be marked as tainted. As a consequence
	// these resources will be deleted and recreated on the next apply call. This pattern
	// is preferred over deleting the resource directly in post_create_failure hooks.
	taint_resource_on_failed_create bool

	// Add a deprecation message for a resource that's been deprecated in the API.
	deprecation_message string
}
