package cai2hcl

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/google/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/google/converters/common"
	computeConverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/google/converters/compute"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta"
)

var ConverterNames = map[string]string{
	converters.ProjectAssetType:                      "google_project",
	converters.ProjectBillingAssetType:               "google_project",
	computeConverters.ComputeInstanceAssetType:       "google_compute_instance",
	computeConverters.ComputeForwardingRuleAssetType: "google_compute_forwarding_rule",
}

var converterFactories = map[string]func(name string, schema map[string]*schema.Schema) common.Converter{
	"google_project":                 converters.NewProjectConverter,
	"google_compute_instance":        computeConverters.NewComputeInstanceConverter,
	"google_compute_forwarding_rule": computeConverters.NewComputeForwardingRuleConverter,
}

var ConverterMap map[string]common.Converter

func init() {
	tpgProvider := tpg.Provider()

	ConverterMap = make(map[string]common.Converter, len(converterFactories))
	for name, factory := range converterFactories {
		ConverterMap[name] = factory(name, tpgProvider.ResourcesMap[name].Schema)
	}
}
