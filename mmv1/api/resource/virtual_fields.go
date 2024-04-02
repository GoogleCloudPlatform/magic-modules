// Copyright 2017 Google Inc.
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

// Virtual fields are Terraform-only fields that control Terraform's
// behaviour. They don't map to underlying API fields (although they
// may map to parameters), and will require custom code to be added to
// control them.
//
// Virtual fields are similar to url_param_only fields in that they create
// a schema entry which is not read from or submitted to the API. However
// virtual fields are meant to provide toggles for Terraform-specific behavior in a resource
// (eg: delete_contents_on_destroy) whereas url_param_only fields _should_
// be used for url construction.
//
// Both are resource level fields and do not make sense, and are also not
// supported, for nested fields. Nested fields that shouldn't be included
// in API payloads are better handled with custom expand/encoder logic.
type VirtualFields struct {
	//< Google::YamlValidator

	// The name of the field in lower snake case.
	Name string

	// The description / docs for the field.
	Description string

	// The API type of the field (defaults to boolean)
	Type string

	// The default value for the field (defaults to false)
	DefaultValue bool `yaml:"default_value"`

	// If set to true, changes in the field's value require recreating the
	// resource.
	Immutable bool
}

// def validate
//   super
//   check :name, type: String, required: true
//   check :description, type: String, required: true
//   check :type, type: Class, default: Api::Type::Boolean
//   check :default_value, default: false
//   check :immutable, default: false
// end
