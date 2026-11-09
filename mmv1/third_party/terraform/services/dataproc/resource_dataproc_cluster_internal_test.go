package dataproc_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/services/dataproc"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestDataprocExtractInitTimeout(t *testing.T) {
	t.Parallel()

	actual, err := dataproc.ExtractInitTimeout("500s")
	expected := 500
	if err != nil {
		t.Fatalf("Expected %d, but got error %v", expected, err)
	}
	if actual != expected {
		t.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

func TestDataprocExtractInitTimeout_nonSeconds(t *testing.T) {
	t.Parallel()

	actual, err := dataproc.ExtractInitTimeout("5m")
	expected := 300
	if err != nil {
		t.Fatalf("Expected %d, but got error %v", expected, err)
	}
	if actual != expected {
		t.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

func TestDataprocExtractInitTimeout_empty(t *testing.T) {
	t.Parallel()

	_, err := dataproc.ExtractInitTimeout("")
	expected := "time: invalid duration"
	if err != nil && err.Error() != expected {
		return
	}
	t.Fatalf("Expected an error with message '%s', but got %v", expected, err.Error())
}

func TestDataprocParseImageVersion(t *testing.T) {
	t.Parallel()

	testCases := map[string]dataproc.DataprocImageVersion{
		"1.2":             {"1", "2", "", ""},
		"1.2.3":           {"1", "2", "3", ""},
		"1.2.3rc":         {"1", "2", "3rc", ""},
		"1.2-debian9":     {"1", "2", "", "debian9"},
		"1.2.3-debian9":   {"1", "2", "3", "debian9"},
		"1.2.3rc-debian9": {"1", "2", "3rc", "debian9"},
	}

	for v, expected := range testCases {
		actual, err := dataproc.ParseDataprocImageVersion(v)
		if actual.Major != expected.Major {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if err != nil {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.Minor != expected.Minor {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.Subminor != expected.Subminor {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.OsName != expected.OsName {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
	}

	errorTestCases := []string{
		"",
		"1",
		"notaversion",
		"1-debian",
	}
	for _, v := range errorTestCases {
		if _, err := dataproc.ParseDataprocImageVersion(v); err == nil {
			t.Errorf("expected parsing invalid version %q to return error", v)
		}
	}
}

func TestDataprocDiffSuppress(t *testing.T) {
	t.Parallel()

	doSuppress := [][]string{
		{"1.3.10-debian9", "1.3"},
		{"1.3.10-debian9", "1.3-debian9"},
		{"1.3.10", "1.3"},
		{"1.3-debian9", "1.3"},
		{"1.3.10-debian9", "1.3.prodprevious-debian9"},
		{"1.3.10-debian9", "1.3.prodcurrent-debian9"},
	}

	noSuppress := [][]string{
		{"1.3.10-debian9", "1.3.10-ubuntu"},
		{"1.3.10-debian9", "1.3.9-debian9"},
		{"1.3.10-debian9", "1.3-ubuntu"},
		{"1.3.10-debian9", "1.3.9"},
		{"1.3.10-debian9", "1.4"},
		{"1.3.10-debian9", "2.3"},
		{"1.3.10", "1.3.10-debian9"},
		{"1.3", "1.3.10"},
		{"1.3", "1.3.10-debian9"},
		{"1.3", "1.3-debian9"},
		{"1.3.prodprevious-debian9", "1.3.10-debian9"},
		{"1.3.prodcurrent-debian9", "1.3.10-debian9"},
		{"1.3.10-debian9", "1.3.randomstring-debian9"},
	}

	for _, tup := range doSuppress {
		if !dataproc.DataprocImageVersionDiffSuppress("", tup[0], tup[1], nil) {
			t.Errorf("expected (old: %q, new: %q) to be suppressed", tup[0], tup[1])
		}
	}
	for _, tup := range noSuppress {
		if dataproc.DataprocImageVersionDiffSuppress("", tup[0], tup[1], nil) {
			t.Errorf("expected (old: %q, new: %q) to not be suppressed", tup[0], tup[1])
		}
	}
}

func TestDataprocSecondaryWorkerFlexibleConfigDiff(t *testing.T) {
	t.Parallel()

	res := dataproc.ResourceDataprocCluster()

	flexType := cty.Object(map[string]cty.Type{
		"num_instances":     cty.Number,
		"core_count":        cty.Number,
		"hardware_class":    cty.String,
		"boot_disk_size_gb": cty.Number,
	})
	autoType := cty.Object(map[string]cty.Type{
		"policy_uri": cty.String,
	})

	makeRawConfigWorkerNum := func(sec []cty.Value, auto []cty.Value, workerNum cty.Value) cty.Value {
		secList := cty.ListValEmpty(flexType)
		if len(sec) > 0 {
			secList = cty.ListVal(sec)
		}
		autoList := cty.ListValEmpty(autoType)
		if len(auto) > 0 {
			autoList = cty.ListVal(auto)
		}
		return cty.ObjectVal(map[string]cty.Value{
			"name":            cty.StringVal("test-cluster"),
			"project":         cty.StringVal("p"),
			"region":          cty.StringVal("us-central1"),
			"deletion_policy": cty.NullVal(cty.String),
			"labels":          cty.NullVal(cty.Map(cty.String)),
			"cluster_config": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"master_flexible_config": cty.ListVal([]cty.Value{
						cty.ObjectVal(map[string]cty.Value{
							"num_instances":     cty.NumberIntVal(1),
							"core_count":        cty.NumberIntVal(4),
							"hardware_class":    cty.StringVal("STANDARD"),
							"boot_disk_size_gb": cty.NumberIntVal(100),
						}),
					}),
					"worker_flexible_config": cty.ListVal([]cty.Value{
						cty.ObjectVal(map[string]cty.Value{
							"num_instances":     workerNum,
							"core_count":        cty.NumberIntVal(8),
							"hardware_class":    cty.StringVal("STANDARD"),
							"boot_disk_size_gb": cty.NumberIntVal(100),
						}),
					}),
					"secondary_worker_flexible_config": secList,
					"autoscaling_config":               autoList,
				}),
			}),
		})
	}
	makeRawConfig := func(sec []cty.Value, auto []cty.Value) cty.Value {
		return makeRawConfigWorkerNum(sec, auto, cty.NumberIntVal(2))
	}

	baseStateAttrs := func() map[string]string {
		return map[string]string{
			"id":                                    "projects/p/regions/us-central1/clusters/test-cluster",
			"name":                                  "test-cluster",
			"project":                               "p",
			"region":                                "us-central1",
			"cluster_config.#":                      "1",
			"cluster_config.0.temp_bucket":          "dataproc-temp-bucket",
			"cluster_config.0.bucket":               "dataproc-staging-bucket",
			"cluster_config.0.cluster_type":         "STANDARD",
			"cluster_config.0.gce_cluster_config.#": "1",
			"cluster_config.0.gce_cluster_config.0.zone":                                                   "us-central1-a",
			"cluster_config.0.gce_cluster_config.0.network":                                                "https://www.googleapis.com/compute/v1/projects/p/global/networks/default",
			"cluster_config.0.gce_cluster_config.0.internal_ip_only":                                       "true",
			"cluster_config.0.gce_cluster_config.0.shielded_instance_config.#":                             "1",
			"cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_integrity_monitoring": "true",
			"cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_secure_boot":          "true",
			"cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_vtpm":                 "true",
			"cluster_config.0.master_config.#":                                                             "1",
			"cluster_config.0.master_config.0.num_instances":                                               "1",
			"cluster_config.0.master_config.0.machine_type":                                                "n4-standard-2",
			"cluster_config.0.master_config.0.min_cpu_platform":                                            "Automatic",
			"cluster_config.0.master_config.0.image_uri":                                                   "https://www.googleapis.com/compute/v1/projects/cloud-dataproc/global/images/dataproc-2-3-deb12",
			"cluster_config.0.master_config.0.disk_config.#":                                               "1",
			"cluster_config.0.master_config.0.disk_config.0.boot_disk_size_gb":                             "100",
			"cluster_config.0.master_config.0.disk_config.0.boot_disk_type":                                "hyperdisk-balanced",
			"cluster_config.0.master_config.0.disk_config.0.num_local_ssds":                                "0",
			"cluster_config.0.master_config.0.disk_config.0.local_ssd_interface":                           "SCSI",
			"cluster_config.0.master_config.0.instance_flexibility_policy.#":                               "1",
			"cluster_config.0.worker_config.#":                                                             "1",
			"cluster_config.0.worker_config.0.num_instances":                                               "2",
			"cluster_config.0.worker_config.0.machine_type":                                                "n4-standard-4",
			"cluster_config.0.worker_config.0.min_cpu_platform":                                            "Automatic",
			"cluster_config.0.worker_config.0.image_uri":                                                   "https://www.googleapis.com/compute/v1/projects/cloud-dataproc/global/images/dataproc-2-3-deb12",
			"cluster_config.0.worker_config.0.disk_config.#":                                               "1",
			"cluster_config.0.worker_config.0.disk_config.0.boot_disk_size_gb":                             "100",
			"cluster_config.0.worker_config.0.disk_config.0.boot_disk_type":                                "hyperdisk-balanced",
			"cluster_config.0.worker_config.0.disk_config.0.num_local_ssds":                                "0",
			"cluster_config.0.worker_config.0.disk_config.0.local_ssd_interface":                           "SCSI",
			"cluster_config.0.worker_config.0.instance_flexibility_policy.#":                               "1",
			"cluster_config.0.preemptible_worker_config.#":                                                 "1",
			"cluster_config.0.preemptible_worker_config.0.num_instances":                                   "2",
			"cluster_config.0.preemptible_worker_config.0.preemptibility":                                  "SPOT",
			"cluster_config.0.preemptible_worker_config.0.disk_config.#":                                   "1",
			"cluster_config.0.preemptible_worker_config.0.disk_config.0.boot_disk_size_gb":                 "100",
			"cluster_config.0.preemptible_worker_config.0.disk_config.0.boot_disk_type":                    "hyperdisk-balanced",
			"cluster_config.0.preemptible_worker_config.0.disk_config.0.num_local_ssds":                    "0",
			"cluster_config.0.preemptible_worker_config.0.disk_config.0.local_ssd_interface":               "SCSI",
			"cluster_config.0.preemptible_worker_config.0.instance_flexibility_policy.#":                   "1",
			"cluster_config.0.software_config.#":                                                           "1",
			"cluster_config.0.software_config.0.image_version":                                             "2.3.19-debian12",
			"cluster_config.0.endpoint_config.#":                                                           "1",
			"cluster_config.0.endpoint_config.0.enable_http_port_access":                                   "false",
			"cluster_config.0.master_flexible_config.#":                                                    "1",
			"cluster_config.0.master_flexible_config.0.num_instances":                                      "1",
			"cluster_config.0.master_flexible_config.0.core_count":                                         "4",
			"cluster_config.0.master_flexible_config.0.hardware_class":                                     "STANDARD",
			"cluster_config.0.master_flexible_config.0.boot_disk_size_gb":                                  "100",
			"cluster_config.0.worker_flexible_config.#":                                                    "1",
			"cluster_config.0.worker_flexible_config.0.num_instances":                                      "2",
			"cluster_config.0.worker_flexible_config.0.core_count":                                         "8",
			"cluster_config.0.worker_flexible_config.0.hardware_class":                                     "STANDARD",
			"cluster_config.0.worker_flexible_config.0.boot_disk_size_gb":                                  "100",
		}
	}

	baseConfigMap := func(sec []interface{}, auto []interface{}, includeWorker bool) map[string]interface{} {
		cc := map[string]interface{}{
			"master_flexible_config": []interface{}{
				map[string]interface{}{
					"num_instances":     1,
					"core_count":        4,
					"hardware_class":    "STANDARD",
					"boot_disk_size_gb": 100,
				},
			},
		}
		if includeWorker {
			cc["worker_flexible_config"] = []interface{}{
				map[string]interface{}{
					"num_instances":     2,
					"core_count":        8,
					"hardware_class":    "STANDARD",
					"boot_disk_size_gb": 100,
				},
			}
		}
		if sec != nil {
			cc["secondary_worker_flexible_config"] = sec
		}
		if auto != nil {
			cc["autoscaling_config"] = auto
		}
		return map[string]interface{}{
			"name":           "test-cluster",
			"project":        "p",
			"region":         "us-central1",
			"cluster_config": []interface{}{cc},
		}
	}

	runDiff := func(t *testing.T, stateAttrs map[string]string, cfgMap map[string]interface{}, rawCfg cty.Value) *terraform.InstanceDiff {
		t.Helper()
		state := &terraform.InstanceState{
			ID:         stateAttrs["id"],
			Attributes: stateAttrs,
			RawConfig:  rawCfg,
			RawPlan:    rawCfg,
		}
		rc := terraform.NewResourceConfigRaw(cfgMap)
		d, err := res.Diff(context.Background(), state, rc, &transport_tpg.Config{Project: "p", Region: "us-central1"})
		if err != nil {
			t.Fatalf("Diff returned unexpected error: %v", err)
		}
		return d
	}

	hasSecDiff := func(d *terraform.InstanceDiff) bool {
		if d == nil {
			return false
		}
		for k := range d.Attributes {
			if strings.HasPrefix(k, "cluster_config.0.secondary_worker_flexible_config") {
				return true
			}
		}
		return false
	}

	t.Run("post_create_add_matching_worker_updates_in_place", func(t *testing.T) {
		attrs := baseStateAttrs()
		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":     2,
				"core_count":        8,
				"hardware_class":    "STANDARD",
				"boot_disk_size_gb": 100,
			},
		}
		rawCfg := makeRawConfig([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.NumberIntVal(2),
				"core_count":        cty.NumberIntVal(8),
				"hardware_class":    cty.StringVal("STANDARD"),
				"boot_disk_size_gb": cty.NumberIntVal(100),
			}),
		}, nil)
		d := runDiff(t, attrs, baseConfigMap(secCfg, nil, true), rawCfg)
		if !hasSecDiff(d) {
			t.Fatalf("expected secondary_worker_flexible_config diff, got none")
		}
		if d.RequiresNew() {
			t.Fatalf("expected in-place update when secondary_worker_flexible_config matches worker_flexible_config, but RequiresNew=true: %#v", d.Attributes)
		}
	})

	t.Run("post_create_add_mismatched_shape_forces_new", func(t *testing.T) {
		attrs := baseStateAttrs()
		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":     2,
				"core_count":        8,
				"hardware_class":    "PERFORMANCE",
				"boot_disk_size_gb": 100,
			},
		}
		rawCfg := makeRawConfig([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.NumberIntVal(2),
				"core_count":        cty.NumberIntVal(8),
				"hardware_class":    cty.StringVal("PERFORMANCE"),
				"boot_disk_size_gb": cty.NumberIntVal(100),
			}),
		}, nil)
		d := runDiff(t, attrs, baseConfigMap(secCfg, nil, true), rawCfg)
		if !d.RequiresNew() {
			t.Fatalf("expected RequiresNew=true when hardware_class differs from worker_flexible_config")
		}
	})

	t.Run("post_create_add_zero_instances_forces_new", func(t *testing.T) {
		attrs := baseStateAttrs()
		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":     0,
				"core_count":        8,
				"hardware_class":    "STANDARD",
				"boot_disk_size_gb": 100,
			},
		}
		rawCfg := makeRawConfig([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.NumberIntVal(0),
				"core_count":        cty.NumberIntVal(8),
				"hardware_class":    cty.StringVal("STANDARD"),
				"boot_disk_size_gb": cty.NumberIntVal(100),
			}),
		}, nil)
		d := runDiff(t, attrs, baseConfigMap(secCfg, nil, true), rawCfg)
		if !d.RequiresNew() {
			t.Fatalf("expected RequiresNew=true when num_instances == 0 on initial post-create addition")
		}
	})

	t.Run("post_create_add_unknown_num_instances_updates_in_place", func(t *testing.T) {
		attrs := baseStateAttrs()
		// hcl2shim.UnknownVariableValue: how the legacy ResourceConfig represents an
		// unknown value (e.g. num_instances from a resource not yet created).
		const unknown = "74D93920-ED26-11E3-AC10-0800200C9A66"
		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":     unknown,
				"core_count":        8,
				"hardware_class":    "STANDARD",
				"boot_disk_size_gb": 100,
			},
		}
		rawCfg := makeRawConfig([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.UnknownVal(cty.Number),
				"core_count":        cty.NumberIntVal(8),
				"hardware_class":    cty.StringVal("STANDARD"),
				"boot_disk_size_gb": cty.NumberIntVal(100),
			}),
		}, nil)
		d := runDiff(t, attrs, baseConfigMap(secCfg, nil, true), rawCfg)
		if !hasSecDiff(d) {
			t.Fatalf("expected secondary_worker_flexible_config diff, got none")
		}
		if d.RequiresNew() {
			t.Fatalf("expected in-place update when num_instances is unknown on initial post-create addition, got RequiresNew=true: %#v", d.Attributes)
		}
	})

	t.Run("post_create_add_omitted_boot_disk_with_worker_set_updates_in_place", func(t *testing.T) {
		attrs := baseStateAttrs()
		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":  2,
				"core_count":     8,
				"hardware_class": "STANDARD",
			},
		}
		rawCfg := makeRawConfig([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.NumberIntVal(2),
				"core_count":        cty.NumberIntVal(8),
				"hardware_class":    cty.StringVal("STANDARD"),
				"boot_disk_size_gb": cty.NullVal(cty.Number),
			}),
		}, nil)
		d := runDiff(t, attrs, baseConfigMap(secCfg, nil, true), rawCfg)
		if !hasSecDiff(d) {
			t.Fatalf("expected secondary_worker_flexible_config diff, got none")
		}
		if d.RequiresNew() {
			t.Fatalf("expected in-place update when boot_disk_size_gb is omitted and worker_flexible_config sets it, got RequiresNew=true: %#v", d.Attributes)
		}
	})

	t.Run("post_create_add_mismatched_boot_disk_forces_new", func(t *testing.T) {
		attrs := baseStateAttrs()
		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":     2,
				"core_count":        8,
				"hardware_class":    "STANDARD",
				"boot_disk_size_gb": 200,
			},
		}
		rawCfg := makeRawConfig([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.NumberIntVal(2),
				"core_count":        cty.NumberIntVal(8),
				"hardware_class":    cty.StringVal("STANDARD"),
				"boot_disk_size_gb": cty.NumberIntVal(200),
			}),
		}, nil)
		d := runDiff(t, attrs, baseConfigMap(secCfg, nil, true), rawCfg)
		if !d.RequiresNew() {
			t.Fatalf("expected RequiresNew=true when boot_disk_size_gb is set and differs from worker_flexible_config")
		}
	})

	t.Run("post_create_add_without_worker_flexible_config_forces_new", func(t *testing.T) {
		attrs := baseStateAttrs()
		delete(attrs, "cluster_config.0.worker_flexible_config.#")
		delete(attrs, "cluster_config.0.worker_flexible_config.0.num_instances")
		delete(attrs, "cluster_config.0.worker_flexible_config.0.core_count")
		delete(attrs, "cluster_config.0.worker_flexible_config.0.hardware_class")
		delete(attrs, "cluster_config.0.worker_flexible_config.0.boot_disk_size_gb")
		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":     2,
				"core_count":        8,
				"hardware_class":    "STANDARD",
				"boot_disk_size_gb": 100,
			},
		}
		rawCfg := makeRawConfig([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.NumberIntVal(2),
				"core_count":        cty.NumberIntVal(8),
				"hardware_class":    cty.StringVal("STANDARD"),
				"boot_disk_size_gb": cty.NumberIntVal(100),
			}),
		}, nil)
		d := runDiff(t, attrs, baseConfigMap(secCfg, nil, false), rawCfg)
		if !d.RequiresNew() {
			t.Fatalf("expected RequiresNew=true when worker_flexible_config is absent")
		}
	})

	t.Run("existing_block_shape_change_forces_new", func(t *testing.T) {
		attrs := baseStateAttrs()
		attrs["cluster_config.0.secondary_worker_flexible_config.#"] = "1"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.num_instances"] = "2"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.core_count"] = "8"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.hardware_class"] = "PERFORMANCE"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.boot_disk_size_gb"] = "100"

		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":     2,
				"core_count":        16,
				"hardware_class":    "PERFORMANCE",
				"boot_disk_size_gb": 100,
			},
		}
		rawCfg := makeRawConfig([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.NumberIntVal(2),
				"core_count":        cty.NumberIntVal(16),
				"hardware_class":    cty.StringVal("PERFORMANCE"),
				"boot_disk_size_gb": cty.NumberIntVal(100),
			}),
		}, nil)
		d := runDiff(t, attrs, baseConfigMap(secCfg, nil, true), rawCfg)
		if !d.RequiresNew() {
			t.Fatalf("expected RequiresNew=true when core_count changes on existing secondary_worker_flexible_config")
		}
	})

	t.Run("block_removed_with_positive_instances_scales_to_zero_in_place", func(t *testing.T) {
		attrs := baseStateAttrs()
		attrs["cluster_config.0.secondary_worker_flexible_config.#"] = "1"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.num_instances"] = "2"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.core_count"] = "8"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.hardware_class"] = "PERFORMANCE"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.boot_disk_size_gb"] = "100"

		rawCfg := makeRawConfig(nil, nil)
		d := runDiff(t, attrs, baseConfigMap(nil, nil, true), rawCfg)
		if !hasSecDiff(d) {
			t.Fatalf("expected secondary_worker_flexible_config diff when block is removed with num_instances=2")
		}
		if d.RequiresNew() {
			t.Fatalf("expected in-place update when removing secondary_worker_flexible_config, got RequiresNew=true")
		}
		numDiff := d.Attributes["cluster_config.0.secondary_worker_flexible_config.0.num_instances"]
		if numDiff == nil || numDiff.Old != "2" || numDiff.New != "0" {
			t.Fatalf("expected secondary_worker_flexible_config.0.num_instances diff 2 -> 0, got %#v", numDiff)
		}
	})

	// Legacy config map with worker_flexible_config.0.num_instances unknown
	// (hcl2shim.UnknownVariableValue).
	withUnknownWorkerNum := func(cfg map[string]interface{}) map[string]interface{} {
		cc := cfg["cluster_config"].([]interface{})[0].(map[string]interface{})
		wrk := cc["worker_flexible_config"].([]interface{})[0].(map[string]interface{})
		wrk["num_instances"] = "74D93920-ED26-11E3-AC10-0800200C9A66"
		return cfg
	}
	secStateAttrs := func() map[string]string {
		attrs := baseStateAttrs()
		attrs["cluster_config.0.secondary_worker_flexible_config.#"] = "1"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.num_instances"] = "2"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.core_count"] = "8"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.hardware_class"] = "STANDARD"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.boot_disk_size_gb"] = "100"
		return attrs
	}

	t.Run("block_removed_with_unknown_sibling_errors", func(t *testing.T) {
		rawCfg := makeRawConfigWorkerNum(nil, nil, cty.UnknownVal(cty.Number))
		state := &terraform.InstanceState{ID: "projects/p/regions/us-central1/clusters/test-cluster", Attributes: secStateAttrs(), RawConfig: rawCfg, RawPlan: rawCfg}
		cfg := withUnknownWorkerNum(baseConfigMap(nil, nil, true))
		_, err := res.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), &transport_tpg.Config{Project: "p", Region: "us-central1"})
		if err == nil {
			t.Fatalf("expected a plan-time error when the block is removed alongside an unknown cluster_config value")
		}
		for _, want := range []string{
			"cluster_config.secondary_worker_flexible_config was removed in the same change as other cluster_config values that are unknown until apply",
			"cluster_config.0.worker_flexible_config.0.num_instances",
			"keep the block with num_instances = 0",
		} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("expected error to contain %q, got: %v", want, err)
			}
		}
	})

	t.Run("block_present_with_unknown_sibling_no_error", func(t *testing.T) {
		secCfg := []interface{}{
			map[string]interface{}{
				"num_instances":     0,
				"core_count":        8,
				"hardware_class":    "STANDARD",
				"boot_disk_size_gb": 100,
			},
		}
		rawCfg := makeRawConfigWorkerNum([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"num_instances":     cty.NumberIntVal(0),
				"core_count":        cty.NumberIntVal(8),
				"hardware_class":    cty.StringVal("STANDARD"),
				"boot_disk_size_gb": cty.NumberIntVal(100),
			}),
		}, nil, cty.UnknownVal(cty.Number))
		state := &terraform.InstanceState{ID: "projects/p/regions/us-central1/clusters/test-cluster", Attributes: secStateAttrs(), RawConfig: rawCfg, RawPlan: rawCfg}
		cfg := withUnknownWorkerNum(baseConfigMap(secCfg, nil, true))
		d, err := res.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), &transport_tpg.Config{Project: "p", Region: "us-central1"})
		if err != nil {
			t.Fatalf("expected no error when the block is present (explicit num_instances = 0), got: %v", err)
		}
		if d.RequiresNew() {
			t.Fatalf("expected in-place update, got RequiresNew=true: %#v", d.Attributes)
		}
		numDiff := d.Attributes["cluster_config.0.secondary_worker_flexible_config.0.num_instances"]
		if numDiff == nil || numDiff.Old != "2" || numDiff.New != "0" {
			t.Fatalf("expected secondary_worker_flexible_config.0.num_instances diff 2 -> 0, got %#v", numDiff)
		}
		if wrk := d.Attributes["cluster_config.0.worker_flexible_config.0.num_instances"]; wrk == nil || !wrk.NewComputed {
			t.Fatalf("expected worker_flexible_config.0.num_instances to stay unknown, got %#v", wrk)
		}
	})

	t.Run("block_removed_when_already_zero_is_suppressed", func(t *testing.T) {
		attrs := baseStateAttrs()
		attrs["cluster_config.0.secondary_worker_flexible_config.#"] = "1"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.num_instances"] = "0"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.core_count"] = "8"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.hardware_class"] = "PERFORMANCE"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.boot_disk_size_gb"] = "100"

		rawCfg := makeRawConfig(nil, nil)
		d := runDiff(t, attrs, baseConfigMap(nil, nil, true), rawCfg)
		if hasSecDiff(d) {
			t.Fatalf("expected no secondary_worker_flexible_config diff when old.num_instances == 0, got %#v", d.Attributes)
		}
	})

	t.Run("block_omitted_with_autoscaling_configured_is_suppressed", func(t *testing.T) {
		attrs := baseStateAttrs()
		attrs["cluster_config.0.secondary_worker_flexible_config.#"] = "1"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.num_instances"] = "2"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.core_count"] = "8"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.hardware_class"] = "STANDARD"
		attrs["cluster_config.0.secondary_worker_flexible_config.0.boot_disk_size_gb"] = "100"
		attrs["cluster_config.0.autoscaling_config.#"] = "1"
		attrs["cluster_config.0.autoscaling_config.0.policy_uri"] = "projects/p/regions/us-central1/autoscalingPolicies/pol1"

		autoCfg := []interface{}{
			map[string]interface{}{
				"policy_uri": "projects/p/regions/us-central1/autoscalingPolicies/pol1",
			},
		}
		rawCfg := makeRawConfig(nil, []cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"policy_uri": cty.StringVal("projects/p/regions/us-central1/autoscalingPolicies/pol1"),
			}),
		})
		d := runDiff(t, attrs, baseConfigMap(nil, autoCfg, true), rawCfg)
		if hasSecDiff(d) {
			t.Fatalf("expected no secondary_worker_flexible_config diff when autoscaling_config is configured, got %#v", d.Attributes)
		}
	})
}

func TestDataprocSecondaryWorkerFlexibleConfigScaleDownKeepsDiffSuppress(t *testing.T) {
	t.Parallel()

	res := dataproc.ResourceDataprocCluster()

	flexType := cty.Object(map[string]cty.Type{
		"num_instances":     cty.Number,
		"core_count":        cty.Number,
		"hardware_class":    cty.String,
		"boot_disk_size_gb": cty.Number,
	})
	rawCfg := cty.ObjectVal(map[string]cty.Value{
		"name":            cty.StringVal("test-cluster"),
		"project":         cty.StringVal("p"),
		"region":          cty.StringVal("us-central1"),
		"deletion_policy": cty.NullVal(cty.String),
		"labels":          cty.NullVal(cty.Map(cty.String)),
		"cluster_config": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"secondary_worker_flexible_config": cty.ListValEmpty(flexType),
				"autoscaling_config":               cty.ListValEmpty(cty.Object(map[string]cty.Type{"policy_uri": cty.String})),
				"lifecycle_config": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{"auto_delete_time": cty.StringVal("2030-01-01T00:00:00Z")}),
				}),
			}),
		}),
	})

	// Complete state for every block the test touches, as written by Read.
	attrs := map[string]string{
		"id":                                  "projects/p/regions/us-central1/clusters/test-cluster",
		"name":                                "test-cluster",
		"project":                             "p",
		"region":                              "us-central1",
		"cluster_config.#":                    "1",
		"cluster_config.0.lifecycle_config.#": "1",
		"cluster_config.0.lifecycle_config.0.idle_delete_ttl": "",
		"cluster_config.0.lifecycle_config.0.idle_start_time": "",
		// Same instant as the config value, in the API's format; suppressed by TimestampDiffSuppress.
		"cluster_config.0.lifecycle_config.0.auto_delete_time": "2030-01-01T00:00:00.000Z",
		"cluster_config.0.lifecycle_config.0.idle_stop_ttl":    "",
		"cluster_config.0.lifecycle_config.0.auto_stop_time":   "",
	}
	for _, blk := range []struct {
		name string
		num  string
		core string
	}{
		{"master_flexible_config", "1", "4"},
		{"worker_flexible_config", "2", "8"},
		{"secondary_worker_flexible_config", "2", "8"},
	} {
		p := "cluster_config.0." + blk.name
		attrs[p+".#"] = "1"
		attrs[p+".0.num_instances"] = blk.num
		attrs[p+".0.core_count"] = blk.core
		attrs[p+".0.hardware_class"] = "STANDARD"
		attrs[p+".0.boot_disk_size_gb"] = "100"
	}

	flex := func(num, core int) []interface{} {
		return []interface{}{map[string]interface{}{
			"num_instances":     num,
			"core_count":        core,
			"hardware_class":    "STANDARD",
			"boot_disk_size_gb": 100,
		}}
	}
	cfg := map[string]interface{}{
		"name":    "test-cluster",
		"project": "p",
		"region":  "us-central1",
		"cluster_config": []interface{}{map[string]interface{}{
			"master_flexible_config": flex(1, 4),
			"worker_flexible_config": flex(2, 8),
			"lifecycle_config": []interface{}{map[string]interface{}{
				"auto_delete_time": "2030-01-01T00:00:00Z",
			}},
		}},
	}

	state := &terraform.InstanceState{ID: attrs["id"], Attributes: attrs, RawConfig: rawCfg, RawPlan: rawCfg}
	d, err := res.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), &transport_tpg.Config{Project: "p", Region: "us-central1"})
	if err != nil {
		t.Fatalf("Diff returned unexpected error: %v", err)
	}
	if d == nil {
		t.Fatalf("expected a scale-down diff, got none")
	}
	if d.RequiresNew() {
		t.Fatalf("expected in-place update, got RequiresNew=true: %#v", d.Attributes)
	}
	numDiff := d.Attributes["cluster_config.0.secondary_worker_flexible_config.0.num_instances"]
	if numDiff == nil || numDiff.Old != "2" || numDiff.New != "0" {
		t.Fatalf("expected secondary_worker_flexible_config.0.num_instances diff 2 -> 0, got %#v", numDiff)
	}
	for k, v := range d.Attributes {
		if strings.HasPrefix(k, "cluster_config.0.lifecycle_config") {
			t.Errorf("expected lifecycle_config diff to stay suppressed during scale-down, got %s: %#v", k, v)
		}
	}

	t.Run("changed_block_preserves_nested_diff_suppress", func(t *testing.T) {
		nestedRawCfg := cty.ObjectVal(map[string]cty.Value{
			"name":            cty.StringVal("test-cluster"),
			"project":         cty.StringVal("p"),
			"region":          cty.StringVal("us-central1"),
			"deletion_policy": cty.NullVal(cty.String),
			"labels":          cty.NullVal(cty.Map(cty.String)),
			"cluster_config": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"secondary_worker_flexible_config": cty.ListValEmpty(flexType),
					"autoscaling_config":               cty.ListValEmpty(cty.Object(map[string]cty.Type{"policy_uri": cty.String})),
					"lifecycle_config": cty.ListVal([]cty.Value{
						cty.ObjectVal(map[string]cty.Value{
							"idle_delete_ttl":  cty.StringVal("1200s"),
							"auto_delete_time": cty.StringVal("2030-01-01T00:00:00Z"),
						}),
					}),
				}),
			}),
		})

		nestedAttrs := make(map[string]string, len(attrs))
		for k, v := range attrs {
			nestedAttrs[k] = v
		}
		nestedAttrs["cluster_config.0.lifecycle_config.0.idle_delete_ttl"] = "600s"

		nestedCfg := map[string]interface{}{
			"name":    "test-cluster",
			"project": "p",
			"region":  "us-central1",
			"cluster_config": []interface{}{map[string]interface{}{
				"master_flexible_config": flex(1, 4),
				"worker_flexible_config": flex(2, 8),
				"lifecycle_config": []interface{}{map[string]interface{}{
					"idle_delete_ttl":  "1200s",
					"auto_delete_time": "2030-01-01T00:00:00Z",
				}},
			}},
		}

		nestedState := &terraform.InstanceState{ID: nestedAttrs["id"], Attributes: nestedAttrs, RawConfig: nestedRawCfg, RawPlan: nestedRawCfg}
		nd, err := res.Diff(context.Background(), nestedState, terraform.NewResourceConfigRaw(nestedCfg), &transport_tpg.Config{Project: "p", Region: "us-central1"})
		if err != nil {
			t.Fatalf("Diff returned unexpected error: %v", err)
		}
		if nd == nil {
			t.Fatalf("expected a diff, got none")
		}
		if nd.RequiresNew() {
			t.Fatalf("expected in-place update, got RequiresNew=true: %#v", nd.Attributes)
		}
		secDiff := nd.Attributes["cluster_config.0.secondary_worker_flexible_config.0.num_instances"]
		if secDiff == nil || secDiff.Old != "2" || secDiff.New != "0" {
			t.Fatalf("expected secondary_worker_flexible_config.0.num_instances diff 2 -> 0, got %#v", secDiff)
		}
		ttlDiff := nd.Attributes["cluster_config.0.lifecycle_config.0.idle_delete_ttl"]
		if ttlDiff == nil || ttlDiff.Old != "600s" || ttlDiff.New != "1200s" {
			t.Fatalf("expected lifecycle_config.0.idle_delete_ttl diff 600s -> 1200s, got %#v", ttlDiff)
		}
		if autoDiff, ok := nd.Attributes["cluster_config.0.lifecycle_config.0.auto_delete_time"]; ok {
			t.Errorf("expected lifecycle_config.0.auto_delete_time diff to stay suppressed during scale-down, got %#v", autoDiff)
		}
	})
}
