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

// require 'api/object'
// require 'api/timeout'

// Base class from which other Async classes can inherit.
type Async struct {
	// < Google::YamlValidator

	// Describes an operation
	// attr_reader :
	Operation Operation

	// The list of methods where operations are used.
	// attr_reader :
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
	// < Google::YamlValidator

	// Contains information about an long-running operation, to make
	// requests for the state of an operation.

	// attr_reader :
	Timeouts

	// attr_reader :
	Result
}

// def validate
//   check :result, type: Result
//   check :timeouts, type: Api::Timeouts
// end

// Base result class
type Result struct {
	// < Google::YamlValidator

	// Contains information about the result of an Operation

	// attr_reader :
	ResourceInsideResponse bool
}

// def validate
//   super
//   check :resource_inside_response, type: :boolean, default: false
// end

// Represents an asynchronous operation definition
type OpAsync struct {
	// < Async

	// attr_reader :
	Operation OpAsyncOperation

	// attr_reader :
	Result OpAsyncResult

	// attr_reader :
	Status OpAsyncStatus

	// attr_reader :
	Error OpAsyncError

	// If true, include project as an argument to OperationWaitTime.
	// It is intended for resources that calculate project/region from a selflink field
	// attr_reader :
	IncludeProject bool

	// The list of methods where operations are used.
	// attr_reader :
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
	// < Async::Operation
	// attr_reader :
	Kind string

	// attr_reader :
	Path string

	// attr_reader :
	BaseUrl string

	// attr_reader :
	WaitMs int

	// attr_reader :
	Timeouts

	// Use this if the resource includes the full operation url.
	// attr_reader :
	FullUrl string
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
	Result `yaml:",inline"`

	// attr_reader :
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
	// < Google::YamlValidator

	// attr_reader :
	Path string

	// attr_reader :
	Complete bool

	// attr_reader :
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
	// < Google::YamlValidator

	// attr_reader :
	Path string

	// attr_reader :
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
