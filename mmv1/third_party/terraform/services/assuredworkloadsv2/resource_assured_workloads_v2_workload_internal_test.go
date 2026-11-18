package assuredworkloadsv2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// configuredPaths returns a function reporting whether a path is one of paths,
// standing in for the raw configuration.
func configuredPaths(paths ...string) func(string) bool {
	return func(p string) bool {
		for _, c := range paths {
			if p == c {
				return true
			}
		}
		return false
	}
}

func paramValuePath(control, param int) string {
	return fmt.Sprintf("cloud_control_configs.%d.parameters.%d.parameter_value.0", control, param)
}

// emptyParamValue is a parameter_value block as Terraform reads it: unset
// attributes hold their zero values.
func emptyParamValue() map[string]interface{} {
	return map[string]interface{}{
		"string_value":      "",
		"bool_value":        false,
		"number_value":      0.0,
		"string_list_value": []interface{}{},
		"oneof_value":       []interface{}{},
	}
}

func withValue(key string, value interface{}) map[string]interface{} {
	m := emptyParamValue()
	m[key] = value
	return m
}

func TestAssuredWorkloadsV2ExpandParamValue(t *testing.T) {
	t.Parallel()

	base := paramValuePath(0, 0)
	nested := emptyParamValue()
	delete(nested, "oneof_value")
	nested["bool_value"] = false

	cases := []struct {
		name       string
		raw        map[string]interface{}
		configured func(string) bool
		want       map[string]interface{}
		wantErr    bool
	}{
		{
			name:       "string",
			raw:        withValue("string_value", "US_REGIONAL_DATA_BOUNDARY"),
			configured: configuredPaths(base + ".string_value"),
			want:       map[string]interface{}{"stringValue": "US_REGIONAL_DATA_BOUNDARY"},
		},
		{
			name:       "configured empty string",
			raw:        emptyParamValue(),
			configured: configuredPaths(base + ".string_value"),
			want:       map[string]interface{}{"stringValue": ""},
		},
		{
			name:       "bool true",
			raw:        withValue("bool_value", true),
			configured: configuredPaths(),
			want:       map[string]interface{}{"boolValue": true},
		},
		{
			name:       "configured bool false",
			raw:        emptyParamValue(),
			configured: configuredPaths(base + ".bool_value"),
			want:       map[string]interface{}{"boolValue": false},
		},
		{
			name:       "number",
			raw:        withValue("number_value", 3.5),
			configured: configuredPaths(),
			want:       map[string]interface{}{"numberValue": 3.5},
		},
		{
			name:       "configured number zero",
			raw:        emptyParamValue(),
			configured: configuredPaths(base + ".number_value"),
			want:       map[string]interface{}{"numberValue": 0.0},
		},
		{
			name:       "string list",
			raw:        withValue("string_list_value", []interface{}{map[string]interface{}{"values": []interface{}{"TLS_VERSION_1"}}}),
			configured: configuredPaths(),
			want:       map[string]interface{}{"stringListValue": map[string]interface{}{"values": []interface{}{"TLS_VERSION_1"}}},
		},
		{
			name: "oneof with configured nested bool false",
			raw: withValue("oneof_value", []interface{}{map[string]interface{}{
				"name":            "nested",
				"parameter_value": []interface{}{nested},
			}}),
			configured: configuredPaths(base + ".oneof_value.0.parameter_value.0.bool_value"),
			want: map[string]interface{}{"oneofValue": map[string]interface{}{
				"name":           "nested",
				"parameterValue": map[string]interface{}{"boolValue": false},
			}},
		},
		{
			name:       "nothing set",
			raw:        emptyParamValue(),
			configured: configuredPaths(),
			want:       map[string]interface{}{},
		},
		{
			name: "multiple values",
			raw: func() map[string]interface{} {
				m := withValue("string_value", "a")
				m["bool_value"] = true
				return m
			}(),
			configured: configuredPaths(),
			wantErr:    true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := assuredWorkloadsV2ExpandParamValue(tc.raw, base, tc.configured, true)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %#v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestAssuredWorkloadsV2ExpandCloudControlConfigs(t *testing.T) {
	t.Parallel()

	v := []interface{}{
		map[string]interface{}{
			"cloud_control":          "organizations/123/locations/us-east4/cloudControls/builtin-enable-axt",
			"cloud_control_revision": 1,
			"enforcement_mode":       "PREVENTIVE",
			"parameters":             []interface{}{},
		},
		map[string]interface{}{
			"cloud_control":          "organizations/123/locations/us-east4/cloudControls/example",
			"cloud_control_revision": 2,
			"enforcement_mode":       "DETECTIVE",
			"parameters": []interface{}{
				map[string]interface{}{"name": "a", "parameter_value": []interface{}{withValue("string_value", "x")}},
				map[string]interface{}{"name": "b", "parameter_value": []interface{}{emptyParamValue()}},
			},
		},
	}
	got, err := assuredWorkloadsV2ExpandCloudControlConfigs(v, configuredPaths(paramValuePath(1, 1)+".bool_value"))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := []interface{}{
		map[string]interface{}{
			"cloudControl":         "organizations/123/locations/us-east4/cloudControls/builtin-enable-axt",
			"cloudControlRevision": 1,
			"enforcementMode":      "PREVENTIVE",
		},
		map[string]interface{}{
			"cloudControl":         "organizations/123/locations/us-east4/cloudControls/example",
			"cloudControlRevision": 2,
			"enforcementMode":      "DETECTIVE",
			"parameters": []interface{}{
				map[string]interface{}{"name": "a", "parameterValue": map[string]interface{}{"stringValue": "x"}},
				map[string]interface{}{"name": "b", "parameterValue": map[string]interface{}{"boolValue": false}},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestAssuredWorkloadsV2ExpandCloudControlConfigsError(t *testing.T) {
	t.Parallel()

	m := withValue("string_value", "x")
	m["number_value"] = 1.0
	v := []interface{}{map[string]interface{}{
		"cloud_control":          "organizations/123/locations/us-east4/cloudControls/example",
		"cloud_control_revision": 1,
		"enforcement_mode":       "DETECTIVE",
		"parameters":             []interface{}{map[string]interface{}{"name": "p", "parameter_value": []interface{}{m}}},
	}}
	_, err := assuredWorkloadsV2ExpandCloudControlConfigs(v, configuredPaths())
	if err == nil || !strings.Contains(err.Error(), `parameter "p"`) {
		t.Fatalf("expected an error naming parameter p, got %v", err)
	}
}

func TestAssuredWorkloadsV2ExpandBillingAccount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "prefixed billing account",
			input: "billingAccounts/012345-6789AB-CDEF01",
			want:  "012345-6789AB-CDEF01",
		},
		{
			name:  "bare billing account",
			input: "012345-6789AB-CDEF01",
			want:  "012345-6789AB-CDEF01",
		},
		{
			name:  "nil value",
			input: nil,
			want:  nil,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := expandAssuredWorkloadsV2WorkloadResourceConfigNewProjectConfigBillingAccount(tc.input, nil, nil)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestAssuredWorkloadsV2BillingAccountDiffSuppress(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{
			name: "prefixed and bare match",
			old:  "billingAccounts/012345-6789AB-CDEF01",
			new:  "012345-6789AB-CDEF01",
			want: true,
		},
		{
			name: "bare and prefixed match",
			old:  "012345-6789AB-CDEF01",
			new:  "billingAccounts/012345-6789AB-CDEF01",
			want: true,
		},
		{
			name: "different billing accounts",
			old:  "billingAccounts/012345-6789AB-CDEF01",
			new:  "billingAccounts/999999-6789AB-CDEF01",
			want: false,
		},
		{
			name: "empty old and prefixed new",
			old:  "",
			new:  "billingAccounts/012345-6789AB-CDEF01",
			want: false,
		},
		{
			name: "empty old and prefix-only new",
			old:  "",
			new:  "billingAccounts/",
			want: false,
		},
		{
			name: "both empty",
			old:  "",
			new:  "",
			want: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := assuredWorkloadsV2BillingAccountDiffSuppress("billing_account", tc.old, tc.new, nil); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAssuredWorkloadsV2WorkloadIdFromName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  string
	}{
		{"organizations/123/locations/us-east4/workloads/my-workload", "my-workload"},
		{"organizations/123/locations/europe/workloads/uuid-1234-5678", "uuid-1234-5678"},
		{"invalid/name", ""},
		{"", ""},
	}
	for _, tc := range cases {
		if got := assuredWorkloadsV2WorkloadIdFromName(tc.input); got != tc.want {
			t.Errorf("assuredWorkloadsV2WorkloadIdFromName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestAssuredWorkloadsV2WorkloadImport(t *testing.T) {
	t.Parallel()

	cases := []string{
		"organizations/123/locations/us-east4/workloads/my-workload",
		"123/us-east4/my-workload",
	}
	for _, importId := range cases {
		importId := importId
		t.Run(importId, func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceAssuredWorkloadsV2Workload().Schema, map[string]interface{}{})
			d.SetId(importId)
			res, err := resourceAssuredWorkloadsV2WorkloadImport(d, &transport_tpg.Config{})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if len(res) != 1 {
				t.Fatalf("got %d resources, want 1", len(res))
			}
			want := map[string]string{
				"organization": "123",
				"location":     "us-east4",
				"workload_id":  "my-workload",
				"name":         "organizations/123/locations/us-east4/workloads/my-workload",
			}
			for k, v := range want {
				if got := res[0].Get(k).(string); got != v {
					t.Errorf("%s = %q, want %q", k, got, v)
				}
			}
			if got := res[0].Id(); got != want["name"] {
				t.Errorf("id = %q, want %q", got, want["name"])
			}
		})
	}

	d := schema.TestResourceDataRaw(t, ResourceAssuredWorkloadsV2Workload().Schema, map[string]interface{}{})
	d.SetId("not/a/valid/workload/id/format")
	if _, err := resourceAssuredWorkloadsV2WorkloadImport(d, &transport_tpg.Config{}); err == nil {
		t.Error("expected an error for an invalid import ID")
	}
}

func TestDataSourceAssuredWorkloadsV2WorkloadSchema(t *testing.T) {
	t.Parallel()

	ds := DataSourceAssuredWorkloadsV2Workload()
	if err := ds.InternalValidate(nil, false); err != nil {
		t.Fatalf("data source schema InternalValidate failed: %s", err)
	}
	if _, exists := ds.Schema["resource_config"]; exists {
		t.Errorf("expected %q to be excluded from data source schema", "resource_config")
	}
	if _, exists := ds.Schema["deletion_policy"]; !exists {
		t.Errorf("expected %q to remain in data source schema so resource read can populate virtual fields", "deletion_policy")
	}

	cases := []struct {
		name    string
		raw     map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid name",
			raw: map[string]interface{}{
				"name": "organizations/123/locations/us-east4/workloads/my-workload",
			},
		},
		{
			name: "valid org location workload_id",
			raw: map[string]interface{}{
				"organization": "123",
				"location":     "us-east4",
				"workload_id":  "my-workload",
			},
		},
		{
			name:    "neither specified",
			raw:     map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "invalid name format",
			raw: map[string]interface{}{
				"name": "invalid-name",
			},
			wantErr: true,
		},
		{
			name: "name conflicts with organization",
			raw: map[string]interface{}{
				"name":         "organizations/123/locations/us-east4/workloads/my-workload",
				"organization": "123",
			},
			wantErr: true,
		},
		{
			name: "incomplete components",
			raw: map[string]interface{}{
				"organization": "123",
				"workload_id":  "my-workload",
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			diags := ds.Validate(terraform.NewResourceConfigRaw(tc.raw))
			if tc.wantErr && !diags.HasError() {
				t.Errorf("expected validation error for %#v, got none", tc.raw)
			}
			if !tc.wantErr && diags.HasError() {
				t.Errorf("unexpected validation errors for %#v: %v", tc.raw, diags)
			}
		})
	}
}

func TestDataSourceAssuredWorkloadsV2WorkloadResolveId(t *testing.T) {
	t.Parallel()

	dsSchema := DataSourceAssuredWorkloadsV2Workload().Schema
	wantName := "organizations/123/locations/us-east4/workloads/my-workload"

	t.Run("by name populates regional URL components", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, dsSchema, map[string]interface{}{
			"name": wantName,
		})
		id, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if id != wantName {
			t.Errorf("id = %q, want %q", id, wantName)
		}
		for k, want := range map[string]string{
			"organization": "123",
			"location":     "us-east4",
			"workload_id":  "my-workload",
			"name":         wantName,
		} {
			if got := d.Get(k).(string); got != want {
				t.Errorf("%s = %q, want %q", k, got, want)
			}
		}
	})

	t.Run("by components populates name", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, dsSchema, map[string]interface{}{
			"organization": "123",
			"location":     "us-east4",
			"workload_id":  "my-workload",
		})
		id, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if id != wantName {
			t.Errorf("id = %q, want %q", id, wantName)
		}
		if got := d.Get("name").(string); got != wantName {
			t.Errorf("name = %q, want %q", got, wantName)
		}
	})

	t.Run("invalid name returns error", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, dsSchema, map[string]interface{}{
			"name": "not-a-valid-name",
		})
		if _, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d); err == nil {
			t.Error("expected an error for invalid name")
		}
	})

	t.Run("missing components returns error", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, dsSchema, map[string]interface{}{
			"organization": "123",
		})
		if _, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d); err == nil {
			t.Error("expected an error for missing components")
		}
	})
}

func sampleAssuredWorkloadsV2WorkloadResponse(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":                      name,
		"description":               "Sample workload",
		"computedTargetResource":    "projects/123456789",
		"targetResourceDisplayName": "sample-project",
		"kajEnrollmentType":         "KEY_ACCESS_TRANSPARENCY_DISABLED",
		"state":                     "ACTIVE",
		"etag":                      "etag-1",
		"createTime":                "2026-11-18T00:00:00Z",
		"updateTime":                "2026-11-18T01:00:00Z",
		"framework": map[string]interface{}{
			"framework":       "organizations/123/locations/us-east4/frameworks/builtin-us-data-boundary",
			"majorRevisionId": "62",
		},
		"cloudControlConfigs": []interface{}{
			map[string]interface{}{
				"cloudControl":         "organizations/123/locations/us-east4/cloudControls/builtin-apigee-regulatory-settings",
				"cloudControlRevision": "2",
				"enforcementMode":      "DETECTIVE",
				"parameters": []interface{}{
					map[string]interface{}{
						"name": "framework",
						"parameterValue": map[string]interface{}{
							"stringValue": "US_REGIONAL_DATA_BOUNDARY",
						},
					},
				},
			},
		},
	}
}

func TestDataSourceAssuredWorkloadsV2WorkloadRead(t *testing.T) {
	t.Parallel()

	workloadName := "organizations/123/locations/us-east4/workloads/my-workload"
	var lastUserProject string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		lastUserProject = req.Header.Get("X-Goog-User-Project")
		if req.Method == "GET" && req.URL.Path == "/us-east4/v2/"+workloadName {
			_ = json.NewEncoder(w).Encode(sampleAssuredWorkloadsV2WorkloadResponse(workloadName))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    404,
				"message": "not found",
			},
		})
	}))
	defer srv.Close()

	cfg := &transport_tpg.Config{
		Client:              srv.Client(),
		UserProjectOverride: true,
		BillingProject:      "billing-proj",
		CustomEndpoints: map[string]string{
			Product.CustomEndpointField: srv.URL + "/{{location}}/",
		},
	}

	ds := DataSourceAssuredWorkloadsV2Workload()

	checkState := func(t *testing.T, d *schema.ResourceData) {
		t.Helper()
		if lastUserProject != "billing-proj" {
			t.Errorf("X-Goog-User-Project = %q, want %q", lastUserProject, "billing-proj")
		}
		if d.Id() != workloadName {
			t.Errorf("Id() = %q, want %q", d.Id(), workloadName)
		}
		for k, want := range map[string]interface{}{
			"name":                         workloadName,
			"organization":                 "123",
			"location":                     "us-east4",
			"workload_id":                  "my-workload",
			"description":                  "Sample workload",
			"computed_target_resource":     "projects/123456789",
			"target_resource_display_name": "sample-project",
			"kaj_enrollment_type":          "KEY_ACCESS_TRANSPARENCY_DISABLED",
			"state":                        "ACTIVE",
			"etag":                         "etag-1",
			"deletion_policy":              "DELETE",
		} {
			if got := d.Get(k); got != want {
				t.Errorf("%s = %#v, want %#v", k, got, want)
			}
		}
		if got := d.Get("framework.0.major_revision_id"); got != 62 {
			t.Errorf("framework.0.major_revision_id = %#v, want 62", got)
		}
		if got := d.Get("cloud_control_configs.0.parameters.0.parameter_value.0.string_value"); got != "US_REGIONAL_DATA_BOUNDARY" {
			t.Errorf("string_value = %#v, want US_REGIONAL_DATA_BOUNDARY", got)
		}
	}

	t.Run("lookup by name", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"name": workloadName,
		})
		if err := dataSourceAssuredWorkloadsV2WorkloadRead(d, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkState(t, d)
	})

	t.Run("lookup by components", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"organization": "123",
			"location":     "us-east4",
			"workload_id":  "my-workload",
		})
		if err := dataSourceAssuredWorkloadsV2WorkloadRead(d, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkState(t, d)
	})

	t.Run("not found returns error", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"name": "organizations/123/locations/us-east4/workloads/missing",
		})
		err := dataSourceAssuredWorkloadsV2WorkloadRead(d, cfg)
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("expected not found error, got %v", err)
		}
	})
}

func TestResourceAssuredWorkloadsV2WorkloadLifecycle(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		inputWorkloadID    string
		wantQueryHasID     bool
		wantQueryID        string
		serverWorkloadName string
		wantWorkloadID     string
	}{
		{
			name:               "server-generated workload_id omits query param and populates workload_id",
			inputWorkloadID:    "",
			wantQueryHasID:     false,
			wantQueryID:        "",
			serverWorkloadName: "organizations/123/locations/us-east4/workloads/server-gen-uuid",
			wantWorkloadID:     "server-gen-uuid",
		},
		{
			name:               "user-assigned workload_id sends query param and strips billing_account prefix",
			inputWorkloadID:    "user-assigned-id",
			wantQueryHasID:     true,
			wantQueryID:        "user-assigned-id",
			serverWorkloadName: "organizations/123/locations/us-east4/workloads/user-assigned-id",
			wantWorkloadID:     "user-assigned-id",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				createQueryHasID    bool
				createQueryID       string
				createBody          map[string]interface{}
				createOpPolled      bool
				createOpUserProject string
				deleteCalled        bool
				deleteOpPolled      bool
				deleteOpUserProject string
			)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case req.Method == "POST" && req.URL.Path == "/us-east4/v2/organizations/123/locations/us-east4/workloads":
					createQueryHasID = req.URL.Query().Has("workloadId")
					createQueryID = req.URL.Query().Get("workloadId")
					if err := json.NewDecoder(req.Body).Decode(&createBody); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"name": "organizations/123/locations/us-east4/operations/create-op-1",
						"done": false,
					})
				case req.Method == "GET" && req.URL.Path == "/us-east4/v1/organizations/123/locations/us-east4/operations/create-op-1":
					createOpPolled = true
					createOpUserProject = req.Header.Get("X-Goog-User-Project")
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"name":     "organizations/123/locations/us-east4/operations/create-op-1",
						"done":     true,
						"response": sampleAssuredWorkloadsV2WorkloadResponse(tc.serverWorkloadName),
					})
				case req.Method == "GET" && req.URL.Path == "/us-east4/v2/"+tc.serverWorkloadName:
					_ = json.NewEncoder(w).Encode(sampleAssuredWorkloadsV2WorkloadResponse(tc.serverWorkloadName))
				case req.Method == "DELETE" && req.URL.Path == "/us-east4/v2/"+tc.serverWorkloadName:
					deleteCalled = true
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"name": "organizations/123/locations/us-east4/operations/delete-op-1",
						"done": false,
					})
				case req.Method == "GET" && req.URL.Path == "/us-east4/v1/organizations/123/locations/us-east4/operations/delete-op-1":
					deleteOpPolled = true
					deleteOpUserProject = req.Header.Get("X-Goog-User-Project")
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"name": "organizations/123/locations/us-east4/operations/delete-op-1",
						"done": true,
					})
				default:
					http.NotFound(w, req)
				}
			}))
			defer srv.Close()

			cfg := &transport_tpg.Config{
				Client:              srv.Client(),
				PollInterval:        time.Millisecond,
				UserProjectOverride: true,
				BillingProject:      "billing-proj",
				CustomEndpoints: map[string]string{
					Product.CustomEndpointField: srv.URL + "/{{location}}/",
				},
			}

			r := ResourceAssuredWorkloadsV2Workload()
			d := r.Data(nil)
			if err := d.Set("organization", "123"); err != nil {
				t.Fatal(err)
			}
			if err := d.Set("location", "us-east4"); err != nil {
				t.Fatal(err)
			}
			if tc.inputWorkloadID != "" {
				if err := d.Set("workload_id", tc.inputWorkloadID); err != nil {
					t.Fatal(err)
				}
			}
			if err := d.Set("resource_config", []interface{}{
				map[string]interface{}{
					"new_project_config": []interface{}{
						map[string]interface{}{
							"parent":          "folders/999",
							"display_name":    "sample-project",
							"billing_account": "billingAccounts/012345-6789AB-CDEF01",
						},
					},
				},
			}); err != nil {
				t.Fatal(err)
			}
			if err := d.Set("framework", []interface{}{
				map[string]interface{}{
					"framework":         "organizations/123/locations/us-east4/frameworks/builtin-us-data-boundary",
					"major_revision_id": 62,
				},
			}); err != nil {
				t.Fatal(err)
			}
			if err := d.Set("cloud_control_configs", []interface{}{
				map[string]interface{}{
					"cloud_control":          "organizations/123/locations/us-east4/cloudControls/builtin-apigee-regulatory-settings",
					"cloud_control_revision": 2,
					"enforcement_mode":       "DETECTIVE",
				},
			}); err != nil {
				t.Fatal(err)
			}

			if err := resourceAssuredWorkloadsV2WorkloadCreate(d, cfg); err != nil {
				t.Fatalf("Create failed: %v", err)
			}
			if createQueryHasID != tc.wantQueryHasID || createQueryID != tc.wantQueryID {
				t.Errorf("workloadId query param = (present=%v, val=%q), want (present=%v, val=%q)", createQueryHasID, createQueryID, tc.wantQueryHasID, tc.wantQueryID)
			}
			if !createOpPolled {
				t.Error("expected Create LRO to be polled")
			}
			if createOpUserProject != "billing-proj" {
				t.Errorf("Create LRO X-Goog-User-Project = %q, want %q", createOpUserProject, "billing-proj")
			}
			resCfg, _ := createBody["resourceConfig"].(map[string]interface{})
			newProjCfg, _ := resCfg["newProjectConfig"].(map[string]interface{})
			if got := newProjCfg["billingAccount"]; got != "012345-6789AB-CDEF01" {
				t.Errorf("billingAccount sent = %#v, want %q", got, "012345-6789AB-CDEF01")
			}
			if got := d.Get("workload_id").(string); got != tc.wantWorkloadID {
				t.Errorf("workload_id = %q, want %q", got, tc.wantWorkloadID)
			}
			if got := d.Id(); got != tc.serverWorkloadName {
				t.Errorf("Id() = %q, want %q", got, tc.serverWorkloadName)
			}

			if err := resourceAssuredWorkloadsV2WorkloadDelete(d, cfg); err != nil {
				t.Fatalf("Delete failed: %v", err)
			}
			if !deleteCalled || !deleteOpPolled {
				t.Errorf("expected Delete and Delete LRO poll, got deleteCalled=%v deleteOpPolled=%v", deleteCalled, deleteOpPolled)
			}
			if deleteOpUserProject != "billing-proj" {
				t.Errorf("Delete LRO X-Goog-User-Project = %q, want %q", deleteOpUserProject, "billing-proj")
			}
		})
	}
}
