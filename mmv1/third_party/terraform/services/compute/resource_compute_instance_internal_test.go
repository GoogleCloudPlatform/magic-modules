package compute

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
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

func buildInstanceFromTemplateGracefulShutdownCty(res *schema.Resource, id string, hasMaxDuration bool, seconds, nanos cty.Value) cty.Value {
	rootType := res.CoreConfigSchema().ImpliedType()
	attrs := make(map[string]cty.Value, len(rootType.AttributeTypes()))
	for k, ty := range rootType.AttributeTypes() {
		attrs[k] = cty.NullVal(ty)
	}
	if id != "" {
		attrs["id"] = cty.StringVal(id)
	}
	attrs["name"] = cty.StringVal("tf-test-vm")
	attrs["zone"] = cty.StringVal("us-central1-a")
	attrs["source_instance_template"] = cty.StringVal("projects/my-project/global/instanceTemplates/tpl")

	schedElemType := rootType.AttributeType("scheduling").ElementType()
	schedAttrs := make(map[string]cty.Value, len(schedElemType.AttributeTypes()))
	for k, ty := range schedElemType.AttributeTypes() {
		schedAttrs[k] = cty.NullVal(ty)
	}

	gsElemType := schedElemType.AttributeType("graceful_shutdown").ElementType()
	mdElemType := gsElemType.AttributeType("max_duration").ElementType()
	var mdList cty.Value
	if hasMaxDuration {
		mdVal := cty.ObjectVal(map[string]cty.Value{
			"seconds": seconds,
			"nanos":   nanos,
		})
		mdList = cty.ListVal([]cty.Value{mdVal})
	} else {
		mdList = cty.ListValEmpty(mdElemType)
	}
	gsVal := cty.ObjectVal(map[string]cty.Value{
		"enabled":      cty.False,
		"max_duration": mdList,
	})
	schedAttrs["graceful_shutdown"] = cty.ListVal([]cty.Value{gsVal})
	attrs["scheduling"] = cty.ListVal([]cty.Value{cty.ObjectVal(schedAttrs)})
	return cty.ObjectVal(attrs)
}

func TestComputeInstanceFromTemplate_gracefulShutdownCustomizeDiff(t *testing.T) {
	t.Parallel()

	res := ResourceComputeInstanceFromTemplate()
	schemaBlock := res.CoreConfigSchema()
	meta := &transport_tpg.Config{Project: "my-project", Zone: "us-central1-a"}

	priorVal := buildInstanceFromTemplateGracefulShutdownCty(
		res,
		"projects/my-project/zones/us-central1-a/instances/tf-test-vm",
		true,
		cty.NumberIntVal(120),
		cty.NumberIntVal(1000),
	)

	cases := map[string]struct {
		proposedHasMD   bool
		configHasMD     bool
		proposedNanos   cty.Value
		configNanos     cty.Value
		checkKey        string
		wantChange      bool
		wantOld         string
		wantNew         string
		wantNewComputed bool
	}{
		"standalone nanos removal zeroes computed nanos": {
			proposedHasMD: true,
			configHasMD:   true,
			proposedNanos: cty.NumberIntVal(1000),
			configNanos:   cty.NullVal(cty.Number),
			checkKey:      "scheduling.0.graceful_shutdown.0.max_duration.0.nanos",
			wantChange:    true,
			wantOld:       "1000",
			wantNew:       "0",
		},
		"in-place max_duration removal clears max_duration list": {
			proposedHasMD: true,
			configHasMD:   false,
			proposedNanos: cty.NumberIntVal(0),
			configNanos:   cty.NullVal(cty.Number),
			checkKey:      "scheduling.0.graceful_shutdown.0.max_duration.#",
			wantChange:    true,
			wantOld:       "1",
			wantNew:       "0",
		},
		"retained nanos produces no diff": {
			proposedHasMD: true,
			configHasMD:   true,
			proposedNanos: cty.NumberIntVal(1000),
			configNanos:   cty.NumberIntVal(1000),
			checkKey:      "scheduling.0.graceful_shutdown.0.max_duration.0.nanos",
			wantChange:    false,
		},
		"unknown nanos in config stays computed": {
			proposedHasMD:   true,
			configHasMD:     true,
			proposedNanos:   cty.UnknownVal(cty.Number),
			configNanos:     cty.UnknownVal(cty.Number),
			checkKey:        "scheduling.0.graceful_shutdown.0.max_duration.0.nanos",
			wantChange:      true,
			wantNewComputed: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			priorState, err := res.ShimInstanceStateFromValue(priorVal)
			if err != nil {
				t.Fatalf("ShimInstanceStateFromValue: %v", err)
			}
			proposedVal := buildInstanceFromTemplateGracefulShutdownCty(res, priorState.ID, tc.proposedHasMD, cty.NumberIntVal(120), tc.proposedNanos)
			configVal := buildInstanceFromTemplateGracefulShutdownCty(res, "", tc.configHasMD, cty.NumberIntVal(120), tc.configNanos)
			priorState.RawState = priorVal
			priorState.RawPlan = proposedVal
			priorState.RawConfig = configVal

			cfg := terraform.NewResourceConfigShimmed(proposedVal, schemaBlock)
			diff, err := res.SimpleDiff(context.Background(), priorState, cfg, meta)
			if err != nil {
				t.Fatalf("SimpleDiff error: %v", err)
			}

			var attrDiff *terraform.ResourceAttrDiff
			if diff != nil {
				attrDiff = diff.Attributes[tc.checkKey]
			}
			if !tc.wantChange {
				if attrDiff != nil {
					t.Fatalf("unexpected diff on %s: %+v", tc.checkKey, attrDiff)
				}
				return
			}
			if attrDiff == nil {
				t.Fatalf("expected diff on %s, got nil (full diff: %+v)", tc.checkKey, diff)
			}
			if tc.wantNewComputed {
				if !attrDiff.NewComputed {
					t.Fatalf("expected %s NewComputed=true, got %+v", tc.checkKey, attrDiff)
				}
				return
			}
			if attrDiff.Old != tc.wantOld || attrDiff.New != tc.wantNew {
				t.Fatalf("expected %s diff %s -> %s, got %s -> %s", tc.checkKey, tc.wantOld, tc.wantNew, attrDiff.Old, attrDiff.New)
			}
		})
	}
}
