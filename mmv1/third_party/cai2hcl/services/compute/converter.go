package compute

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/common/matchers"
)

var forwardingRuleMatchers = []matchers.ConverterMatcher{
	matchers.ByAssetName(ComputeForwardingRuleAssetNameRegex, "google_compute_forwarding_rule"),
	matchers.ByAssetName(ComputeGlobalForwardingRuleAssetNameRegex, "google_compute_global_forwarding_rule"),
}

var backendServiceMatchers = []matchers.ConverterMatcher{
	matchers.ByAssetName(ComputeRegionBackendServiceAssetNameRegex, "google_compute_region_backend_service"),
	matchers.ByAssetName(ComputeBackendServiceAssetNameRegex, "google_compute_backend_service"),
}

var UberConverter = common.UberConverter{
	ConverterByAssetType: map[string]string{
		ComputeInstanceAssetType:    "google_compute_instance",
		ComputeHealthCheckAssetType: "google_compute_health_check",
	},
	ConverterMatchersByAssetType: map[string][]matchers.ConverterMatcher{
		ComputeForwardingRuleAssetType:       forwardingRuleMatchers,
		ComputeGlobalForwardingRuleAssetType: forwardingRuleMatchers,

		ComputeBackendServiceAssetType:       backendServiceMatchers,
		ComputeRegionBackendServiceAssetType: backendServiceMatchers,
	},
	Converters: common.CreateConverterMap(map[string]common.ConverterFactory{
		"google_compute_instance":               NewComputeInstanceConverter,
		"google_compute_health_check":           NewComputeHealthCheckConverter,
		"google_compute_forwarding_rule":        NewComputeForwardingRuleConverter,
		"google_compute_global_forwarding_rule": NewComputeGlobalForwardingRuleConverter,
		"google_compute_region_backend_service": NewComputeRegionBackendServiceConverter,
		"google_compute_backend_service":        NewComputeBackendServiceConverter,
	}),
}
