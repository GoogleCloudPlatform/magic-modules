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

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

var ORDER = []string{"ga", "beta", "nightly", "alpha", "private", "internal"}

// A version of the API for a given product / API group
// In GCP, different product versions are generally ordered where alpha is
// a superset of beta, and beta a superset of GA. Each version will have a
// different version url.
type Version struct {
	CaiBaseUrl       string `yaml:"cai_base_url,omitempty"`
	CaiLegacyBaseUrl string `yaml:"cai_legacy_base_url,omitempty"`
	BaseUrl          string `yaml:"base_url"`
	Name             string
	RepUrl           string `yaml:"rep_url,omitempty"`

	// EXPERIMENTAL: RPC settings are not fully implemented, and should not be
	// used at this time.
	RPCAddress string `yaml:"rpc_address,omitempty"`
	RPCPackage string `yaml:"rpc_package,omitempty"`
}

func (v *Version) Validate(pName string) {
	if v.Name == "" {
		log.Fatalf("Missing `name` in `version` for product %s", pName)
	}
	if v.BaseUrl == "" {
		log.Fatalf("Missing `base_url` in `version` for product %s", pName)
	}
	if v.RepUrl != "" {
		if n := len(google.ExtractTemplateVariables(v.RepUrl)); n != 1 {
			log.Fatalf("`rep_url` %q in `version` %s for product %s must contain exactly one template variable (e.g. {{location}} or {{region}}), found %d", v.RepUrl, v.Name, pName, n)
		}
	}
	baseVars := google.ExtractTemplateVariables(v.BaseUrl)
	if n := len(baseVars); n > 1 {
		log.Fatalf("`base_url` %q in `version` %s for product %s must contain at most one template variable (e.g. {{location}} or {{region}}), found %d", v.BaseUrl, v.Name, pName, n)
	}
	// Generated code substitutes a single variable name into whichever URL
	// transport_tpg.BaseUrl returns, so base_url and rep_url must agree on it.
	if v.RepUrl != "" && len(baseVars) == 1 {
		if repVar := google.ExtractTemplateVariables(v.RepUrl)[0]; repVar != baseVars[0] {
			log.Fatalf("`base_url` %q and `rep_url` %q in `version` %s for product %s must use the same template variable, found %q and %q", v.BaseUrl, v.RepUrl, v.Name, pName, baseVars[0], repVar)
		}
	}
}

func (v *Version) CompareTo(other *Version) int {
	return slices.Index(ORDER, v.Name) - slices.Index(ORDER, other.Name)
}

// HasRegionalUrl returns whether this version has a regionalized endpoint:
// either a regional endpoint policy (REP) URL, or a base URL that is itself
// templated on a location (for example
// https://{{location}}-documentai.googleapis.com/v1/). The default of
// regional vs global for REP is controlled at the product level.
func (v *Version) HasRegionalUrl() bool {
	return v.RepUrl != "" || len(google.ExtractTemplateVariables(v.BaseUrl)) > 0
}

// RegionalUrlVariable returns the name of the template variable in the
// regional URL (for example "location" for
// https://foo.{{location}}.rep.googleapis.com/v1/). RepUrl is preferred,
// falling back to BaseUrl. Validate ensures each contains at most one variable
// and that they match. Returns an empty string if this version has no
// regional URL.
func (v *Version) RegionalUrlVariable() string {
	if vars := google.ExtractTemplateVariables(v.RepUrl); len(vars) > 0 {
		return vars[0]
	}
	if vars := google.ExtractTemplateVariables(v.BaseUrl); len(vars) > 0 {
		return vars[0]
	}
	return ""
}
