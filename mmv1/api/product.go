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
	"reflect"
	"strings"
	"unicode"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"golang.org/x/exp/slices"
)

// Represents a product to be managed
type Product struct {
	// The name of the product's API capitalised in the appropriate places.
	// This isn't just the API name because it doesn't meaningfully separate
	// words in the api name - "accesscontextmanager" vs "AccessContextManager"
	// Example inputs: "Compute", "AccessContextManager"
	Name string

	// original value of :name before the provider override happens
	// same as :name if not overridden in provider
	ApiName string `yaml:"api_name"`

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

func (p *Product) UnmarshalYAML(unmarshal func(any) error) error {
	type productAlias Product
	aliasObj := (*productAlias)(p)

	if err := unmarshal(aliasObj); err != nil {
		return err
	}

	p.SetApiName()
	p.SetDisplayName()

	return nil
}

func (p *Product) Validate() {
	if len(p.Name) == 0 {
		log.Fatalf("Missing `name` for product")
	}

	// product names must start with a capital
	for i, ch := range p.Name {
		if !unicode.IsUpper(ch) {
			log.Fatalf("product name `%s` must start with a capital letter.", p.Name)
		}
		if i == 0 {
			break
		}
	}

	if len(p.Scopes) == 0 {
		log.Fatalf("Missing `scopes` for product %s", p.Name)
	}

	if p.Versions == nil {
		log.Fatalf("Missing `versions` for product %s", p.Name)
	}

	for _, v := range p.Versions {
		v.Validate(p.Name)
	}

	if p.Async != nil {
		p.Async.Validate()
	}
}

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

// Prints a dot notation path to where the field is nested within the parent
// object when called on a property. eg: parent.meta.label.foo
// Redefined on Product to terminate the calls up the parent chain.
func (p Product) Lineage() string {
	return p.Name
}

func Merge(self, otherObj reflect.Value) {

	selfObj := reflect.Indirect(self)
	for i := 0; i < selfObj.NumField(); i++ {

		// skip if the override is the "empty" value
		emptyOverrideValue := reflect.DeepEqual(reflect.Zero(otherObj.Field(i).Type()).Interface(), otherObj.Field(i).Interface())

		if emptyOverrideValue && selfObj.Type().Field(i).Name != "Required" {
			continue
		}

		if selfObj.Field(i).Kind() == reflect.Slice {
			DeepMerge(selfObj.Field(i), otherObj.Field(i))
		} else {
			selfObj.Field(i).Set(otherObj.Field(i))
		}
	}
}

func DeepMerge(arr1, arr2 reflect.Value) {
	if arr1.Len() == 0 {
		arr1.Set(arr2)
		return
	}
	if arr2.Len() == 0 {
		return
	}

	// Scopes is an array of standard strings. In which case return the
	// version in the overrides. This allows scopes to be removed rather
	// than allowing for a merge of the two arrays
	if arr1.Index(0).Kind() == reflect.String {
		arr1.Set(arr2)
		return
	}

	// Merge any elements that exist in both
	for i := 0; i < arr1.Len(); i++ {
		currentVal := arr1.Index(i)
		pointer := currentVal.Kind() == reflect.Ptr
		if pointer {
			currentVal = currentVal.Elem()
		}
		var otherVal reflect.Value
		for j := 0; j < arr2.Len(); j++ {
			currentName := currentVal.FieldByName("Name").Interface()
			tempOtherVal := arr2.Index(j)
			if pointer {
				tempOtherVal = tempOtherVal.Elem()
			}
			otherName := tempOtherVal.FieldByName("Name").Interface()

			if otherName == currentName {
				otherVal = tempOtherVal
				break
			}
		}
		if otherVal.IsValid() {
			Merge(currentVal, otherVal)
		}
	}

	// Add any elements of arr2 that don't exist in arr1
	for i := 0; i < arr2.Len(); i++ {
		otherVal := arr2.Index(i)
		pointer := otherVal.Kind() == reflect.Ptr
		if pointer {
			otherVal = otherVal.Elem()
		}

		found := false
		for j := 0; j < arr1.Len(); j++ {
			currentVal := arr1.Index(j)
			if pointer {
				currentVal = currentVal.Elem()
			}
			currentName := currentVal.FieldByName("Name").Interface()
			otherName := otherVal.FieldByName("Name").Interface()

			if otherName == currentName {
				found = true
				break
			}
		}
		if !found {
			arr1.Set(reflect.Append(arr1, arr2.Index(i)))
		}
	}
}
