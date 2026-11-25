// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
package bigquery

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
)

func TestUnitBigQueryDataset_isBackfillError(t *testing.T) {
	t.Parallel()

	retryableErr := &googleapi.Error{
		Code:    400,
		Message: "The secondary replica 'us-east4' has not finished backfill. The backfill must be complete before associating a failover reservation to the dataset.",
	}
	if ok, _ := IsBigQueryDatasetBackfillError(retryableErr); !ok {
		t.Fatalf("expected backfill 400 error to be retryable")
	}

	other400Err := &googleapi.Error{
		Code:    400,
		Message: "Invalid dataset ID",
	}
	if ok, _ := IsBigQueryDatasetBackfillError(other400Err); ok {
		t.Fatalf("expected non-backfill 400 error not to be retryable")
	}

	if ok, _ := IsBigQueryDatasetBackfillError(errors.New("plain error")); ok {
		t.Fatalf("expected non-googleapi error not to be retryable")
	}
}

func TestUnitBigQueryDataset_findPrimaryReplica(t *testing.T) {
	t.Parallel()

	replicas := []interface{}{
		map[string]interface{}{
			"id":           "sec-replica",
			"location":     "us-east4",
			"replicaState": "SECONDARY",
		},
		map[string]interface{}{
			"id":           "pri-replica",
			"location":     "us-central1",
			"replicaState": "PRIMARY",
		},
	}
	id, loc, found := FindBigQueryDatasetPrimaryReplica(replicas, "us-central1")
	if !found || id != "pri-replica" || loc != "us-central1" {
		t.Fatalf("expected primary replica (pri-replica, us-central1), got (%q, %q, %v)", id, loc, found)
	}

	fallbackReplicas := []interface{}{
		map[string]interface{}{
			"id":       "us-primary",
			"location": "US",
		},
	}
	id, loc, found = FindBigQueryDatasetPrimaryReplica(fallbackReplicas, "us")
	if !found || id != "us-primary" || loc != "US" {
		t.Fatalf("expected location fallback replica (us-primary, US), got (%q, %q, %v)", id, loc, found)
	}
}

func TestUnitBigQueryDataset_replicasSetHashCaseInsensitiveLocation(t *testing.T) {
	t.Parallel()

	replicasSchema := ResourceBigQueryDataset().Schema["replicas"]
	if replicasSchema == nil || replicasSchema.Set == nil {
		t.Fatalf("expected replicas schema with custom Set hash function")
	}
	h1 := replicasSchema.Set(map[string]interface{}{
		"id":       "primary",
		"location": "US",
	})
	h2 := replicasSchema.Set(map[string]interface{}{
		"id":       "primary",
		"location": "us",
	})
	if h1 != h2 {
		t.Fatalf("expected case-insensitive location hash equality, got %d != %d", h1, h2)
	}
}

var bigQueryDatasetReplicaOptsCtyType = cty.Object(map[string]cty.Type{
	"default_kms_key_name": cty.String,
})

var bigQueryDatasetReplicaCtyType = cty.Object(map[string]cty.Type{
	"id":              cty.String,
	"location":        cty.String,
	"replica_options": cty.List(bigQueryDatasetReplicaOptsCtyType),
})

func buildBigQueryDatasetRawCty(replicas ...cty.Value) cty.Value {
	var repVal cty.Value
	if len(replicas) == 0 {
		repVal = cty.SetValEmpty(bigQueryDatasetReplicaCtyType)
	} else {
		repVal = cty.SetVal(replicas)
	}
	return cty.ObjectVal(map[string]cty.Value{
		"replicas": repVal,
	})
}

func buildBigQueryDatasetReplicaCty(id, location, kmsKey string) cty.Value {
	var optsVal cty.Value
	if kmsKey == "" {
		optsVal = cty.ListValEmpty(bigQueryDatasetReplicaOptsCtyType)
	} else {
		optsVal = cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"default_kms_key_name": cty.StringVal(kmsKey),
			}),
		})
	}
	return cty.ObjectVal(map[string]cty.Value{
		"id":              cty.StringVal(id),
		"location":        cty.StringVal(location),
		"replica_options": optsVal,
	})
}

func TestUnitBigQueryDataset_extractReplicasFromCty(t *testing.T) {
	t.Parallel()

	t.Run("populated with mixed casing and optional replica_options", func(t *testing.T) {
		t.Parallel()
		raw := buildBigQueryDatasetRawCty(
			buildBigQueryDatasetReplicaCty("pri", "US", ""),
			buildBigQueryDatasetReplicaCty("sec", "US-EAST4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
		)
		got, hasUnknown, ok := extractBigQueryDatasetReplicasFromCty(raw)
		if !ok || hasUnknown {
			t.Fatalf("expected ok=true, hasUnknown=false, got ok=%v, hasUnknown=%v", ok, hasUnknown)
		}
		want := map[bigQueryDatasetReplicaIdentity]bigQueryDatasetReplicaConfig{
			{id: "pri", location: "us"}: {
				id:       "pri",
				location: "US",
				kmsKey:   "",
				known:    true,
			},
			{id: "sec", location: "us-east4"}: {
				id:       "sec",
				location: "US-EAST4",
				kmsKey:   "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2",
				known:    true,
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("extractBigQueryDatasetReplicasFromCty mismatch:\ngot:  %#v\nwant: %#v", got, want)
		}
	})

	t.Run("null replicas returns empty map", func(t *testing.T) {
		t.Parallel()
		raw := cty.ObjectVal(map[string]cty.Value{
			"replicas": cty.NullVal(cty.Set(bigQueryDatasetReplicaCtyType)),
		})
		got, hasUnknown, ok := extractBigQueryDatasetReplicasFromCty(raw)
		if !ok || hasUnknown || len(got) != 0 {
			t.Fatalf("expected empty map with ok=true, hasUnknown=false; got %#v, hasUnknown=%v, ok=%v", got, hasUnknown, ok)
		}
	})

	t.Run("unknown replicas attribute reports hasUnknown", func(t *testing.T) {
		t.Parallel()
		raw := cty.ObjectVal(map[string]cty.Value{
			"replicas": cty.UnknownVal(cty.Set(bigQueryDatasetReplicaCtyType)),
		})
		_, hasUnknown, ok := extractBigQueryDatasetReplicasFromCty(raw)
		if !ok || !hasUnknown {
			t.Fatalf("expected ok=true, hasUnknown=true, got ok=%v, hasUnknown=%v", ok, hasUnknown)
		}
	})

	t.Run("unknown default_kms_key_name marks replica known=false", func(t *testing.T) {
		t.Parallel()
		raw := buildBigQueryDatasetRawCty(
			cty.ObjectVal(map[string]cty.Value{
				"id":       cty.StringVal("sec"),
				"location": cty.StringVal("us-east4"),
				"replica_options": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"default_kms_key_name": cty.UnknownVal(cty.String),
					}),
				}),
			}),
		)
		got, hasUnknown, ok := extractBigQueryDatasetReplicasFromCty(raw)
		if !ok || !hasUnknown {
			t.Fatalf("expected ok=true, hasUnknown=true, got ok=%v, hasUnknown=%v", ok, hasUnknown)
		}
		info, exists := got[bigQueryDatasetReplicaIdentity{id: "sec", location: "us-east4"}]
		if !exists || info.known {
			t.Fatalf("expected replica entry with known=false, got exists=%v, info=%#v", exists, info)
		}
	})

	t.Run("marked cty values are unmarked cleanly", func(t *testing.T) {
		t.Parallel()
		markedKms := cty.StringVal("projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2").Mark("sensitive")
		elem := cty.ObjectVal(map[string]cty.Value{
			"id":       cty.StringVal("sec").Mark("id-mark"),
			"location": cty.StringVal("us-east4").Mark("loc-mark"),
			"replica_options": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"default_kms_key_name": markedKms,
				}).Mark("opt-elem-mark"),
			}).Mark("opts-mark"),
		})
		raw := cty.ObjectVal(map[string]cty.Value{
			"replicas": cty.ListVal([]cty.Value{elem.Mark("elem-mark")}).Mark("replicas-mark"),
		}).Mark("outer")
		got, hasUnknown, ok := extractBigQueryDatasetReplicasFromCty(raw)
		if !ok || hasUnknown {
			t.Fatalf("expected ok=true, hasUnknown=false for marked value, got ok=%v, hasUnknown=%v", ok, hasUnknown)
		}
		info := got[bigQueryDatasetReplicaIdentity{id: "sec", location: "us-east4"}]
		if info.kmsKey != "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2" || !info.known {
			t.Fatalf("unexpected unmarked replica info: %#v", info)
		}
	})
}

type bigQueryDatasetDiffMock struct {
	*tpgresource.ResourceDiffMock
	id             string
	newComputed    map[string]bool
	rawState       cty.Value
	rawPlan        cty.Value
	changedKeys    []string
	forcedNewKeys  []string
	useRawCty      bool
	useChangedKeys bool
}

func (m *bigQueryDatasetDiffMock) Id() string {
	return m.id
}

func (m *bigQueryDatasetDiffMock) SetNewComputed(key string) error {
	if m.newComputed == nil {
		m.newComputed = map[string]bool{}
	}
	m.newComputed[key] = true
	return nil
}

func (m *bigQueryDatasetDiffMock) ForceNew(key string) error {
	m.forcedNewKeys = append(m.forcedNewKeys, key)
	return m.ResourceDiffMock.ForceNew(key)
}

type bigQueryDatasetRawCtyDiffMock struct {
	*bigQueryDatasetDiffMock
}

func (m *bigQueryDatasetRawCtyDiffMock) GetRawState() cty.Value {
	return m.rawState
}

func (m *bigQueryDatasetRawCtyDiffMock) GetRawPlan() cty.Value {
	return m.rawPlan
}

func (m *bigQueryDatasetRawCtyDiffMock) GetChangedKeysPrefix(prefix string) []string {
	var matched []string
	for _, k := range m.changedKeys {
		if strings.HasPrefix(k, prefix) {
			matched = append(matched, k)
		}
	}
	return matched
}

func TestUnitBigQueryDataset_replicasCustomDiff(t *testing.T) {
	t.Parallel()

	replicasSchema := ResourceBigQueryDataset().Schema["replicas"]
	if replicasSchema == nil || replicasSchema.Set == nil {
		t.Fatalf("expected replicas schema with custom Set hash function")
	}

	makeReplica := func(id, location, kmsKey string) map[string]interface{} {
		r := map[string]interface{}{
			"id":       id,
			"location": location,
		}
		if kmsKey != "" {
			r["replica_options"] = []interface{}{
				map[string]interface{}{
					"default_kms_key_name": kmsKey,
				},
			}
		}
		return r
	}
	makeSet := func(items ...map[string]interface{}) *schema.Set {
		raw := make([]interface{}, len(items))
		for i, item := range items {
			raw[i] = item
		}
		return schema.NewSet(replicasSchema.Set, raw)
	}

	// (a) Adding a brand-new secondary CMEK replica does NOT force replacement.
	addReplicaDiff := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				makeReplica("us-east4", "us-east4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
			),
		},
		After: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				makeReplica("us-east4", "us-east4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
				makeReplica("us-west1", "us-west1", "projects/p/locations/us-west1/keyRings/kr/cryptoKeys/k3"),
			),
		},
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(addReplicaDiff); err != nil {
		t.Fatalf("unexpected error when adding CMEK replica: %v", err)
	}
	if addReplicaDiff.IsForceNew {
		t.Fatalf("expected IsForceNew=false when adding a new CMEK replica, got true")
	}

	// (b) Dropping an existing secondary CMEK replica does NOT force replacement.
	dropReplicaDiff := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				makeReplica("us-east4", "us-east4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
				makeReplica("us-west1", "us-west1", "projects/p/locations/us-west1/keyRings/kr/cryptoKeys/k3"),
			),
		},
		After: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				makeReplica("us-east4", "us-east4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
			),
		},
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(dropReplicaDiff); err != nil {
		t.Fatalf("unexpected error when dropping secondary CMEK replica: %v", err)
	}
	if dropReplicaDiff.IsForceNew {
		t.Fatalf("expected IsForceNew=false when dropping a secondary CMEK replica, got true")
	}

	// (c) Populating write-only default_kms_key_name from "" -> key (e.g. after terraform import
	// or adding replica_options to an existing primary replica) does NOT force replacement
	// on either the fallback *schema.Set path or the raw cty path.
	populateKmsFallbackDiff := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", ""),
				makeReplica("us-east4", "us-east4", ""),
			),
		},
		After: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				makeReplica("us-east4", "us-east4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
			),
		},
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(populateKmsFallbackDiff); err != nil {
		t.Fatalf("unexpected error when populating write-only KMS key (fallback path): %v", err)
	}
	if populateKmsFallbackDiff.IsForceNew {
		t.Fatalf("expected IsForceNew=false when populating write-only KMS key from empty state (fallback path), got true")
	}

	populateKmsCtyDiff := &bigQueryDatasetRawCtyDiffMock{
		bigQueryDatasetDiffMock: &bigQueryDatasetDiffMock{
			ResourceDiffMock: &tpgresource.ResourceDiffMock{},
			id:               "projects/proj/datasets/ds",
			rawState: buildBigQueryDatasetRawCty(
				buildBigQueryDatasetReplicaCty("us-central1", "us-central1", ""),
				buildBigQueryDatasetReplicaCty("us-east4", "us-east4", ""),
			),
			rawPlan: buildBigQueryDatasetRawCty(
				buildBigQueryDatasetReplicaCty("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				buildBigQueryDatasetReplicaCty("us-east4", "us-east4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
			),
		},
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(populateKmsCtyDiff); err != nil {
		t.Fatalf("unexpected error when populating write-only KMS key (raw cty path): %v", err)
	}
	if populateKmsCtyDiff.IsForceNew || len(populateKmsCtyDiff.forcedNewKeys) != 0 {
		t.Fatalf("expected no ForceNew when populating write-only KMS key from empty state (raw cty path), got IsForceNew=%v keys=%v",
			populateKmsCtyDiff.IsForceNew, populateKmsCtyDiff.forcedNewKeys)
	}

	// (d) Mutating a non-empty default_kms_key_name on an existing (id, location) replica
	// triggers ForceNew on the fallback path.
	mutateKmsDiff := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				makeReplica("us-east4", "us-east4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
			),
		},
		After: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				makeReplica("us-east4", "US-EAST4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2-new"),
			),
		},
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(mutateKmsDiff); err != nil {
		t.Fatalf("unexpected error when mutating existing replica KMS key: %v", err)
	}
	if !mutateKmsDiff.IsForceNew {
		t.Fatalf("expected IsForceNew=true when changing default_kms_key_name on an existing replica, got false")
	}

	// (e) Mutating a non-empty default_kms_key_name via raw cty with unchanged set length
	// (2 -> 2) exercises extractBigQueryDatasetReplicasFromCty and the GetChangedKeysPrefix
	// leaf-key ForceNew branch in forceNewBigQueryDatasetReplicas.
	mutateKmsCtyDiff := &bigQueryDatasetRawCtyDiffMock{
		bigQueryDatasetDiffMock: &bigQueryDatasetDiffMock{
			ResourceDiffMock: &tpgresource.ResourceDiffMock{
				Before: map[string]interface{}{
					"replicas":              "old",
					"replicas.222.id":       "us-east4",
					"replicas.222.location": "us-east4",
					"replicas.222.replica_options.0.default_kms_key_name": "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2",
					"replicas.111.id": "",
				},
				After: map[string]interface{}{
					"replicas":              "new",
					"replicas.222.id":       "",
					"replicas.222.location": "",
					"replicas.222.replica_options.0.default_kms_key_name": "",
					"replicas.111.id": "us-east4",
				},
			},
			id: "projects/proj/datasets/ds",
			rawState: buildBigQueryDatasetRawCty(
				buildBigQueryDatasetReplicaCty("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				buildBigQueryDatasetReplicaCty("us-east4", "us-east4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2"),
			),
			rawPlan: buildBigQueryDatasetRawCty(
				buildBigQueryDatasetReplicaCty("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
				buildBigQueryDatasetReplicaCty("us-east4", "US-EAST4", "projects/p/locations/us-east4/keyRings/kr/cryptoKeys/k2-new"),
			),
			changedKeys: []string{
				"replicas.222.replica_options.0.default_kms_key_name",
				"replicas.222.location",
				"replicas.222.id",
				"replicas.111.id",
			},
		},
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(mutateKmsCtyDiff); err != nil {
		t.Fatalf("unexpected error on raw cty KMS mutation diff: %v", err)
	}
	if !mutateKmsCtyDiff.IsForceNew {
		t.Fatalf("expected IsForceNew=true on raw cty KMS mutation diff")
	}
	wantForcedKeys := []string{"replicas.111.id", "replicas"}
	if !reflect.DeepEqual(mutateKmsCtyDiff.forcedNewKeys, wantForcedKeys) {
		t.Fatalf("expected ForceNew keys %v (deterministic sorted leaf key + replicas), got %v",
			wantForcedKeys, mutateKmsCtyDiff.forcedNewKeys)
	}

	// Removing a non-empty default_kms_key_name on an existing replica also triggers ForceNew.
	removeKmsDiff := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", "projects/p/locations/us-central1/keyRings/kr/cryptoKeys/k1"),
			),
		},
		After: map[string]interface{}{
			"replicas": makeSet(
				makeReplica("us-central1", "us-central1", ""),
			),
		},
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(removeKmsDiff); err != nil {
		t.Fatalf("unexpected error when removing existing replica KMS key: %v", err)
	}
	if !removeKmsDiff.IsForceNew {
		t.Fatalf("expected IsForceNew=true when removing default_kms_key_name on an existing replica, got false")
	}

	// (f) Changing failover_reservation_id without changing replicas marks replicas computed via SetNewComputed("replicas").
	unchangedSet := makeSet(
		makeReplica("us-central1", "us-central1", ""),
		makeReplica("us-east4", "us-east4", ""),
	)
	failoverDiff := &bigQueryDatasetDiffMock{
		ResourceDiffMock: &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"failover_reservation_id": "proj.res1",
				"replicas":                unchangedSet,
			},
			After: map[string]interface{}{
				"failover_reservation_id": "proj.res2",
				"replicas":                unchangedSet,
			},
		},
		id: "projects/proj/datasets/ds",
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(failoverDiff); err != nil {
		t.Fatalf("unexpected error when changing failover_reservation_id: %v", err)
	}
	if failoverDiff.IsForceNew {
		t.Fatalf("expected IsForceNew=false when changing failover_reservation_id, got true")
	}
	if !failoverDiff.newComputed["replicas"] {
		t.Fatalf("expected SetNewComputed(\"replicas\") when failover_reservation_id changes without replicas changing")
	}

	// Create (Id() == "") is a no-op.
	createDiff := &bigQueryDatasetDiffMock{
		ResourceDiffMock: &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{},
			After: map[string]interface{}{
				"failover_reservation_id": "proj.res1",
				"replicas":                unchangedSet,
			},
		},
		id: "",
	}
	if err := ResourceBigQueryDatasetReplicasCustomDiffFunc(createDiff); err != nil {
		t.Fatalf("unexpected error on create diff: %v", err)
	}
	if createDiff.IsForceNew || createDiff.newComputed["replicas"] {
		t.Fatalf("expected no ForceNew or SetNewComputed on create diff")
	}
}

func TestUnitBigQueryDataset_cleanupReplicationBeforeDelete(t *testing.T) {
	t.Parallel()

	newResourceData := func(t *testing.T, state map[string]interface{}) *schema.ResourceData {
		t.Helper()
		d := schema.TestResourceDataRaw(t, ResourceBigQueryDataset().Schema, state)
		d.SetId("projects/test-proj/datasets/test_ds")
		return d
	}

	t.Run("non-replicated dataset skips GET and PATCH entirely", func(t *testing.T) {
		t.Parallel()
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		config := &transport_tpg.Config{
			Client: srv.Client(),
			CustomEndpoints: map[string]string{
				Product.CustomEndpointField: srv.URL + "/bigquery/v2/",
			},
		}

		for _, state := range []map[string]interface{}{
			{
				"project":    "test-proj",
				"dataset_id": "test_ds",
				"location":   "us-central1",
			},
			{
				"project":    "test-proj",
				"dataset_id": "test_ds",
				"location":   "us-central1",
				"replicas": []interface{}{
					map[string]interface{}{
						"id":       "us-central1",
						"location": "us-central1",
					},
				},
			},
		} {
			d := newResourceData(t, state)
			if err := cleanupBigQueryDatasetReplicationBeforeDelete(d, config, "test-proj", "test-ua"); err != nil {
				t.Fatalf("unexpected error for non-replicated dataset: %v", err)
			}
		}
		if requests != 0 {
			t.Fatalf("expected 0 HTTP requests for non-replicated datasets, got %d", requests)
		}
	})

	for _, contentCase := range []struct {
		name         string
		endpointPath string
		responseJSON string
		wantErrPart  string
	}{
		{
			name:         "tables present blocks cleanup when delete_contents_on_destroy=false",
			endpointPath: "/bigquery/v2/projects/test-proj/datasets/test_ds/tables",
			responseJSON: `{"tables":[{"id":"test-proj:test_ds.t1"}],"totalItems":1}`,
			wantErrPart:  "contains tables",
		},
		{
			name:         "routines present blocks cleanup when delete_contents_on_destroy=false",
			endpointPath: "/bigquery/v2/projects/test-proj/datasets/test_ds/routines",
			responseJSON: `{"routines":[{"routineReference":{"routineId":"r1"}}]}`,
			wantErrPart:  "contains routines",
		},
		{
			name:         "models present blocks cleanup when delete_contents_on_destroy=false",
			endpointPath: "/bigquery/v2/projects/test-proj/datasets/test_ds/models",
			responseJSON: `{"models":[{"modelReference":{"modelId":"m1"}}]}`,
			wantErrPart:  "contains models",
		},
	} {
		contentCase := contentCase
		t.Run(contentCase.name, func(t *testing.T) {
			t.Parallel()
			patchCalls := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch {
					patchCalls++
				}
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/bigquery/v2/projects/test-proj/datasets/test_ds":
					_, _ = w.Write([]byte(`{
						"location": "us-central1",
						"failoverReservationId": "test-proj.res1",
						"replicas": [
							{"id": "us-central1", "location": "us-central1", "replicaState": "PRIMARY"},
							{"id": "us-east4", "location": "us-east4", "replicaState": "SECONDARY"}
						]
					}`))
				case contentCase.endpointPath:
					_, _ = w.Write([]byte(contentCase.responseJSON))
				default:
					_, _ = w.Write([]byte(`{}`))
				}
			}))
			defer srv.Close()

			config := &transport_tpg.Config{
				Client: srv.Client(),
				CustomEndpoints: map[string]string{
					Product.CustomEndpointField: srv.URL + "/bigquery/v2/",
				},
			}
			d := newResourceData(t, map[string]interface{}{
				"project":                    "test-proj",
				"dataset_id":                 "test_ds",
				"location":                   "us-central1",
				"delete_contents_on_destroy": false,
				"failover_reservation_id":    "test-proj.res1",
				"replicas": []interface{}{
					map[string]interface{}{"id": "us-central1", "location": "us-central1"},
					map[string]interface{}{"id": "us-east4", "location": "us-east4"},
				},
			})
			err := cleanupBigQueryDatasetReplicationBeforeDelete(d, config, "test-proj", "test-ua")
			if err == nil || !strings.Contains(err.Error(), contentCase.wantErrPart) {
				t.Fatalf("expected error containing %q, got %v", contentCase.wantErrPart, err)
			}
			if patchCalls != 0 {
				t.Fatalf("expected 0 PATCH requests when dataset has contents and delete_contents_on_destroy=false, got %d", patchCalls)
			}
		})
	}

	for _, tc := range []struct {
		name             string
		deleteContents   bool
		wantListChecks   int
		wantPatchActions []string
	}{
		{
			name:             "empty replicated dataset with delete_contents_on_destroy=false verifies contents then cleans up",
			deleteContents:   false,
			wantListChecks:   3,
			wantPatchActions: []string{"disassociate", "drop_secondary"},
		},
		{
			name:             "replicated dataset with delete_contents_on_destroy=true skips content checks and cleans up",
			deleteContents:   true,
			wantListChecks:   0,
			wantPatchActions: []string{"disassociate", "drop_secondary"},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			listChecks := 0
			var patchActions []string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodGet {
					switch r.URL.Path {
					case "/bigquery/v2/projects/test-proj/datasets/test_ds":
						_, _ = w.Write([]byte(`{
							"location": "us-central1",
							"failoverReservationId": "test-proj.res1",
							"replicas": [
								{"id": "us-central1", "location": "us-central1", "replicaState": "PRIMARY"},
								{"id": "us-east4", "location": "us-east4", "replicaState": "SECONDARY"}
							]
						}`))
					case "/bigquery/v2/projects/test-proj/datasets/test_ds/tables",
						"/bigquery/v2/projects/test-proj/datasets/test_ds/routines",
						"/bigquery/v2/projects/test-proj/datasets/test_ds/models":
						listChecks++
						_, _ = w.Write([]byte(`{}`))
					default:
						w.WriteHeader(http.StatusNotFound)
					}
					return
				}
				if r.Method == http.MethodPatch {
					var body map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("failed to decode PATCH body: %v", err)
					}
					if _, ok := body["failoverReservationId"]; ok {
						patchActions = append(patchActions, "disassociate")
					}
					if reps, ok := body["replicas"].([]interface{}); ok && len(reps) == 1 {
						patchActions = append(patchActions, "drop_secondary")
					}
					_, _ = w.Write([]byte(`{}`))
					return
				}
				w.WriteHeader(http.StatusMethodNotAllowed)
			}))
			defer srv.Close()

			config := &transport_tpg.Config{
				Client: srv.Client(),
				CustomEndpoints: map[string]string{
					Product.CustomEndpointField: srv.URL + "/bigquery/v2/",
				},
			}
			d := newResourceData(t, map[string]interface{}{
				"project":                    "test-proj",
				"dataset_id":                 "test_ds",
				"location":                   "us-central1",
				"delete_contents_on_destroy": tc.deleteContents,
				"failover_reservation_id":    "test-proj.res1",
				"replicas": []interface{}{
					map[string]interface{}{"id": "us-central1", "location": "us-central1"},
					map[string]interface{}{"id": "us-east4", "location": "us-east4"},
				},
			})
			if err := cleanupBigQueryDatasetReplicationBeforeDelete(d, config, "test-proj", "test-ua"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if listChecks != tc.wantListChecks {
				t.Fatalf("expected %d content list checks, got %d", tc.wantListChecks, listChecks)
			}
			if !reflect.DeepEqual(patchActions, tc.wantPatchActions) {
				t.Fatalf("expected PATCH actions %v, got %v", tc.wantPatchActions, patchActions)
			}
		})
	}
}
