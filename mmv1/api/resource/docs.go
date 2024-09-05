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

// Inserts custom strings into terraform resource docs.
type Docs struct {
	// google.YamlValidator

	// All these values should be strings, which will be inserted
	// directly into the terraform resource documentation.  The
	// strings should _not_ be the names of template files
	// (This should be reconsidered if we find ourselves repeating
	// any string more than ones), but rather the actual text
	// (including markdown) which needs to be injected into the
	// template.
	// The text will be injected at the bottom of the specified
	// section.
	// attr_reader :
	Warning string

	// attr_reader :
	Note string

	// attr_reader :
	RequiredProperties string `yaml:"required_properties"`

	// attr_reader :
	OptionalProperties string `yaml:"optional_properties"`

	// attr_reader :
	Attributes string
}
