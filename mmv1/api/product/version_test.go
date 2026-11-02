// Copyright 2026 Google Inc.
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

import "testing"

func TestVersionRegionalUrl(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		version      Version
		wantRegional bool
		wantVariable string
	}{
		{
			name:         "global base url only",
			version:      Version{BaseUrl: "https://foo.googleapis.com/v1/"},
			wantRegional: false,
			wantVariable: "",
		},
		{
			name: "rep url",
			version: Version{
				BaseUrl: "https://foo.googleapis.com/v1/",
				RepUrl:  "https://foo.{{region}}.rep.googleapis.com/v1/",
			},
			wantRegional: true,
			wantVariable: "region",
		},
		{
			name:         "locational base url",
			version:      Version{BaseUrl: "https://{{location}}-foo.googleapis.com/v1/"},
			wantRegional: true,
			wantVariable: "location",
		},
		{
			name: "rep url and locational base url",
			version: Version{
				BaseUrl: "https://{{location}}-foo.googleapis.com/v1/",
				RepUrl:  "https://foo.{{location}}.rep.googleapis.com/v1/",
			},
			wantRegional: true,
			wantVariable: "location",
		},
		{
			name:         "url-encoded variables are ignored",
			version:      Version{BaseUrl: "https://foo.googleapis.com/v1/{{%name}}"},
			wantRegional: false,
			wantVariable: "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.version.HasRegionalUrl(); got != tc.wantRegional {
				t.Errorf("HasRegionalUrl() = %v, want %v", got, tc.wantRegional)
			}
			if got := tc.version.RegionalUrlVariable(); got != tc.wantVariable {
				t.Errorf("RegionalUrlVariable() = %q, want %q", got, tc.wantVariable)
			}
		})
	}
}
