package compute

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestComputeInstance_networkIPCustomizedDiff(t *testing.T) {
	t.Parallel()

	d := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"network_interface.#": 0,
		},
		After: map[string]interface{}{
			"network_interface.#": 1,
		},
	}

	err := forceNewIfNetworkIPNotUpdatableFunc(d)
	if err != nil {
		t.Error(err)
	}

	if d.IsForceNew {
		t.Errorf("Expected not force new if network_interface array size changes")
	}

	type NetworkInterface struct {
		Network           string
		Subnetwork        string
		SubnetworkProject string
		NetworkIP         string
	}
	NIBefore := NetworkInterface{
		Network:           "a",
		Subnetwork:        "a",
		SubnetworkProject: "a",
		NetworkIP:         "a",
	}

	cases := map[string]struct {
		ExpectedForceNew bool
		Before           NetworkInterface
		After            NetworkInterface
	}{
		"NetworkIP only change": {
			ExpectedForceNew: true,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and Network change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "b",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and Subnetwork change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "b",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and SubnetworkProject change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "b",
				NetworkIP:         "b",
			},
		},
		"All change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "b",
				Subnetwork:        "b",
				SubnetworkProject: "b",
				NetworkIP:         "b",
			},
		},
		"No change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "a",
			},
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"network_interface.#":                    1,
				"network_interface.0.network":            tc.Before.Network,
				"network_interface.0.subnetwork":         tc.Before.Subnetwork,
				"network_interface.0.subnetwork_project": tc.Before.SubnetworkProject,
				"network_interface.0.network_ip":         tc.Before.NetworkIP,
			},
			After: map[string]interface{}{
				"network_interface.#":                    1,
				"network_interface.0.network":            tc.After.Network,
				"network_interface.0.subnetwork":         tc.After.Subnetwork,
				"network_interface.0.subnetwork_project": tc.After.SubnetworkProject,
				"network_interface.0.network_ip":         tc.After.NetworkIP,
			},
		}
		err := forceNewIfNetworkIPNotUpdatableFunc(d)
		if err != nil {
			t.Error(err)
		}
		if tc.ExpectedForceNew != d.IsForceNew {
			t.Errorf("%v: expected d.IsForceNew to be %v, but was %v", tn, tc.ExpectedForceNew, d.IsForceNew)
		}
	}
}

func TestPerformanceMonitoringUnitDiffSuppress(t *testing.T) {
	t.Parallel()

	newResource := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
	existingResource := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
	existingResource.SetId("existing")

	cases := map[string]struct {
		old, new string
		data     *schema.ResourceData
		want     bool
	}{
		"new resource preserves standard":                       {old: "", new: "STANDARD", data: newResource, want: false},
		"existing resource suppresses empty to standard":        {old: "", new: "STANDARD", data: existingResource, want: true},
		"existing resource suppresses standard to empty":        {old: "STANDARD", new: "", data: existingResource, want: true},
		"existing resource preserves empty to architectural":    {old: "", new: "ARCHITECTURAL", data: existingResource, want: false},
		"existing resource preserves architectural to standard": {old: "ARCHITECTURAL", new: "STANDARD", data: existingResource, want: false},
		"nil resource does not suppress":                        {old: "", new: "STANDARD", data: nil, want: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := performanceMonitoringUnitDiffSuppress("performance_monitoring_unit", tc.old, tc.new, tc.data); got != tc.want {
				t.Fatalf("performanceMonitoringUnitDiffSuppress(%q, %q) = %t, want %t", tc.old, tc.new, got, tc.want)
			}
		})
	}
}

func TestComputeInstanceUpdate_internalIpv6PrefixLength(t *testing.T) {
	t.Parallel()

	var patchBodies []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "PATCH" && strings.Contains(r.URL.Path, "/updateNetworkInterface") {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("failed to decode PATCH body: %v", err)
			}
			patchBodies = append(patchBodies, body)
			w.Write([]byte(`{"name": "op-1", "status": "DONE", "zone": "https://www.googleapis.com/compute/v1/projects/my-proj/zones/us-central1-a"}`))
			return
		}
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/operations/") {
			w.Write([]byte(`{"name": "op-1", "status": "DONE", "zone": "https://www.googleapis.com/compute/v1/projects/my-proj/zones/us-central1-a"}`))
			return
		}
		if strings.Contains(r.URL.Path, "/stop") || strings.Contains(r.URL.Path, "/start") {
			w.Write([]byte(`{"name": "op-1", "status": "DONE", "zone": "https://www.googleapis.com/compute/v1/projects/my-proj/zones/us-central1-a"}`))
			return
		}
		w.Write([]byte(`{
			"name": "test-instance",
			"machineType": "https://www.googleapis.com/compute/v1/projects/my-proj/zones/us-central1-a/machineTypes/e2-medium",
			"status": "RUNNING",
			"zone": "https://www.googleapis.com/compute/v1/projects/my-proj/zones/us-central1-a",
			"cpuPlatform": "Intel Broadwell",
			"networkInterfaces": [{
				"name": "nic0",
				"network": "https://www.googleapis.com/compute/v1/projects/my-proj/global/networks/default",
				"subnetwork": "https://www.googleapis.com/compute/v1/projects/my-proj/regions/us-central1/subnetworks/default",
				"fingerprint": "abc123",
				"stackType": "IPV4_IPV6",
				"ipv6AccessType": "INTERNAL",
				"internalIpv6PrefixLength": 96
			}]
		}`))
	}))
	defer ts.Close()

	res := ResourceComputeInstance()
	state := &terraform.InstanceState{
		ID: "projects/my-proj/zones/us-central1-a/instances/test-instance",
		Attributes: map[string]string{
			"id":                                   "projects/my-proj/zones/us-central1-a/instances/test-instance",
			"name":                                 "test-instance",
			"project":                              "my-proj",
			"zone":                                 "us-central1-a",
			"machine_type":                         "e2-medium",
			"network_interface.#":                  "1",
			"network_interface.0.name":             "nic0",
			"network_interface.0.network":          "https://www.googleapis.com/compute/v1/projects/my-proj/global/networks/default",
			"network_interface.0.subnetwork":       "https://www.googleapis.com/compute/v1/projects/my-proj/regions/us-central1/subnetworks/default",
			"network_interface.0.stack_type":       "IPV4_IPV6",
			"network_interface.0.ipv6_access_type": "INTERNAL",
			"network_interface.0.internal_ipv6_prefix_length": "0",
		},
	}

	rawCfg := map[string]interface{}{
		"name":         "test-instance",
		"project":      "my-proj",
		"zone":         "us-central1-a",
		"machine_type": "e2-medium",
		"network_interface": []interface{}{
			map[string]interface{}{
				"network":                     "https://www.googleapis.com/compute/v1/projects/my-proj/global/networks/default",
				"subnetwork":                  "https://www.googleapis.com/compute/v1/projects/my-proj/regions/us-central1/subnetworks/default",
				"stack_type":                  "IPV4_IPV6",
				"ipv6_access_type":            "INTERNAL",
				"internal_ipv6_prefix_length": 96,
			},
		},
		"service_account": []interface{}{nil},
		"scheduling":      []interface{}{nil},
	}

	cfg := &transport_tpg.Config{
		Project: "my-proj",
		Zone:    "us-central1-a",
		Client:  ts.Client(),
		CustomEndpoints: map[string]string{
			"compute_custom_endpoint": ts.URL + "/",
		},
	}

	diff, err := schema.InternalMap(res.Schema).Diff(t.Context(), state, terraform.NewResourceConfigRaw(rawCfg), nil, cfg, false)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("Data failed: %v", err)
	}

	if err := resourceComputeInstanceUpdate(d, cfg); err != nil {
		t.Fatalf("resourceComputeInstanceUpdate failed: %v", err)
	}

	if len(patchBodies) != 1 {
		t.Fatalf("expected 1 updateNetworkInterface PATCH request, got %d", len(patchBodies))
	}
	if got := patchBodies[0]["internalIpv6PrefixLength"]; got != float64(96) {
		t.Errorf("expected internalIpv6PrefixLength=96 in PATCH body, got %v (%T)", got, got)
	}

	// Also verify updateDuringStop populates internalIpv6PrefixLength (not ipv6Address)
	patchBodies = nil
	rawCfgStop := map[string]interface{}{
		"name":                      "test-instance",
		"project":                   "my-proj",
		"zone":                      "us-central1-a",
		"machine_type":              "e2-medium",
		"allow_stopping_for_update": true,
		"network_interface": []interface{}{
			map[string]interface{}{
				"network":                     "https://www.googleapis.com/compute/v1/projects/my-proj/global/networks/other-net",
				"subnetwork":                  "https://www.googleapis.com/compute/v1/projects/my-proj/regions/us-central1/subnetworks/other-subnet",
				"stack_type":                  "IPV4_IPV6",
				"ipv6_access_type":            "INTERNAL",
				"internal_ipv6_prefix_length": 96,
			},
		},
	}

	diffStop, err := schema.InternalMap(res.Schema).Diff(t.Context(), state, terraform.NewResourceConfigRaw(rawCfgStop), nil, cfg, false)
	if err != nil {
		t.Fatalf("Diff (updateDuringStop) failed: %v", err)
	}

	dStop, err := schema.InternalMap(res.Schema).Data(state, diffStop)
	if err != nil {
		t.Fatalf("Data (updateDuringStop) failed: %v", err)
	}

	if err := resourceComputeInstanceUpdate(dStop, cfg); err != nil {
		t.Fatalf("resourceComputeInstanceUpdate (updateDuringStop) failed: %v", err)
	}

	if len(patchBodies) != 1 {
		t.Fatalf("expected 1 updateNetworkInterface PATCH request during stop, got %d", len(patchBodies))
	}
	if got := patchBodies[0]["internalIpv6PrefixLength"]; got != float64(96) {
		t.Errorf("expected internalIpv6PrefixLength=96 in stopped PATCH body, got %v (%T)", got, got)
	}
	if _, hasIpv6Addr := patchBodies[0]["ipv6Address"]; hasIpv6Addr {
		t.Errorf("did not expect ipv6Address in stopped PATCH body when only internal_ipv6_prefix_length changed: %v", patchBodies[0])
	}
}
