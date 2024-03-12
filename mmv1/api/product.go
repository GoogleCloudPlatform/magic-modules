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
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"golang.org/x/exp/slices"
)

// require 'api/object'
// require 'api/product/version'
// require 'google/logger'
// require 'compile/core'
// require 'json'

// Represents a product to be managed
type Product struct {
	NamedObject `yaml:",inline"`

	// include Compile::Core

	// Inherited:
	// The name of the product's API capitalised in the appropriate places.
	// This isn't just the API name because it doesn't meaningfully separate
	// words in the api name - "accesscontextmanager" vs "AccessContextManager"
	// Example inputs: "Compute", "AccessContextManager"
	// Name string

	// Display Name: The full name of the GCP product; eg "Cloud Bigtable"
	DisplayName string `yaml:"display_name"`

	Objects []*Resource

	// The list of permission scopes available for the service
	// For example: `https://www.googleapis.com/auth/compute`

	Scopes []string

	// The API versions of this product

	Versions []product.Version

	// The base URL for the service API endpoint
	// For example: `https://www.googleapis.com/compute/v1/`

	BaseUrl string `yaml:"base_url"`

	// A function reference designed for the rare case where you
	// need to use retries in operation calls. Used for the service api
	// as it enables itself (self referential) and can result in occasional
	// failures on operation_get. see github.com/hashicorp/terraform-provider-google/issues/9489

	OperationRetry string `yaml:"operation_retry"`

	Async OpAsync

	LegacyName string `yaml:"legacy_name"`

	ClientName string `yaml:"client_name"`
}

func (p *Product) Validate() {
	// TODO Q1 Rewrite super
	//     super
	for _, o := range p.Objects {
		o.ProductMetadata = p
	}
}

// def validate
//     super
//     set_variables @objects, :__product

//     // name comes from Named, and product names must start with a capital
//     caps = ('A'..'Z').to_a
//     unless caps.include? @name[0]
//       raise "product name `//{@name}` must start with a capital letter."
//     end

//     check :display_name, type: String
//     check :objects, type: Array, item_type: Api::Resource
//     check :scopes, type: Array, item_type: String, required: true
//     check :operation_retry, type: String

//     check :async, type: Api::Async
//     check :legacy_name, type: String
//     check :client_name, type: String

//     check :versions, type: Array, item_type: Api::Product::Version, required: true
//   end

//   // ====================
//   // Custom Getters
//   // ====================

//   // The name of the product's API; "compute", "accesscontextmanager"
//   def api_name
//     name.downcase
//   end

// The product full name is the "display name" in string form intended for
// users to read in documentation; "Google Compute Engine", "Cloud Bigtable"
func (p Product) GetDisplayName() string {
	if p.DisplayName == "" {
		return google.SpaceSeparated(p.Name)
	}

	return p.DisplayName
}

//   // Most general version that exists for the product
//   // If GA is present, use that, else beta, else alpha
//   def lowest_version
//     Version::ORDER.each do |ordered_version_name|
//       @versions.each do |product_version|
//         return product_version if ordered_version_name == product_version.name
//       end
//     end
//     raise "Unable to find lowest version for product //{display_name}"
//   end

//   def version_obj(name)
//     @versions.each do |v|
//       return v if v.name == name
//     end

//     raise "API version '//{name}' does not exist for product '//{@name}'"
//   end

//   // Get the version of the object specified by the version given if present
//   // Or else fall back to the closest version in the chain defined by Version::ORDER
//   def version_obj_or_closest(name)
//     return version_obj(name) if exists_at_version(name)

//     // versions should fall back to the closest version to them that exists
//     name ||= Version::ORDER[0]
//     lower_versions = Version::ORDER[0..Version::ORDER.index(name)]

//     lower_versions.reverse_each do |version|
//       return version_obj(version) if exists_at_version(version)
//     end

//     raise "Could not find object for version //{name} and product //{display_name}"
//   end

//   def exists_at_version_or_lower(name)
//     // Versions aren't normally going to be empty since products need a
//     // base_url. This nil check exists for atypical products, like _bundle.
//     return true if @versions.nil?

//     name ||= Version::ORDER[0]
//     return false unless Version::ORDER.include?(name)

//     (0..Version::ORDER.index(name)).each do |i|
//       return true if exists_at_version(Version::ORDER[i])
//     end
//     false
//   end

func (p *Product) ExistsAtVersionOrLower(name string) bool {
	if !slices.Contains(product.ORDER, name) {
		return false
	}

	for i := 0; i <= slices.Index(product.ORDER, name); i++ {
		if p.ExistsAtVersion(product.ORDER[i]) {
			return true
		}
	}

	return false
}

func (p *Product) ExistsAtVersion(name string) bool {
	for _, v := range p.Versions {
		if v.Name == name {
			return true
		}
	}
	return false
}

//   def exists_at_version(name)
//     // Versions aren't normally going to be empty since products need a
//     // base_url. This nil check exists for atypical products, like _bundle.
//     return true if @versions.nil?

//     @versions.any? { |v| v.name == name }
//   end

//   // Not a conventional setter, so ignore rubocop's warning
//   // rubocop:disable Naming/AccessorMethodName
//   def set_properties_based_on_version(version)
//     @base_url = version.base_url
//   end
//   // rubocop:enable Naming/AccessorMethodName

//   // ====================
//   // Debugging Methods
//   // ====================

//   def to_s
//     // relies on the custom to_json definitions
//     JSON.pretty_generate(self)
//   end

//   // Prints a dot notation path to where the field is nested within the parent
//   // object when called on a property. eg: parent.meta.label.foo
//   // Redefined on Product to terminate the calls up the parent chain.
//   def lineage
//     name
//   end

//   def to_json(opts = nil)
//     json_out = {}

//     instance_variables.each do |v|
//       if v == :@objects
//         json_out['@resources'] = objects.to_h { |o| [o.name, o] }
//       elsif instance_variable_get(v) == false || instance_variable_get(v).nil?
//         // ignore false or missing because omitting them cleans up result
//         // and both are the effective defaults of their types
//       else
//         json_out[v] = instance_variable_get(v)
//       end
//     end

//     JSON.generate(json_out, opts)
//   end
