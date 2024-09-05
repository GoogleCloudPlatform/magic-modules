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
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"golang.org/x/exp/slices"
)

// Base class from which other Async classes can inherit.
type Async struct {
	// Embed YamlValidator object
	// google.YamlValidator

	// Describes an operation
	Operation *Operation

	// The list of methods where operations are used.
	Actions []string

	// Describes an operation, one of "OpAsync", "PollAsync"
	Type string

	OpAsync `yaml:",inline"`

	PollAsync `yaml:",inline"`
}

// def allow?(method)
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

// def initialize(path, base_url, wait_ms, timeouts)
func NewOperation() *Operation {
	//   super()
	op := new(Operation)
	op.Timeouts = NewTimeouts()
	return op
}

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
	Result OpAsyncResult

	Status OpAsyncStatus

	Error OpAsyncError

	// If true, include project as an argument to OperationWaitTime.
	// It is intended for resources that calculate project/region from a selflink field
	IncludeProject bool `yaml:"include_project"`
}

// def initialize(operation, result, status, error)
//   super()
//   @operation = operation
//   @result = result
//   @status = status
//   @error = error
// end

type OpAsyncOperation struct {
	Kind string

	Path string

	BaseUrl string `yaml:"base_url"`

	WaitMs int `yaml:"wait_ms"`

	// Use this if the resource includes the full operation url.
	FullUrl string `yaml:"full_url"`
}

// def validate
//   super

//   check :kind, type: String
//   check :path, type: String
//   check :base_url, type: String
//   check :wait_ms, type: Integer

//   check :full_url, type: String

//   conflicts %i[base_url full_url]
// end

// Represents the results of an Operation request
type OpAsyncResult struct {
	ResourceInsideResponse bool `yaml:"resource_inside_response"`

	Path string
}

// def initialize(path = nil, resource_inside_response = nil)
//   super()
//   @path = path
//   @resource_inside_response = resource_inside_response
// end

// Provides information to parse the result response to check operation
// status
type OpAsyncStatus struct {
	// google.YamlValidator

	Path string

	Complete bool

	Allowed []bool
}

// def initialize(path, complete, allowed)
//   super()
//   @path = path
//   @complete = complete
//   @allowed = allowed
// end

// Provides information on how to retrieve errors of the executed operations
type OpAsyncError struct {
	google.YamlValidator

	Path string

	Message string
}

// def initialize(path, message)
//   super()
//   @path = path
//   @message = message
// end

// Async implementation for polling in Terraform
type PollAsync struct {
	// Details how to poll for an eventually-consistent resource state.

	// Function to call for checking the Poll response for
	// creating and updating a resource
	CheckResponseFuncExistence string `yaml:"check_response_func_existence"`

	// Function to call for checking the Poll response for
	// deleting a resource
	CheckResponseFuncAbsence string `yaml:"check_response_func_absence"`

	// Custom code to get a poll response, if needed.
	// Will default to same logic as Read() to get current resource
	CustomPollRead string `yaml:"custom_poll_read"`

	// If true, will suppress errors from polling and default to the
	// result of the final Read()
	SuppressError bool `yaml:"suppress_error"`

	// Number of times the desired state has to occur continuously
	// during polling before returning a success
	TargetOccurrences int `yaml:"target_occurrences"`
}

func (a *Async) UnmarshalYAML(unmarshal func(any) error) error {
	a.Actions = []string{"create", "delete", "update"}
	type asyncAlias Async
	aliasObj := (*asyncAlias)(a)

	err := unmarshal(aliasObj)
	if err != nil {
		return err
	}

	if a.Type == "PollAsync" && a.TargetOccurrences == 0 {
		a.TargetOccurrences = 1
	}

	return nil
}

func (a *Async) Validate() {
	if a.Type == "OpAsync" {
		if a.Operation == nil {
			log.Fatalf("Missing `Operation` for OpAsync")
		} else {
			if a.Operation.BaseUrl != "" && a.Operation.FullUrl != "" {
				log.Fatalf("`base_url` and `full_url` cannot be set at the same time in OpAsync operation.")
			}
		}
	}
}
