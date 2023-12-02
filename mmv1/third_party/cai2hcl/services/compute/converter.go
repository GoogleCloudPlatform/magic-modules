package compute

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
)

var UberConverter = common.UberConverter{
	ConverterByAssetType: map[string]string{
		ComputeInstanceAssetType:       "google_compute_instance",
		ComputeForwardingRuleAssetType: "google_compute_forwarding_rule",
	},
	Converters: common.CreateConverterMap(map[string]common.ConverterFactory{
		"google_compute_instance":        NewComputeInstanceConverter,
		"google_compute_forwarding_rule": NewComputeForwardingRuleConverter,
	}),
}
