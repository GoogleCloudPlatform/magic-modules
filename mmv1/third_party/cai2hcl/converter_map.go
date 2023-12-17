package cai2hcl

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/services/compute"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg_provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
)

var provider *schema.Provider = tpg_provider.Provider()

var ConverterMap = common.ConverterMap{
	AssetTypeToConverterName: map[string]string{
		compute.ComputeInstanceAssetType:       compute.ComputeInstanceSchemaName,
		compute.ComputeForwardingRuleAssetType: compute.ComputeForwardingRuleSchemaName,

		resourcemanager.ProjectAssetType:        resourcemanager.ProjectSchemaName,
		resourcemanager.ProjectBillingAssetType: resourcemanager.ProjectSchemaName,
	},

	ConverterNameToConverter: map[string]common.Converter{
		compute.ComputeInstanceSchemaName:       compute.NewComputeInstanceConverter(provider),
		compute.ComputeForwardingRuleSchemaName: compute.NewComputeForwardingRuleConverter(provider),

		resourcemanager.ProjectSchemaName: resourcemanager.NewProjectConverter(provider),
	},
}
