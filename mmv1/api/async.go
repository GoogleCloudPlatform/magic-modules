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
	"reflect"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/utils"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

// Base class from which other Async classes can inherit.
type Async struct {
	// Describes an operation, one of "OpAsync", "PollAsync"
	Type string `yaml:"type,omitempty"`

	// Describes an operation
	Operation *Operation `yaml:"operation,omitempty"`

	// The list of methods where operations are used.
	Actions []string `yaml:"actions,omitempty"`

	OpAsync   `yaml:",inline"`
	PollAsync `yaml:",inline"`
}

func (a Async) Allow(method string) bool {
	return slices.Contains(a.Actions, strings.ToLower(method))
}

func (a Async) IsA(asyncType string) bool {
	return a.Type == asyncType
}

// The main implementation of Operation,
// corresponding to common GCP Operation resources.
type Operation struct {
	Timeouts         *Timeouts
	OpAsyncOperation `yaml:",inline"`
}

func NewOperation() *Operation {
	op := new(Operation)
	op.Timeouts = NewTimeouts()
	return op
}

// It is only used in openapi-generate
func NewAsync() *Async {
	oa := &Async{
		Actions:   []string{"create", "delete", "update"},
		Type:      "OpAsync",
		Operation: NewOperation(),
	}
	return oa
}

// Represents an asynchronous operation definition
type OpAsync struct {
	Result OpAsyncResult `yaml:"result,omitempty"`

	// If true, include project as an argument to OperationWaitTime.
	// It is intended for resources that calculate project/region from a selflink field
	IncludeProject bool `yaml:"include_project,omitempty"`
}

type OpAsyncOperation struct {
	BaseUrl string `yaml:"base_url,omitempty"`
}

// Represents the results of an Operation request
type OpAsyncResult struct {
	ResourceInsideResponse bool `yaml:"resource_inside_response,omitempty"`
}

// Async implementation for polling in Terraform
type PollAsync struct {
	// Details how to poll for an eventually-consistent resource state.

	// Function to call for checking the Poll response for
	// creating and updating a resource
	CheckResponseFuncExistence string `yaml:"check_response_func_existence,omitempty"`

	// Function to call for checking the Poll response for
	// deleting a resource
	CheckResponseFuncAbsence string `yaml:"check_response_func_absence,omitempty"`

	// If true, will suppress errors from polling and default to the
	// result of the final Read()
	SuppressError bool `yaml:"suppress_error,omitempty"`

	// Number of times the desired state has to occur continuously
	// during polling before returning a success
	TargetOccurrences int `yaml:"target_occurrences,omitempty"`
}

// newAsyncWithDefaults returns an Async object with default values set.
func newAsyncWithDefaults() Async {
	a := Async{
		Actions: []string{"create", "delete", "update"},
		Type:    "OpAsync",
	}
	return a
}

func (a *Async) UnmarshalYAML(value *yaml.Node) error {
	// Start with a struct containing all the default values.
	*a = newAsyncWithDefaults()

	type asyncAlias Async
	aliasObj := (*asyncAlias)(a)

	if err := value.Decode(aliasObj); err != nil {
		return err
	}

	if a.Type == "PollAsync" && a.TargetOccurrences == 0 {
		a.TargetOccurrences = 1
	}

	return nil
}

// MarshalYAML implements a custom marshaller for the Async struct.
// It omits fields that are set to their default values.
func (a *Async) MarshalYAML() (interface{}, error) {
	// Use a type alias to prevent infinite recursion during marshaling.
	type asyncAlias Async

	// Create a defaults object that reflects the defaults for the current object's state.
	defaults := newAsyncWithDefaults()

	// Use the generic helper for simple types. It returns a pointer to a clone.
	clone, err := utils.OmitDefaultsForMarshaling(*a, defaults)
	if err != nil {
		return nil, err
	}
	clonePtr := clone.(*Async)

	// The helper ignores slices, so we handle `Actions` manually on the clone.
	if reflect.DeepEqual(clonePtr.Actions, defaults.Actions) {
		clonePtr.Actions = nil
	}

	return (*asyncAlias)(clonePtr), nil
}
