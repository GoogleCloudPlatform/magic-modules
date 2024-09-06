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

package product

import (
	"log"

	"golang.org/x/exp/slices"
)

var ORDER = []string{"ga", "beta", "alpha", "private"}

// A version of the API for a given product / API group
// In GCP, different product versions are generally ordered where alpha is
// a superset of beta, and beta a superset of GA. Each version will have a
// different version url.
type Version struct {
	CaiBaseUrl string `yaml:"cai_base_url"`
	BaseUrl    string `yaml:"base_url"`
	Name       string
}

func (v *Version) Validate(pName string) {
	if v.Name == "" {
		log.Fatalf("Missing `name` in `version` for product %s", pName)
	}
	if v.BaseUrl == "" {
		log.Fatalf("Missing `base_url` in `version` for product %s", pName)
	}
}

func (v *Version) CompareTo(other *Version) int {
	return slices.Index(ORDER, v.Name) - slices.Index(ORDER, other.Name)
}
