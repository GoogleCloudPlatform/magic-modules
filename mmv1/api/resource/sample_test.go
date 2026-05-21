package resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
)

func TestSample_TestServiceDependencies(t *testing.T) {
	cases := []struct {
		name                     string
		sample                   resource.Sample
		resourcePrefixServiceMap map[string]string
		want                     map[string]string
	}{
		{
			name: "empty",
			sample: resource.Sample{
				Steps: []*resource.Step{
					{
						TestContextVars: map[string]string{},
						TestHCLText:     "",
					},
				},
			},
			resourcePrefixServiceMap: map[string]string{},
			want:                     map[string]string{},
		},
		{
			name: "bootstrap iam",
			sample: resource.Sample{
				BootstrapIam: []resource.IamMember{
					{
						Member: "whatever",
						Role:   "role",
					},
				},
			},
			resourcePrefixServiceMap: map[string]string{},
			want: map[string]string{
				"resourcemanager": "",
			},
		},
		{
			name: "no conflict",
			sample: resource.Sample{
				Steps: []*resource.Step{
					{
						TestContextVars: map[string]string{
							"network": "compute.BootstrapSubnet",
						},
						TestHCLText: "",
					},
					{
						TestContextVars: map[string]string{
							"network": "compute.BootstrapSubnet",
						},
						TestHCLText: "",
					},
				},
			},
			resourcePrefixServiceMap: map[string]string{},
			want: map[string]string{
				"compute": "",
			},
		},
		{
			name: "underscore vs empty string conflict",
			sample: resource.Sample{
				Steps: []*resource.Step{
					{
						TestContextVars: map[string]string{
							"network": "compute.BootstrapSubnet",
						},
						TestHCLText: "",
					},
					{
						TestContextVars: map[string]string{},
						TestHCLText:     `resource "google_compute_instance"`,
					},
				},
			},
			resourcePrefixServiceMap: map[string]string{
				"google_compute_": "compute",
			},
			want: map[string]string{
				"compute": "",
			},
		},
		{
			name: "merge boostrapped iam",
			sample: resource.Sample{
				BootstrapIam: []resource.IamMember{
					{
						Member: "whatever",
						Role:   "role",
					},
				},
				Steps: []*resource.Step{
					{
						TestContextVars: map[string]string{},
						TestHCLText:     `resource "google_project"`,
					},
				},
			},
			resourcePrefixServiceMap: map[string]string{
				"google_project": "resourcemanager",
			},
			want: map[string]string{
				"resourcemanager": "",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.sample.TestServiceDependencies(tc.resourcePrefixServiceMap)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TestServiceDependencies() mismatch (-want +got:\n%s", diff)
			}
		})
	}
}
