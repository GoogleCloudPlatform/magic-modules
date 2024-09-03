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

// Default timeout for all operation types is 20, the Terraform default (https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts)
// minutes. This can be overridden for each resource.
const DEFAULT_INSERT_TIMEOUT_MINUTES = 20
const DEFAULT_UPDATE_TIMEOUT_MINUTES = 20
const DEFAULT_DELETE_TIMEOUT_MINUTES = 20

// Provides timeout information for the different operation types
type Timeouts struct {
	// google.YamlValidator
	InsertMinutes int `yaml:"insert_minutes"`
	UpdateMinutes int `yaml:"update_minutes"`
	DeleteMinutes int `yaml:"delete_minutes"`
}

// def initialize
//   super

//   validate
// end

func NewTimeouts() *Timeouts {
	return &Timeouts{
		InsertMinutes: DEFAULT_INSERT_TIMEOUT_MINUTES,
		UpdateMinutes: DEFAULT_UPDATE_TIMEOUT_MINUTES,
		DeleteMinutes: DEFAULT_DELETE_TIMEOUT_MINUTES,
	}
}
