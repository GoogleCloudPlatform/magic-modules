package lustre

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	"testing"
)

func TestLustreInstanceTargetVersionDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		New              string
		AvailableVersion string
		EffectiveVersion string
		ShouldSuppress   bool
	}{
		"latest with nothing available is a no-op": {
			New:              "latest",
			AvailableVersion: "",
			EffectiveVersion: "2.0.0",
			ShouldSuppress:   true,
		},
		"latest is case insensitive": {
			New:              "LATEST",
			AvailableVersion: "",
			EffectiveVersion: "2.0.0",
			ShouldSuppress:   true,
		},
		"latest with a newer version available should upgrade": {
			New:              "latest",
			AvailableVersion: "2.1.0",
			EffectiveVersion: "2.0.0",
			ShouldSuppress:   false,
		},
		"re-applying the running version is a no-op": {
			New:              "2.0.0",
			AvailableVersion: "",
			EffectiveVersion: "2.0.0",
			ShouldSuppress:   true,
		},
		"an explicit newer version should upgrade": {
			New:              "2.1.0",
			AvailableVersion: "2.1.0",
			EffectiveVersion: "2.0.0",
			ShouldSuppress:   false,
		},
		"an explicit older version is a downgrade and is ignored": {
			New:              "1.9.0",
			AvailableVersion: "2.1.0",
			EffectiveVersion: "2.0.0",
			ShouldSuppress:   true,
		},
		"an explicit version with no known running version should not be suppressed": {
			New:              "2.1.0",
			AvailableVersion: "",
			EffectiveVersion: "",
			ShouldSuppress:   false,
		},
	}
	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"available_version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"effective_version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"target_version": {
					Type:     schema.TypeString,
					Optional: true,
				},
			}, map[string]interface{}{
				"available_version": tc.AvailableVersion,
				"effective_version": tc.EffectiveVersion,
			})
			// old is always "" in practice: the API clears target_version once the
			// upgrade completes, so the previous state value carries no information.
			if got := lustreInstanceTargetVersionDiffSuppress("target_version", "", tc.New, d); got != tc.ShouldSuppress {
				t.Errorf("lustreInstanceTargetVersionDiffSuppress(target_version=%q, available=%q, effective=%q) = %t, want %t",
					tc.New, tc.AvailableVersion, tc.EffectiveVersion, got, tc.ShouldSuppress)
			}
		})
	}
}
func TestLustreInstanceVersionUpgradeCustomDiff(t *testing.T) {
	t.Parallel()
	// maintenance_policy is a nested object, so it surfaces as a list. The mock's
	// HasChange compares the two values directly, so use nil to mean "absent".
	policy := []interface{}{map[string]interface{}{}}
	cases := map[string]struct {
		BeforeTargetVersion interface{}
		AfterTargetVersion  interface{}
		BeforeCapacityGib   interface{}
		AfterCapacityGib    interface{}
		BeforeMaintenance   interface{}
		AfterMaintenance    interface{}
		ExpectError         bool
	}{
		"no version change is always allowed": {
			BeforeCapacityGib: "18000",
			AfterCapacityGib:  "27000",
			BeforeMaintenance: nil,
			AfterMaintenance:  policy,
		},
		"version change on its own is allowed": {
			AfterTargetVersion: "2.1.0",
			BeforeCapacityGib:  "18000",
			AfterCapacityGib:   "18000",
		},
		"version change with unchanged capacity is allowed": {
			AfterTargetVersion: "2.1.0",
			BeforeCapacityGib:  "18000",
			AfterCapacityGib:   "18000",
			BeforeMaintenance:  nil,
			AfterMaintenance:   nil,
		},
		"version change with a capacity change is rejected": {
			AfterTargetVersion: "2.1.0",
			BeforeCapacityGib:  "18000",
			AfterCapacityGib:   "27000",
			ExpectError:        true,
		},
		"version change with a maintenance policy change is rejected": {
			AfterTargetVersion: "2.1.0",
			BeforeCapacityGib:  "18000",
			AfterCapacityGib:   "18000",
			BeforeMaintenance:  nil,
			AfterMaintenance:   policy,
			ExpectError:        true,
		},
		"version change with both is rejected": {
			AfterTargetVersion: "2.1.0",
			BeforeCapacityGib:  "18000",
			AfterCapacityGib:   "27000",
			BeforeMaintenance:  nil,
			AfterMaintenance:   policy,
			ExpectError:        true,
		},
	}
	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			d := &tpgresource.ResourceDiffMock{
				Before: map[string]interface{}{
					"target_version":     tc.BeforeTargetVersion,
					"capacity_gib":       tc.BeforeCapacityGib,
					"maintenance_policy": tc.BeforeMaintenance,
				},
				After: map[string]interface{}{
					"target_version":     tc.AfterTargetVersion,
					"capacity_gib":       tc.AfterCapacityGib,
					"maintenance_policy": tc.AfterMaintenance,
				},
			}
			err := lustreInstanceVersionUpgradeCustomDiffFunc(d)
			if tc.ExpectError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Errorf("expected no error but got: %s", err)
			}
		})
	}
}
