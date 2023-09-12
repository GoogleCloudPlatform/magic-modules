package cai2hcl

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/converters/google/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/converters/google/resources/compute"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta"
)

var ConverterNames = map[string]string{
	resources.ProjectAssetType:             "google_project",
	resources.ProjectBillingAssetType:      "google_project",
	compute.ComputeInstanceAssetType:       "google_compute_instance",
	compute.ComputeForwardingRuleAssetType: "google_compute_forwarding_rule",
}

var converterFactories = map[string]func(name string, schema map[string]*schema.Schema) common.Converter{
	"google_project":                 resources.NewProjectConverter,
	"google_compute_instance":        compute.NewComputeInstanceConverter,
	"google_compute_forwarding_rule": compute.NewComputeForwardingRuleConverter,
}

var ConverterMap map[string]common.Converter

func init() {
	tpgProvider := tpg.Provider()

	ConverterMap = make(map[string]common.Converter, len(converterFactories))
	for name, factory := range converterFactories {
		ConverterMap[name] = factory(name, tpgProvider.ResourcesMap[name].Schema)
	}
}
