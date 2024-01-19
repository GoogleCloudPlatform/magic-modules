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

// require 'api/product'
// require 'api/resource'
// require 'api/type'
// require 'google/yaml_validator'

import (
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

// Process <product>.yaml and produces output module
type Compiler struct {
	Catalog []byte
	Object  interface{}
}

func NewCompiler(catalog []byte, obj interface{}) *Compiler {
	c := Compiler{
		Catalog: catalog,
		Object:  obj,
	}
	return &c
}

func (c *Compiler) Run() {
	// Compile step //1: compile with generic class to instantiate target class
	yamlValidator := google.YamlValidator{}
	yamlValidator.Parse(c.Catalog, c.Object)
	// unless config.class <= Api::Product || config.class <= Api::Resource
	//   raise StandardError, "//{@catalog} is //{config.class} instead of Api::Product"
	// end
}
