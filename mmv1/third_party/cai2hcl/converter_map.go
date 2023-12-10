package cai2hcl

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/services/compute"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/services/resourcemanager"
)

var ConverterMap = common.ConverterMap{
	AssetTypeToConverter: map[string]string{
		// Compute
		compute.ComputeInstanceAssetType:       "google_compute_instance",
		compute.ComputeForwardingRuleAssetType: "google_compute_forwarding_rule",

		// ResourceManager
		resourcemanager.ProjectAssetType:        "google_project",
		resourcemanager.ProjectBillingAssetType: "google_project",
	},
	Converters: common.CreateConverters(map[string]common.ConverterFactory{
		// Compute
		"google_compute_instance":        compute.NewComputeInstanceConverter,
		"google_compute_forwarding_rule": compute.NewComputeForwardingRuleConverter,

		// ResourceManager
		"google_project": resourcemanager.NewProjectConverter,
	}),
}
