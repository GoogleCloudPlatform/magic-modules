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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
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

	// This is the name of the package path relative to mmv1 root repo
	PackagePath string

	// original value of :name before the provider override happens
	// same as :name if not overridden in provider
	ApiName string `yaml:"api_name,omitempty"`

	// Display Name: The full name of the GCP product; eg "Cloud Bigtable"
	DisplayName string `yaml:"display_name,omitempty"`

	Objects []*Resource `yaml:"objects,omitempty"`

	// The list of permission scopes available for the service
	// For example: `https://www.googleapis.com/auth/compute`
	Scopes []string

	// The API versions of this product
	Versions []*product.Version

	// The base URL for the service API endpoint
	// For example: `https://www.googleapis.com/compute/v1/`
	BaseUrl string `yaml:"base_url,omitempty"`

	// The validator "relative URI" of a resource, relative to the product
	// base URL. Specific to defining the resource as a CAI asset.
	CaiBaseUrl string

	// CaiResourceType of resources that already have an AssetType constant defined in the product.
	ResourcesWithCaiAssetType map[string]struct{}

	// A function reference designed for the rare case where you
	// need to use retries in operation calls. Used for the service api
	// as it enables itself (self referential) and can result in occasional
	// failures on operation_get. see github.com/hashicorp/terraform-provider-google/issues/9489
	OperationRetry string `yaml:"operation_retry,omitempty"`

	Async *Async `yaml:"async,omitempty"`

	LegacyName string `yaml:"legacy_name,omitempty"`

	ClientName string `yaml:"client_name,omitempty"`

	// The compiler to generate the downstream files, for example "terraformgoogleconversion-codegen".
	Compiler string `yaml:"-"`
}

// Load compiles a product with all its resources from the given path and optional overrides
// This loads the product configuration and all its resources into memory without generating any files
func (p *Product) Load(productName string, version string, overrideDirectory string) error {
	productYamlPath := filepath.Join(productName, "product.yaml")

	var productOverridePath string
	if overrideDirectory != "" {
		productOverridePath = filepath.Join(overrideDirectory, productName, "product.yaml")
	}

	_, baseProductErr := os.Stat(productYamlPath)
	baseProductExists := !errors.Is(baseProductErr, os.ErrNotExist)

	_, overrideProductErr := os.Stat(productOverridePath)
	overrideProductExists := !errors.Is(overrideProductErr, os.ErrNotExist)

	if !(baseProductExists || overrideProductExists) {
		return fmt.Errorf("%s does not contain a product.yaml file", productName)
	}

	// Compile the product configuration
	if overrideProductExists {
		if baseProductExists {
			Compile(productYamlPath, p, overrideDirectory)
			overrideApiProduct := &Product{}
			Compile(productOverridePath, overrideApiProduct, overrideDirectory)
			Merge(reflect.ValueOf(p).Elem(), reflect.ValueOf(*overrideApiProduct), version)
		} else {
			Compile(productOverridePath, p, overrideDirectory)
		}
	} else {
		Compile(productYamlPath, p, overrideDirectory)
	}

	// Check if product exists at the requested version
	if !p.ExistsAtVersionOrLower(version) {
		return &ErrProductVersionNotFound{ProductName: productName, Version: version}
	}

	// Compile all resources
	resources, err := p.loadResources(productName, version, overrideDirectory)
	if err != nil {
		return err
	}

	p.Objects = resources
	p.PackagePath = productName
	p.Validate()

	return nil
}

// loadResources loads all resources for a product
func (p *Product) loadResources(productName string, version string, overrideDirectory string) ([]*Resource, error) {
	var resources []*Resource = make([]*Resource, 0)

	// Get base resource files
	resourceFiles, err := filepath.Glob(fmt.Sprintf("%s/*", productName))
	if err != nil {
		return nil, fmt.Errorf("cannot get resource files: %v", err)
	}

	// Compile base resources (skip those that will be merged with overrides)
	for _, resourceYamlPath := range resourceFiles {
		if filepath.Base(resourceYamlPath) == "product.yaml" || filepath.Ext(resourceYamlPath) != ".yaml" {
			continue
		}

		// Skip if resource will be merged in the override loop
		if overrideDirectory != "" {
			resourceOverridePath := filepath.Join(overrideDirectory, resourceYamlPath)
			_, overrideResourceErr := os.Stat(resourceOverridePath)
			if !errors.Is(overrideResourceErr, os.ErrNotExist) {
				continue
			}
		}

		resource := p.loadResource(resourceYamlPath, "", version, overrideDirectory)
		resources = append(resources, resource)
	}

	// Compile override resources
	if overrideDirectory != "" {
		resources, err = p.reconcileOverrideResources(productName, version, overrideDirectory, resources)
		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}

// reconcileOverrideResources handles resolution of override resources
func (p *Product) reconcileOverrideResources(productName string, version string, overrideDirectory string, resources []*Resource) ([]*Resource, error) {
	productOverridePath := filepath.Join(overrideDirectory, productName, "product.yaml")
	productOverrideDir := filepath.Dir(productOverridePath)

	overrideFiles, err := filepath.Glob(fmt.Sprintf("%s/*", productOverrideDir))
	if err != nil {
		return nil, fmt.Errorf("cannot get override files: %v", err)
	}

	for _, overrideYamlPath := range overrideFiles {
		if filepath.Base(overrideYamlPath) == "product.yaml" || filepath.Ext(overrideYamlPath) != ".yaml" {
			continue
		}

		baseResourcePath := filepath.Join(productName, filepath.Base(overrideYamlPath))
		resource := p.loadResource(baseResourcePath, overrideYamlPath, version, overrideDirectory)
		resources = append(resources, resource)
	}

	// Sort resources by name for consistent output
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})

	return resources, nil
}

// loadResource loads a single resource with optional override
func (p *Product) loadResource(baseResourcePath string, overrideResourcePath string, version string, overrideDirectory string) *Resource {
	resource := &Resource{}

	// Check if base resource exists
	_, baseResourceErr := os.Stat(baseResourcePath)
	baseResourceExists := !errors.Is(baseResourceErr, os.ErrNotExist)

	if overrideResourcePath != "" {
		if baseResourceExists {
			// Merge base and override
			Compile(baseResourcePath, resource, overrideDirectory)
			overrideResource := &Resource{}
			Compile(overrideResourcePath, overrideResource, overrideDirectory)
			Merge(reflect.ValueOf(resource).Elem(), reflect.ValueOf(*overrideResource), version)
			resource.SourceYamlFile = baseResourcePath
		} else {
			// Override only
			Compile(overrideResourcePath, resource, overrideDirectory)
		}
	} else {
		// Base only
		Compile(baseResourcePath, resource, overrideDirectory)
		resource.SourceYamlFile = baseResourcePath
	}

	// Set resource defaults and validate
	resource.TargetVersionName = version
	// SetDefault before AddExtraFields to ensure relevant metadata is available on existing fields
	resource.SetDefault(p)
	resource.Properties = resource.AddExtraFields(resource.PropertiesWithExcluded(), nil)
	// SetDefault after AddExtraFields to ensure relevant metadata is available for the newly generated fields
	resource.SetDefault(p)
	resource.Validate()

	return resource
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

func (p *Product) SetCompiler(t string) {
	p.Compiler = fmt.Sprintf("%s-codegen", strings.ToLower(t))
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
	p.CaiBaseUrl = version.CaiBaseUrl
}

func (p *Product) TerraformName() string {
	if p.LegacyName != "" {
		return google.Underscore(p.LegacyName)
	}
	return google.Underscore(p.Name)
}

func (p *Product) ServiceBaseUrl() string {
	if p.CaiBaseUrl != "" {
		return p.CaiBaseUrl
	}
	return p.BaseUrl
}

func (p *Product) ServiceName() string {
	parts := strings.Split(p.ServiceBaseUrl(), "/")
	// remove location prefix if present
	trimmed, _ := strings.CutPrefix(parts[2], "{{location}}-")
	return trimmed
}

var versionRegexp = regexp.MustCompile(`^v[0-9]+|beta`)

func (p *Product) ServiceVersion() string {
	parts := strings.Split(p.ServiceBaseUrl(), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		// stop when we get to the domain name
		if strings.Contains(part, ".com") {
			break
		}
		v := versionRegexp.FindString(part)
		if v != "" {
			return part
		}
	}
	return ""
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

func Merge(self, otherObj reflect.Value, version string) {
	selfObj := reflect.Indirect(self)

	// Skip merge if otherObj targets a higher version than what is being generated
	for i := 0; i < otherObj.NumField(); i++ {
		if otherObj.Type().Field(i).Name == "MinVersion" {
			for j := slices.Index(product.ORDER, version) + 1; j < len(product.ORDER); j++ {
				if otherObj.Field(i).String() == product.ORDER[j] {
					return
				}
			}
		}
	}

	for i := 0; i < selfObj.NumField(); i++ {

		// skip if the override is the "empty" value
		emptyOverrideValue := reflect.DeepEqual(reflect.Zero(otherObj.Field(i).Type()).Interface(), otherObj.Field(i).Interface())

		if emptyOverrideValue && selfObj.Type().Field(i).Name != "Required" {
			continue
		}

		if selfObj.Field(i).Kind() == reflect.Slice {
			DeepMerge(selfObj.Field(i), otherObj.Field(i), version)
		} else {
			selfObj.Field(i).Set(otherObj.Field(i))
		}
	}
}

func DeepMerge(arr1, arr2 reflect.Value, version string) {
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
			Merge(currentVal, otherVal, version)
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
