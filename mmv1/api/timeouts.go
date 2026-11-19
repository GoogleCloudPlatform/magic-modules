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

import "fmt"

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

// SetDefault populates any unset (0) timeout fields with the 20-minute default.
// This runs after YAML unmarshaling and override merging (api.Merge) so that a
// partial override timeouts block does not clobber non-default base timeouts.
func (t *Timeouts) SetDefault() {
	if t == nil {
		return
	}
	if t.InsertMinutes == 0 {
		t.InsertMinutes = DEFAULT_INSERT_TIMEOUT_MINUTES
	}
	if t.UpdateMinutes == 0 {
		t.UpdateMinutes = DEFAULT_UPDATE_TIMEOUT_MINUTES
	}
	if t.DeleteMinutes == 0 {
		t.DeleteMinutes = DEFAULT_DELETE_TIMEOUT_MINUTES
	}
}

// withDefaults returns a copy of t with any unset (0) timeout fields populated
// with their default values, or nil if t is nil.
func (t *Timeouts) withDefaults() *Timeouts {
	if t == nil {
		return nil
	}
	resolved := *t
	resolved.SetDefault()
	return &resolved
}

// Validate checks that no timeout fields are negative.
func (t *Timeouts) Validate(rName string) (es []error) {
	if t == nil {
		return nil
	}
	if t.InsertMinutes < 0 {
		es = append(es, fmt.Errorf("`insert_minutes` on %s must be > 0, got %d", rName, t.InsertMinutes))
	}
	if t.UpdateMinutes < 0 {
		es = append(es, fmt.Errorf("`update_minutes` on %s must be > 0, got %d", rName, t.UpdateMinutes))
	}
	if t.DeleteMinutes < 0 {
		es = append(es, fmt.Errorf("`delete_minutes` on %s must be > 0, got %d", rName, t.DeleteMinutes))
	}
	return es
}

type timeoutsYAML struct {
	InsertMinutes *int `yaml:"insert_minutes,omitempty"`
	UpdateMinutes *int `yaml:"update_minutes,omitempty"`
	DeleteMinutes *int `yaml:"delete_minutes,omitempty"`
}

// UnmarshalYAML validates that any explicitly configured timeout values are
// positive (> 0) while leaving omitted fields as 0 until SetDefault / GetTimeouts
// so that partial override blocks merge cleanly in api.Merge.
func (t *Timeouts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw timeoutsYAML
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*t = Timeouts{}
	if raw.InsertMinutes != nil {
		if *raw.InsertMinutes <= 0 {
			return fmt.Errorf("insert_minutes must be > 0, got %d", *raw.InsertMinutes)
		}
		t.InsertMinutes = *raw.InsertMinutes
	}
	if raw.UpdateMinutes != nil {
		if *raw.UpdateMinutes <= 0 {
			return fmt.Errorf("update_minutes must be > 0, got %d", *raw.UpdateMinutes)
		}
		t.UpdateMinutes = *raw.UpdateMinutes
	}
	if raw.DeleteMinutes != nil {
		if *raw.DeleteMinutes <= 0 {
			return fmt.Errorf("delete_minutes must be > 0, got %d", *raw.DeleteMinutes)
		}
		t.DeleteMinutes = *raw.DeleteMinutes
	}

	return nil
}

// IsZero enables the omitempty tag on the parent Resource struct.
// If Timeouts matches the default values, it is considered zero and omitted entirely.
func (t *Timeouts) IsZero() bool {
	if t == nil {
		return true
	}
	defaults := NewTimeouts()
	return *t.withDefaults() == *defaults
}

// MarshalYAML implements a custom marshaller for the Timeouts struct.
func (t *Timeouts) MarshalYAML() (interface{}, error) {
	if t == nil {
		return nil, nil
	}
	// Use a type alias to prevent infinite recursion.
	type timeoutsAlias Timeouts

	defaults := NewTimeouts()
	resolved := t.withDefaults()

	// TEMP: Retain legacy behavior where we only strip the block if
	// ALL values are default. If any value differs, we write the full block.
	// This prevents partial objects like { insert: 40 } (implicitly update=20)
	// from being written if the intention is { insert: 40, update: 20, delete: 20 }.
	//
	// If it matches defaults exactly, return nil to print 'null' (or omit via IsZero).
	if *resolved == *defaults {
		return nil, nil
	}

	// Return the resolved struct (cast to timeoutsAlias).
	// Because defaults are 20 (non-zero), all three fields print explicitly.
	return (*timeoutsAlias)(resolved), nil
}
