package compute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
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

func TestComputeInstance_networkInterfaceDynamicNicCustomizedDiff(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		before  map[string]interface{}
		after   map[string]interface{}
		wantNew bool
	}{
		"new resource creation does not force new": {
			before: map[string]interface{}{
				"network_interface.#": 0,
			},
			after: map[string]interface{}{
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       2,
			},
			wantNew: false,
		},
		"add single trailing dynamic VLAN NIC updates in place": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       2,
			},
			wantNew: false,
		},
		"add multiple trailing dynamic VLAN NICs updates in place": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            3,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       2,
				"network_interface.2.subnetwork": "sub3",
				"network_interface.2.vlan":       255,
			},
			wantNew: false,
		},
		"add trailing dynamic VLAN NIC with network_attachment updates in place": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
			},
			after: map[string]interface{}{
				"network_interface":                      []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":                    2,
				"network_interface.0.subnetwork":         "sub1",
				"network_interface.0.vlan":               0,
				"network_interface.1.network_attachment": "projects/p/regions/r/networkAttachments/na1",
				"network_interface.1.service_class_id":   "sc-1",
				"network_interface.1.vlan":               2,
			},
			wantNew: false,
		},
		"remove single trailing dynamic VLAN NIC updates in place": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       2,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
			},
			wantNew: false,
		},
		"remove multiple trailing dynamic VLAN NICs updates in place": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            3,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       2,
				"network_interface.2.subnetwork": "sub3",
				"network_interface.2.vlan":       3,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
			},
			wantNew: false,
		},
		"add non-dynamic NIC forces replacement": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       0,
			},
			wantNew: true,
		},
		"remove non-dynamic NIC forces replacement": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       0,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
			},
			wantNew: true,
		},
		"insert non-dynamic NIC before existing dynamic NIC forces replacement": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       2,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            3,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub3",
				"network_interface.1.vlan":       0,
				"network_interface.2.subnetwork": "sub2",
				"network_interface.2.vlan":       2,
			},
			wantNew: true,
		},
		"remove non-trailing dynamic NIC forces replacement": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            3,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       2,
				"network_interface.2.subnetwork": "sub3",
				"network_interface.2.vlan":       3,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub3",
				"network_interface.1.vlan":       3,
			},
			wantNew: true,
		},
		"change vlan on existing dynamic NIC forces replacement": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       2,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "sub2",
				"network_interface.1.vlan":       3,
			},
			wantNew: true,
		},
		"reorder non-dynamic NICs forces replacement": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "https://www.googleapis.com/compute/v1/projects/p/regions/r/subnetworks/sub1",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "https://www.googleapis.com/compute/v1/projects/p/regions/r/subnetworks/sub2",
				"network_interface.1.vlan":       0,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":            2,
				"network_interface.0.subnetwork": "projects/p/regions/r/subnetworks/sub2",
				"network_interface.0.vlan":       0,
				"network_interface.1.subnetwork": "projects/p/regions/r/subnetworks/sub1",
				"network_interface.1.vlan":       0,
			},
			wantNew: true,
		},
		"migrate single non-dynamic NIC to new subnetwork does not force new": {
			before: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "https://www.googleapis.com/compute/v1/projects/p/regions/r/subnetworks/sub1",
				"network_interface.0.vlan":       0,
			},
			after: map[string]interface{}{
				"network_interface":              []interface{}{map[string]interface{}{}},
				"network_interface.#":            1,
				"network_interface.0.subnetwork": "projects/p/regions/r/subnetworks/sub2",
				"network_interface.0.vlan":       0,
			},
			wantNew: false,
		},
		"change network_attachment on existing NIC forces replacement": {
			before: map[string]interface{}{
				"network_interface":                      []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":                    2,
				"network_interface.0.network":            "default",
				"network_interface.0.vlan":               0,
				"network_interface.1.network_attachment": "projects/p/regions/r/networkAttachments/na1",
				"network_interface.1.vlan":               0,
			},
			after: map[string]interface{}{
				"network_interface":                      []interface{}{map[string]interface{}{}, map[string]interface{}{}},
				"network_interface.#":                    2,
				"network_interface.0.network":            "default",
				"network_interface.0.vlan":               0,
				"network_interface.1.network_attachment": "projects/p/regions/r/networkAttachments/na2",
				"network_interface.1.vlan":               0,
			},
			wantNew: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			d := &tpgresource.ResourceDiffMock{
				Before: tc.before,
				After:  tc.after,
			}
			if err := forceNewIfNetworkInterfaceNotUpdatableFunc(d); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if d.IsForceNew != tc.wantNew {
				t.Fatalf("got IsForceNew = %v, want %v", d.IsForceNew, tc.wantNew)
			}
		})
	}
}

func TestResolveDynamicNicParentName(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		instNics    []interface{}
		targetIndex int
		want        string
		wantErr     bool
	}{
		"single parent nic0": {
			instNics: []interface{}{
				map[string]interface{}{"name": "nic0"},
			},
			targetIndex: 1,
			want:        "nic0",
		},
		"multiple dynamic NICs added onto nic0 in one apply": {
			instNics: []interface{}{
				map[string]interface{}{"name": "nic0"},
			},
			targetIndex: 2,
			want:        "nic0",
		},
		"existing dynamic NIC on nic0": {
			instNics: []interface{}{
				map[string]interface{}{"name": "nic0"},
				map[string]interface{}{"name": "nic0.2", "parentNicName": "nic0", "vlan": float64(2)},
			},
			targetIndex: 2,
			want:        "nic0",
		},
		"two regular NICs resolves to nic1": {
			instNics: []interface{}{
				map[string]interface{}{"name": "nic0"},
				map[string]interface{}{"name": "nic0.2", "parentNicName": "nic0", "vlan": float64(2)},
				map[string]interface{}{"name": "nic1"},
			},
			targetIndex: 3,
			want:        "nic1",
		},
		"empty instNics errors": {
			instNics:    []interface{}{},
			targetIndex: 1,
			wantErr:     true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := resolveDynamicNicParentName(tc.instNics, tc.targetIndex)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("resolveDynamicNicParentName() = %q, want %q", got, tc.want)
			}
		})
	}
}
