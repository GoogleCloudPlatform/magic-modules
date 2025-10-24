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

package resource

import (
	// "bytes"
	"fmt"
	"log"
	// "net/url"
	// "os"
	// "path/filepath"
	// "regexp"
	"slices"
	// "strings"
	// "text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	// "github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	// "github.com/golang/glog"
)

////////////////////////////////////////////////
// TODO: comment out after example.go is removed
////////////////////////////////////////////////
// type IamMember struct {
// 	Member, Role string
// }

type Sample struct {
	Name string

	// If the example should be skipped during VCR testing.
	// This is the case when something about the resource or config causes VCR to fail for example
	// a resource with a unique identifier generated within the resource via id.UniqueId()
	// Or a config with two fine grained resources that have a race condition during create
	SkipVcr bool `yaml:"skip_vcr,omitempty"`

	// The reason to skip a test. For example, a link to a ticket explaining the issue that needs to be resolved before
	// unskipping the test. If this is not empty, the test will be skipped.
	SkipTest string `yaml:"skip_test,omitempty"`

	// Specify which external providers are needed for the testcase.
	// Think before adding as there is latency and adds an external dependency to
	// your test so avoid if you can.
	ExternalProviders []string `yaml:"external_providers,omitempty"`

	// BootstrapIam will automatically bootstrap the given member/role pairs.
	// This should be used in cases where specific IAM permissions must be
	// present on the default test project, to avoid race conditions between
	// tests. Permissions attached to resources created in a test should instead
	// be provisioned with standard terraform resources.
	BootstrapIam []IamMember `yaml:"bootstrap_iam,omitempty"`

	// The version name of of the example's version if it's different than the
	// resource version, eg. `beta`
	//
	// This should be the highest version of all the features used in the
	// example; if there's a single beta field in an example, the example's
	// min_version is beta. This is only needed if an example uses features
	// with a different version than the resource; a beta resource's examples
	// are all automatically versioned at beta.
	//
	// When an example has a version of beta, each resource must use the
	// `google-beta` provider in the config. If the `google` provider is
	// implicitly used, the test will fail.
	//
	// NOTE: Until Terraform 0.12 is released and is used in the OiCS tests, an
	// explicit provider block should be defined. While the tests @ 0.12 will
	// use `google-beta` automatically, past Terraform versions required an
	// explicit block.
	MinVersion string `yaml:"min_version,omitempty"`

	// The version name provided by the user through CI
	TargetVersionName string `yaml:"-"`

	// The id of the "primary" resource in an example. Used in import tests.
	// This is the value that will appear in the Terraform config url. For
	// example:
	// resource "google_compute_address" {{primary_resource_id}} {
	//   ...
	// }
	PrimaryResourceId string `yaml:"primary_resource_id"`

	PrimaryResourceType string `yaml:"primary_resource_type,omitempty"`

	// The name of the primary resource for use in IAM tests. IAM tests need
	// a reference to the primary resource to create IAM policies for
	PrimaryResourceName string `yaml:"primary_resource_name,omitempty"`

	ExcludeTest bool `yaml:"exclude_test,omitempty"`

	Steps []Step

	NewConfigFuncs []Step `yaml:"-"`

	// The name of the location/region override for use in IAM tests. IAM
	// tests may need this if the location is not inherited on the resource
	// for one reason or another
	RegionOverride string `yaml:"region_override,omitempty"`

	// ====================
	// TGC
	// ====================
	// Extra properties to ignore test.
	// These properties are present in Terraform resources schema, but not in CAI assets.
	// Virtual Fields and url parameters are already ignored by default and do not need to be duplicated here.
	TGCTestIgnoreExtra []string `yaml:"tgc_test_ignore_extra,omitempty"`
	// The properties ignored in CAI assets. It is rarely used and only used
	// when the nested field has sent_empty_value: true.
	// But its parent field is C + O and not specified in raw_config.
	// Example: ['RESOURCE.cdnPolicy.signedUrlCacheMaxAgeSec'].
	// "RESOURCE" means that the property is for resource data in CAI asset.
	TGCTestIgnoreInAsset []string `yaml:"tgc_test_ignore_in_asset,omitempty"`
	// The reason to skip a test. For example, a link to a ticket explaining the issue that needs to be resolved before
	// unskipping the test. If this is not empty, the test will be skipped.
	TGCSkipTest string `yaml:"tgc_skip_test,omitempty"`
}

// Set default value for fields
func (s *Sample) UnmarshalYAML(unmarshal func(any) error) error {
	type sampleAlias Sample
	aliasObj := (*sampleAlias)(s)

	err := unmarshal(aliasObj)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sample) TestSampleSlug(productName, resourceName string) string {
	ret := fmt.Sprintf("%s%s_%sExample", productName, resourceName, google.Camelize(s.Name, "lower"))
	return ret
}

func (s *Sample) TestSteps() []Step {
	return google.Reject(s.Steps, func(st Step) bool {
		return st.MinVersion != "" && slices.Index(product.ORDER, s.TargetVersionName) < slices.Index(product.ORDER, st.MinVersion)
	})
}

func (s *Sample) ResourceType(terraformName string) string {
	if s.PrimaryResourceType != "" {
		return s.PrimaryResourceType
	}
	return terraformName
}

func (s *Sample) Validate(rName string) {
	if s.Name == "" {
		log.Fatalf("Missing `name` for one example in resource %s", rName)
	}
	s.ValidateExternalProviders()

	for _, step := range s.Steps {
		step.Validate(rName, s.Name)
	}
}

func (s *Sample) ValidateExternalProviders() {
	// Official providers supported by HashiCorp
	// https://registry.terraform.io/search/providers?namespace=hashicorp&tier=official
	HASHICORP_PROVIDERS := []string{"aws", "random", "null", "template", "azurerm", "kubernetes", "local",
		"external", "time", "vault", "archive", "tls", "helm", "azuread", "http", "cloudinit", "tfe", "dns",
		"consul", "vsphere", "nomad", "awscc", "googleworkspace", "hcp", "boundary", "ad", "azurestack", "opc",
		"oraclepaas", "hcs", "salesforce"}

	var unallowedProviders []string
	for _, p := range s.ExternalProviders {
		if !slices.Contains(HASHICORP_PROVIDERS, p) {
			unallowedProviders = append(unallowedProviders, p)
		}
	}

	if len(unallowedProviders) > 0 {
		log.Fatalf("Providers %#v are not allowed. Only providers published by HashiCorp are allowed.", unallowedProviders)
	}
}
