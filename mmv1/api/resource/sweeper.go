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

// Sweeper provides configuration for the test sweeper to clean up test resources
type Sweeper struct {
	// IdentifierField specifies which field in the resource object should be used
	// to identify resources for deletion (typically "name" or "id")
	IdentifierField string `yaml:"identifier_field"`

	// Regions defines which regions to run the sweeper in
	// If empty, defaults to just us-central1
	Regions []string `yaml:"regions,omitempty"`

	// Prefixes specifies name prefixes that identify resources eligible for sweeping
	// Resources whose names start with any of these prefixes will be deleted
	Prefixes []string `yaml:"prefixes,omitempty"`

	// URLSubstitutions allows customizing URL parameters when listing resources
	// Each map entry represents a set of key-value pairs to substitute in the URL template
	URLSubstitutions []map[string]interface{} `yaml:"url_substitutions,omitempty"`

	// Dependencies lists other resource types that must be swept before this one
	Dependencies []string `yaml:"dependencies,omitempty"`

	// Parent defines the parent-child relationship for hierarchical resources
	// When specified, the sweeper will first collect parent resources before listing child resources
	Parent *ParentResource `yaml:"parent,omitempty"`
}

// ParentResource specifies how to handle parent-child resource dependencies
type ParentResource struct {
	// ResourceType is the parent resource type that will be used to find the parent sweeper
	// Example: "GoogleContainerCluster"
	ResourceType string `yaml:"resource_type"`

	// ParentField specifies which field to extract from the parent resource
	// Example: "name" or "id"
	// Required unless Template is provided
	ParentField string `yaml:"parent_field"`

	// ParentFieldRegex is a regex pattern to apply to the parent field value
	// The first capture group will be used as the final value
	ParentFieldRegex string `yaml:"parent_field_regex"`

	// ParentFieldExtractName when true indicates the parent field contains a self-link
	// and only the resource name (portion after the last slash) should be used
	ParentFieldExtractName bool `yaml:"parent_field_extract_name"`

	// ChildField is the field in the child resource that needs to reference the parent
	// Example: "cluster", "instance", etc.
	ChildField string `yaml:"child_field"`

	// Template provides a format string to construct the parent reference
	// Variables in {{brackets}} will be replaced with values from the parent resource
	// The special placeholder {{value}} is populated with the processed parent field value
	// Example: "projects/{{project}}/locations/{{location}}/clusters/{{value}}"
	// If specified, takes precedence over direct field mapping
	Template string `yaml:"template"`
}
