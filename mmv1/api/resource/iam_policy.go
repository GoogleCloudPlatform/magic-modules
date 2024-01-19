// Copyright 2019 Google Inc.
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

// Information about the IAM policy for this resource
// Several GCP resources have IAM policies that are scoped to
// and accessed via their parent resource
// See: https://cloud.google.com/iam/docs/overview
type IamPolicy struct {
	// boolean of if this binding should be generated
	exclude bool

	// boolean of if this binding should be generated
	excludeTgc bool

	// Boolean of if tests for IAM resources should exclude import test steps
	// Used to handle situations where typical generated IAM tests cannot import
	// due to the parent resource having an API-generated id
	skipImportTest bool

	// Character that separates resource identifier from method call in URL
	// For example, PubSub subscription uses {resource}:getIamPolicy
	// While Compute subnetwork uses {resource}/getIamPolicy
	methodNameSeparator string

	// The terraform type of the parent resource if it is not the same as the
	// IAM resource. The IAP product needs these as its IAM policies refer
	// to compute resources
	parentResourceType string

	// Some resources allow retrieving the IAM policy with GET requests,
	// others expect POST requests
	fetchIamPolicyVerb string

	// Last part of URL for fetching IAM policy.
	fetchIamPolicyMethod string

	// Some resources allow setting the IAM policy with POST requests,
	// others expect PUT requests
	setIamPolicyVerb string

	// Last part of URL for setting IAM policy.
	setIamPolicyMethod string

	// Whether the policy JSON is contained inside of a 'policy' object.
	wrappedPolicyObj bool

	// Certain resources allow different sets of roles to be set with IAM policies
	// This is a role that is acceptable for the given IAM policy resource for use in tests
	allowedIamRole string

	// This is a role that grants create/read/delete for the parent resource for use in tests.
	// If set, the test runner will receive a binding to this role in _policy tests in order to
	// avoid getting locked out of the resource.
	adminIamRole string

	// Certain resources need an attribute other than "id" from their parent resource
	// Especially when a parent is not the same type as the IAM resource
	parentResourceAttribute string

	// If the IAM resource test needs a new project to be created, this is the name of the project
	testProjectName string

	// Resource name may need a custom diff suppress function. Default is to use
	// CompareSelfLinkOrResourceName
	customDiffSuppress *string

	// Some resources (IAP) use fields named differently from the parent resource.
	// We need to use the parent's attributes to create an IAM policy, but they may not be
	// named as the IAM IAM resource expects.
	// This allows us to specify a file (relative to MM root) containing a partial terraform
	// config with the test/example attributes of the IAM resource.
	exampleConfigBody string

	// How the API supports IAM conditions
	iamConditionsRequestType string

	// Allows us to override the base_url of the resource. This is required for Cloud Run as the
	// IAM resources use an entirely different base URL from the actual resource
	baseUrl string

	// Allows us to override the import format of the resource. Useful for Cloud Run where we need
	// variables that are outside of the base_url qualifiers.
	importFormat []string

	// Allows us to override the self_link of the resource. This is required for Artifact Registry
	// to prevent breaking changes
	selfLink string

	// [Optional] Version number in the request payload.
	// if set, it overrides the default IamPolicyVersion
	iamPolicyVersion string

	// [Optional] Min version to make IAM resources available at
	// If unset, defaults to 'ga'
	minVersion string

	// [Optional] Check to see if zone value should be replaced with GOOGLE_ZONE in iam tests
	// Defaults to true
	substituteZoneValue bool
}

// func (p *IamPolicy) validate() {

// }
