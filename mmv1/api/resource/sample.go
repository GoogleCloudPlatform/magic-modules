// Copyright 2025 Google Inc.
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
	"fmt"
	"slices"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

type IamMember struct {
	Member, Role string
}

type Sample struct {
	Name string

	// If the test should be skipped during VCR testing.
	// This is the case when something about the resource or config causes VCR to fail for example
	// a resource with a unique identifier generated within the resource via id.UniqueId()
	// Or a config with two fine grained resources that have a race condition during create
	SkipVcr bool `yaml:"skip_vcr,omitempty"`

	// The reason to skip a test. For example, a link to a ticket explaining the issue that needs to be resolved before
	// unskipping the test. If this is not empty, the test will be skipped.
	SkipTest string `yaml:"skip_test,omitempty"`

	// Whether to skip generating tests for this resource
	ExcludeTest bool `yaml:"exclude_test,omitempty"`

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

	// The version name of the sample's version if it's different than the
	// resource version, eg. `beta`
	MinVersion string `yaml:"min_version,omitempty"`

	// The id of the "primary" resource in a Test. Used in import test steps.
	// This is the value that will appear in the Terraform config url. For
	// example:
	// resource "google_compute_address" {{primary_resource_id}} {
	//   ...
	// }
	PrimaryResourceId string `yaml:"primary_resource_id"`

	// Optional resource type of the "primary" resource. Used in import tests.
	// If set, this will override the default resource type implied from the
	// object parent
	PrimaryResourceType string `yaml:"primary_resource_type,omitempty"`

	// The name of the primary resource for use in IAM tests. IAM tests need
	// a reference to the primary resource to create IAM policies for
	PrimaryResourceName string `yaml:"primary_resource_name,omitempty"`

	// Steps
	Steps []*Step

	// The version name provided by the user through CI
	TargetVersionName string `yaml:"-"`

	// Step configs that first appears
	NewConfigFuncs []*Step `yaml:"-"`

	// The name of the location/region override for use in IAM tests. IAM
	// tests may need this if the location is not inherited on the resource
	// for one reason or another
	RegionOverride string `yaml:"region_override,omitempty"`

	// ====================
	// TGC
	// ====================
	// The reason to skip a test. For example, a link to a ticket explaining the issue that needs to be resolved before
	// unskipping the test. If this is not empty, the test will be skipped.
	TGCSkipTest string `yaml:"tgc_skip_test,omitempty"`
}

func (s *Sample) TestSampleSlug(productName, resourceName string) string {
	ret := fmt.Sprintf("%s%s_%sExample", productName, resourceName, google.Camelize(s.Name, "lower"))
	return ret
}

func (s *Sample) TestSteps() []*Step {
	return google.Reject(s.Steps, func(st *Step) bool {
		return st.MinVersion != "" && slices.Index(product.ORDER, s.TargetVersionName) < slices.Index(product.ORDER, st.MinVersion)
	})
}

func (s *Sample) ResourceType(terraformName string) string {
	if s.PrimaryResourceType != "" {
		return s.PrimaryResourceType
	}
	return terraformName
}

func (s *Sample) Validate(rName string) (es []error) {
	if s.Name == "" {
		es = append(es, fmt.Errorf("missing `name` for one sample in resource %s", rName))
	}
	es = append(es, s.ValidateExternalProviders()...)

	for _, step := range s.Steps {
		es = append(es, step.Validate(rName, s.Name)...)
	}

	return es
}

func (s *Sample) ValidateExternalProviders() (es []error) {
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
		es = append(es, fmt.Errorf("providers %#v are not allowed. Only providers published by HashiCorp are allowed.", unallowedProviders))
	}

	return es
}
