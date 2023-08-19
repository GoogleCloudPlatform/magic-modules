package generated

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/generated/converters/common"
	computeConverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/generated/converters/compute"
)

var ConverterNames = map[string]string{
	computeConverters.ComputeForwardingRuleAssetType: "google_compute_forwarding_rule",
	computeConverters.ComputeBackendServiceAssetType: "google_compute_backend_service",
	computeConverters.ComputeHealthCheckAssetType:    "google_compute_health_check",
}

var ConverterMap map[string]common.Converter

func init() {
	var factoryMap = map[string]func(name string) common.Converter{
		"google_compute_forwarding_rule": computeConverters.NewComputeForwardingRuleConverter,
		"google_compute_backend_service": computeConverters.NewComputeBackendServiceConverter,
		"google_compute_health_check":    computeConverters.NewComputeHealthCheckConverter,
	}

	ConverterMap = make(map[string]common.Converter, len(factoryMap))
	for name, factory := range factoryMap {
		ConverterMap[name] = factory(name)
	}
}
