package tpgdclresource

import (
	containeraws "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containeraws"
	containerazure "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containerazure"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Returns the terraform representation of a three-state boolean value represented by a pointer to bool in DCL.
func FlattenEnumBool(v interface{}) string {
	b, ok := v.(*bool)
	if !ok || b == nil {
		return ""
	}
	if *b {
		return "TRUE"
	}
	return "FALSE"
}

func FlattenContainerAwsNodePoolManagement(obj *containeraws.NodePoolManagement, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if obj == nil {
		return nil
	}
	transformed := make(map[string]interface{})

	if obj.AutoRepair == nil || obj.Empty() {
		transformed["auto_repair"] = false
	} else {
		transformed["auto_repair"] = obj.AutoRepair
	}

	return []interface{}{transformed}
}

func FlattenContainerAzureNodePoolManagement(obj *containerazure.NodePoolManagement, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if obj == nil {
		return nil
	}
	transformed := make(map[string]interface{})

	if obj.AutoRepair == nil || obj.Empty() {
		transformed["auto_repair"] = false
	} else {
		transformed["auto_repair"] = obj.AutoRepair
	}

	return []interface{}{transformed}
}
