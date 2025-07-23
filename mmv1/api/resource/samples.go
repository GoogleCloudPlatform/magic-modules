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

package resource

import (
	// "bytes"
	"fmt"
	// "log"
	// "net/url"
	// "os"
	// "path/filepath"
	// "regexp"
	// "slices"
	// "strings"
	// "text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	// "github.com/golang/glog"
)

type Steps struct {
	Config string `yaml:"config,omitempty"`

	ConfigPath string `yaml:"config_path,omitempty"`

	// IdVars map[string]string

	Vars map[string]string

	GenerateDoc bool `yaml:"generate_doc,omitempty"`

	TestHCLText string `yaml:"-"`
}

type Samples struct {
	Name string

	// The id of the "primary" resource in an example. Used in import tests.
	// This is the value that will appear in the Terraform config url. For
	// example:
	// resource "google_compute_address" {{primary_resource_id}} {
	//   ...
	// }
	PrimaryResourceId string `yaml:"primary_resource_id"`

	PrimaryResourceType string `yaml:"primary_resource_type,omitempty"`

	ExcludeTest bool `yaml:"exclude_test,omitempty"`

	MinVersion string `yaml:"min_version,omitempty"`

	Steps []Steps
}

// Set default value for fields
func (s *Samples) UnmarshalYAML(unmarshal func(any) error) error {
	type sampleAlias Samples
	aliasObj := (*sampleAlias)(s)

	err := unmarshal(aliasObj)
	if err != nil {
		return err
	}

	return nil
}

func (s *Steps) UnmarshalYAML(unmarshal func(any) error) error {
	type stepAlias Steps
	aliasObj := (*stepAlias)(s)

	err := unmarshal(aliasObj)
	if err != nil {
		return err
	}

	return nil
}

func (s *Samples) TestSampleSlug(productName, resourceName string) string {
	ret := fmt.Sprintf("%s%s_%sExample", productName, resourceName, google.Camelize(s.Name, "lower"))
	return ret
}

func (s *Steps) TestStepSlug(productName, resourceName string) string {
	ret := fmt.Sprintf("%s%s_%sExample", productName, resourceName, google.Camelize(s.Config, "lower"))
	return ret
}

func (e *Samples) ResourceType(terraformName string) string {
	if e.PrimaryResourceType != "" {
		return e.PrimaryResourceType
	}
	return terraformName
}
