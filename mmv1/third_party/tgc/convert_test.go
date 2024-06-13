package google

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/ancestrymanager"
	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/services/resourcemanager"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/tfdata"
	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const testProject = "test-project"

type bufferWriteSyncer struct {
	*bytes.Buffer
}

func (bws bufferWriteSyncer) Sync() error {
	return nil
}

// This returns a logger to allow deterministic testing of the output
// by omitting the timestamp and calling function.
func newTestErrorLogger() (*zap.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	syncer := bufferWriteSyncer{Buffer: buf}

	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), syncer, zap.DebugLevel)
	return zap.New(core), syncer.Buffer
}

func newTestConverter(convertUnchanged bool) (*Converter, *bytes.Buffer, error) {
	ctx := context.Background()
	project := testProject
	offline := true
	cfg, err := resources.NewConfig(ctx, project, "", "", offline, "", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("constructing configuration: %w", err)
	}
	errorLogger, buf := newTestErrorLogger()
	c := NewConverter(cfg, &ancestrymanager.NoOpAncestryManager{}, offline, convertUnchanged, errorLogger)

	return c, buf, nil
}

func TestSortByName(t *testing.T) {
	cases := []struct {
		name           string
		unsorted       []caiasset.Asset
		expectedSorted []caiasset.Asset
	}{
		{
			name:           "Empty",
			unsorted:       []caiasset.Asset{},
			expectedSorted: []caiasset.Asset{},
		},
		{
			name: "BCAtoABC",
			unsorted: []caiasset.Asset{
				{
					Name: "b",
					Type: "b-type",
				},
				{
					Name: "c",
					Type: "c-type",
				},
				{
					Name: "a",
					Type: "a-type",
				},
			},
			expectedSorted: []caiasset.Asset{
				{
					Name: "a",
					Type: "a-type",
				},
				{
					Name: "b",
					Type: "b-type",
				},
				{
					Name: "c",
					Type: "c-type",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assets := c.unsorted
			sort.Sort(byName(assets))
			assert.EqualValues(t, c.expectedSorted, assets)
		})
	}
}

func TestAddResourceChanges_nonGoogleResourceIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.aws_api_gateway_account.foo",
		Mode:         "managed",
		Type:         "aws_api_gateway_account",
		Name:         "foo",
		ProviderName: "aws",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"change"},
			Before:  nil,
			After:   nil,
		},
	}
	c, buf, err := newTestConverter(false)
	assert.Nil(t, err)
	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.EqualValues(t, map[string]Asset{}, c.assets)
	assert.Equal(t, "", buf.String())
}

func TestAddResourceChanges_unknownResourceIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.google_unknown.foo",
		Mode:         "managed",
		Type:         "google_unknown",
		Name:         "foo",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"change"},
			Before:  nil,
			After:   nil,
		},
	}
	c, buf, err := newTestConverter(false)
	assert.Nil(t, err)
	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.EqualValues(t, map[string]Asset{}, c.assets)
	assert.Contains(t, buf.String(), "resource type not found")
	assert.Contains(t, buf.String(), rc.Address)
}

func TestAddResourceChanges_unsupportedResourceIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.google_unknown.foo",
		Mode:         "managed",
		Type:         "google_unsupported",
		Name:         "foo",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"change"},
			Before:  nil,
			After:   nil,
		},
	}
	c, buf, err := newTestConverter(false)
	assert.Nil(t, err)

	// fake that this resource is known to the provider; it will never be "supported" by the
	// converter.
	c.schema.ResourcesMap[rc.Type] = c.schema.ResourcesMap["google_compute_disk"]

	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.EqualValues(t, map[string]Asset{}, c.assets)
	assert.Contains(t, buf.String(), "resource type cannot be converted")
	assert.Contains(t, buf.String(), rc.Address)
}

func TestAddResourceChanges_noopIgnoredWhenConvertUnchangedFalse(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.google_compute_disk.foo",
		Mode:         "managed",
		Type:         "google_compute_disk",
		Name:         "foo",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"no-op"},
			Before:  nil,
			After:   nil,
		},
	}
	convertUnchanged := false
	c, buf, err := newTestConverter(convertUnchanged)
	assert.Nil(t, err)

	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.EqualValues(t, map[string]Asset{}, c.assets)
	assert.Equal(t, "", buf.String())
}

func TestAddResourceChanges_deleteProcessed(t *testing.T) {
	cases := []struct {
		name             string
		convertUnchanged bool
	}{
		{
			name:             "Delete when convertUnchanged is false",
			convertUnchanged: false,
		},
		{
			name:             "Delete when convertUnchanged is true",
			convertUnchanged: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rc := tfjson.ResourceChange{
				Address:      "whatever.google_compute_disk.foo",
				Mode:         "managed",
				Type:         "google_compute_disk",
				Name:         "foo",
				ProviderName: "google",
				Change: &tfjson.Change{
					Actions: tfjson.Actions{"delete"},
					Before: map[string]interface{}{
						"project": testProject,
						"name":    "test-disk",
						"type":    "pd-ssd",
						"zone":    "us-central1-a",
						"image":   "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
						"labels": map[string]interface{}{
							"environment": "dev",
						},
						"physical_block_size_bytes": 4096,
					},
					After: nil,
				},
			}
			c, buf, err := newTestConverter(tc.convertUnchanged)
			assert.Nil(t, err)

			err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
			assert.Nil(t, err)
			assert.EqualValues(t, map[string]Asset{}, c.assets)
			assert.Equal(t, "", buf.String())
		})
	}
}

func TestAddResourceChanges_betaResourcesLogged(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address:      "whatever.google_compute_disk.foo",
		Mode:         "managed",
		Type:         "google_compute_disk",
		Name:         "foo",
		ProviderName: "registry.terraform.io/hashicorp/google-beta",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"create"},
			Before:  nil, // Ignore Before because it's unused
			After: map[string]interface{}{
				"project": testProject,
				"name":    "test-disk",
				"type":    "pd-ssd",
				"zone":    "us-central1-a",
				"image":   "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
				"labels": map[string]interface{}{
					"environment": "dev",
				},
				"physical_block_size_bytes": 4096,
			},
		},
	}
	c, buf, err := newTestConverter(false)
	assert.Nil(t, err)

	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)

	caiKey := "compute.googleapis.com/Disk//compute.googleapis.com/projects/test-project/zones/us-central1-a/disks/test-disk"
	assert.Contains(t, c.assets, caiKey)

	assert.Contains(t, buf.String(), "resource uses the google-beta provider")
	assert.Contains(t, buf.String(), rc.Address)
}

func TestAddResourceChanges_createOrUpdateOrDeleteCreateOrNoopProcessed(t *testing.T) {
	cases := []struct {
		name             string
		actions          tfjson.Actions
		convertUnchanged bool
	}{
		{
			name:             "Create when convertUnchanged is false",
			actions:          tfjson.Actions{"create"},
			convertUnchanged: false,
		},
		{
			name:             "Create when convertUnchanged is true",
			actions:          tfjson.Actions{"create"},
			convertUnchanged: true,
		},
		{
			name:             "Update when convertUnchanged is false",
			actions:          tfjson.Actions{"update"},
			convertUnchanged: false,
		},
		{
			name:             "Update when convertUnchanged is true",
			actions:          tfjson.Actions{"update"},
			convertUnchanged: true,
		},
		{
			name:             "DeleteCreate when convertUnchanged is false",
			actions:          tfjson.Actions{"delete", "create"},
			convertUnchanged: false,
		},
		{
			name:             "DeleteCreate when convertUnchanged is true",
			actions:          tfjson.Actions{"delete", "create"},
			convertUnchanged: true,
		},
		{
			name:             "Noop when convertUnchanged is true",
			actions:          tfjson.Actions{"no-op"},
			convertUnchanged: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rc := tfjson.ResourceChange{
				Address:      "whatever.google_compute_disk.foo",
				Mode:         "managed",
				Type:         "google_compute_disk",
				Name:         "foo",
				ProviderName: "google",
				Change: &tfjson.Change{
					Actions: c.actions,
					Before:  nil, // Ignore Before because it's unused
					After: map[string]interface{}{
						"project": testProject,
						"name":    "test-disk",
						"type":    "pd-ssd",
						"zone":    "us-central1-a",
						"image":   "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
						"labels": map[string]interface{}{
							"environment": "dev",
						},
						"physical_block_size_bytes": 4096,
					},
				},
			}
			c, buf, err := newTestConverter(c.convertUnchanged)
			assert.Nil(t, err)

			err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
			assert.Nil(t, err)

			caiKey := "compute.googleapis.com/Disk//compute.googleapis.com/projects/test-project/zones/us-central1-a/disks/test-disk"
			assert.Contains(t, c.assets, caiKey)
			assert.Equal(t, "", buf.String())
			assert.NotContains(t, buf.String(), "resource uses the google-beta provider")
		})
	}
}

func TestAddDuplicatedResources(t *testing.T) {
	rcb1 := tfjson.ResourceChange{
		Address:      "google_billing_budget.budget1",
		Mode:         "managed",
		Type:         "google_billing_budget",
		Name:         "budget1",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"create"},
			Before:  nil,
			After: map[string]interface{}{
				"all_updates_rule": []map[string]interface{}{},
				"amount": []map[string]interface{}{
					{
						"last_period_amount": nil,
						"specified_amount": []map[string]interface{}{
							{
								"currency_code": "USD",
								"nanos":         nil,
								"units":         "100",
							},
						},
					},
				},
				"billing_account": "000000-000000-000000",
				"budget_filter": []map[string]interface{}{
					{
						"credit_types_treatment": "INCLUDE_ALL_CREDITS",
					},
				},
				"display_name": "Example Billing Budget 1",
				"threshold_rules": []map[string]interface{}{
					{
						"spend_basis":       "CURRENT_SPEND",
						"threshold_percent": 0.5,
					},
				},
				"timeouts": nil,
			},
		},
	}
	rcb2 := tfjson.ResourceChange{
		Address:      "google_billing_budget.budget2",
		Mode:         "managed",
		Type:         "google_billing_budget",
		Name:         "budget2",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"create"},
			Before:  nil,
			After: map[string]interface{}{
				"all_updates_rule": []map[string]interface{}{},
				"amount": []map[string]interface{}{
					{
						"last_period_amount": nil,
						"specified_amount": []map[string]interface{}{
							{
								"currency_code": "USD",
								"nanos":         nil,
								"units":         "100",
							},
						},
					},
				},
				"billing_account": "000000-000000-000000",
				"budget_filter": []map[string]interface{}{
					{
						"credit_types_treatment": "INCLUDE_ALL_CREDITS",
					},
				},
				"display_name": "Example Billing Budget 2",
				"threshold_rules": []map[string]interface{}{
					{
						"spend_basis":       "CURRENT_SPEND",
						"threshold_percent": 0.5,
					},
				},
				"timeouts": nil,
			},
		},
	}
	rcp1 := tfjson.ResourceChange{
		Address:      "google_project.my_project1",
		Mode:         "managed",
		Type:         "google_project",
		Name:         "my_project1",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"create"},
			Before:  nil,
			After: map[string]interface{}{
				"auto_create_network": true,
				"billing_account":     "000000-000000-000000",
				"labels":              nil,
				"name":                "My Project1",
				"org_id":              "00000000000000",
				"timeouts":            nil,
			},
		},
	}
	rcp2 := tfjson.ResourceChange{
		Address:      "google_project.my_project2",
		Mode:         "managed",
		Type:         "google_project",
		Name:         "my_project2",
		ProviderName: "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"create"},
			Before:  nil,
			After: map[string]interface{}{
				"auto_create_network": true,
				"billing_account":     "000000-000000-000000",
				"labels":              nil,
				"name":                "My Project2",
				"org_id":              "00000000000000",
				"timeouts":            nil,
			},
		},
	}
	c, buf, err := newTestConverter(false)
	assert.Nil(t, err)

	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rcb1, &rcb2, &rcp1, &rcp2})
	assert.Nil(t, err)

	caiKeyBilling := "cloudbilling.googleapis.com/ProjectBillingInfo//cloudbilling.googleapis.com/projects/test-project/billingInfo"
	assert.Contains(t, c.assets, caiKeyBilling)

	caiKeyProject := "cloudresourcemanager.googleapis.com/Project//cloudresourcemanager.googleapis.com/projects/test-project"
	assert.Contains(t, c.assets, caiKeyProject)

	assert.Contains(t, buf.String(), "duplicate asset")
}

func TestAddStorageModuleAfterUnknown(t *testing.T) {
	var nilValue map[string]interface{} = nil
	rc := tfjson.ResourceChange{
		Address:       "module.gcs_buckets.google_storage_bucket.buckets[0]",
		ModuleAddress: "module.gcs_buckets",
		Mode:          "managed",
		Type:          "google_storage_bucket",
		Name:          "buckets",
		Index:         0,
		ProviderName:  "google",
		Change: &tfjson.Change{
			Actions: tfjson.Actions{"create"},
			Before:  nil,
			After: map[string]interface{}{
				"cors": []interface{}{
					nilValue,
				},
				"default_event_based_hold": nil,
				"encryption": []interface{}{
					nilValue,
				},
				"lifecycle_rule":   []interface{}{},
				"location":         "US",
				"logging":          []interface{}{},
				"project":          "test-project",
				"requester_pays":   nil,
				"retention_policy": []interface{}{},
				"storage_class":    "MULTI_REGIONAL",
				"versioning": []interface{}{
					nilValue,
				},
				"website": []interface{}{
					nilValue,
				},
			},
		},
	}
	c, buf, err := newTestConverter(false)
	assert.Nil(t, err)

	err = c.AddResourceChanges([]*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.Len(t, c.assets, 1)
	for key := range c.assets {
		assert.EqualValues(t, c.assets[key].Type, "storage.googleapis.com/Bucket")
	}

	assert.Equal(t, "", buf.String())
}

func TestTimestampMarshalJSON(t *testing.T) {
	expectedJSON := []byte("\"2021-04-14T15:16:17Z\"")
	date := time.Date(2021, time.April, 14, 15, 16, 17, 0, time.UTC)
	ts := caiasset.Timestamp{
		Seconds: int64(date.Unix()),
		Nanos:   int64(date.UnixNano()),
	}
	json, err := ts.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	assert.EqualValues(t, json, expectedJSON)
}

func TestTimestampUnmarshalJSON(t *testing.T) {
	expectedDate := time.Date(2021, time.April, 14, 15, 16, 17, 0, time.UTC)
	expected := caiasset.Timestamp{
		Seconds: int64(expectedDate.Unix()),
		Nanos:   int64(expectedDate.UnixNano()),
	}
	json := []byte("\"2021-04-14T15:16:17Z\"")
	ts := caiasset.Timestamp{}
	err := ts.UnmarshalJSON(json)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	assert.EqualValues(t, ts, expected)
}

func TestConvertWrapper(t *testing.T) {
	values := map[string]interface{}{
		"name": "test-disk",
	}
	d := tfdata.NewFakeResourceData(
		"google_compute_disk",
		provider.Provider().ResourcesMap["google_compute_disk"].Schema,
		values,
	)

	panicConvertFunc := resources.ConvertFunc(func(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]resources.Asset, error) {
		// should panic
		_ = d.Get("abc").(string)
		return nil, nil
	})

	convertFunc := resources.ConvertFunc(func(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]resources.Asset, error) {
		_ = d.Get("name").(string)
		return nil, nil
	})

	tests := []struct {
		name       string
		converter  resources.ResourceConverter
		rd         tpgresource.TerraformResourceData
		cfg        *transport_tpg.Config
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "field missing",
			converter: resources.ResourceConverter{
				Convert: panicConvertFunc,
			},
			rd:         d,
			wantErr:    true,
			wantErrMsg: "interface conversion",
		},
		{
			name: "field exists",
			converter: resources.ResourceConverter{
				Convert: convertFunc,
			},
			rd: d,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := convertWrapper(test.converter, test.rd, test.cfg)
			if !test.wantErr {
				if err != nil {
					t.Errorf("convertWrapper() = %q, want = nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("convertWrapper() = nil, want = %v", test.wantErrMsg)
				} else if !strings.Contains(err.Error(), test.wantErrMsg) {
					t.Errorf("convertWrapper() = %q, want containing %q", err, test.wantErrMsg)
				}
			}
		})
	}
}

func TestConvertMergingWithExistAsset(t *testing.T) {
	tests := []struct {
		name    string
		changes []*tfjson.ResourceChange
		want    []caiasset.Asset
	}{
		{
			name: "CreateOrUpdateOrNoop",
			changes: []*tfjson.ResourceChange{
				{
					Address:      "google_project_iam_binding.test",
					Mode:         "managed",
					Type:         "google_project_iam_binding",
					Name:         "test",
					ProviderName: "registry.terraform.io/hashicorp/google-beta",
					Change: &tfjson.Change{
						Actions: tfjson.Actions{"create"},
						Before:  nil,
						After: map[string]interface{}{
							"members": []interface{}{
								"user:jane@example.com",
							},
							"project": "example-project",
							"role":    "roles/editor",
						},
					},
				},
			},
			want: []caiasset.Asset{
				{
					Name: "//cloudresourcemanager.googleapis.com/projects/example-project",
					Type: "cloudresourcemanager.googleapis.com/Project",
					IAMPolicy: &caiasset.IAMPolicy{
						Bindings: []caiasset.IAMBinding{
							{
								Role:    "roles/editor",
								Members: []string{"user:jane@example.com"},
							},
						},
					},
				},
			},
		},
		{
			name: "Delete",
			changes: []*tfjson.ResourceChange{
				{
					Address:      "google_project_iam_binding.test",
					Mode:         "managed",
					Type:         "google_project_iam_binding",
					Name:         "test",
					ProviderName: "registry.terraform.io/hashicorp/google-beta",
					Change: &tfjson.Change{
						Actions: tfjson.Actions{"delete"},
						After:   nil,
						Before: map[string]interface{}{
							"members": []interface{}{
								"user:jane@example.com",
							},
							"project": "example-project",
							"role":    "roles/editor",
						},
					},
				},
			},
			want: []caiasset.Asset{
				{
					Name:      "//cloudresourcemanager.googleapis.com/projects/example-project",
					Type:      "cloudresourcemanager.googleapis.com/Project",
					IAMPolicy: &caiasset.IAMPolicy{},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			converter, _, err := newTestConverter(false)
			assert.Nil(t, err)

			fetchFuncCalled := 0
			fetchFunc := func(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
				fetchFuncCalled++
				return resourcemanager.FetchProjectIamPolicy(d, config)
			}

			converter.converters["google_project_iam_binding"] = []cai.ResourceConverter{
				{
					AssetType:         "cloudresourcemanager.googleapis.com/Project",
					Convert:           resourcemanager.GetProjectIamBindingCaiObject,
					FetchFullResource: fetchFunc,
					MergeCreateUpdate: resourcemanager.MergeProjectIamBinding,
					MergeDelete:       resourcemanager.MergeProjectIamBindingDelete,
				},
			}

			// pre-populate the converter's asset cache
			asset, err := converter.augmentAsset(nil, converter.cfg, cai.Asset{
				Name: "//cloudresourcemanager.googleapis.com/projects/example-project",
				Type: "cloudresourcemanager.googleapis.com/Project",
				IAMPolicy: &cai.IAMPolicy{
					Bindings: []cai.IAMBinding{
						{
							Role:    "roles/editor",
							Members: []string{"user:abc@example.com"},
						},
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			converter.assets = map[string]Asset{
				"cloudresourcemanager.googleapis.com/Project//cloudresourcemanager.googleapis.com/projects/example-project": asset,
			}

			err = converter.AddResourceChanges(test.changes)
			if err != nil {
				t.Fatalf("AddResourceChanges() = %s, want = nil", err)
			}
			if diff := cmp.Diff(test.want, converter.Assets()); diff != "" {
				t.Errorf("AddResourceChanges() returned unexpected diff (-want +got):\n%s", diff)
			}
			assert.Equal(t, 0, fetchFuncCalled)
		})
	}
}
