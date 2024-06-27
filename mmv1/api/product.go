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
	"log"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"golang.org/x/exp/slices"
)

// Represents a product to be managed
type Product struct {
	NamedObject `yaml:",inline"`

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

	Versions []*product.Version

	// The base URL for the service API endpoint
	// For example: `https://www.googleapis.com/compute/v1/`

	BaseUrl string `yaml:"base_url"`

	// A function reference designed for the rare case where you
	// need to use retries in operation calls. Used for the service api
	// as it enables itself (self referential) and can result in occasional
	// failures on operation_get. see github.com/hashicorp/terraform-provider-google/issues/9489

	OperationRetry string `yaml:"operation_retry"`

	Async *Async

	LegacyName string `yaml:"legacy_name"`

	ClientName string `yaml:"client_name"`
}

func (p *Product) UnmarshalYAML(n *yaml.Node) error {
	type productAlias Product
	aliasObj := (*productAlias)(p)

	err := n.Decode(&aliasObj)
	if err != nil {
		return err
	}

	p.SetApiName()
	p.SetDisplayName()

	return nil
}

func (p *Product) Validate() {
	// TODO Q2 Rewrite super
	//     super
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

// ====================
// Custom Setters
// ====================

func (p *Product) SetApiName() {
	// The name of the product's API; "compute", "accesscontextmanager"
	p.ApiName = strings.ToLower(p.Name)
}

// The product full name is the "display name" in string form intended for
// users to read in documentation; "Google Compute Engine", "Cloud Bigtable"
func (p *Product) SetDisplayName() {
	if p.DisplayName == "" {
		p.DisplayName = google.SpaceSeparated(p.Name)
	}
}

// ====================
// Version-related methods
// ====================

// Most general version that exists for the product
// If GA is present, use that, else beta, else alpha
func (p Product) lowestVersion() *product.Version {
	for _, orderedVersionName := range product.ORDER {
		for _, productVersion := range p.Versions {
			if orderedVersionName == productVersion.Name {
				return productVersion
			}
		}
	}

	log.Fatalf("Unable to find lowest version for product %s", p.DisplayName)
	return nil
}

func (p Product) versionObj(name string) *product.Version {
	for _, v := range p.Versions {
		if v.Name == name {
			return v
		}
	}

	log.Fatalf("API version '%s' does not exist for product '%s'", name, p.Name)
	return nil
}

// Get the version of the object specified by the version given if present
// Or else fall back to the closest version in the chain defined by product.ORDER
func (p Product) VersionObjOrClosest(name string) *product.Version {
	if p.ExistsAtVersion(name) {
		return p.versionObj(name)
	}

	// versions should fall back to the closest version to them that exists
	if name == "" {
		name = product.ORDER[0]
	}

	lowerVersions := make([]string, 0)
	for _, v := range product.ORDER {
		lowerVersions = append(lowerVersions, v)
		if v == name {
			break
		}
	}

	for i := len(lowerVersions) - 1; i >= 0; i-- {
		if p.ExistsAtVersion(lowerVersions[i]) {
			return p.versionObj(lowerVersions[i])
		}
	}

	log.Fatalf("Could not find object for version %s and product %s", name, p.DisplayName)
	return nil
}

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

func (p *Product) SetPropertiesBasedOnVersion(version *product.Version) {
	p.BaseUrl = version.BaseUrl
}

func (p *Product) TerraformName() string {
	if p.LegacyName != "" {
		return google.Underscore(p.LegacyName)
	}
	return google.Underscore(p.Name)
}

// ====================
// Debugging Methods
// ====================

//   def to_s
//     // relies on the custom to_json definitions
//     JSON.pretty_generate(self)
//   end

// Prints a dot notation path to where the field is nested within the parent
// object when called on a property. eg: parent.meta.label.foo
// Redefined on Product to terminate the calls up the parent chain.
func (p Product) Lineage() string {
	return p.Name
}

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
