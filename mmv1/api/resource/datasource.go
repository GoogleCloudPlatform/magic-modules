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

type Datasource struct {
	// boolean to determine whether the datasource file should be generated
	Generate bool `yaml:"generate"`
	// boolean to determine whether tests should be generated for a datasource
	ExcludeTest bool `yaml:"exclude_test,omitempty"`
	// URI or schema fields forced to Required (overrides default optional scope fields)
	RequiredFields []string `yaml:"required_fields,omitempty"`
	// Additional fields forced to Optional on the datasource schema
	OptionalFields []string `yaml:"optional_fields,omitempty"`
	// Top-level fields to remove from the generated datasource schema
	DeleteFields []string `yaml:"delete_fields,omitempty"`
	// Custom template path executed before constructing the ID and calling resource Read
	PreRead string `yaml:"pre_read,omitempty"`
	// Custom template path executed after resource Read and state post-processing
	PostRead string `yaml:"post_read,omitempty"`
}
