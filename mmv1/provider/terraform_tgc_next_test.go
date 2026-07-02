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

package provider

import (
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
)

func TestFindIdentityParams(t *testing.T) {
	cases := []struct {
		name     string
		input    []ResourceIdentifier
		expected [][]string // Comparing IdentityParams slices in return sequence
	}{
		{
			name: "single resource with single distinct segment",
			input: []ResourceIdentifier{
				{
					CaiAssetNameFormat: "folders/{{folder}}/feeds/{{name}}",
				},
				{
					CaiAssetNameFormat: "organizations/{{org_id}}/feeds/{{name}}",
				},
				{
					CaiAssetNameFormat: "projects/{{project}}/feeds/{{name}}",
				},
			},
			expected: [][]string{
				{"folders"},
				{"organizations"},
				{"projects"},
			},
		},
		{
			name: "complex multi-segment firewall collision logic",
			input: []ResourceIdentifier{
				{
					// NetworkFirewallPolicy
					CaiAssetNameFormat: "projects/{{project}}/global/firewallPolicies/{{name}}",
				},
				{
					// RegionNetworkFirewallPolicy
					CaiAssetNameFormat: "projects/{{project}}/regions/{{region}}/firewallPolicies/{{name}}",
				},
				{
					// FirewallPolicy
					CaiAssetNameFormat: "locations/global/firewallPolicies/{{name}}",
				},
			},
			expected: [][]string{
				{"projects", "global"},
				{"projects", "regions"},
				{"locations", "global"},
			},
		},
		{
			name: "fallback to import formats when there is a collision",
			input: []ResourceIdentifier{
				{
					CaiAssetNameFormat: "projects/{{project}}/regions/{{region}}/forwardingRules/{{name}}",
					ImportFormats:      []string{"projects/{{project}}/regions/{{region}}/forwardingRules/{{name}}"},
				},
				{
					// Forced collision through CaiAssetNameFormat identical structure
					CaiAssetNameFormat: "projects/{{project}}/regions/{{region}}/forwardingRules/{{name}}",
					ImportFormats:      []string{"projects/{{project}}/global/forwardingRules/{{name}}"},
				},
			},
			expected: [][]string{
				{"regions"},
				{"global"},
			},
		},
		{
			name: "empty identify params grouped at end",
			input: []ResourceIdentifier{
				{
					CaiAssetNameFormat: "projects/{{project}}/global/backendServices/{{name}}",
					ImportFormats:      []string{"projects/{{project}}/global/backendServices/{{name}}"},
				},
				{
					CaiAssetNameFormat: "projects/{{project}}/global/backendServices/{{name}}",
					ImportFormats:      []string{"projects/{{project}}/global/backendServices/{{name}}"},
				},
			},
			expected: [][]string{
				nil,
				nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// FindIdentityParams modifies the input array directly and returns it
			inputCopy := make([]ResourceIdentifier, len(c.input))
			copy(inputCopy, c.input)

			result := FindIdentityParams(inputCopy)
			if len(result) != len(c.expected) {
				t.Fatalf("expected length %d, got %d", len(c.expected), len(result))
			}

			for i, exp := range c.expected {
				if len(result[i].IdentityParams) == 0 && len(exp) == 0 {
					continue // Both represent an empty/nil slice successfully
				}

				if !reflect.DeepEqual(result[i].IdentityParams, exp) {
					t.Errorf("at index %d: expected IdentityParams %v, got %v", i, exp, result[i].IdentityParams)
				}
			}
		})
	}
}

func TestAddTestsFromHandwrittenTests(t *testing.T) {
	mockFS := fstest.MapFS{
		"third_party/terraform/services/dummy/resource_dummy_dummy_test.go": &fstest.MapFile{
			Data: []byte(`func TestAccDummyDummyResource_basic(t *testing.T) {}`),
		},
		"third_party/terraform/services/dummy/resource_dummy_dummy_extra_test.go": &fstest.MapFile{
			Data: []byte(`func TestAccDummyDummyResource_extra(t *testing.T) {}`),
		},
		"third_party/terraform/services/dummy/resource_dummy_dummy_association_test.go": &fstest.MapFile{
			Data: []byte(`func TestAccDummyAssociation_basic(t *testing.T) {}`),
		},
	}

	dummyRes := &api.Resource{
		Name:             "DummyResource",
		FilenameOverride: "dummy",
		ProductMetadata: &api.Product{
			Name: "Dummy",
		},
	}

	dummyAssocRes := &api.Resource{
		Name:             "DummyResourceAssociation",
		FilenameOverride: "dummy_association",
		ProductMetadata: &api.Product{
			Name: "Dummy",
		},
	}

	tgc := TerraformGoogleConversionNext{
		Product: &api.Product{
			Name: "Dummy",
			Objects: []*api.Resource{
				dummyRes,
				dummyAssocRes,
			},
		},
		templateFS: mockFS,
	}

	err := tgc.addTestsFromHandwrittenTests(dummyRes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTests := []string{"TestAccDummyDummyResource_basic", "TestAccDummyDummyResource_extra"}
	if len(dummyRes.TGCTests) != len(expectedTests) {
		t.Errorf("expected %d tests, got %d", len(expectedTests), len(dummyRes.TGCTests))
	}

	foundBasic := false
	foundExtra := false
	foundAssoc := false
	for _, test := range dummyRes.TGCTests {
		if test.Name == "TestAccDummyDummyResource_basic" {
			foundBasic = true
		}
		if test.Name == "TestAccDummyDummyResource_extra" {
			foundExtra = true
		}
		if test.Name == "TestAccDummyAssociation_basic" {
			foundAssoc = true
		}
	}

	if !foundBasic || !foundExtra {
		t.Errorf("did not find all expected tests. dummyRes.TGCTests: %v", dummyRes.TGCTests)
	}
	if foundAssoc {
		t.Errorf("found test from association resource in dummyRes.TGCTests: %v", dummyRes.TGCTests)
	}
}
