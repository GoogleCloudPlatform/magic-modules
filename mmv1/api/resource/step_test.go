package resource_test

import (
	"strings"
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

func TestStep_Validate_UnusedVars(t *testing.T) {
	t.Parallel()

	sysfs := fstest.MapFS{
		"templates/terraform/samples/services/pubsub/used.tf.tmpl": &fstest.MapFile{
			Data: []byte(`resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  name    = "{{index $.ResourceIdVars "topic_name"}}"
  project = "{{index $.TestEnvVars "project"}}"
  labels = {
    env = "{{index $.Vars "env"}}"
  }
  source = "path/to/index.zip"
}
`),
		},
	}

	t.Run("all variables used passes validation", func(t *testing.T) {
		t.Parallel()
		st := &resource.Step{
			Name:       "used",
			ConfigPath: "templates/terraform/samples/services/pubsub/used.tf.tmpl",
			ResourceIdVars: map[string]string{
				"topic_name": "my-topic",
			},
			Vars: resource.StepVars{
				"env": "dev",
			},
			TestEnvVars: map[string]string{
				"project":          "PROJECT_NAME",
				"consumer_project": "PROJECT_NAME",
			},
			TestVarsOverrides: map[string]string{
				"topic_name": `"custom-topic"`,
				"zip_path":   `"./test-fixtures/index.zip"`,
			},
			OicsVarsOverrides: map[string]string{
				"env": "oics-dev",
			},
		}
		st.PrimaryResourceId = "primary"
		st.SetHCLText(sysfs)
		st.IamHCLText = `consumer_project = "%{consumer_project}"`

		if errs := st.Validate("Topic", "pubsub_topic_used"); len(errs) != 0 {
			t.Errorf("Validate() returned unexpected errors: %v", errs)
		}
	})

	t.Run("unused variables across maps fail validation", func(t *testing.T) {
		t.Parallel()
		st := &resource.Step{
			Name:       "used",
			ConfigPath: "templates/terraform/samples/services/pubsub/used.tf.tmpl",
			ResourceIdVars: map[string]string{
				"topic_name": "my-topic",
				"unused_id":  "unused",
				"paired_id":  "paired",
			},
			Vars: resource.StepVars{
				"env":        "dev",
				"unused_var": "val",
			},
			TestEnvVars: map[string]string{
				"project":    "PROJECT_NAME",
				"unused_env": "ORG_ID",
			},
			TestVarsOverrides: map[string]string{
				"unused_override": `"val"`,
				"paired_id":       `"val"`,
			},
			OicsVarsOverrides: map[string]string{
				"orphan_oics": "val",
				"paired_id":   "val",
			},
		}
		st.PrimaryResourceId = "primary"
		st.SetHCLText(sysfs)

		errs := st.Validate("Topic", "pubsub_topic_unused")
		got := make([]string, len(errs))
		for i, err := range errs {
			got[i] = err.Error()
		}
		want := []string{
			`unused variable "orphan_oics" in 'oics_vars_overrides' (not defined in 'vars' or 'resource_id_vars') for step "used" in sample "pubsub_topic_unused" of resource "Topic"`,
			`unused variable "unused_var" in 'vars' for step "used" in sample "pubsub_topic_unused" of resource "Topic"`,
			`unused variable "paired_id" in 'resource_id_vars' for step "used" in sample "pubsub_topic_unused" of resource "Topic"`,
			`unused variable "unused_id" in 'resource_id_vars' for step "used" in sample "pubsub_topic_unused" of resource "Topic"`,
			`unused variable "unused_env" in 'test_env_vars' for step "used" in sample "pubsub_topic_unused" of resource "Topic"`,
			`unused variable "paired_id" in 'test_vars_overrides' for step "used" in sample "pubsub_topic_unused" of resource "Topic"`,
			`unused variable "unused_override" in 'test_vars_overrides' for step "used" in sample "pubsub_topic_unused" of resource "Topic"`,
			`unused variable "paired_id" in 'oics_vars_overrides' for step "used" in sample "pubsub_topic_unused" of resource "Topic"`,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Validate() errors mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("unloaded step skips template usage checks", func(t *testing.T) {
		t.Parallel()
		st := &resource.Step{
			Name:       "beta_only_step",
			MinVersion: "beta",
			ResourceIdVars: map[string]string{
				"topic_name": "my-topic",
			},
			Vars: resource.StepVars{
				"env": "dev",
			},
			TestEnvVars: map[string]string{
				"project": "PROJECT_NAME",
			},
			TestVarsOverrides: map[string]string{
				"topic_name": `"custom"`,
			},
		}
		if errs := st.Validate("Topic", "pubsub_topic_beta"); len(errs) != 0 {
			t.Errorf("Validate() on unloaded step returned unexpected errors: %v", errs)
		}
	})
}

func TestStep_TestStepSlug(t *testing.T) {
	t.Parallel()

	st := &resource.Step{Name: "pubsub_topic_basic"}
	if got, want := st.TestStepSlug("Pubsub", "Topic"), "PubsubTopic_pubsubTopicBasicExample"; got != want {
		t.Errorf("TestStepSlug() = %q, want %q", got, want)
	}

	st.StepSlug = "pubsub_topic_exhaustive_step2"
	if got, want := st.TestStepSlug("Pubsub", "Topic"), "PubsubTopic_pubsubTopicExhaustiveStep2Example"; got != want {
		t.Errorf("TestStepSlug() with StepSlug = %q, want %q", got, want)
	}
}

func TestStep_ExecuteTemplate_StrictIndexAndMissingKey(t *testing.T) {
	t.Parallel()

	sysfs := fstest.MapFS{
		"templates/missing_var.tf.tmpl": &fstest.MapFile{
			Data: []byte(`name = "{{printf "%s" (index $.Vars "missing_key")}}"`),
		},
		"templates/missing_dot_key.tf.tmpl": &fstest.MapFile{
			Data: []byte(`name = "{{$.Vars.missing_key}}"`),
		},
		"templates/nil_map.tf.tmpl": &fstest.MapFile{
			Data: []byte(`name = "{{index $.ResourceIdVars "id"}}"`),
		},
	}

	t.Run("nested index on missing Vars key returns error", func(t *testing.T) {
		t.Parallel()
		st := &resource.Step{
			Name:       "s",
			ConfigPath: "templates/missing_var.tf.tmpl",
			Vars:       resource.StepVars{"present": "ok"},
		}
		_, err := st.ExecuteTemplate(sysfs)
		if err == nil || !strings.Contains(err.Error(), `map has no entry for key "missing_key"`) {
			t.Errorf("ExecuteTemplate() err = %v, want missing_key error", err)
		}
	})

	t.Run("dot lookup on missing map key returns error via missingkey=error", func(t *testing.T) {
		t.Parallel()
		st := &resource.Step{
			Name:       "s",
			ConfigPath: "templates/missing_dot_key.tf.tmpl",
			Vars:       resource.StepVars{"present": "ok"},
		}
		_, err := st.ExecuteTemplate(sysfs)
		if err == nil || !strings.Contains(err.Error(), `map has no entry for key "missing_key"`) {
			t.Errorf("ExecuteTemplate() err = %v, want missing_key error", err)
		}
	})

	t.Run("index on nil map returns error", func(t *testing.T) {
		t.Parallel()
		st := &resource.Step{
			Name:       "s",
			ConfigPath: "templates/nil_map.tf.tmpl",
		}
		_, err := st.ExecuteTemplate(sysfs)
		if err == nil || !strings.Contains(err.Error(), `map has no entry for key "id"`) {
			t.Errorf("ExecuteTemplate() err = %v, want missing id error", err)
		}
	})
}

func TestStep_Validate_QuoteWrappedVars(t *testing.T) {
	t.Parallel()

	t.Run("rejects quote-wrapped values across vars, resource_id_vars, and oics_vars_overrides", func(t *testing.T) {
		t.Parallel()
		st := &resource.Step{
			Name: "sample_step",
			Vars: resource.StepVars{
				"double_quoted":    `"DELETE"`,
				"single_quoted":    `'DELETE'`,
				"empty_double":     `""`,
				"empty_single":     `''`,
				"whitespace_pad":   " \"DELETE\"\n",
				"escaped_inner":    `"foo\"bar"`,
				"escaped_bslash":   `"foo\\"`,
				"regex_like":       `"\d+"`,
				"valid_plain":      "DELETE",
				"valid_empty":      "",
				"valid_null":       "null",
				"valid_hcl_list":   `["-c", "echo hello"]`,
				"valid_hcl_map":    `{"key" = "value"}`,
				"valid_hcl_ref":    `data.google_container_engine_versions.versions.release_channel_latest_version["REGULAR"]`,
				"valid_json_array": `[\"Default Environment\"]`,
				"valid_pair":       `"a" = "b"`,
				"valid_single_chr": `"`,
				"valid_unclosed":   `"unterminated`,
				"valid_esc_end":    `"foo\"`,
			},
			ResourceIdVars: map[string]string{
				"quoted_id":  `"custom-nam11-config"`,
				"single_id":  `'custom-id'`,
				"valid_id":   "custom-nam11-config",
				"valid_cidr": "136.124.3.120/32",
			},
			TestVarsOverrides: map[string]string{
				"valid_id": `"custom-tf-test-nam11-config"`,
			},
			OicsVarsOverrides: map[string]string{
				"valid_id":    `"DELETE"`,
				"valid_plain": `'DELETE'`,
				"valid_cidr":  "DELETE",
			},
		}

		errs := st.Validate("InstanceConfig", "spanner_instance_config_basic")
		got := make([]string, len(errs))
		for i, err := range errs {
			got[i] = err.Error()
		}
		want := []string{
			`variable "double_quoted" in 'vars' has quote-wrapped value "\"DELETE\"" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "empty_double" in 'vars' has quote-wrapped value "\"\"" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "empty_single" in 'vars' has quote-wrapped value "''" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "escaped_bslash" in 'vars' has quote-wrapped value "\"foo\\\\\"" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "escaped_inner" in 'vars' has quote-wrapped value "\"foo\\\"bar\"" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "regex_like" in 'vars' has quote-wrapped value "\"\\d+\"" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "single_quoted" in 'vars' has quote-wrapped value "'DELETE'" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "whitespace_pad" in 'vars' has quote-wrapped value " \"DELETE\"\n" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "quoted_id" in 'resource_id_vars' has quote-wrapped value "\"custom-nam11-config\"" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "single_id" in 'resource_id_vars' has quote-wrapped value "'custom-id'" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "valid_id" in 'oics_vars_overrides' has quote-wrapped value "\"DELETE\"" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "valid_plain" in 'oics_vars_overrides' has quote-wrapped value "'DELETE'" for step "sample_step" in sample "spanner_instance_config_basic" of resource "InstanceConfig"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Validate() errors mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("yaml unmarshaling distinguishes YAML quotes from literal HCL quotes", func(t *testing.T) {
		t.Parallel()
		rawYAML := `
name: sample_step
resource_id_vars:
  yaml_single_quoted: 'my-id'
  yaml_double_quoted: "my-other-id"
  unquoted_id: my-plain-id
vars:
  empty_str: ""
  null_val: null
  yaml_quoted_str: "DELETE"
  literal_quoted_str: '"DELETE"'
test_vars_overrides:
  yaml_single_quoted: '"override-id"'
oics_vars_overrides:
  yaml_single_quoted: "DELETE"
  yaml_double_quoted: '"DELETE"'
`
		var st resource.Step
		if err := yaml.Unmarshal([]byte(rawYAML), &st); err != nil {
			t.Fatalf("yaml.Unmarshal failed: %v", err)
		}

		errs := st.Validate("Instance", "sample_name")
		got := make([]string, len(errs))
		for i, err := range errs {
			got[i] = err.Error()
		}
		want := []string{
			`variable "literal_quoted_str" in 'vars' has quote-wrapped value "\"DELETE\"" for step "sample_step" in sample "sample_name" of resource "Instance"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
			`variable "yaml_double_quoted" in 'oics_vars_overrides' has quote-wrapped value "\"DELETE\"" for step "sample_step" in sample "sample_name" of resource "Instance"; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes`,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Validate() errors mismatch (-want +got):\n%s", diff)
		}
	})
}
