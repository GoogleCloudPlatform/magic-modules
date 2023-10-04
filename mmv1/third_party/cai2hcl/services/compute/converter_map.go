package compute

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/common"
)

var ConverterNames = map[string]string{
	ComputeInstanceAssetType:       "google_compute_instance",
	ComputeForwardingRuleAssetType: "google_compute_forwarding_rule",
}

var ConverterMap = common.CreateConverterMap(map[string]common.ConverterFactory{
	"google_compute_instance":        NewComputeInstanceConverter,
	"google_compute_forwarding_rule": NewComputeForwardingRuleConverter,
})
