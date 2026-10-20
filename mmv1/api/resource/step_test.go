package resource_test

import (
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
)

func TestStep_TestDependencies(t *testing.T) {
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
				"services/compute": "",
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
				"services/kms": "",
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
				"services/servicenetworking": "",
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
				"google_compute_": "services/compute",
				"google_kms_":     "services/kms",
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
				"google_compute_": "services/compute",
				"google_kms_":     "services/kms",
			},
			want: map[string]string{
				"services/compute": "_",
				"services/kms":     "_",
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
				"google_compute_": "services/compute",
				"google_kms_":     "services/kms",
			},
			want: map[string]string{
				"services/compute": "",
				"services/kms":     "",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.step.TestDependencies(tc.resourcePrefixServiceMap)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TestDependencies() mismatch (-want +got:\n%s", diff)
			}
		})
	}
}

func TestStep_SetHCLText_NullVars(t *testing.T) {
	t.Parallel()

	sysfs := fstest.MapFS{
		"templates/terraform/samples/services/appengine/sample.tf.tmpl": &fstest.MapFile{
			Data: []byte(`resource "google_app_engine_standard_app_version" "{{$.PrimaryResourceId}}" {
  service        = "{{index $.ResourceIdVars "service_id"}}"
  vpc_egress     = "{{index $.Vars "vpc_egress"}}"
  egress_setting = "{{index $.Vars "egress_setting"}}"
  empty_str      = "{{index $.Vars "empty_str"}}"
  max_instances  = {{index $.Vars "max_instances"}}
  delete_service = {{index $.Vars "delete_service"}}
  tags           = ["primary", "{{index $.Vars "optional_tag"}}"]
  counts         = [1, {{index $.Vars "optional_count"}}]
}
`),
		},
	}

	rawYAML := `
- name: sample
  config_path: templates/terraform/samples/services/appengine/sample.tf.tmpl
  resource_id_vars:
    service_id: my-service
  vars:
    vpc_egress: "ALL_TRAFFIC"
    egress_setting: null
    empty_str: ""
    max_instances: 10
    delete_service: true
    optional_tag: "v1"
    optional_count: 2
- name: sample
  config_path: templates/terraform/samples/services/appengine/sample.tf.tmpl
  resource_id_vars:
    service_id: my-service
  vars:
    vpc_egress: ~
    egress_setting: "ALL_TRAFFIC"
    empty_str: ""
    max_instances: null
    delete_service: "null"
    optional_tag: null
    optional_count: ~
`

	var steps []*resource.Step
	if err := yaml.Unmarshal([]byte(rawYAML), &steps); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}
	if len(steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(steps))
	}

	for _, st := range steps {
		st.PrimaryResourceId = "primary"
		st.SetHCLText(sysfs)
		st.SetOiCSHCLText(sysfs)
	}

	wantTestHCL := `resource "google_app_engine_standard_app_version" "primary" {
  service        = "%{service_id}"
  vpc_egress     = "%{vpc_egress}"
  egress_setting = "%{egress_setting}"
  empty_str      = "%{empty_str}"
  max_instances  = %{max_instances}
  delete_service = %{delete_service}
  tags           = ["primary", "%{optional_tag}"]
  counts         = [1, %{optional_count}]
}
`

	// Step 0: vpc_egress="ALL_TRAFFIC", egress_setting=null, empty_str="", max_instances=10, delete_service=true, optional_tag="v1", optional_count=2
	wantStep0Context := map[string]string{
		"service_id":     `"tf-test-my-service"+randomSuffix`,
		"vpc_egress":     `"ALL_TRAFFIC"`,
		"egress_setting": "nil",
		"empty_str":      `""`,
		"max_instances":  "10",
		"delete_service": "true",
		"optional_tag":   `"v1"`,
		"optional_count": "2",
	}
	wantStep0Doc := `resource "google_app_engine_standard_app_version" "primary" {
  service        = "my-service"
  vpc_egress     = "ALL_TRAFFIC"
  empty_str      = ""
  max_instances  = 10
  delete_service = true
  tags           = ["primary", "v1"]
  counts         = [1, 2]
}
`
	wantStep0OiCS := `resource "google_app_engine_standard_app_version" "primary" {
  service        = "my-service-${local.name_suffix}"
  vpc_egress     = "ALL_TRAFFIC"
  empty_str      = ""
  max_instances  = 10
  delete_service = true
  tags           = ["primary", "v1"]
  counts         = [1, 2]
}
`

	if diff := cmp.Diff(wantStep0Context, steps[0].TestContextVars); diff != "" {
		t.Errorf("Step 0 TestContextVars mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantTestHCL, steps[0].TestHCLText); diff != "" {
		t.Errorf("Step 0 TestHCLText mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantStep0Doc, steps[0].DocumentationHCLText); diff != "" {
		t.Errorf("Step 0 DocumentationHCLText mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantStep0OiCS, steps[0].OicsHCLText); diff != "" {
		t.Errorf("Step 0 OicsHCLText mismatch (-want +got):\n%s", diff)
	}

	// Step 1: vpc_egress=null (~), egress_setting="ALL_TRAFFIC", empty_str="", max_instances=null, delete_service="null", optional_tag=null, optional_count=null (~)
	wantStep1Context := map[string]string{
		"service_id":     `"tf-test-my-service"+randomSuffix`,
		"vpc_egress":     "nil",
		"egress_setting": `"ALL_TRAFFIC"`,
		"empty_str":      `""`,
		"max_instances":  "nil",
		"delete_service": "nil",
		"optional_tag":   "nil",
		"optional_count": "nil",
	}
	wantStep1Doc := `resource "google_app_engine_standard_app_version" "primary" {
  service        = "my-service"
  egress_setting = "ALL_TRAFFIC"
  empty_str      = ""
  tags           = ["primary", null]
  counts         = [1, null]
}
`
	wantStep1OiCS := `resource "google_app_engine_standard_app_version" "primary" {
  service        = "my-service-${local.name_suffix}"
  egress_setting = "ALL_TRAFFIC"
  empty_str      = ""
  tags           = ["primary", null]
  counts         = [1, null]
}
`

	if diff := cmp.Diff(wantStep1Context, steps[1].TestContextVars); diff != "" {
		t.Errorf("Step 1 TestContextVars mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantTestHCL, steps[1].TestHCLText); diff != "" {
		t.Errorf("Step 1 TestHCLText mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantStep1Doc, steps[1].DocumentationHCLText); diff != "" {
		t.Errorf("Step 1 DocumentationHCLText mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantStep1OiCS, steps[1].OicsHCLText); diff != "" {
		t.Errorf("Step 1 OicsHCLText mismatch (-want +got):\n%s", diff)
	}
}
