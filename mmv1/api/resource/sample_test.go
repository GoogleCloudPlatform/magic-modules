package resource_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
)

func TestSample_TestDependencies(t *testing.T) {
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
				"services/resourcemanager": "",
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
				"services/compute": "",
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
				"google_compute_": "services/compute",
			},
			want: map[string]string{
				"services/compute": "",
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
				"google_project": "services/resourcemanager",
			},
			want: map[string]string{
				"services/resourcemanager": "",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.sample.TestDependencies(tc.resourcePrefixServiceMap)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TestDependencies() mismatch (-want +got:\n%s", diff)
			}
		})
	}
}

func TestSample_ShouldRunImportStep(t *testing.T) {
	t.Parallel()

	collect := func(s *resource.Sample, total int, steps []*resource.Step) []bool {
		got := make([]bool, total)
		for i := 0; i < total; i++ {
			st := &resource.Step{}
			if i < len(steps) && steps[i] != nil {
				st = steps[i]
			}
			got[i] = s.ShouldRunImportStep(i, total, st)
		}
		return got
	}

	if diff := cmp.Diff([]bool{true, true, true}, collect(&resource.Sample{}, 3, nil)); diff != "" {
		t.Errorf("3-step default mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]bool{true, true, false, false, true}, collect(&resource.Sample{}, 5, nil)); diff != "" {
		t.Errorf("5-step default mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(
		[]bool{false, true, true, false, true},
		collect(&resource.Sample{}, 5, []*resource.Step{
			{ExcludeImportTest: true},
			nil,
			{IncludeImportTest: true},
			nil,
			nil,
		}),
	); diff != "" {
		t.Errorf("5-step overrides mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(
		[]bool{true, true, true, false, true},
		collect(&resource.Sample{ImportAllSteps: true}, 5, []*resource.Step{
			nil, nil, nil, {ExcludeImportTest: true}, nil,
		}),
	); diff != "" {
		t.Errorf("ImportAllSteps mismatch (-want +got):\n%s", diff)
	}

	badStep := &resource.Step{Name: "s", ExcludeImportTest: true, IncludeImportTest: true}
	if errs := badStep.Validate("Topic", "samp"); len(errs) != 1 || !strings.Contains(errs[0].Error(), "cannot set both") {
		t.Errorf("Validate() mutual exclusion = %v, want cannot set both error", errs)
	}
}
