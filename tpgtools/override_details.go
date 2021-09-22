// Copyright 2021 Google LLC. All Rights Reserved.
//
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

// -----------BEFORE ADDING NEW OVERRIDES------------
// Any new override added to tpgtools should have an associated design doc
// See go/tpgtools-new-feature for the template.
// This allows us to further document why overrides were added and
// how they should be used.
// Add these docs to go/tpgtools-overrides when the override is submitted.

// VirtualFieldDetails are the details used to construct a virtual field, a
// Terraform-only field that represents Terraform-specific resource behaviour
// that deviates from the DCL.
type VirtualFieldDetails struct {
	// Name is the name of the property in Terraform, in snake_case.
	Name string
	// Name of the field's type eg "string", "boolean", "integer"
	Type string
	// If set to true, the field is an output-only field.
	Output bool
}

type PreCreateFunctionDetails struct {
	// Function is the name of the function to call. Arguments are decided by
	// the generated code based on the resource's identity, with the format
	// (d, config, res).
	// This function is expected to return an error.
	Function string
}

type PostCreateFunctionDetails struct {
	// Function is the name of the function to call. Arguments are decided by
	// the generated code based on the resource's identity, with the format
	// (d, config, res).
	// This function is expected to return an error.
	Function string
}

type PreDeleteFunctionDetails struct {
	// Function is the name of the function to call. Arguments are decided by
	// the generated code based on the resource's identity, with the format
	// (d, config, res).
	// This function is expected to return an error.
	Function string
}

type CustomImportFunctionDetails struct {
	// Function is the name of the function to call. Arguments are decided by
	// the generated code based on the resource's identity, with the format
	// (d, config, {{identity}}). For example:
	// (d *schema.ResourceData, config *Config, project, location, name string)
	// This function is expected to return an error.
	Function string
}

type AppendToBasePathDetails struct {
	// Append to base path appends this string to the end of the resource's
	// base path.
	String string
}

type CustomizeDiffDetails struct {
	// Functions is a list of CustomizeDiffFunc to use with
	// customdiff.All(...).
	Functions []string
}

type CustomConfigModeDetails struct {
	Mode string
}

type CustomDescriptionDetails struct {
	// Formatted CommonMark description for the property.
	Description string
}

type CustomDiffSuppressFuncDetails struct {
	// Name of the DSF to apply to a property
	DiffSuppressFunc string
}

type CustomIDDetails struct {
	// The pattern string of the Terraform resource's id
	ID string
}

type CustomNameDetails struct {
	// The overriding name of the property in Terraform, in snake_case.
	Name string
}

type CustomValidationDetails struct {
	// Function is the name of a ValidationFunc to apply to a property
	Function string
}

type SetHashFuncDetails struct {
	// Name of the function to determine the unique ID of an item in the set
	Function string
}

type RemovedDetails struct {
	// Message describing a removed field
	Message string
}

type DeprecatedDetails struct {
	// Message describing a deprecated field
	Message string
}

// A CustomIdentityGetter is used to replace the default getX function for
// fields that are inferred from multiple places.
// In the future, this override could be used to attach getX functions to fields
// that don't match the standard name exactly as well.
type CustomIdentityGetterDetails struct {
	// The name of the function to call to retrieve the value. An error is
	// expected to be returned.
	Function string
}

type CustomDefaultDetails struct {
	Default string
}

type CustomListSizeConstraintDetails struct {
	Min int64
	Max int64
}

type CustomRequiredDetails struct {
	Required bool
	Optional bool
	Computed bool
	ForceNew bool
}

type ImportFormatDetails struct {
	// List of import format pattern strings
	Formats []string
}

type MutexDetails struct {
	// The pattern string for the mutex lock name preventing concurrent calls to
	// the resource
	Mutex string
}

type CustomStateGetterDetails struct {
	// The function that is used as the StateGetter.
	Function string
}

type CustomStateSetterDetails struct {
	// The function that is used as the StateSetter.
	Function string
}

type CustomResourceNameDetails struct {
	// The name of the resource in terraform.
	Title string
}

type CustomCreateDirectiveDetails struct {
	// The name of function that will return the create directive.
	Function string
}

type CustomSerializerDetails struct {
	// The name of the function that will serialize this resource.
	Function string
}

type SkipDeleteFunctionDetails struct {
	// The name of the function that determines if we should skip delete.
	Function string
}

type TerraformProductNameDetails struct {
	// The name of the product in terraform.
	Product string
}

type ProductBasePathDetails struct {
	// If set to true, generating the product base path should be skipped
	// This is the case when mmv1 already generates base path support
	Skip bool
	// Alternative product base path name to allow for DCL-based resources to use
	// different default base paths than mmv1 generated resources
	BasePathIdentifier string
}

type ProductTitleDetails struct {
	// alternative name to be used for the product resources
	Title string
}
