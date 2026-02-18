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

import (
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/utils"
)

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
// It uses a generic helper to omit fields that are set to their default values.
func (t *Timeouts) MarshalYAML() (interface{}, error) {
	// Use a type alias to prevent infinite recursion.
	type Alias Timeouts

	defaults := NewTimeouts()
	omitted, err := utils.OmitDefaultsForMarshaling(*t, *defaults)
	if err != nil {
		return nil, err
	}

	// If the resulting struct is empty (all fields match defaults), return nil.
	// This ensures we get 'null' instead of '{}' if IsZero wasn't used.
	if utils.IsEmpty(omitted) {
		return nil, nil
	}

	return (*Alias)(omitted.(*Timeouts)), nil
}
