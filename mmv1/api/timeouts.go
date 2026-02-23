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

package api

// Default timeout for all operation types is 20, the Terraform default
// (https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts)
// minutes. This can be overridden for each resource.
const DEFAULT_INSERT_TIMEOUT_MINUTES = 20
const DEFAULT_UPDATE_TIMEOUT_MINUTES = 20
const DEFAULT_DELETE_TIMEOUT_MINUTES = 20

// Provides timeout information for the different operation types
type Timeouts struct {
	InsertMinutes int `yaml:"insert_minutes,omitempty"`
	UpdateMinutes int `yaml:"update_minutes,omitempty"`
	DeleteMinutes int `yaml:"delete_minutes,omitempty"`
}

func NewTimeouts() *Timeouts {
	return &Timeouts{
		InsertMinutes: DEFAULT_INSERT_TIMEOUT_MINUTES,
		UpdateMinutes: DEFAULT_UPDATE_TIMEOUT_MINUTES,
		DeleteMinutes: DEFAULT_DELETE_TIMEOUT_MINUTES,
	}
}

// IsZero enables the omitempty tag on the parent Resource struct.
// If Timeouts matches the default values, it is considered zero and omitted entirely.
func (t *Timeouts) IsZero() bool {
	defaults := NewTimeouts()
	return *t == *defaults
}

// MarshalYAML implements a custom marshaller for the Timeouts struct.
func (t *Timeouts) MarshalYAML() (interface{}, error) {
	// Use a type alias to prevent infinite recursion.
	type Alias Timeouts

	defaults := NewTimeouts()

	// TEMP: Retain legacy behavior where we only strip the block if
	// ALL values are default. If any value differs, we write the full block.
	// This prevents partial objects like { insert: 40 } (implicitly update=20)
	// from being written if the intention is { insert: 40, update: 20, delete: 20 }.
	//
	// If it matches defaults exactly, return nil to print 'null' (or omit via IsZero).
	if *t == *defaults {
		return nil, nil
	}

	// Return the struct as is (cast to Alias).
	// Because the fields have `omitempty` but are `int`, they will print
	// unless they are 0. Since defaults are 20, they will print explicitly.
	return (*Alias)(t), nil
}
