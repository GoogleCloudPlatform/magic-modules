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
	"fmt"
	"strings"
)

// Sweeper provides configuration for the test sweeper to clean up test resources.
// Sweepers are a testing infrastructure mechanism that automatically clean up
// resources created during tests. They run before tests start and can be run
// manually to clean up dangling resources. Sweepers help prevent test failures
// due to resource quota limits and reduce cloud infrastructure costs by removing
// test resources that were not properly cleaned up.
//
// Sweeper generation is enabled by default, except for resources with custom
// deletion code, parent-child relationships (unless configured via Parent), or
// complex URL parameters. Defining the sweeper block overrides these exclusions.
type Sweeper struct {
	// IdentifierField specifies which field in the resource object should be used
	// to identify resources for deletion. If not specified, defaults to "name"
	// if present in the resource, otherwise falls back to "id".
	IdentifierField string `yaml:"identifier_field"`

	// Regions (deprecated - use url_substitutions) defines which regions to run
	// the sweeper in. If empty, defaults to just us-central1. Note that
	// URLSubstitutions provides more granular control over list request parameters,
	// including regions.
	Regions []string `yaml:"regions,omitempty"`

	// Prefixes specifies name prefixes that identify resources eligible for sweeping.
	// Resources whose names start with any of these prefixes will be deleted.
	// By default, resources with the "tf-test-" prefix are automatically eligible
	// for sweeping even if no prefixes are specified here.
	Prefixes []string `yaml:"prefixes,omitempty"`

	// URLSubstitutions allows customizing URL parameters when listing resources.
	// Each map entry represents a set of key-value pairs to substitute in the
	// base_url template when listing resources. This is commonly used to specify
	// regions or other parameters required for the list API call. If not specified,
	// the sweeper will typically only run in the default region (us-central1) and
	// zone (us-central1-a), depending on the resource's base_url structure.
	URLSubstitutions []map[string]string `yaml:"url_substitutions,omitempty"`

	// Dependencies lists other resource types (e.g., "google_compute_instance")
	// that must be swept *before* this resource type. This ensures proper cleanup
	// order for resources with dependencies. If not specified, no dependencies
	// are assumed.
	Dependencies []string `yaml:"dependencies,omitempty"`

	// Parent configures sweeping for resources that depend on parent resources
	// (like a nodepool that belongs to a cluster). When specified, the sweeper
	// will first collect parent resources before listing and deleting child resources.
	// See the ParentResource struct for configuration details.
	Parent *ParentResource `yaml:"parent,omitempty"`

	// QueryString allows appending additional query parameters to the resource's
	// delete URL when performing delete operations. Format should include the
	// starting character, e.g., "?force=true" or "&verbose=true". If not specified,
	// no additional query parameters are added to the delete request.
	QueryString string `yaml:"query_string,omitempty"`

	// EnsureValue specifies a field that must be set to a specific value before
	// deletion can occur. This is used for resources that have fields like
	// 'deletionProtection' that must be explicitly disabled before the API allows
	// deletion. The sweeper automatically handles checking the current value and
	// updating it if necessary before attempting deletion. See the EnsureValue
	// struct for configuration details.
	EnsureValue *EnsureValue `yaml:"ensure_value,omitempty"`
}

// EnsureValue specifies a field and value that must be set before a resource can be deleted.
// Used for resources that have fields like 'deletionProtectionEnabled' that must be
// explicitly disabled before the resource can be deleted.
type EnsureValue struct {
	// Field is the API field name that needs to be updated before deletion.
	// Can include dot notation for nested fields (e.g., "settings.deletionProtectionEnabled").
	Field string `yaml:"field,omitempty"`

	// Value is the required value that Field must be set to before deletion.
	// For boolean fields use "true" or "false", for integers use string representation,
	// for string fields use the exact string value required. The template automatically
	// converts this string to the appropriate type in the API request.
	// Example values: "false", "0", "DISABLED".
	Value string `yaml:"value,omitempty"`

	// IncludeFullResource determines whether to send the entire resource object
	// with the updated field (true) or to send just the field that needs updating (false)
	// in the update request payload. Some APIs require the full resource to be sent
	// in update operations. Defaults to false if not specified.
	IncludeFullResource bool `yaml:"include_full_resource,omitempty"`
}

// ParentResource specifies how to handle parent-child resource dependencies during sweeping.
// It defines how to identify and reference the parent resource when listing or processing
// child resources.
type ParentResource struct {
	// ResourceType is the type name of the parent resource (e.g., "google_container_cluster")
	// used to find the corresponding parent sweeper logic.
	ResourceType string `yaml:"resource_type"`

	// ParentField specifies which field to extract from the parent resource object.
	// This value is then used either directly as the parent identifier or as input
	// for the Template. Example: "name" or "id".
	// Required unless Template is provided.
	ParentField string `yaml:"parent_field"`

	// ParentFieldRegex is a regex pattern with at least one capture group used to
	// extract a specific portion of the ParentField value. The first capture group's
	// match will be used as the final value for substitution in the Template's
	// {{value}} placeholder or as the parent identifier if Template is not used.
	ParentFieldRegex string `yaml:"parent_field_regex"`

	// ParentFieldExtractName, when true, indicates the ParentField contains a self_link
	// URL (e.g., "projects/p/zones/z/instances/i"). It extracts just the final
	// resource name component ("i") from the URL path. This extracted name is then used
	// for substitution in the Template's {{value}} placeholder or as the parent
	// identifier if Template is not used.
	ParentFieldExtractName bool `yaml:"parent_field_extract_name"`

	// ChildField is the field name within the *child* resource's API list/get URL
	// (defined in base_url/self_link) that needs to be populated with the identifier
	// derived from the parent resource (via ParentField/Template).
	// Example: "cluster", "instance".
	ChildField string `yaml:"child_field"`

	// Template provides a format string to construct the parent reference identifier
	// needed in the child resource's URL. Variables in {{curly_braces}} are replaced
	// with values from the parent resource object (e.g., {{project}}, {{location}}).
	// The special placeholder {{value}} is populated with the processed parent field
	// value (obtained from ParentField, potentially modified by ParentFieldRegex or
	// ParentFieldExtractName).
	// Example: "projects/{{project}}/locations/{{location}}/clusters/{{value}}"
	// If specified, Template takes precedence over using the raw ParentField value.
	// At least one of ParentField or Template is required.
	Template string `yaml:"template"`
}

// EnvVarInterpolate takes a string and replaces any environment variable patterns
// with their corresponding function calls, returning a valid Go expression.
func (s Sweeper) EnvVarInterpolate(value string) string {
	// For exact matches, return the function directly
	switch value {
	case "ORG_ID":
		return "envvar.GetTestOrgFromEnv(t)"
	case "ORG_DOMAIN":
		return "envvar.GetTestOrgDomainFromEnv(t)"
	case "CREDENTIALS":
		return "envvar.GetTestCredsFromEnv(t)"
	case "REGION":
		return "envvar.GetTestRegionFromEnv()"
	case "ORG_TARGET":
		return "envvar.GetTestOrgTargetFromEnv(t)"
	case "BILLING_ACCT":
		return "envvar.GetTestBillingAccountFromEnv(t)"
	case "MASTER_BILLING_ACCT":
		return "envvar.GetTestMasterBillingAccountFromEnv(t)"
	case "SERVICE_ACCT":
		return "envvar.GetTestServiceAccountFromEnv(t)"
	case "PROJECT_NAME":
		return "envvar.GetTestProjectFromEnv()"
	case "PROJECT_NUMBER":
		return "envvar.GetTestProjectNumberFromEnv()"
	case "CUST_ID":
		return "envvar.GetTestCustIdFromEnv(t)"
	case "IDENTITY_USER":
		return "envvar.GetTestIdentityUserFromEnv(t)"
	case "PAP_DESCRIPTION":
		return "envvar.GetTestPublicAdvertisedPrefixDescriptionFromEnv(t)"
	case "CHRONICLE_ID":
		return "envvar.GetTestChronicleInstanceIdFromEnv(t)"
	case "VMWAREENGINE_PROJECT":
		return "envvar.GetTestVmwareengineProjectFromEnv(t)"
	case "ZONE":
		return "envvar.GetTestZoneFromEnv()"
	}

	// Check if the string contains any patterns that need to be replaced
	hasPattern := false
	for _, pattern := range []string{
		"${ORG_ID}", "${ORG_DOMAIN}", "${CREDENTIALS}", "${REGION}",
		"${ORG_TARGET}", "${BILLING_ACCT}", "${MASTER_BILLING_ACCT}",
		"${SERVICE_ACCT}", "${PROJECT_NAME}", "${PROJECT_NUMBER}",
		"${CUST_ID}", "${IDENTITY_USER}", "${PAP_DESCRIPTION}",
		"${CHRONICLE_ID}", "${VMWAREENGINE_PROJECT}", "${ZONE}",
	} {
		if strings.Contains(value, pattern) {
			hasPattern = true
			break
		}
	}

	// If no patterns found, return as a string literal
	if !hasPattern {
		return fmt.Sprintf("%q", value)
	}

	// Start with the string as a literal
	result := fmt.Sprintf("%q", value)

	// Replace each pattern with the corresponding function call
	replacements := map[string]string{
		"${ORG_ID}":               "\" + envvar.GetTestOrgFromEnv(t) + \"",
		"${ORG_DOMAIN}":           "\" + envvar.GetTestOrgDomainFromEnv(t) + \"",
		"${CREDENTIALS}":          "\" + envvar.GetTestCredsFromEnv(t) + \"",
		"${REGION}":               "\" + envvar.GetTestRegionFromEnv() + \"",
		"${ORG_TARGET}":           "\" + envvar.GetTestOrgTargetFromEnv(t) + \"",
		"${BILLING_ACCT}":         "\" + envvar.GetTestBillingAccountFromEnv(t) + \"",
		"${MASTER_BILLING_ACCT}":  "\" + envvar.GetTestMasterBillingAccountFromEnv(t) + \"",
		"${SERVICE_ACCT}":         "\" + envvar.GetTestServiceAccountFromEnv(t) + \"",
		"${PROJECT_NAME}":         "\" + envvar.GetTestProjectFromEnv() + \"",
		"${PROJECT_NUMBER}":       "\" + envvar.GetTestProjectNumberFromEnv() + \"",
		"${CUST_ID}":              "\" + envvar.GetTestCustIdFromEnv(t) + \"",
		"${IDENTITY_USER}":        "\" + envvar.GetTestIdentityUserFromEnv(t) + \"",
		"${PAP_DESCRIPTION}":      "\" + envvar.GetTestPublicAdvertisedPrefixDescriptionFromEnv(t) + \"",
		"${CHRONICLE_ID}":         "\" + envvar.GetTestChronicleInstanceIdFromEnv(t) + \"",
		"${VMWAREENGINE_PROJECT}": "\" + envvar.GetTestVmwareengineProjectFromEnv(t) + \"",
		"${ZONE}":                 "\" + envvar.GetTestZoneFromEnv() + \"",
	}

	for pattern, replacement := range replacements {
		result = strings.Replace(result, pattern, replacement, -1)
	}

	// Clean up unnecessary concatenations like "" + and + ""
	result = strings.Replace(result, "\"\" + ", "", -1)
	result = strings.Replace(result, " + \"\"", "", -1)

	return result
}
