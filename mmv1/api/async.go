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
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

// require 'api/object'
// require 'api/timeout'

// Base class from which other Async classes can inherit.
type Async struct {
	// Embed YamlValidator object
	google.YamlValidator

	// Describes an operation
	Operation *Operation

	// The list of methods where operations are used.
	Actions []string
}

// def validate
//   super

//   check :operation, type: Operation
//   check :actions, default: %w[create delete update], type: ::Array, item_type: ::String
// end

// def allow?(method)
//   @actions.include?(method.downcase)
// end

// Base async operation type
type Operation struct {
	google.YamlValidator

	// Contains information about an long-running operation, to make
	// requests for the state of an operation.

	Timeouts *Timeouts

	Result Result
}

// def validate
//   check :result, type: Result
//   check :timeouts, type: Api::Timeouts
// end

// Base result class
type Result struct {
	google.YamlValidator

	// Contains information about the result of an Operation

	ResourceInsideResponse bool `yaml:"resource_inside_response"`
}

// def validate
//   super
//   check :resource_inside_response, type: :boolean, default: false
// end

// Represents an asynchronous operation definition
type OpAsync struct {
	// TODO: Should embed Async or not?
	// < Async

	Operation *OpAsyncOperation

	Result OpAsyncResult

	Status OpAsyncStatus

	Error OpAsyncError

	// If true, include project as an argument to OperationWaitTime.
	// It is intended for resources that calculate project/region from a selflink field
	IncludeProject bool `yaml:"include_project"`

	// The list of methods where operations are used.
	Actions []string
}

// def initialize(operation, result, status, error)
//   super()
//   @operation = operation
//   @result = result
//   @status = status
//   @error = error
// end

// def validate
//   super

//   check :operation, type: Operation, required: true
//   check :result, type: Result, default: Result.new
//   check :status, type: Status
//   check :error, type: Error
//   check :actions, default: %w[create delete update], type: ::Array, item_type: ::String
//   check :include_project, type: :boolean, default: false
// end

// The main implementation of Operation,
// corresponding to common GCP Operation resources.
type OpAsyncOperation struct {
	// TODO: Should embed Operation or not?
	// < Async::Operation
	Kind string

	Path string

	BaseUrl string `yaml:"base_url"`

	WaitMs int `yaml:"wait_ms"`

	Timeouts *Timeouts

	// Use this if the resource includes the full operation url.
	FullUrl string `yaml:"full_url"`
}

// def initialize(path, base_url, wait_ms, timeouts)
//   super()
//   @path = path
//   @base_url = base_url
//   @wait_ms = wait_ms
//   @timeouts = timeouts
// end

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
	Result Result `yaml:",inline"`

	Path string
}

// def initialize(path = nil, resource_inside_response = nil)
//   super()
//   @path = path
//   @resource_inside_response = resource_inside_response
// end

// def validate
//   super

//   check :path, type: String
// end

// Provides information to parse the result response to check operation
// status
type OpAsyncStatus struct {
	google.YamlValidator

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

// def validate
//   super
//   check :path, type: String
//   check :allowed, type: Array, item_type: [::String, :boolean]
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

// def validate
//   super
//   check :path, type: String
//   check :message, type: String
// end
