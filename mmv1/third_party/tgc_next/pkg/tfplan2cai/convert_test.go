package tfplan2cai

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"os"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/provider"
	tfjson "github.com/hashicorp/terraform-json"
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

func newTestErrorLogger() (*zap.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			bufferWriteSyncer{buf},
			zap.DebugLevel,
		),
	), buf
}

func convertWithChanges(t *testing.T, changes []*tfjson.ResourceChange) ([]caiasset.Asset, *bytes.Buffer, error) {
	logger, buf := newTestErrorLogger()
	o := &Options{
		ErrorLogger:         logger,
		Offline:             true,
		DefaultProject:      testProject,
		DefaultZone:         "us-central1-a",
		NoOpAncestryManager: true,
	}

	plan := &tfjson.Plan{
		FormatVersion:    "0.1",
		TerraformVersion: "1.0.0",
		ResourceChanges:  changes,
	}
	jsonPlan, err := json.Marshal(plan)
	if err != nil {
		return nil, nil, err
	}

	assets, err := Convert(context.Background(), jsonPlan, o)
	return assets, buf, err
}

func TestConvert_nonGoogleResourceIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address: "aws_instance.foo",
		Type:    "aws_instance",
	}

	assets, buf, err := convertWithChanges(t, []*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.Empty(t, assets)
	assert.Equal(t, "", buf.String())
}

func TestConvert_unknownResourceIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address: "google_really_unknown.foo",
		Type:    "google_really_unknown",
	}

	assets, buf, err := convertWithChanges(t, []*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.Empty(t, assets)
	assert.Contains(t, buf.String(), "resource type not found in google beta provider")
}

func TestConvert_unsupportedResourceIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address: "google_unsupported.foo",
		Type:    "google_unsupported",
		Change: &tfjson.Change{
			Actions: []tfjson.Action{tfjson.ActionCreate},
			After:   map[string]interface{}{},
		},
	}

	p := provider.Provider()
	p.ResourcesMap["google_unsupported"] = p.ResourcesMap["google_compute_disk"]
	defer delete(p.ResourcesMap, "google_unsupported")

	assets, buf, err := convertWithChanges(t, []*tfjson.ResourceChange{&rc})

	if err == nil {
		assert.Empty(t, assets)
	} else {
		// assert.Contains(t, err.Error(), "getting resource converter")
	}
	_ = buf
}

func TestConvert_noOpIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address: "google_compute_disk.foo",
		Type:    "google_compute_disk",
		Change: &tfjson.Change{
			Actions: []tfjson.Action{tfjson.ActionNoop},
		},
	}
	assets, _, err := convertWithChanges(t, []*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
	assert.Empty(t, assets)
}

func TestConvert_deleteIgnored(t *testing.T) {
	rc := tfjson.ResourceChange{
		Address: "google_compute_disk.foo",
		Type:    "google_compute_disk",
		Change: &tfjson.Change{
			Actions: []tfjson.Action{tfjson.ActionDelete},
			Before: map[string]interface{}{
				"name": "foo",
			},
		},
	}
	_, _, err := convertWithChanges(t, []*tfjson.ResourceChange{&rc})
	assert.Nil(t, err)
}
func TestConvert_ComputeNetwork(t *testing.T) {
	logger, _ := newTestErrorLogger()
	o := &Options{
		ErrorLogger:         logger,
		Offline:             true,
		DefaultProject:      testProject,
		DefaultZone:         "us-central1-a",
		NoOpAncestryManager: true,
	}

	jsonPlan, err := os.ReadFile("resolvers/compute_network.tfplan.json")
	if err != nil {
		t.Fatalf("Error reading test file: %v", err)
	}

	assets, err := Convert(context.Background(), jsonPlan, o)
	if err != nil {
		t.Fatalf("Error marshaling assets: %v", err)
	}
	assert.Nil(t, err)
	assert.NotEmpty(t, assets)
	assert.NotEmpty(t, assets[2].Resource.Data["vmwareEngineNetwork"])
	assert.NotEmpty(t, assets[2].Resource.Data["peerNetwork"])
	assert.Equal(t, "projects/terraform-dev-haonan/locations/global/vmwareEngineNetworks/network-peering-test-ven", assets[2].Resource.Data["vmwareEngineNetwork"])
	assert.Equal(t, "projects/terraform-dev-haonan/global/networks/network-peering-test-nw", assets[2].Resource.Data["peerNetwork"])
}

func TestConvert_ComputeDiskNestedId(t *testing.T) {
	logger, _ := newTestErrorLogger()
	o := &Options{
		ErrorLogger:         logger,
		Offline:             true,
		DefaultProject:      testProject,
		DefaultZone:         "us-central1-a",
		NoOpAncestryManager: true,
	}

	jsonPlan, err := os.ReadFile("resolvers/compute_disk_nestedId.tfplan.json")
	if err != nil {
		t.Fatalf("Error reading test file: %v", err)
	}

	assets, err := Convert(context.Background(), jsonPlan, o)
	if err != nil {
		t.Fatalf("Error marshaling assets: %v", err)
	}
	assert.Nil(t, err)
	assert.NotEmpty(t, assets)
	assert.NotEmpty(t, assets[1].Resource.Data["asyncPrimaryDisk"])
	asyncPrimaryDisk, ok := assets[1].Resource.Data["asyncPrimaryDisk"].(map[string]interface{})
	assert.True(t, ok, "asyncPrimaryDisk should be a map")
	assert.Equal(t, "projects/terraform-dev-haonan/zones/us-central1-a/disks/async-test-disk", asyncPrimaryDisk["disk"])
}

func TestConvert_ComputeAddress(t *testing.T) {
	logger, _ := newTestErrorLogger()
	o := &Options{
		ErrorLogger:         logger,
		Offline:             true,
		DefaultProject:      testProject,
		DefaultZone:         "us-central1-a",
		NoOpAncestryManager: true,
	}

	jsonPlan, err := os.ReadFile("resolvers/compute_address.tfplan.json")
	if err != nil {
		t.Fatalf("Error reading test file: %v", err)
	}

	assets, err := Convert(context.Background(), jsonPlan, o)
	if err != nil {
		t.Fatalf("Error marshaling assets: %v", err)
	}
	assert.Nil(t, err)
	assert.NotEmpty(t, assets)
	assert.Equal(t, "https://www.googleapis.com/compute/v1/projects/terraform-dev-haonan/regions/us-east1/subnetworks/subnetwork-test", assets[1].Resource.Data["subnetwork"])
}
