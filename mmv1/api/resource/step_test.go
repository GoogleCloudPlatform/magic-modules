package resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
)

func TestStep_TestServiceDependencies(t *testing.T) {
	cases := []struct {
		name                     string
		step                     resource.Step
		resourcePrefixServiceMap map[string]string
		want                     map[string]string
	}{
		{
			name: "empty",
			step: resource.Step{
				TestContextVars: map[string]string{},
				TestHCLText:     "",
			},
			resourcePrefixServiceMap: map[string]string{},
			want:                     map[string]string{},
		},
		{
			name: "compute bootstrapped",
			step: resource.Step{
				TestContextVars: map[string]string{
					"network": "compute.BootstrapSubnet",
				},
				TestHCLText: "",
			},
			resourcePrefixServiceMap: map[string]string{},
			want: map[string]string{
				"compute": "",
			},
		},
		{
			name: "kms bootstrapped",
			step: resource.Step{
				TestContextVars: map[string]string{
					"kms_key": "kms.BootstrapKMSKey",
				},
				TestHCLText: "",
			},
			resourcePrefixServiceMap: map[string]string{},
			want: map[string]string{
				"kms": "",
			},
		},
		{
			name: "servicenetworking bootstrapped",
			step: resource.Step{
				TestContextVars: map[string]string{
					"network_name": "servicenetworking.BootstrapSharedServiceNetworkingConnection",
				},
				TestHCLText: "",
			},
			resourcePrefixServiceMap: map[string]string{},
			want: map[string]string{
				"servicenetworking": "",
			},
		},
		{
			name: "hcl without prefixes",
			step: resource.Step{
				TestContextVars: map[string]string{},
				TestHCLText: `
resource "google_compute_instance" "foobar" {
}
resource "google_kms_crypto_key" "foobar" {
}`,
			},
			resourcePrefixServiceMap: map[string]string{},
			want:                     map[string]string{},
		},
		{
			name: "prefixes without hcl",
			step: resource.Step{
				TestContextVars: map[string]string{},
				TestHCLText:     "",
			},
			resourcePrefixServiceMap: map[string]string{
				"google_compute_": "compute",
				"google_kms_":     "kms",
			},
			want: map[string]string{},
		},
		{
			name: "hcl with prefixes",
			step: resource.Step{
				TestContextVars: map[string]string{},
				TestHCLText: `
resource "google_compute_instance" "foobar" {
}
resource "google_kms_crypto_key" "foobar" {
}`,
			},
			resourcePrefixServiceMap: map[string]string{
				"google_compute_": "compute",
				"google_kms_":     "kms",
			},
			want: map[string]string{
				"compute": "_",
				"kms":     "_",
			},
		},
		{
			name: "hcl with prefixes plus bootstrapped",
			step: resource.Step{
				TestContextVars: map[string]string{
					"network": "compute.BootstrapSubnet",
					"kms_key": "kms.BootstrapKMSKey",
				},
				TestHCLText: `
resource "google_compute_instance" "foobar" {
}
resource "google_kms_crypto_key" "foobar" {
}`,
			},
			resourcePrefixServiceMap: map[string]string{
				"google_compute_": "compute",
				"google_kms_":     "kms",
			},
			want: map[string]string{
				"compute": "",
				"kms":     "",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.step.TestServiceDependencies(tc.resourcePrefixServiceMap)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TestServiceDependencies() mismatch (-want +got:\n%s", diff)
			}
		})
	}
}
