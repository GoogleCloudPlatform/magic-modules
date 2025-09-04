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
	"log"
	"slices"
)

// Information about the IAM policy for this resource
// Several GCP resources have IAM policies that are scoped to
// and accessed via their parent resource
// See: https://cloud.google.com/iam/docs/overview
type IamPolicy struct {
	// boolean of if this binding should be generated
	Exclude bool

	// boolean of if this binding should be generated
	ExcludeTgc bool `yaml:"exclude_tgc"`

	// Boolean of if tests for IAM resources should exclude import test steps
	// Used to handle situations where typical generated IAM tests cannot import
	// due to the parent resource having an API-generated id
	ExcludeImportTest bool `yaml:"exclude_import_test"`

	// Character that separates resource identifier from method call in URL
	// For example, PubSub subscription uses {resource}:getIamPolicy
	// While Compute subnetwork uses {resource}/getIamPolicy
	MethodNameSeparator string `yaml:"method_name_separator"`

	// The terraform type (e.g. 'google_endpoints_service') of the parent resource
	// if it is not the same as the IAM resource. The IAP product needs these
	// as its IAM policies refer to compute resources.
	ParentResourceType string `yaml:"parent_resource_type"`

	// Some resources allow retrieving the IAM policy with GET requests,
	// others expect POST requests
	FetchIamPolicyVerb string `yaml:"fetch_iam_policy_verb"`

	// Last part of URL for fetching IAM policy.
	FetchIamPolicyMethod string `yaml:"fetch_iam_policy_method"`

	// Some resources allow setting the IAM policy with POST requests,
	// others expect PUT requests
	SetIamPolicyVerb string `yaml:"set_iam_policy_verb"`

	// Last part of URL for setting IAM policy.
	SetIamPolicyMethod string `yaml:"set_iam_policy_method"`

	// Whether the policy JSON is contained inside of a 'policy' object.
	WrappedPolicyObj bool `yaml:"wrapped_policy_obj"`

	// Certain resources allow different sets of roles to be set with IAM policies
	// This is a role that is acceptable for the given IAM policy resource for use in tests
	AllowedIamRole string `yaml:"allowed_iam_role"`

	// This is a role that grants create/read/delete for the parent resource for use in tests.
	// If set, the test runner will receive a binding to this role in _policy tests in order to
	// avoid getting locked out of the resource.
	AdminIamRole string `yaml:"admin_iam_role"`

	// Certain resources need an attribute other than "id" from their parent resource
	// Especially when a parent is not the same type as the IAM resource
	ParentResourceAttribute string `yaml:"parent_resource_attribute"`

	// If the IAM resource test needs a new project to be created, this is the name of the project
	TestProjectName string `yaml:"test_project_name"`

	// Resource name may need a custom diff suppress function. Default is to use
	// CompareSelfLinkOrResourceName
	CustomDiffSuppress *string `yaml:"custom_diff_suppress"`

	// Some resources (IAP) use fields named differently from the parent resource.
	// We need to use the parent's attributes to create an IAM policy, but they may not be
	// named as the IAM resource expects.
	// This allows us to specify a file (relative to MM root) containing a partial terraform
	// config with the test/example attributes of the IAM resource.
	ExampleConfigBody string `yaml:"example_config_body"`

	// How the API supports IAM conditions
	IamConditionsRequestType string `yaml:"iam_conditions_request_type"`

	// Allows us to override the base_url of the resource. This is required for Cloud Run as the
	// IAM resources use an entirely different base URL from the actual resource
	BaseUrl string `yaml:"base_url"`

	// Allows us to override the import format of the resource. Useful for Cloud Run where we need
	// variables that are outside of the base_url qualifiers.
	ImportFormat []string `yaml:"import_format"`

	// Allows us to override the self_link of the resource. This is required for Artifact Registry
	// to prevent breaking changes
	SelfLink string `yaml:"self_link"`

	// [Optional] Version number in the request payload.
	// if set, it overrides the default IamPolicyVersion
	IamPolicyVersion string `yaml:"iam_policy_version"`

	// [Optional] Min version to make IAM resources available at
	// If unset, defaults to 'ga'
	MinVersion string `yaml:"min_version"`

	// [Optional] Check to see if zone value should be replaced with GOOGLE_ZONE in iam tests
	// Defaults to true
	SubstituteZoneValue bool `yaml:"substitute_zone_value"`
}

func (p *IamPolicy) UnmarshalYAML(unmarshal func(any) error) error {
	p.MethodNameSeparator = "/"
	p.FetchIamPolicyVerb = "GET"
	p.FetchIamPolicyMethod = "getIamPolicy"
	p.SetIamPolicyVerb = "POST"
	p.SetIamPolicyMethod = "setIamPolicy"
	p.WrappedPolicyObj = true
	p.AllowedIamRole = "roles/viewer"
	p.ParentResourceAttribute = "id"
	p.ExampleConfigBody = "templates/terraform/iam/iam_attributes.go.tmpl"
	p.SubstituteZoneValue = true

	type iamPolicyAlias IamPolicy
	aliasObj := (*iamPolicyAlias)(p)

	err := unmarshal(aliasObj)
	if err != nil {
		return err
	}

	return nil
}

func (p *IamPolicy) Validate(rName string) {
	allowed := []string{"GET", "POST"}
	if !slices.Contains(allowed, p.FetchIamPolicyVerb) {
		log.Fatalf("Value on `fetch_iam_policy_verb` should be one of %#v in resource %s", allowed, rName)
	}

	allowed = []string{"POST", "PUT"}
	if !slices.Contains(allowed, p.SetIamPolicyVerb) {
		log.Fatalf("Value on `set_iam_policy_verb` should be one of %#v in resource %s", allowed, rName)
	}

	allowed = []string{"REQUEST_BODY", "QUERY_PARAM", "QUERY_PARAM_NESTED"}
	if p.IamConditionsRequestType != "" && !slices.Contains(allowed, p.IamConditionsRequestType) {
		log.Fatalf("Value on `iam_conditions_request_type` should be one of %#v in resource %s", allowed, rName)
	}
}
