package resource_test

import (
	"strings"
	"testing"
	"testing/fstest"

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

func TestSample_InheritFromSample(t *testing.T) {
	t.Parallel()

	sysfs := fstest.MapFS{
		"templates/default.tf.tmpl": &fstest.MapFile{
			Data: []byte(`resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  project = "{{index $.TestEnvVars "project"}}"
  name    = "{{index $.ResourceIdVars "topic_id"}}"
  env     = "{{index $.Vars "env"}}"
}
`),
		},
		"templates/override.tf.tmpl": &fstest.MapFile{
			Data: []byte(`resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  project = "{{index $.TestEnvVars "project"}}"
  name    = "{{index $.ResourceIdVars "topic_id"}}"
  env     = "{{index $.Vars "env"}}"
  extra   = "step1"
}
`),
		},
	}

	s := &resource.Sample{
		Name:              "topic_inherit",
		PrimaryResourceId: "primary",
		ConfigPath:        "templates/default.tf.tmpl",
		ResourceIdVars:    map[string]string{"topic_id": "my-topic"},
		Vars:              resource.StepVars{"env": "v1"},
		TestEnvVars:       map[string]string{"project": "PROJECT_NAME"},
		TestVarsOverrides: map[string]string{"topic_id": `"custom-topic"`},
		OicsVarsOverrides: map[string]string{"env": "oics-v1"},
		IgnoreReadExtra:   []string{"etag", "secret"},
		Steps: []*resource.Step{
			{},
			{
				ConfigPath:      "templates/override.tf.tmpl",
				Vars:            resource.StepVars{"env": "v2"},
				IgnoreReadExtra: []string{"secret", "step_extra"},
			},
		},
	}

	for _, st := range s.Steps {
		st.Name = s.Name
		st.PrimaryResourceId = s.PrimaryResourceId
		st.InheritFromSample(s)
		st.SetHCLText(sysfs)
		st.SetOiCSHCLText(sysfs)
	}

	if errs := s.Validate("Topic"); len(errs) != 0 {
		t.Fatalf("Sample.Validate() unexpected errors: %v", errs)
	}
	if got := s.Steps[0].ConfigPath; got != "templates/default.tf.tmpl" {
		t.Errorf("step 0 ConfigPath = %q, want %q", got, "templates/default.tf.tmpl")
	}
	if got := s.Steps[1].ConfigPath; got != "templates/override.tf.tmpl" {
		t.Errorf("step 1 ConfigPath = %q, want %q", got, "templates/override.tf.tmpl")
	}
	if got := s.Steps[0].TestContextVars["topic_id"]; got != `"custom-topic"` {
		t.Errorf("step 0 topic_id TestContextVars = %q, want %q", got, `"custom-topic"`)
	}
	if got := s.Steps[1].TestContextVars["env"]; got != `"v2"` {
		t.Errorf("step 1 env TestContextVars = %q, want %q", got, `"v2"`)
	}
	if !strings.Contains(s.Steps[0].OicsHCLText, `env     = "oics-v1"`) {
		t.Errorf("step 0 OicsHCLText missing inherited oics_vars_overrides:\n%s", s.Steps[0].OicsHCLText)
	}
	wantIgnore := []string{"etag", "secret", "step_extra"}
	if diff := cmp.Diff(wantIgnore, s.Steps[1].IgnoreReadExtra); diff != "" {
		t.Errorf("step 1 IgnoreReadExtra mismatch (-want +got):\n%s", diff)
	}

	// Verify an unused sample-level variable is caught by Step.Validate.
	sUnused := &resource.Sample{
		Name:              "topic_unused_sample_var",
		PrimaryResourceId: "primary",
		ConfigPath:        "templates/default.tf.tmpl",
		ResourceIdVars:    map[string]string{"topic_id": "my-topic"},
		Vars:              resource.StepVars{"env": "v1", "never_used_anywhere": "dead"},
		TestEnvVars:       map[string]string{"project": "PROJECT_NAME"},
		Steps:             []*resource.Step{{Name: "topic_unused_sample_var"}},
	}
	sUnused.Steps[0].PrimaryResourceId = "primary"
	sUnused.Steps[0].InheritFromSample(sUnused)
	sUnused.Steps[0].SetHCLText(sysfs)
	errs := sUnused.Validate("Topic")
	if len(errs) != 1 || !strings.Contains(errs[0].Error(), `unused variable "never_used_anywhere"`) {
		t.Errorf("Validate() = %v, want unused never_used_anywhere error", errs)
	}

	// Verify cross-level vars vs resource_id_vars collision is caught by Step.Validate.
	sCollision := &resource.Sample{
		Name:              "topic_collision",
		PrimaryResourceId: "primary",
		ResourceIdVars:    map[string]string{"topic_id": "my-topic"},
		Steps: []*resource.Step{
			{
				Name: "topic_collision",
				Vars: resource.StepVars{"topic_id": "override-in-wrong-map"},
			},
		},
	}
	sCollision.Steps[0].InheritFromSample(sCollision)
	errs = sCollision.Validate("Topic")
	if len(errs) != 1 || !strings.Contains(errs[0].Error(), "cannot exist in both 'vars' and 'resource_id_vars'") {
		t.Errorf("Validate() = %v, want cross-level vars/resource_id_vars collision error", errs)
	}
}
