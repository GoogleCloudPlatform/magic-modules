package toolkit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/tfplan"
	legacytfplan2cai "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai"
	legacy "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/converters/google/resources"
	tfjson "github.com/hashicorp/terraform-json"
)

// IsSupported returns true if the resource is migrated or supported in legacy TGC
// resource is expected to be a Terraform resource type (e.g. "google_compute_instance")
func IsSupported(resource string, migratedResourceMap map[string]bool) bool {
	if migratedResourceMap[resource] {
		return true
	}
	_, tgcExists := legacy.ResourceConverters()[resource]
	return tgcExists
}

// Convert converts a Terraform plan to CAI assets.
// It splits resources into migrated (pkg) and legacy (tfplan2cai) buckets based on the provided map.
// It uses the appropriate converter for each bucket and merges the results, normalizing legacy assets to the format in TGC pkg.
func Convert(ctx context.Context, jsonPlan []byte, o *tfplan2cai.Options, migratedResourceMap map[string]bool) ([]caiasset.Asset, error) {
	changes, err := tfplan.ReadResourceChanges(jsonPlan)
	if err != nil {
		return nil, err
	}

	var migratedChanges, legacyChanges []*tfjson.ResourceChange
	for _, change := range changes {
		if migratedResourceMap[change.Type] {
			migratedChanges = append(migratedChanges, change)
		} else {
			legacyChanges = append(legacyChanges, change)
		}
	}

	migratedAssets, err := tfplan2cai.ConvertChanges(ctx, jsonPlan, migratedChanges, o)
	if err != nil {
		return nil, fmt.Errorf("converting migrated resources: %w", err)
	}

	legacyOptions := &legacytfplan2cai.Options{
		ErrorLogger:    o.ErrorLogger,
		Offline:        o.Offline,
		DefaultProject: o.DefaultProject,
		DefaultRegion:  o.DefaultRegion,
		DefaultZone:    o.DefaultZone,
		UserAgent:      o.UserAgent,
		HTTPClient:     o.HTTPClient,
		AncestryCache:  o.AncestryCache,
	}

	legacyAssets, err := legacytfplan2cai.ConvertChanges(ctx, legacyChanges, legacyOptions)
	if err != nil {
		return nil, fmt.Errorf("converting legacy resources: %w", err)
	}

	// Convert legacy assets to TGC format
	convertedLegacyAssets := make([]caiasset.Asset, 0, len(legacyAssets))
	for _, legacyAsset := range legacyAssets {
		marshaledAsset, err := json.Marshal(legacyAsset)
		if err != nil {
			return nil, fmt.Errorf("marshaling legacy asset: %w", err)
		}
		var asset caiasset.Asset
		if err := json.Unmarshal(marshaledAsset, &asset); err != nil {
			return nil, fmt.Errorf("unmarshaling legacy asset: %w", err)
		}
		convertedLegacyAssets = append(convertedLegacyAssets, asset)
	}

	return append(migratedAssets, convertedLegacyAssets...), nil
}
