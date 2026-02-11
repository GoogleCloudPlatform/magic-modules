package bigtable

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var BigtableTableIamStateUpgraders = []schema.StateUpgrader{
	{
		Type:    resourceBigtableTableIAMV0().CoreConfigSchema().ImpliedType(),
		Upgrade: ResourceBigtableTableIAMUpgradeV0,
		Version: 0,
	},
}

func resourceBigtableTableIAMV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"table": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func ResourceBigtableTableIAMUpgradeV0(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)

	if _, ok := rawState["instance"]; ok {
		rawState["instance_name"] = rawState["instance"]
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}
