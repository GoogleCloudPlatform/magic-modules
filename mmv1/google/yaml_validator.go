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

package google

import (
	"log"

	"gopkg.in/yaml.v2"
)

// A helper class to validate contents coming from YAML files.
type YamlValidator struct{}

func (v *YamlValidator) Parse(content []byte, obj interface{}, yamlPath string) {
	if err := yaml.UnmarshalStrict(content, obj); err != nil {
		log.Fatalf("Cannot unmarshal data from file %s: %v", yamlPath, err)
	}
}
