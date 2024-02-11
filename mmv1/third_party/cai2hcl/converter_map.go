package cai2hcl

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/services/compute"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg_provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
)

var provider *schema.Provider = tpg_provider.Provider()

// AssetTypeToConverter is a mapping from Asset Type to converter instance.
var AssetTypeToConverter = map[string]string{
	compute.ComputeInstanceAssetType:       "google_compute_instance",
	compute.ComputeForwardingRuleAssetType: "google_compute_forwarding_rule",

	compute.ComputeBackendServiceAssetType: "google_compute_backend_service",

	resourcemanager.ProjectAssetType:        "google_project",
	resourcemanager.ProjectBillingAssetType: "google_project",
}

// ConverterMap is a collection of converters instances, indexed by name.
var ConverterMap = map[string]common.Converter{
	"google_compute_instance":        compute.NewComputeInstanceConverter(provider),
	"google_compute_forwarding_rule": compute.NewComputeForwardingRuleConverter(provider),
	"google_compute_backend_service": compute.NewComputeBackendServiceConverter(provider),

	"google_project": resourcemanager.NewProjectConverter(provider),
}
