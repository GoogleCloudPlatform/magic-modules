package cai2hcl

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/services/compute"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/services/networksecurity"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg_provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
)

var provider *schema.Provider = tpg_provider.Provider()

// AssetTypeToConverter is a mapping from Asset Type to converter instance.
var AssetTypeToConverter = map[string]string{
	compute.ComputeInstanceAssetType:       "google_compute_instance",
	compute.ComputeForwardingRuleAssetType: "google_compute_forwarding_rule",

	compute.ComputeBackendServiceAssetType:       "google_compute_backend_service",
	compute.ComputeRegionBackendServiceAssetType: "google_compute_region_backend_service",

	compute.ComputeRegionHealthCheckAssetType: "google_compute_region_health_check",

	resourcemanager.ProjectAssetType:        "google_project",
	resourcemanager.ProjectBillingAssetType: "google_project",

	networksecurity.ServerTLSPolicyAssetType:             "google_network_security_server_tls_policy",
	networksecurity.BackendAuthenticationConfigAssetType: "google_network_security_backend_authentication_config",
}

// ConverterMap is a collection of converters instances, indexed by name.
var ConverterMap = map[string]common.Converter{
	"google_compute_instance":        compute.NewComputeInstanceConverter(provider),
	"google_compute_forwarding_rule": compute.NewComputeForwardingRuleConverter(provider),

	"google_compute_backend_service":        compute.NewComputeBackendServiceConverter(provider),
	"google_compute_region_backend_service": compute.NewComputeRegionBackendServiceConverter(provider),

	"google_compute_region_health_check": compute.NewComputeRegionHealthCheckConverter(provider),

	"google_project": resourcemanager.NewProjectConverter(provider),

	"google_network_security_server_tls_policy":             networksecurity.NewServerTLSPolicyConverter(provider),
	"google_network_security_backend_authentication_config": networksecurity.NewBackendAuthenticationConfigConverter(provider),
}
